package nullable

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Nullable represents a value that may be null.
// It embeds sql.Null[T] to provide database compatibility while adding JSON marshaling.
type Nullable[T any] struct {
	sql.Null[T]
}

// NewNullable creates a new valid Nullable with the given value.
func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{
		sql.Null[T]{
			V:     value,
			Valid: true,
		},
	}
}

// NewNull creates a new null Nullable.
func NewNull[T any]() Nullable[T] {
	return Nullable[T]{
		sql.Null[T]{
			Valid: false,
		},
	}
}

// Ptr returns a pointer to the value if valid, otherwise nil.
func (n Nullable[T]) Ptr() *T {
	if !n.Valid {
		return nil
	}
	return &n.V
}

// ValueOr returns the value if valid, otherwise returns the default value.
func (n Nullable[T]) ValueOr(defaultValue T) T {
	if !n.Valid {
		return defaultValue
	}
	return n.V
}

// MarshalJSON implements the json.Marshaler interface.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.V)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	n.Valid = true
	return json.Unmarshal(data, &n.V)
}

// Scan implements the sql.Scanner interface.
func (n *Nullable[T]) Scan(value any) error {
	return n.Null.Scan(value)
}

// Value implements the T interface.
func (n Nullable[T]) Value() (T, error) {
	if !n.Valid {
		var zero T
		return zero, nil
	}
	// For driver.Value, we need to return a basic type
	// This might need type assertions depending on T
	return n.V, nil
}

// String returns a string representation of the nullable value.
func (n Nullable[T]) String() string {
	if !n.Valid {
		return "null"
	}
	return fmt.Sprintf("%v", n.V)
}
