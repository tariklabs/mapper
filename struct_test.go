package mapper

import (
	"testing"
)

// Nested Struct Tests
type SrcAddress struct {
	Street  string
	City    string
	ZipCode string
}

type DstAddress struct {
	Street  string
	City    string
	ZipCode string
}

type SrcPerson struct {
	Name    string
	Age     int
	Address SrcAddress
}

type DstPerson struct {
	Name    string
	Age     int
	Address DstAddress
}

func TestNestedStruct_Basic(t *testing.T) {
	src := SrcPerson{
		Name: "John",
		Age:  30,
		Address: SrcAddress{
			Street:  "123 Main St",
			City:    "Boston",
			ZipCode: "02101",
		},
	}
	var dst DstPerson

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "John" {
		t.Errorf("expected Name = 'John', got %q", dst.Name)
	}
	if dst.Age != 30 {
		t.Errorf("expected Age = 30, got %d", dst.Age)
	}
	if dst.Address.Street != "123 Main St" {
		t.Errorf("expected Address.Street = '123 Main St', got %q", dst.Address.Street)
	}
	if dst.Address.City != "Boston" {
		t.Errorf("expected Address.City = 'Boston', got %q", dst.Address.City)
	}
	if dst.Address.ZipCode != "02101" {
		t.Errorf("expected Address.ZipCode = '02101', got %q", dst.Address.ZipCode)
	}
}

func TestNestedStruct_DeepCopy(t *testing.T) {
	src := SrcPerson{
		Name: "John",
		Address: SrcAddress{
			City: "Boston",
		},
	}
	var dst DstPerson

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src.Address.City = "New York"

	if dst.Address.City != "Boston" {
		t.Errorf("nested struct modification affected destination: got %q", dst.Address.City)
	}
}

// Deeply nested structs (3 levels)
type SrcCountry struct {
	Name string
	Code string
}

type DstCountry struct {
	Name string
	Code string
}

type SrcFullAddress struct {
	Street  string
	City    string
	Country SrcCountry
}

type DstFullAddress struct {
	Street  string
	City    string
	Country DstCountry
}

type SrcEmployee struct {
	Name    string
	Address SrcFullAddress
}

type DstEmployee struct {
	Name    string
	Address DstFullAddress
}

func TestNestedStruct_ThreeLevels(t *testing.T) {
	src := SrcEmployee{
		Name: "Alice",
		Address: SrcFullAddress{
			Street: "456 Oak Ave",
			City:   "Seattle",
			Country: SrcCountry{
				Name: "United States",
				Code: "US",
			},
		},
	}
	var dst DstEmployee

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Alice" {
		t.Errorf("expected Name = 'Alice', got %q", dst.Name)
	}
	if dst.Address.Street != "456 Oak Ave" {
		t.Errorf("expected Address.Street = '456 Oak Ave', got %q", dst.Address.Street)
	}
	if dst.Address.Country.Name != "United States" {
		t.Errorf("expected Address.Country.Name = 'United States', got %q", dst.Address.Country.Name)
	}
	if dst.Address.Country.Code != "US" {
		t.Errorf("expected Address.Country.Code = 'US', got %q", dst.Address.Country.Code)
	}
}

// Pointer to Nested Struct Tests
type SrcWithPtrAddress struct {
	Name    string
	Address *SrcAddress
}

type DstWithPtrAddress struct {
	Name    string
	Address *DstAddress
}

func TestNestedStruct_Pointer(t *testing.T) {
	src := SrcWithPtrAddress{
		Name: "Bob",
		Address: &SrcAddress{
			Street: "789 Pine St",
			City:   "Portland",
		},
	}
	var dst DstWithPtrAddress

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Bob" {
		t.Errorf("expected Name = 'Bob', got %q", dst.Name)
	}
	if dst.Address == nil {
		t.Fatal("expected Address to not be nil")
	}
	if dst.Address.Street != "789 Pine St" {
		t.Errorf("expected Address.Street = '789 Pine St', got %q", dst.Address.Street)
	}
	if dst.Address.City != "Portland" {
		t.Errorf("expected Address.City = 'Portland', got %q", dst.Address.City)
	}
}

