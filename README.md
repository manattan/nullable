# Nullable

A Go package that provides a generic nullable type that combines database compatibility with JSON marshaling.

## Features

- Generic `Nullable[T]` type that works with any type
- Database compatibility through `sql.Scanner` and `driver.Valuer` interfaces
- JSON marshaling/unmarshaling support
- Built on top of Go's `sql.Null[T]` for robust database integration
- Helper methods for common operations

## Installation

```bash
go get github.com/manattan/nullable
```

## Usage

### Basic Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/manattan/nullable"
)

func main() {
    // Create a nullable string with a value
    name := nullable.NewNullable("John")
    fmt.Println(name.String()) // "John"
    
    // Create a null nullable string
    emptyName := nullable.NewNull[string]()
    fmt.Println(emptyName.String()) // "null"
    
    // Use ValueOr for safe value access
    actualName := name.ValueOr("Unknown")      // "John"
    defaultName := emptyName.ValueOr("Unknown") // "Unknown"
}
```

### JSON Marshaling

```go
type User struct {
    Name nullable.Nullable[string] `json:"name"`
    Age  nullable.Nullable[int]    `json:"age"`
}

user := User{
    Name: nullable.NewNullable("Alice"),
    Age:  nullable.NewNull[int](),
}

// Marshal to JSON
data, _ := json.Marshal(user)
fmt.Println(string(data)) // {"name":"Alice","age":null}

// Unmarshal from JSON
var decoded User
json.Unmarshal([]byte(`{"name":"Bob","age":25}`), &decoded)
```

### Database Usage

```go
import "database/sql"

type Person struct {
    ID   int64
    Name nullable.Nullable[string]
    Age  nullable.Nullable[int]
}

// Insert into database
_, err := db.Exec("INSERT INTO people (name, age) VALUES (?, ?)", 
    person.Name, person.Age)

// Scan from database
var p Person
err := db.QueryRow("SELECT name, age FROM people WHERE id = ?", id).
    Scan(&p.Name, &p.Age)
```

## API Reference

### Constructor Functions

- `NewNullable[T](value T) Nullable[T]` - Creates a nullable with a valid value
- `NewNull[T]() Nullable[T]` - Creates a null nullable

### Methods

- `Ptr() *T` - Returns pointer to value if valid, nil otherwise
- `ValueOr(defaultValue T) T` - Returns value if valid, otherwise default
- `String() string` - String representation
- `MarshalJSON() ([]byte, error)` - JSON marshaling
- `UnmarshalJSON(data []byte) error` - JSON unmarshaling
- `Scan(value any) error` - Database scanning (sql.Scanner)
- `Value() (T, error)` - Database value (driver.Valuer)

## Testing

Run the test suite:

```bash
go test
```

## License

This project is available under the MIT License.