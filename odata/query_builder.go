package odata

import (
	"fmt"
	"net/url"
	"strings"
)

// QueryOptions builder for OData v2 parameters
type QueryOptions struct {
	params url.Values
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		params: make(url.Values),
	}
}

// Format adds $format parameter (e.g., "json")
func (q *QueryOptions) Format(format string) *QueryOptions {
	q.params.Set("$format", format)
	return q
}

// Filter adds $filter parameter
func (q *QueryOptions) Filter(filter string) *QueryOptions {
	q.params.Set("$filter", filter)
	return q
}

// Select adds $select parameter
func (q *QueryOptions) Select(fields []string) *QueryOptions {
	q.params.Set("$select", strings.Join(fields, ","))
	return q
}

// Expand adds $expand parameter
func (q *QueryOptions) Expand(entities []string) *QueryOptions {
	q.params.Set("$expand", strings.Join(entities, ","))
	return q
}

// OrderBy adds $orderby parameter
func (q *QueryOptions) OrderBy(field string, asc bool) *QueryOptions {
	direction := "asc"
	if !asc {
		direction = "desc"
	}
	// Append if multiple orderby? V2 usually supports one string like "Name asc, Date desc"
	// For simplicity, this helper sets one. Users can pass the full string if needed or we can append.
	// Let's check if it exists to append
	current := q.params.Get("$orderby")
	clause := fmt.Sprintf("%s %s", field, direction)
	if current != "" {
		q.params.Set("$orderby", current+","+clause)
	} else {
		q.params.Set("$orderby", clause)
	}
	return q
}

// Top adds $top parameter (pagination)
func (q *QueryOptions) Top(n int) *QueryOptions {
	q.params.Set("$top", fmt.Sprintf("%d", n))
	return q
}

// Skip adds $skip parameter (pagination)
func (q *QueryOptions) Skip(n int) *QueryOptions {
	q.params.Set("$skip", fmt.Sprintf("%d", n))
	return q
}

// InlineCount adds $inlinecount parameter (allpages or none)
func (q *QueryOptions) InlineCount(allPages bool) *QueryOptions {
	val := "none"
	if allPages {
		val = "allpages"
	}
	q.params.Set("$inlinecount", val)
	return q
}

// Build returns the map of query parameters for Resty
func (q *QueryOptions) Build() map[string]string {
	m := make(map[string]string)
	for k, v := range q.params {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}