func TestNestedStruct_NilPointer(t *testing.T) {
	src := SrcWithPtrAddress{
		Name:    "Charlie",
		Address: nil,
	}
	var dst DstWithPtrAddress

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Charlie" {
		t.Errorf("expected Name = 'Charlie', got %q", dst.Name)
	}
	if dst.Address != nil {
		t.Errorf("expected Address to be nil, got %v", dst.Address)
	}
}

func TestNestedStruct_Pointer_DeepCopy(t *testing.T) {
	src := SrcWithPtrAddress{
		Name: "Dave",
		Address: &SrcAddress{
			City: "Denver",
		},
	}
	var dst DstWithPtrAddress

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src.Address.City = "Dallas"

	if dst.Address.City != "Denver" {
		t.Errorf("pointer struct modification affected destination: got %q", dst.Address.City)
	}
}

// Struct in Slice Tests
type SrcItem struct {
	ID    int
	Name  string
	Price float64
}

type DstItem struct {
	ID    int
	Name  string
	Price float64
}

type SrcOrder struct {
	OrderID string
	Items   []SrcItem
}

type DstOrder struct {
	OrderID string
	Items   []DstItem
}

func TestStructInSlice_Basic(t *testing.T) {
	src := SrcOrder{
		OrderID: "ORD-001",
		Items: []SrcItem{
			{ID: 1, Name: "Widget", Price: 9.99},
			{ID: 2, Name: "Gadget", Price: 19.99},
			{ID: 3, Name: "Gizmo", Price: 29.99},
		},
	}
	var dst DstOrder

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.OrderID != "ORD-001" {
		t.Errorf("expected OrderID = 'ORD-001', got %q", dst.OrderID)
	}
	if len(dst.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(dst.Items))
	}
	if dst.Items[0].Name != "Widget" {
		t.Errorf("expected Items[0].Name = 'Widget', got %q", dst.Items[0].Name)
	}
	if dst.Items[1].Price != 19.99 {
		t.Errorf("expected Items[1].Price = 19.99, got %f", dst.Items[1].Price)
	}
	if dst.Items[2].ID != 3 {
		t.Errorf("expected Items[2].ID = 3, got %d", dst.Items[2].ID)
	}
}

func TestStructInSlice_Empty(t *testing.T) {
	src := SrcOrder{
		OrderID: "ORD-002",
		Items:   []SrcItem{},
	}
	var dst DstOrder

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Items == nil {
		t.Error("expected Items to be empty slice, got nil")
	}
	if len(dst.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(dst.Items))
	}
}

func TestStructInSlice_Nil(t *testing.T) {
	src := SrcOrder{
		OrderID: "ORD-003",
		Items:   nil,
	}
	var dst DstOrder

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Items != nil {
		t.Errorf("expected Items to be nil, got %v", dst.Items)
	}
}

func TestStructInSlice_DeepCopy(t *testing.T) {
	src := SrcOrder{
		OrderID: "ORD-004",
		Items: []SrcItem{
			{ID: 1, Name: "Original"},
		},
	}
	var dst DstOrder

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src.Items[0].Name = "Modified"

	if dst.Items[0].Name != "Original" {
		t.Errorf("struct in slice modification affected destination: got %q", dst.Items[0].Name)
	}
}

// Pointer to struct in slice
type SrcOrderWithPtrs struct {
	OrderID string
	Items   []*SrcItem
}

type DstOrderWithPtrs struct {
	OrderID string
	Items   []*DstItem
}

func TestPointerStructInSlice_Basic(t *testing.T) {
	src := SrcOrderWithPtrs{
		OrderID: "ORD-005",
		Items: []*SrcItem{
			{ID: 1, Name: "Widget", Price: 9.99},
			{ID: 2, Name: "Gadget", Price: 19.99},
		},
	}
	var dst DstOrderWithPtrs

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(dst.Items))
	}
	if dst.Items[0] == nil {
		t.Fatal("expected Items[0] to not be nil")
	}
	if dst.Items[0].Name != "Widget" {
		t.Errorf("expected Items[0].Name = 'Widget', got %q", dst.Items[0].Name)
	}
}

