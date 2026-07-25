package models

import "testing"

func TestMarshalJSONColumn_Slice(t *testing.T) {
	got := MarshalJSONColumn([]string{"a", "b"})
	want := `["a","b"]`
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestMarshalJSONColumn_Map(t *testing.T) {
	got := MarshalJSONColumn(map[string]int{"synced_count": 3})
	want := `{"synced_count":3}`
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestMarshalJSONColumn_Unmarshalable(t *testing.T) {
	// A channel cannot be marshaled to JSON; MarshalJSONColumn should
	// swallow the error and return "" rather than panicking.
	got := MarshalJSONColumn(make(chan int))
	if got != "" {
		t.Errorf("expected empty string for an unmarshalable value, got %q", got)
	}
}
