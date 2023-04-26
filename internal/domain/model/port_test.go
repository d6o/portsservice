package model_test

import (
	"encoding/json"
	"github.com/d6o/portsservice/internal/domain/model"
	"testing"
)

func TestPortDecoding(t *testing.T) {
	jsonData := `{
		"name": "Ajman",
		"city": "Ajman",
		"country": "United Arab Emirates",
		"alias": [],
		"regions": [],
		"coordinates": [55.5136433, 25.4052165],
		"province": "Ajman",
		"timezone": "Asia/Dubai",
		"unlocs": ["AEAJM"],
		"code": "52000"
	}`

	var p model.Port
	err := json.Unmarshal([]byte(jsonData), &p)
	if err != nil {
		t.Errorf("Failed to decode JSON: %v", err)
	}

	if p.Name != "Ajman" {
		t.Errorf("Expected Name to be 'Ajman', got: %s", p.Name)
	}
	if p.City != "Ajman" {
		t.Errorf("Expected City to be 'Ajman', got: %s", p.City)
	}
	if p.Country != "United Arab Emirates" {
		t.Errorf("Expected Country to be 'United Arab Emirates', got: %s", p.Country)
	}
}