func TestPointerStructInSlice_WithNil(t *testing.T) {
	src := SrcOrderWithPtrs{
		OrderID: "ORD-006",
		Items: []*SrcItem{
			{ID: 1, Name: "Widget"},
			nil,
			{ID: 3, Name: "Gizmo"},
		},
	}
	var dst DstOrderWithPtrs

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(dst.Items))
	}
	if dst.Items[0].Name != "Widget" {
		t.Errorf("expected Items[0].Name = 'Widget', got %q", dst.Items[0].Name)
	}
	if dst.Items[1] != nil {
		t.Errorf("expected Items[1] to be nil, got %v", dst.Items[1])
	}
	if dst.Items[2].Name != "Gizmo" {
		t.Errorf("expected Items[2].Name = 'Gizmo', got %q", dst.Items[2].Name)
	}
}

func TestPointerStructInSlice_DeepCopy(t *testing.T) {
	src := SrcOrderWithPtrs{
		OrderID: "ORD-007",
		Items: []*SrcItem{
			{ID: 1, Name: "Original"},
		},
	}
	var dst DstOrderWithPtrs

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src.Items[0].Name = "Modified"

	if dst.Items[0].Name != "Original" {
		t.Errorf("pointer struct in slice modification affected destination: got %q", dst.Items[0].Name)
	}
}

// Struct in Map Tests
type SrcConfig struct {
	Host string
	Port int
}

type DstConfig struct {
	Host string
	Port int
}

type SrcService struct {
	Name    string
	Configs map[string]SrcConfig
}

type DstService struct {
	Name    string
	Configs map[string]DstConfig
}

func TestStructInMap_Basic(t *testing.T) {
	src := SrcService{
		Name: "api-service",
		Configs: map[string]SrcConfig{
			"database": {Host: "localhost", Port: 5432},
			"cache":    {Host: "redis", Port: 6379},
		},
	}
	var dst DstService

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "api-service" {
		t.Errorf("expected Name = 'api-service', got %q", dst.Name)
	}
	if len(dst.Configs) != 2 {
		t.Fatalf("expected 2 configs, got %d", len(dst.Configs))
	}
	if dst.Configs["database"].Host != "localhost" {
		t.Errorf("expected Configs[database].Host = 'localhost', got %q", dst.Configs["database"].Host)
	}
	if dst.Configs["database"].Port != 5432 {
		t.Errorf("expected Configs[database].Port = 5432, got %d", dst.Configs["database"].Port)
	}
	if dst.Configs["cache"].Host != "redis" {
		t.Errorf("expected Configs[cache].Host = 'redis', got %q", dst.Configs["cache"].Host)
	}
}

func TestStructInMap_Empty(t *testing.T) {
	src := SrcService{
		Name:    "empty-service",
		Configs: map[string]SrcConfig{},
	}
	var dst DstService

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Configs == nil {
		t.Error("expected Configs to be empty map, got nil")
	}
	if len(dst.Configs) != 0 {
		t.Errorf("expected 0 configs, got %d", len(dst.Configs))
	}
}

func TestStructInMap_Nil(t *testing.T) {
	src := SrcService{
		Name:    "nil-service",
		Configs: nil,
	}
	var dst DstService

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Configs != nil {
		t.Errorf("expected Configs to be nil, got %v", dst.Configs)
	}
}

func TestStructInMap_DeepCopy(t *testing.T) {
	src := SrcService{
		Name: "copy-service",
		Configs: map[string]SrcConfig{
			"database": {Host: "original"},
		},
	}
	var dst DstService

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := src.Configs["database"]
	cfg.Host = "modified"
	src.Configs["database"] = cfg

	if dst.Configs["database"].Host != "original" {
		t.Errorf("struct in map modification affected destination: got %q", dst.Configs["database"].Host)
	}
}

