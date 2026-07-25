package models

import "encoding/json"

// MarshalJSONColumn serializes v to a JSON string for storage in a plain
// string column (e.g. Attachments, Metadata). The simple value types these
// columns hold (string slices, maps of primitives) never fail to marshal, so
// the error is swallowed and "" is returned instead of forcing every call
// site to handle an error that in practice never occurs.
func MarshalJSONColumn(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
