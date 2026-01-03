package mapper

import (
	"strings"
	"testing"
)

// Basic Map() function tests
func TestMap_BasicMapping(t *testing.T) {
	type Src struct {
		Name  string
		Age   int
		Email string
	}
	type Dst struct {
		Name  string
		Age   int
		Email string
	}

	src := Src{Name: "Alice", Age: 30, Email: "alice@example.com"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Alice" {
		t.Errorf("expected Name = 'Alice', got %q", dst.Name)
	}
	if dst.Age != 30 {
		t.Errorf("expected Age = 30, got %d", dst.Age)
	}
	if dst.Email != "alice@example.com" {
		t.Errorf("expected Email = 'alice@example.com', got %q", dst.Email)
	}
}

func TestMap_SourceAsPointer(t *testing.T) {
	type Src struct {
		Name string
	}
	type Dst struct {
		Name string
	}

	src := &Src{Name: "Bob"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Bob" {
		t.Errorf("expected Name = 'Bob', got %q", dst.Name)
	}
}

func TestMap_PartialMapping(t *testing.T) {
	type Src struct {
		Name string
		Age  int
	}
	type Dst struct {
		Name    string
		Age     int
		Country string // Not in source
	}

	src := Src{Name: "Charlie", Age: 25}
	dst := Dst{Country: "USA"}

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Charlie" {
		t.Errorf("expected Name = 'Charlie', got %q", dst.Name)
	}
	if dst.Age != 25 {
		t.Errorf("expected Age = 25, got %d", dst.Age)
	}
	if dst.Country != "USA" {
		t.Errorf("expected Country = 'USA' (preserved), got %q", dst.Country)
	}
}