// Pointer to struct in map
type SrcServiceWithPtrs struct {
	Name    string
	Configs map[string]*SrcConfig
}

type DstServiceWithPtrs struct {
	Name    string
	Configs map[string]*DstConfig
}

func TestPointerStructInMap_Basic(t *testing.T) {
	src := SrcServiceWithPtrs{
		Name: "ptr-service",
		Configs: map[string]*SrcConfig{
			"database": {Host: "localhost", Port: 5432},
			"cache":    {Host: "redis", Port: 6379},
		},
	}
	var dst DstServiceWithPtrs

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Configs) != 2 {
		t.Fatalf("expected 2 configs, got %d", len(dst.Configs))
	}
	if dst.Configs["database"] == nil {
		t.Fatal("expected Configs[database] to not be nil")
	}
	if dst.Configs["database"].Host != "localhost" {
		t.Errorf("expected Configs[database].Host = 'localhost', got %q", dst.Configs["database"].Host)
	}
}

func TestPointerStructInMap_WithNil(t *testing.T) {
	src := SrcServiceWithPtrs{
		Name: "nil-ptr-service",
		Configs: map[string]*SrcConfig{
			"database": {Host: "localhost"},
			"cache":    nil,
		},
	}
	var dst DstServiceWithPtrs

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Configs["database"].Host != "localhost" {
		t.Errorf("expected Configs[database].Host = 'localhost', got %q", dst.Configs["database"].Host)
	}
	if dst.Configs["cache"] != nil {
		t.Errorf("expected Configs[cache] to be nil, got %v", dst.Configs["cache"])
	}
}

func TestPointerStructInMap_DeepCopy(t *testing.T) {
	src := SrcServiceWithPtrs{
		Name: "deep-copy-service",
		Configs: map[string]*SrcConfig{
			"database": {Host: "original"},
		},
	}
	var dst DstServiceWithPtrs

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src.Configs["database"].Host = "modified"

	if dst.Configs["database"].Host != "original" {
		t.Errorf("pointer struct in map modification affected destination: got %q", dst.Configs["database"].Host)
	}
}

// Complex Nested Structures Tests
type SrcNestedItem struct {
	Name  string
	Value int
}

type DstNestedItem struct {
	Name  string
	Value int
}

type SrcNestedSlice struct {
	Items []SrcNestedItem
}

type DstNestedSlice struct {
	Items []DstNestedItem
}

type SrcComplex struct {
	Name    string
	Nested  SrcNestedSlice
	ItemMap map[string][]SrcNestedItem
}

type DstComplex struct {
	Name    string
	Nested  DstNestedSlice
	ItemMap map[string][]DstNestedItem
}

func TestComplex_NestedSliceInStruct(t *testing.T) {
	src := SrcComplex{
		Name: "complex",
		Nested: SrcNestedSlice{
			Items: []SrcNestedItem{
				{Name: "item1", Value: 100},
				{Name: "item2", Value: 200},
			},
		},
	}
	var dst DstComplex

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "complex" {
		t.Errorf("expected Name = 'complex', got %q", dst.Name)
	}
	if len(dst.Nested.Items) != 2 {
		t.Fatalf("expected 2 nested items, got %d", len(dst.Nested.Items))
	}
	if dst.Nested.Items[0].Name != "item1" {
		t.Errorf("expected Nested.Items[0].Name = 'item1', got %q", dst.Nested.Items[0].Name)
	}
	if dst.Nested.Items[1].Value != 200 {
		t.Errorf("expected Nested.Items[1].Value = 200, got %d", dst.Nested.Items[1].Value)
	}
}

