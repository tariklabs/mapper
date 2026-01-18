# mapper

[![Go Reference](https://pkg.go.dev/badge/github.com/tariklabs/mapper.svg)](https://pkg.go.dev/github.com/tariklabs/mapper)
[![Go Report Card](https://goreportcard.com/badge/github.com/tariklabs/mapper)](https://goreportcard.com/report/github.com/tariklabs/mapper)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![CI](https://github.com/tariklabs/mapper/actions/workflows/ci.yml/badge.svg)](https://github.com/tariklabs/mapper/actions/workflows/ci.yml)

A zero-boilerplate, reflection-based struct mapper for Go. Automatically maps fields between structs with tags, nested paths, patch semantics, and slice support.

## Features

- **Zero Boilerplate** - No code generation, no manual field assignments
- **Tag-Based Aliasing** - Map fields with different names using struct tags
- **String Conversion** - Automatic string-to-primitive conversion via `mapconv` tag
- **Nested Structs** - Recursive mapping of arbitrarily nested structures
- **Deep Copying** - Slices and maps are deep-copied, not shared
- **Pointer Flexibility** - Seamless conversion between pointer and value types
- **Patch Semantics** - Skip zero values for partial updates
- **Strict Mode** - Ensure all destination fields are populated
- **Thread Safe** - Safe for concurrent use with internal caching
- **Performance Optimized** - Struct metadata caching minimizes reflection overhead

## Installation

```bash
go get github.com/tariklabs/mapper
```

Requires Go 1.21 or later.

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/tariklabs/mapper"
)

type UserDTO struct {
    FullName string `map:"Name"`
    Email    string
    Age      string `mapconv:"int"`
}

type User struct {
    Name  string
    Email string
    Age   int
}

func main() {
    dto := UserDTO{
        FullName: "Alice Smith",
        Email:    "alice@example.com",
        Age:      "30",
    }

    var user User
    if err := mapper.Map(&user, dto); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", user)
    // Output: {Name:Alice Smith Email:alice@example.com Age:30}
}
```

## Documentation

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/tariklabs/mapper).

## Usage

### Basic Mapping

Fields are matched by name (case-sensitive):

```go
type Source struct {
    Name  string
    Email string
    Age   int
}

type Destination struct {
    Name  string
    Email string
    Age   int
}

src := Source{Name: "Alice", Email: "alice@example.com", Age: 30}
var dst Destination

err := mapper.Map(&dst, src)
// dst = {Name: "Alice", Email: "alice@example.com", Age: 30}
```

### Tag-Based Field Aliasing

Use the `map` tag to map fields with different names:

```go
type APIResponse struct {
    UserName    string `map:"Name"`
    UserEmail   string `map:"Email"`
    YearsOld    int    `map:"Age"`
}

type User struct {
    Name  string
    Email string
    Age   int
}

response := APIResponse{
    UserName:  "Bob",
    UserEmail: "bob@example.com",
    YearsOld:  25,
}
var user User

err := mapper.Map(&user, response)
// user = {Name: "Bob", Email: "bob@example.com", Age: 25}
```

### String-to-Type Conversion

Use the `mapconv` tag to convert string fields to typed values:

```go
type FormInput struct {
    Age      string `mapconv:"int"`
    Score    string `mapconv:"float64"`
    Active   string `mapconv:"bool"`
    Quantity string `mapconv:"uint"`
}

type Product struct {
    Age      int
    Score    float64
    Active   bool
    Quantity uint
}

input := FormInput{
    Age:      "42",
    Score:    "95.5",
    Active:   "true",
    Quantity: "100",
}
var product Product

err := mapper.Map(&product, input)
// product = {Age: 42, Score: 95.5, Active: true, Quantity: 100}
```

**Supported conversion types:** `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`

Tags can be combined:

```go
type Input struct {
    UserAge string `map:"Age" mapconv:"int"`
}
```

### Nested Structs

Nested structs are mapped recursively:

```go
type SrcAddress struct {
    Street string
    City   string
}

type SrcPerson struct {
    Name    string
    Address SrcAddress
}

type DstAddress struct {
    Street string
    City   string
}

type DstPerson struct {
    Name    string
    Address DstAddress
}

src := SrcPerson{
    Name: "Alice",
    Address: SrcAddress{
        Street: "123 Main St",
        City:   "Seattle",
    },
}
var dst DstPerson

err := mapper.Map(&dst, src)
// dst.Address = {Street: "123 Main St", City: "Seattle"}
```

### Slices and Maps

Slices and maps are deep-copied with element type conversion:

```go
type Source struct {
    Tags    []string
    Counts  []int32
    Config  map[string]string
    Values  map[string]int32
}

type Destination struct {
    Tags    []string
    Counts  []int64             // Different element type
    Config  map[string]string
    Values  map[string]int64    // Different value type
}

src := Source{
    Tags:   []string{"go", "mapper"},
    Counts: []int32{1, 2, 3},
    Config: map[string]string{"env": "prod"},
    Values: map[string]int32{"a": 100},
}
var dst Destination

err := mapper.Map(&dst, src)
// Deep copy with automatic type conversion
```

**Important:** Modifying the source after mapping does not affect the destination.

### Pointer Handling

Seamless conversion between pointer and value types:

```go
// Value to pointer
type Source struct {
    Name string
}

type Destination struct {
    Name *string
}

src := Source{Name: "Alice"}
var dst Destination

err := mapper.Map(&dst, src)
// *dst.Name == "Alice"
```

```go
// Pointer to value
type Source struct {
    Name *string
}

type Destination struct {
    Name string
}

name := "Bob"
src := Source{Name: &name}
var dst Destination

err := mapper.Map(&dst, src)
// dst.Name == "Bob"
```

Nil pointers are handled gracefully and do not overwrite destination values.

## Options

Use `MapWithOptions` for customized behavior:

### WithTagName

Use a different struct tag for field aliasing:

```go
type Response struct {
    UserName string `json:"name"`
}

type User struct {
    Name string `json:"name"`
}

err := mapper.MapWithOptions(&user, response, mapper.WithTagName("json"))
```

### WithIgnoreZeroSource

Skip zero-value fields for patch/partial update operations:

```go
type User struct {
    Name  string
    Email string
    Age   int
}

existing := User{Name: "Alice", Email: "alice@old.com", Age: 25}
patch := User{Name: "Alicia", Email: "", Age: 0}  // Only update Name

err := mapper.MapWithOptions(&existing, patch, mapper.WithIgnoreZeroSource())
// existing = {Name: "Alicia", Email: "alice@old.com", Age: 25}
```

### WithStrictMode

Return an error if any destination field has no matching source:

```go
type Source struct {
    Name string
}

type Destination struct {
    Name  string
    Email string  // No source field!
}

err := mapper.MapWithOptions(&dst, src, mapper.WithStrictMode())
// Error: no matching source field found for "Email"
```

### WithMaxDepth

Set maximum nesting depth (default: 64):

```go
err := mapper.MapWithOptions(&dst, src, mapper.WithMaxDepth(100))
```

### Combining Options

```go
err := mapper.MapWithOptions(&dst, src,
    mapper.WithTagName("json"),
    mapper.WithIgnoreZeroSource(),
    mapper.WithStrictMode(),
    mapper.WithMaxDepth(100),
)
```

## Error Handling

Errors are returned as `*MappingError` with detailed context:

```go
err := mapper.Map(&dst, src)
if err != nil {
    var mappingErr *mapper.MappingError
    if errors.As(err, &mappingErr) {
        fmt.Printf("Source type: %s\n", mappingErr.SrcType)
        fmt.Printf("Destination type: %s\n", mappingErr.DstType)
        fmt.Printf("Field path: %s\n", mappingErr.FieldPath)
        fmt.Printf("Reason: %s\n", mappingErr.Reason)
    }
}
```

Field paths use dot notation for nested fields (`Address.City`) and bracket notation for slices (`Items[0]`) and maps (`Config[key]`).

### Common Errors

| Reason | Cause |
|--------|-------|
| `dst must be a non-nil pointer to struct` | Destination is not a valid pointer |
| `src must be a struct or pointer to struct` | Source is not a struct type |
| `no matching source field found` | Strict mode: destination field has no source |
| `incompatible field types: X -> Y` | Types cannot be converted |
| `maximum nesting depth exceeded` | Depth limit reached (circular reference protection) |
| `cannot convert "X" to Y` | String conversion failed |

## Performance

The mapper uses struct metadata caching to minimize reflection overhead. First-time mapping of a struct type incurs reflection cost, but subsequent mappings use cached metadata.

```
BenchmarkMap_Flat-12                 847845    1412 ns/op    312 B/op    11 allocs/op
BenchmarkMap_Nested-12               398174    2987 ns/op    720 B/op    24 allocs/op
BenchmarkMap_Slice100-12              27862   42891 ns/op  14832 B/op   409 allocs/op
```

For performance-critical hot paths where nanoseconds matter, consider manual field assignment. For typical application code (API handlers, DTOs, data transformation), the mapper provides an excellent balance of convenience and performance.

See [BENCHMARKS.md](BENCHMARKS.md) for detailed profiling instructions.

## Thread Safety

All functions are safe for concurrent use. The internal metadata cache uses `sync.Map` for thread-safe access.

## Limitations

- **Exported fields only** - Unexported (private) fields cannot be mapped
- **Structs only** - Interface types are not supported as field types
- **Built-in conversions** - No custom converter functions
- **Depth-based protection** - Circular references are protected by depth limit, not runtime detection

## Real-World Examples

### API Request to Domain Model

```go
type CreateOrderRequest struct {
    CustomerName string            `map:"Customer"`
    TotalPrice   string            `map:"Total" mapconv:"float64"`
    IsPriority   string            `mapconv:"bool"`
    Tags         []string          `map:"Labels"`
    Metadata     map[string]string `map:"Info"`
}

type Order struct {
    Customer   string
    Total      float64
    IsPriority bool
    Labels     []string
    Info       map[string]string
}

func CreateOrder(req CreateOrderRequest) (*Order, error) {
    var order Order
    if err := mapper.Map(&order, req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    return &order, nil
}
```

### Database Entity to API Response

```go
type UserEntity struct {
    ID        int64
    Username  string
    Email     string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}

type UserResponse struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func ToResponse(entity UserEntity) UserResponse {
    var response UserResponse
    mapper.Map(&response, entity)  // Only copies matching fields
    return response
}
```

### Partial Update (Patch)

```go
type UpdateUserRequest struct {
    Name  string `json:"name,omitempty"`
    Email string `json:"email,omitempty"`
    Age   int    `json:"age,omitempty"`
}

func UpdateUser(userID int64, req UpdateUserRequest) error {
    user, err := db.GetUser(userID)
    if err != nil {
        return err
    }

    // Only update fields that were provided (non-zero)
    if err := mapper.MapWithOptions(user, req, mapper.WithIgnoreZeroSource()); err != nil {
        return err
    }

    return db.SaveUser(user)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`go test -v -race ./...`)
4. Run linter (`golangci-lint run`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
