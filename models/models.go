package models

import "encoding/json"

// ODataResponse is the generic container for OData v2 JSON responses.
// V2 typically wraps results in a "d" object.
type ODataResponse[T any] struct {
	D DWrapper[T] `json:"d"`
}

// DWrapper handles the "result" vs "results" discrepancy or direct object return.
// In many OData v2 implementations:
// - Collections are in d.results
// - Single entities are directly in d (or d.results depending on system, but usually d)
// We use a custom unmarshaller or pointer strategies to handle this generic.
// However, for strict V2, collections are usually `d: { results: [] }`.
// DWrapper handles the "result" vs "results" discrepancy.
type DWrapper[T any] struct {
	Result T
}

func (w *DWrapper[T]) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Case 1: d.results exists (Common for collections and some single entities)
	if val, ok := raw["results"]; ok {
		return json.Unmarshal(val, &w.Result)
	}

	// Case 2: Direct entity properties in d
	return json.Unmarshal(data, &w.Result)
}

// ODataErrorResponse handles OData error structures
type ODataErrorResponse struct {
	Err ODataError `json:"error"`
}

type ODataError struct {
	Code    string       `json:"code"`
	Message ODataMessage `json:"message"`
}

type ODataMessage struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// ODataError implements the error interface
func (e *ODataErrorResponse) Error() string {
	return e.Err.Message.Value
}
