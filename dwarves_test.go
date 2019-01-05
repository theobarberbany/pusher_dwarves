package main

import "testing"

func TestGetDwarf(t *testing.T) {
	expected_thorin := `{"dwarf":{"name":"Thorin","birth":"TA 2746","death":"TA 2941","culture":"Durin's Folk"}}`
	expected_error := `{"error":"dwarf not found"}`

	out, err := getDwarf("Thorin")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if string(*out) != expected_thorin {
		t.Errorf("Unexpected response, got: %s, expected: %s", string(*out), expected_thorin)
	}

	out2, err := getDwarf("NotADwarf")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if string(*out2) != expected_error {
		t.Errorf("Unexpected response, got: %s, expected: %s", string(*out2), expected_error)
	}
}