func TestComplex_MapOfSlicesOfStructs(t *testing.T) {
	src := SrcComplex{
		Name: "map-slice-complex",
		ItemMap: map[string][]SrcNestedItem{
			"group1": {
				{Name: "a", Value: 1},
				{Name: "b", Value: 2},
			},
			"group2": {
				{Name: "c", Value: 3},
			},
		},
	}
	var dst DstComplex

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.ItemMap) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(dst.ItemMap))
	}
	if len(dst.ItemMap["group1"]) != 2 {
		t.Fatalf("expected 2 items in group1, got %d", len(dst.ItemMap["group1"]))
	}
	if dst.ItemMap["group1"][0].Name != "a" {
		t.Errorf("expected ItemMap[group1][0].Name = 'a', got %q", dst.ItemMap["group1"][0].Name)
	}
	if dst.ItemMap["group2"][0].Value != 3 {
		t.Errorf("expected ItemMap[group2][0].Value = 3, got %d", dst.ItemMap["group2"][0].Value)
	}
}

// Tag Mapping with Nested Structs
type SrcTaggedAddress struct {
	StreetName string `map:"Street"`
	CityName   string `map:"City"`
}

type DstTaggedAddress struct {
	Street string
	City   string
}

type SrcTaggedPerson struct {
	FullName    string           `map:"Name"`
	HomeAddress SrcTaggedAddress `map:"Address"`
}

type DstTaggedPerson struct {
	Name    string
	Address DstTaggedAddress
}

func TestNestedStruct_WithTags(t *testing.T) {
	src := SrcTaggedPerson{
		FullName: "Jane Doe",
		HomeAddress: SrcTaggedAddress{
			StreetName: "100 Tech Blvd",
			CityName:   "San Francisco",
		},
	}
	var dst DstTaggedPerson

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Jane Doe" {
		t.Errorf("expected Name = 'Jane Doe', got %q", dst.Name)
	}
	if dst.Address.Street != "100 Tech Blvd" {
		t.Errorf("expected Address.Street = '100 Tech Blvd', got %q", dst.Address.Street)
	}
	if dst.Address.City != "San Francisco" {
		t.Errorf("expected Address.City = 'San Francisco', got %q", dst.Address.City)
	}
}

// Type Conversion in Nested Structs
type SrcTypedAddress struct {
	Street  string
	ZipCode int32
}

type DstTypedAddress struct {
	Street  string
	ZipCode int64
}

type SrcTypedPerson struct {
	Name    string
	Age     int32
	Address SrcTypedAddress
}

type DstTypedPerson struct {
	Name    string
	Age     int64
	Address DstTypedAddress
}

func TestNestedStruct_TypeConversion(t *testing.T) {
	src := SrcTypedPerson{
		Name: "Type Test",
		Age:  25,
		Address: SrcTypedAddress{
			Street:  "Type St",
			ZipCode: 12345,
		},
	}
	var dst DstTypedPerson

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Age != 25 {
		t.Errorf("expected Age = 25, got %d", dst.Age)
	}
	if dst.Address.ZipCode != 12345 {
		t.Errorf("expected Address.ZipCode = 12345, got %d", dst.Address.ZipCode)
	}
}

// Slice of Structs with Nested Structs
type SrcNestedAddress struct {
	City string
}

type DstNestedAddress struct {
	City string
}

type SrcNestedPerson struct {
	Name    string
	Address SrcNestedAddress
}

type DstNestedPerson struct {
	Name    string
	Address DstNestedAddress
}

type SrcTeam struct {
	Name    string
	Members []SrcNestedPerson
}

type DstTeam struct {
	Name    string
	Members []DstNestedPerson
}

func TestSliceOfStructs_WithNestedStructs(t *testing.T) {
	src := SrcTeam{
		Name: "Engineering",
		Members: []SrcNestedPerson{
			{Name: "Alice", Address: SrcNestedAddress{City: "NYC"}},
			{Name: "Bob", Address: SrcNestedAddress{City: "LA"}},
		},
	}
	var dst DstTeam

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Engineering" {
		t.Errorf("expected Name = 'Engineering', got %q", dst.Name)
	}
	if len(dst.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(dst.Members))
	}
	if dst.Members[0].Name != "Alice" {
		t.Errorf("expected Members[0].Name = 'Alice', got %q", dst.Members[0].Name)
	}
	if dst.Members[0].Address.City != "NYC" {
		t.Errorf("expected Members[0].Address.City = 'NYC', got %q", dst.Members[0].Address.City)
	}
	if dst.Members[1].Address.City != "LA" {
		t.Errorf("expected Members[1].Address.City = 'LA', got %q", dst.Members[1].Address.City)
	}
}

