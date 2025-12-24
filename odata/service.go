package odata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Willias7788/go-odata-v2-sdk/client"
	"github.com/Willias7788/go-odata-v2-sdk/models"
)

// Service represents a specific OData service endpoint
type Service struct {
	client      *client.SAPClient
	servicePath string // e.g. "/sap/opu/odata/IWBEP/GWSAMPLE_BASIC/"
}

// NewService creates a new OData service handler
func NewService(client *client.SAPClient, servicePath string) *Service {
	// Ensure service path has trailing slash
	if !strings.HasSuffix(servicePath, "/") {
		servicePath += "/"
	}
	// Ensure service path has leading slash
	if !strings.HasPrefix(servicePath, "/") {
		servicePath = "/" + servicePath
	}
	return &Service{
		client:      client,
		servicePath: servicePath,
	}
}

func (s *Service) buildURL(entitySet string) string {
	return s.servicePath + entitySet
}

func (s *Service) buildKeyURL(entitySet, key string) string {
	// Simple check: if key doesn't start with (, wrap it? 
	// OData keys can be complicated (prop=val vs 'val'). 
	// We assume user passes valid key predicate like "('123')" or "(Id='123',Type='A')"
	// If the user just passes "123", we might want to be smart, but generic SDKs should prioritize predictability.
	// We'll trust the user passed the predicate.
	if !strings.HasPrefix(key, "(") {
		key = "(" + key + ")"
	}
	return s.servicePath + entitySet + key
}

// GetEntitySet fetches a collection of entities
func GetEntitySet[T any](s *Service, entitySet string, opts *QueryOptions) (*models.ODataResponse[[]T], error) {
	url := s.buildURL(entitySet)
	var qParams map[string]string
	if opts != nil {
		qParams = opts.Build()
	}

	resp, err := s.client.ExecuteRequest(http.MethodGet, url, nil, qParams)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp.Body())
	}

	var result models.ODataResponse[[]T]
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &result, nil
}

// GetEntityByKey fetches a single entity
func GetEntityByKey[T any](s *Service, entitySet, key string, opts *QueryOptions) (*models.ODataResponse[T], error) {
	url := s.buildKeyURL(entitySet, key)
	var qParams map[string]string
	if opts != nil {
		qParams = opts.Build()
	}

	resp, err := s.client.ExecuteRequest(http.MethodGet, url, nil, qParams)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp.Body())
	}

	var result models.ODataResponse[T]
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &result, nil
}

// CreateEntity creates a new entity
func CreateEntity[T any](s *Service, entitySet string, payload interface{}) (*models.ODataResponse[T], error) {
	url := s.buildURL(entitySet)
	
	resp, err := s.client.ExecuteRequest(http.MethodPost, url, payload, nil)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, parseError(resp.Body())
	}

	var result models.ODataResponse[T]
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &result, nil
}

// UpdateEntity updates an existing entity (PUT)
func UpdateEntity(s *Service, entitySet, key string, payload interface{}) error {
	url := s.buildKeyURL(entitySet, key)
	
	resp, err := s.client.ExecuteRequest(http.MethodPut, url, payload, nil)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return parseError(resp.Body())
	}

	return nil
}

// PatchEntity updates an existing entity (PATCH/MERGE)
func PatchEntity(s *Service, entitySet, key string, payload interface{}) error {
	url := s.buildKeyURL(entitySet, key)
	
	resp, err := s.client.ExecuteRequest(http.MethodPatch, url, payload, nil)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return parseError(resp.Body())
	}

	return nil
}

// DeleteEntity deletes an entity
func DeleteEntity(s *Service, entitySet, key string) error {
	url := s.buildKeyURL(entitySet, key)
	
	resp, err := s.client.ExecuteRequest(http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return parseError(resp.Body())
	}

	return nil
}

func parseError(body []byte) error {
	var errResp models.ODataErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("http error and failed to parse odata error: %s", string(body))
	}
	return &errResp
}
