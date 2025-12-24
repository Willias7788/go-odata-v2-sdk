package client

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	CSRFHeader = "X-CSRF-Token"
	CSRFValue  = "Fetch"
)

type SAPClient struct {
	client      *resty.Client
	baseURL     string
	csrfToken   string
	csrfCookies []*http.Cookie
	mu          sync.RWMutex
}

// NewSAPClient initializes the Resty client with basic auth and defaults
func NewSAPClient(baseURL, username, password string) *SAPClient {
	r := resty.New()
	r.SetBaseURL(baseURL)
	r.SetBasicAuth(username, password)
	
	// Set default timeouts and headers
	r.SetTimeout(time.Second * 30)
	r.SetHeader("Accept", "application/json")
	r.SetHeader("Content-Type", "application/json")

	return &SAPClient{
		client:  r,
		baseURL: baseURL,
	}
}

// SetDebug enables resty debug mode
func (s *SAPClient) SetDebug(debug bool) {
	s.client.SetDebug(debug)
}

// GetClient returns the underlying resty client if direct access is needed
func (s *SAPClient) GetClient() *resty.Client {
	return s.client
}

// executeRequest wraps the resty request execution with CSRF handling.
// It takes a function meant to build and execute the request.
func (s *SAPClient) ExecuteRequest(method, url string, body interface{}, queryParams map[string]string) (*resty.Response, error) {
	var resp *resty.Response
	var err error

	// 1. Try with existing token (if we have one, or just try if we don't know it's needed yet)
	// For mutating requests, we check if we need to fetch first.
	isMutating := isMutatingMethod(method)

	// If we anticipate needing a token but don't have one, fetch it now to save a round trip failure.
	// However, standard flow is: Try -> Fail -> Fetch -> Retry
	// We'll optimistically try if we have a token, or if it's GET (doesn't need one usually).
	
	req := s.buildRequest()
	if body != nil {
		req.SetBody(body)
	}
	if len(queryParams) > 0 {
		req.SetQueryParams(queryParams)
	}

	// Attach current token if available
	s.mu.RLock()
	token := s.csrfToken
	s.mu.RUnlock()
	
	if token != "" {
		req.SetHeader(CSRFHeader, token)
	}

	resp, err = req.Execute(method, url)
	if err != nil {
		return nil, err
	}

	// 2. Check for CSRF error
	// SAP usually returns 403 Forbidden with proper header indication, or sometimes generic 403.
	// We detect need for refresh if 403 AND we tried a mutating method.
	if isMutating && (resp.StatusCode() == http.StatusForbidden || resp.Header().Get(CSRFHeader) == "Required") {
		// Log or Debug: "CSRF token invalid or missing, refreshing..."
		if err := s.RefreshCSRFToken(); err != nil {
			return nil, fmt.Errorf("failed to refresh CSRF token: %w", err)
		}

		// 3. Retry with new token
		reqRetry := s.buildRequest()
		if body != nil {
			reqRetry.SetBody(body)
		}
		if len(queryParams) > 0 {
			reqRetry.SetQueryParams(queryParams)
		}
		
		s.mu.RLock()
		newToken := s.csrfToken
		s.mu.RUnlock()
		
		reqRetry.SetHeader(CSRFHeader, newToken)
		
		resp, err = reqRetry.Execute(method, url)
	}

	return resp, err
}

// buildRequest creates a new request and attaches managed cookies
func (s *SAPClient) buildRequest() *resty.Request {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	req := s.client.R()
	// Resty automatically manages cookies in its built-in CookieJar if enabled, 
	// but we might want explicit control if we reset the jar or want to persist generic session cookies manually.
	// For now, let's trust Resty's JAR for session cookies, but valid X-CSRF-Token often requires specific cookies to accompany it.
	// If we manually captured cookies during Fetch, we set them here.
	if len(s.csrfCookies) > 0 {
		req.SetCookies(s.csrfCookies)
	}
	return req
}

// RefreshCSRFToken fetches a new token and updates the client state
func (s *SAPClient) RefreshCSRFToken() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use HEAD or GET to valid endpoint. Service Root "/" is standard.
	req := s.client.R().
		SetHeader(CSRFHeader, CSRFValue)

	resp, err := req.Head("/") // or GET
	if err != nil {
		return err
	}

	if resp.IsError() {
		// Fallback to GET if HEAD fails
		resp, err = req.Get("/")
		if err != nil {
			return err
		}
		if resp.IsError() {
			return fmt.Errorf("csrf fetch failed with status: %d", resp.StatusCode())
		}
	}

	token := resp.Header().Get(CSRFHeader)
	if token == "" {
		return fmt.Errorf("csrf token header not found in response")
	}

	s.csrfToken = token
	s.csrfCookies = resp.Cookies() // Capture cookies explicitly e.g. SAP_SESSIONID
	
	return nil
}

func isMutatingMethod(method string) bool {
	m := strings.ToUpper(method)
	return m == http.MethodPost || m == http.MethodPut || m == http.MethodPatch || m == http.MethodDelete
}