// Test maximum depth limit
func TestMaxDepthLimit_Exceeded(t *testing.T) {
	// Create a deeply nested structure that exceeds a small depth limit
	type Level struct {
		Value int
		Next  *Level
	}

	// Build a chain of 10 levels
	src := &Level{Value: 1}
	current := src
	for i := 2; i <= 10; i++ {
		current.Next = &Level{Value: i}
		current = current.Next
	}

	var dst Level

	// With depth limit of 3, this should fail
	err := MapWithOptions(&dst, *src, WithMaxDepth(3))
	if err == nil {
		t.Fatal("expected error due to depth limit exceeded, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.Reason != "maximum nesting depth exceeded (possible circular reference)" {
		t.Errorf("expected depth exceeded reason, got %q", mappingErr.Reason)
	}
}

func TestMaxDepthLimit_NotExceeded(t *testing.T) {
	// Create a nested structure within limits
	type Level struct {
		Value int
		Next  *Level
	}

	src := Level{
		Value: 1,
		Next: &Level{
			Value: 2,
			Next: &Level{
				Value: 3,
			},
		},
	}

	var dst Level

	// With default depth limit of 64, this should succeed
	err := Map(&dst, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Value != 1 {
		t.Errorf("expected Value = 1, got %d", dst.Value)
	}
	if dst.Next == nil || dst.Next.Value != 2 {
		t.Errorf("expected Next.Value = 2")
	}
	if dst.Next.Next == nil || dst.Next.Next.Value != 3 {
		t.Errorf("expected Next.Next.Value = 3")
	}
}

// Test depth limit with slices containing pointers to structs
func TestMaxDepthLimit_SliceWithPointers(t *testing.T) {
	type Node struct {
		Value    int
		Children []*Node
	}

	// Build a tree-like structure that exceeds depth limit
	src := Node{
		Value: 1,
		Children: []*Node{
			{
				Value: 2,
				Children: []*Node{
					{
						Value: 3,
						Children: []*Node{
							{Value: 4},
						},
					},
				},
			},
		},
	}

	var dst Node

	// With depth limit of 3, this should fail
	err := MapWithOptions(&dst, src, WithMaxDepth(3))
	if err == nil {
		t.Fatal("expected error due to depth limit exceeded in slice, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.Reason != "maximum nesting depth exceeded (possible circular reference)" {
		t.Errorf("expected depth exceeded reason, got %q", mappingErr.Reason)
	}
}

// Test depth limit with maps containing pointers to structs
func TestMaxDepthLimit_MapWithPointers(t *testing.T) {
	type Node struct {
		Value   int
		Related map[string]*Node
	}

	// Build a structure with nested maps that exceeds depth limit
	src := Node{
		Value: 1,
		Related: map[string]*Node{
			"child": {
				Value: 2,
				Related: map[string]*Node{
					"grandchild": {
						Value: 3,
						Related: map[string]*Node{
							"great-grandchild": {Value: 4},
						},
					},
				},
			},
		},
	}

	var dst Node

	// With depth limit of 3, this should fail
	err := MapWithOptions(&dst, src, WithMaxDepth(3))
	if err == nil {
		t.Fatal("expected error due to depth limit exceeded in map, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.Reason != "maximum nesting depth exceeded (possible circular reference)" {
		t.Errorf("expected depth exceeded reason, got %q", mappingErr.Reason)
	}
}

// Test depth limit with deeply nested slices (no structs)
func TestMaxDepthLimit_NestedSlices(t *testing.T) {
	type Container struct {
		Data [][][]int
	}

	src := Container{
		Data: [][][]int{
			{
				{1, 2, 3},
				{4, 5, 6},
			},
		},
	}

	var dst Container

	// With depth limit of 2, nested slices should fail
	err := MapWithOptions(&dst, src, WithMaxDepth(2))
	if err == nil {
		t.Fatal("expected error due to depth limit exceeded in nested slices, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.Reason != "maximum nesting depth exceeded (possible circular reference)" {
		t.Errorf("expected depth exceeded reason, got %q", mappingErr.Reason)
	}
}

// Test depth limit with deeply nested maps (no structs)
func TestMaxDepthLimit_NestedMaps(t *testing.T) {
	type Container struct {
		Data map[string]map[string]map[string]int
	}

	src := Container{
		Data: map[string]map[string]map[string]int{
			"a": {
				"b": {
					"c": 1,
				},
			},
		},
	}

	var dst Container

	// With depth limit of 2, nested maps should fail
	err := MapWithOptions(&dst, src, WithMaxDepth(2))
	if err == nil {
		t.Fatal("expected error due to depth limit exceeded in nested maps, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.Reason != "maximum nesting depth exceeded (possible circular reference)" {
		t.Errorf("expected depth exceeded reason, got %q", mappingErr.Reason)
	}
}

// Test depth limit with pointer chains
func TestMaxDepthLimit_PointerChain(t *testing.T) {
	type Box struct {
		Value int
		Inner **Box
	}

	innermost := &Box{Value: 3}
	inner := &Box{Value: 2, Inner: &innermost}
	src := Box{Value: 1, Inner: &inner}

	var dst Box

	// With depth limit of 2, pointer chain should fail
	err := MapWithOptions(&dst, src, WithMaxDepth(2))
	if err == nil {
		t.Fatal("expected error due to depth limit exceeded in pointer chain, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.Reason != "maximum nesting depth exceeded (possible circular reference)" {
		t.Errorf("expected depth exceeded reason, got %q", mappingErr.Reason)
	}
}

// Test that custom MaxDepth option works correctly
func TestMaxDepthLimit_CustomDepth(t *testing.T) {
	type Level struct {
		Value int
		Next  *Level
	}

	// Build a chain of 5 levels
	src := Level{Value: 1}
	current := &src
	for i := 2; i <= 5; i++ {
		current.Next = &Level{Value: i}
		current = current.Next
	}

	var dst Level

	// With depth limit of 10, this should succeed (5 levels < 10)
	err := MapWithOptions(&dst, src, WithMaxDepth(10))
	if err != nil {
		t.Fatalf("unexpected error with sufficient depth limit: %v", err)
	}

	// Verify the chain was copied correctly
	d := &dst
	for i := 1; i <= 5; i++ {
		if d == nil {
			t.Fatalf("expected level %d, got nil", i)
		}
		if d.Value != i {
			t.Errorf("expected Value = %d, got %d", i, d.Value)
		}
		d = d.Next
	}
}

// Test depth limit with mixed nested structures
func TestMaxDepthLimit_MixedStructures(t *testing.T) {
	type Inner struct {
		Value int
	}

	type Container struct {
		Slices [][]Inner
		Maps   map[string]map[string]Inner
	}

	src := Container{
		Slices: [][]Inner{
			{{Value: 1}, {Value: 2}},
		},
		Maps: map[string]map[string]Inner{
			"outer": {
				"inner": {Value: 3},
			},
		},
	}

	var dst Container

	// With sufficient depth, this should succeed
	err := MapWithOptions(&dst, src, WithMaxDepth(10))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Slices) != 1 || len(dst.Slices[0]) != 2 {
		t.Errorf("expected Slices structure to be preserved")
	}
	if dst.Slices[0][0].Value != 1 {
		t.Errorf("expected Slices[0][0].Value = 1, got %d", dst.Slices[0][0].Value)
	}
	if dst.Maps["outer"]["inner"].Value != 3 {
		t.Errorf("expected Maps[outer][inner].Value = 3, got %d", dst.Maps["outer"]["inner"].Value)
	}
}
