package nullable

import (
	"encoding/json"
	"testing"
)

func TestNewNullable(t *testing.T) {
	n := NewNullable("test")
	if !n.Valid {
		t.Error("Expected Valid to be true")
	}
	if n.V != "test" {
		t.Errorf("Expected V to be 'test', got %v", n.V)
	}
}

func TestNewNull(t *testing.T) {
	n := NewNull[string]()
	if n.Valid {
		t.Error("Expected Valid to be false")
	}
}

func TestPtr(t *testing.T) {
	// Valid nullable
	n1 := NewNullable(42)
	ptr := n1.Ptr()
	if ptr == nil {
		t.Error("Expected non-nil pointer for valid nullable")
	}
	if *ptr != 42 {
		t.Errorf("Expected *ptr to be 42, got %v", *ptr)
	}

	// Null nullable
	n2 := NewNull[int]()
	ptr2 := n2.Ptr()
	if ptr2 != nil {
		t.Error("Expected nil pointer for null nullable")
	}
}

func TestValueOr(t *testing.T) {
	// Valid nullable
	n1 := NewNullable("hello")
	result := n1.ValueOr("default")
	if result != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}

	// Null nullable
	n2 := NewNull[string]()
	result2 := n2.ValueOr("default")
	if result2 != "default" {
		t.Errorf("Expected 'default', got %v", result2)
	}
}

func TestMarshalJSON(t *testing.T) {
	// Valid nullable
	n1 := NewNullable("test")
	data, err := json.Marshal(n1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := `"test"`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}

	// Null nullable
	n2 := NewNull[string]()
	data2, err := json.Marshal(n2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected2 := "null"
	if string(data2) != expected2 {
		t.Errorf("Expected %s, got %s", expected2, string(data2))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	// Valid JSON value
	var n1 Nullable[string]
	err := json.Unmarshal([]byte(`"test"`), &n1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !n1.Valid {
		t.Error("Expected Valid to be true")
	}
	if n1.V != "test" {
		t.Errorf("Expected V to be 'test', got %v", n1.V)
	}

	// Null JSON value
	var n2 Nullable[string]
	err = json.Unmarshal([]byte("null"), &n2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if n2.Valid {
		t.Error("Expected Valid to be false")
	}

	// Number type
	var n3 Nullable[int]
	err = json.Unmarshal([]byte("42"), &n3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !n3.Valid {
		t.Error("Expected Valid to be true")
	}
	if n3.V != 42 {
		t.Errorf("Expected V to be 42, got %v", n3.V)
	}
}

func TestScan(t *testing.T) {
	// Valid value
	var n1 Nullable[string]
	err := n1.Scan("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !n1.Valid {
		t.Error("Expected Valid to be true")
	}
	if n1.V != "test" {
		t.Errorf("Expected V to be 'test', got %v", n1.V)
	}

	// Nil value
	var n2 Nullable[string]
	err = n2.Scan(nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if n2.Valid {
		t.Error("Expected Valid to be false")
	}

	// sql.Null[T] handles type conversion, so test with compatible types
	var n3 Nullable[int64]
	err = n3.Scan(int64(42))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !n3.Valid {
		t.Error("Expected Valid to be true")
	}
	if n3.V != 42 {
		t.Errorf("Expected V to be 42, got %v", n3.V)
	}
}

func TestValue(t *testing.T) {
	// Valid nullable
	n1 := NewNullable("test")
	val, err := n1.Value()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != "test" {
		t.Errorf("Expected 'test', got %v", val)
	}

	// Null nullable
	n2 := NewNull[string]()
	val2, err := n2.Value()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val2 != "" {
		t.Errorf("Expected nil, got %v", val2)
	}
}

func TestString(t *testing.T) {
	// Valid nullable
	n1 := NewNullable("test")
	str := n1.String()
	if str != "test" {
		t.Errorf("Expected 'test', got %s", str)
	}

	// Null nullable
	n2 := NewNull[string]()
	str2 := n2.String()
	if str2 != "null" {
		t.Errorf("Expected 'null', got %s", str2)
	}

	// Number
	n3 := NewNullable(42)
	str3 := n3.String()
	if str3 != "42" {
		t.Errorf("Expected '42', got %s", str3)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	type TestStruct struct {
		Name Nullable[string] `json:"name"`
		Age  Nullable[int]    `json:"age"`
	}

	// Test with valid values
	original := TestStruct{
		Name: NewNullable("John"),
		Age:  NewNullable(30),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Errorf("Marshal error: %v", err)
	}

	var decoded TestStruct
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	if !decoded.Name.Valid || decoded.Name.V != "John" {
		t.Errorf("Name mismatch: expected valid 'John', got %+v", decoded.Name)
	}
	if !decoded.Age.Valid || decoded.Age.V != 30 {
		t.Errorf("Age mismatch: expected valid 30, got %+v", decoded.Age)
	}

	// Test with null values
	original2 := TestStruct{
		Name: NewNull[string](),
		Age:  NewNull[int](),
	}

	data2, err := json.Marshal(original2)
	if err != nil {
		t.Errorf("Marshal error: %v", err)
	}

	var decoded2 TestStruct
	err = json.Unmarshal(data2, &decoded2)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	if decoded2.Name.Valid {
		t.Error("Expected Name to be null")
	}
	if decoded2.Age.Valid {
		t.Error("Expected Age to be null")
	}
}

func TestDatabaseInteraction(t *testing.T) {
	// Test Value method for database storage
	n1 := NewNullable("test")
	val, err := n1.Value()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != "test" {
		t.Errorf("Expected 'test', got %v", val)
	}

	// Test Scan method for database retrieval
	var n2 Nullable[string]
	err = n2.Scan("scanned_value")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !n2.Valid {
		t.Error("Expected Valid to be true")
	}
	if n2.V != "scanned_value" {
		t.Errorf("Expected 'scanned_value', got %v", n2.V)
	}

	// Test null scan
	var n3 Nullable[string]
	err = n3.Scan(nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if n3.Valid {
		t.Error("Expected Valid to be false")
	}
}