// Error cases for Map()
func TestMap_NilDst(t *testing.T) {
	type Src struct {
		Name string
	}

	src := Src{Name: "Test"}
	err := Map(nil, src)

	if err == nil {
		t.Fatal("expected error for nil dst, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if !strings.Contains(mappingErr.Reason, "nil") {
		t.Errorf("error should mention 'nil', got: %q", mappingErr.Reason)
	}
}

func TestMap_NilSrc(t *testing.T) {
	type Dst struct {
		Name string
	}

	var dst Dst
	err := Map(&dst, nil)

	if err == nil {
		t.Fatal("expected error for nil src, got nil")
	}
}

func TestMap_DstNotPointer(t *testing.T) {
	type Src struct {
		Name string
	}
	type Dst struct {
		Name string
	}

	src := Src{Name: "Test"}
	var dst Dst

	err := Map(dst, src) // dst instead of &dst

	if err == nil {
		t.Fatal("expected error for non-pointer dst, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if !strings.Contains(mappingErr.Reason, "pointer") {
		t.Errorf("error should mention 'pointer', got: %q", mappingErr.Reason)
	}
}

func TestMap_SrcNilPointer(t *testing.T) {
	type Src struct {
		Name string
	}
	type Dst struct {
		Name string
	}

	var src *Src = nil
	var dst Dst

	err := Map(&dst, src)

	if err == nil {
		t.Fatal("expected error for nil pointer src, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if !strings.Contains(mappingErr.Reason, "nil") {
		t.Errorf("error should mention 'nil', got: %q", mappingErr.Reason)
	}
}

func TestMap_SrcNotStruct(t *testing.T) {
	type Dst struct {
		Name string
	}

	src := "not a struct"
	var dst Dst

	err := Map(&dst, src)

	if err == nil {
		t.Fatal("expected error for non-struct src, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if !strings.Contains(mappingErr.Reason, "struct") {
		t.Errorf("error should mention 'struct', got: %q", mappingErr.Reason)
	}
}

func TestMap_DstNotStruct(t *testing.T) {
	type Src struct {
		Name string
	}

	src := Src{Name: "Test"}
	dst := "not a struct"

	err := Map(&dst, src)

	if err == nil {
		t.Fatal("expected error for non-struct dst, got nil")
	}
}

// Options tests
func TestMapWithOptions_WithTagName(t *testing.T) {
	type Src struct {
		FullName string `custom:"Name"`
	}
	type Dst struct {
		Name string
	}

	src := Src{FullName: "Diana"}
	var dst Dst

	// Using "custom" as tag name
	if err := MapWithOptions(&dst, src, WithTagName("custom")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Diana" {
		t.Errorf("expected Name = 'Diana', got %q", dst.Name)
	}
}

func TestMapWithOptions_WithIgnoreZeroSource(t *testing.T) {
	type Src struct {
		Name  string
		Age   int
		Email string
	}
	type Dst struct {
		Name  string
		Age   int
		Email string
	}

	src := Src{Name: "Eve", Age: 0, Email: ""} // Age and Email are zero values
	dst := Dst{Name: "Old", Age: 50, Email: "old@example.com"}

	if err := MapWithOptions(&dst, src, WithIgnoreZeroSource()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Eve" {
		t.Errorf("expected Name = 'Eve', got %q", dst.Name)
	}
	if dst.Age != 50 {
		t.Errorf("expected Age = 50 (preserved), got %d", dst.Age)
	}
	if dst.Email != "old@example.com" {
		t.Errorf("expected Email = 'old@example.com' (preserved), got %q", dst.Email)
	}
}

func TestMapWithOptions_WithStrictMode(t *testing.T) {
	type Src struct {
		Name string
	}
	type Dst struct {
		Name  string
		Age   int // Not in source
		Email string
	}

	src := Src{Name: "Frank"}
	var dst Dst

	err := MapWithOptions(&dst, src, WithStrictMode())

	if err == nil {
		t.Fatal("expected error in strict mode for missing field, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if !strings.Contains(mappingErr.Reason, "no matching") {
		t.Errorf("error should mention 'no matching', got: %q", mappingErr.Reason)
	}
}

func TestMapWithOptions_MultipleOptions(t *testing.T) {
	type Src struct {
		UserName string `custom:"Name"`
		UserAge  int    `custom:"Age"`
	}
	type Dst struct {
		Name string
		Age  int
	}

	src := Src{UserName: "Grace", UserAge: 0}
	dst := Dst{Name: "Old", Age: 99}

	err := MapWithOptions(&dst, src, WithTagName("custom"), WithIgnoreZeroSource())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Grace" {
		t.Errorf("expected Name = 'Grace', got %q", dst.Name)
	}
	if dst.Age != 99 {
		t.Errorf("expected Age = 99 (preserved due to zero source), got %d", dst.Age)
	}
}

// =============================================================================
// Pointer field handling tests
// =============================================================================

func TestMap_ValueToPointer(t *testing.T) {
	type Src struct {
		Name  string
		Value int
	}
	type Dst struct {
		Name  *string
		Value *int
	}

	src := Src{Name: "Henry", Value: 42}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name == nil {
		t.Fatal("expected Name to be non-nil")
	}
	if *dst.Name != "Henry" {
		t.Errorf("expected *Name = 'Henry', got %q", *dst.Name)
	}
	if dst.Value == nil {
		t.Fatal("expected Value to be non-nil")
	}
	if *dst.Value != 42 {
		t.Errorf("expected *Value = 42, got %d", *dst.Value)
	}
}

func TestMap_PointerToValue(t *testing.T) {
	type Src struct {
		Name  *string
		Value *int
	}
	type Dst struct {
		Name  string
		Value int
	}

	name := "Ivy"
	value := 100
	src := Src{Name: &name, Value: &value}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Ivy" {
		t.Errorf("expected Name = 'Ivy', got %q", dst.Name)
	}
	if dst.Value != 100 {
		t.Errorf("expected Value = 100, got %d", dst.Value)
	}
}

func TestMap_PointerToValue_NilPointer(t *testing.T) {
	type Src struct {
		Name  *string
		Value *int
	}
	type Dst struct {
		Name  string
		Value int
	}

	src := Src{Name: nil, Value: nil}
	dst := Dst{Name: "Original", Value: 999}

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Nil pointers should not overwrite destination
	if dst.Name != "Original" {
		t.Errorf("expected Name = 'Original' (preserved), got %q", dst.Name)
	}
	if dst.Value != 999 {
		t.Errorf("expected Value = 999 (preserved), got %d", dst.Value)
	}
}

// Type conversion tests
func TestMap_ConvertibleTypes(t *testing.T) {
	type Src struct {
		IntVal   int
		FloatVal float32
	}
	type Dst struct {
		IntVal   int64
		FloatVal float64
	}

	src := Src{IntVal: 42, FloatVal: 3.14}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.IntVal != 42 {
		t.Errorf("expected IntVal = 42, got %d", dst.IntVal)
	}
	if dst.FloatVal < 3.13 || dst.FloatVal > 3.15 {
		t.Errorf("expected FloatVal ≈ 3.14, got %f", dst.FloatVal)
	}
}

func TestMap_IncompatibleTypes(t *testing.T) {
	type Src struct {
		Value string
	}
	type Dst struct {
		Value int // string not convertible to int without mapconv
	}

	src := Src{Value: "not-a-number"}
	var dst Dst

	err := Map(&dst, src)

	if err == nil {
		t.Fatal("expected error for incompatible types, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.FieldPath != "Value" {
		t.Errorf("expected FieldPath = 'Value', got %q", mappingErr.FieldPath)
	}
}

// MappingError tests
func TestMappingError_Error(t *testing.T) {
	err := &MappingError{
		SrcType:   "SourceStruct",
		DstType:   "DestStruct",
		FieldPath: "FieldName",
		Reason:    "test reason",
	}

	errMsg := err.Error()

	if !strings.Contains(errMsg, "SourceStruct") {
		t.Errorf("error message should contain SrcType, got: %q", errMsg)
	}
	if !strings.Contains(errMsg, "DestStruct") {
		t.Errorf("error message should contain DstType, got: %q", errMsg)
	}
	if !strings.Contains(errMsg, "FieldName") {
		t.Errorf("error message should contain FieldPath, got: %q", errMsg)
	}
	if !strings.Contains(errMsg, "test reason") {
		t.Errorf("error message should contain Reason, got: %q", errMsg)
	}
}

func TestMappingError_EmptyFieldPath(t *testing.T) {
	err := &MappingError{
		SrcType:   "Src",
		DstType:   "Dst",
		FieldPath: "",
		Reason:    "nil pointer",
	}

	errMsg := err.Error()

	// Should still be valid error message
	if errMsg == "" {
		t.Error("error message should not be empty")
	}
}

// Edge cases
func TestMap_EmptyStructs(t *testing.T) {
	type Src struct{}
	type Dst struct{}

	src := Src{}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMap_UnexportedFields(t *testing.T) {
	type Src struct {
		Public  string
		private string //nolint:unused
	}
	type Dst struct {
		Public  string
		private string //nolint:unused
	}

	src := Src{Public: "visible"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Public != "visible" {
		t.Errorf("expected Public = 'visible', got %q", dst.Public)
	}
}

func TestMap_DifferentFieldOrder(t *testing.T) {
	type Src struct {
		A string
		B int
		C bool
	}
	type Dst struct {
		C bool
		A string
		B int
	}

	src := Src{A: "alpha", B: 42, C: true}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.A != "alpha" {
		t.Errorf("expected A = 'alpha', got %q", dst.A)
	}
	if dst.B != 42 {
		t.Errorf("expected B = 42, got %d", dst.B)
	}
	if !dst.C {
		t.Errorf("expected C = true, got %v", dst.C)
	}
}

func TestMap_FieldNameCaseSensitive(t *testing.T) {
	type Src struct {
		Name string
		NAME string
	}
	type Dst struct {
		Name string
		NAME string
	}

	src := Src{Name: "lower", NAME: "UPPER"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "lower" {
		t.Errorf("expected Name = 'lower', got %q", dst.Name)
	}
	if dst.NAME != "UPPER" {
		t.Errorf("expected NAME = 'UPPER', got %q", dst.NAME)
	}
}
