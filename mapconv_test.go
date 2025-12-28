package mapper

import (
	"testing"
)

// TestMapconv_StringToInt tests converting string to int using mapconv tag.
func TestMapconv_StringToInt(t *testing.T) {
	type Src struct {
		Age string `mapconv:"int"`
	}
	type Dst struct {
		Age int
	}

	src := Src{Age: "25"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Age != 25 {
		t.Errorf("expected Age = 25, got %d", dst.Age)
	}
}

// TestMapconv_StringToFloat64 tests converting string to float64 using mapconv tag.
func TestMapconv_StringToFloat64(t *testing.T) {
	type Src struct {
		Score string `mapconv:"float64"`
	}
	type Dst struct {
		Score float64
	}

	src := Src{Score: "95.5"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Score != 95.5 {
		t.Errorf("expected Score = 95.5, got %f", dst.Score)
	}
}

// TestMapconv_StringToBool tests converting string to bool using mapconv tag.
func TestMapconv_StringToBool(t *testing.T) {
	type Src struct {
		Active string `mapconv:"bool"`
	}
	type Dst struct {
		Active bool
	}

	src := Src{Active: "true"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Active != true {
		t.Errorf("expected Active = true, got %v", dst.Active)
	}
}

// TestMapconv_StringToBool_False tests converting "false" string to bool.
func TestMapconv_StringToBool_False(t *testing.T) {
	type Src struct {
		Active string `mapconv:"bool"`
	}
	type Dst struct {
		Active bool
	}

	src := Src{Active: "false"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Active != false {
		t.Errorf("expected Active = false, got %v", dst.Active)
	}
}

// TestMapconv_StringToInt64 tests converting string to int64 using mapconv tag.
func TestMapconv_StringToInt64(t *testing.T) {
	type Src struct {
		Count string `mapconv:"int64"`
	}
	type Dst struct {
		Count int64
	}

	src := Src{Count: "9223372036854775807"} // Max int64
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Count != 9223372036854775807 {
		t.Errorf("expected Count = 9223372036854775807, got %d", dst.Count)
	}
}

// TestMapconv_StringToUint tests converting string to uint using mapconv tag.
func TestMapconv_StringToUint(t *testing.T) {
	type Src struct {
		Amount string `mapconv:"uint"`
	}
	type Dst struct {
		Amount uint
	}

	src := Src{Amount: "42"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Amount != 42 {
		t.Errorf("expected Amount = 42, got %d", dst.Amount)
	}
}

// TestMapconv_NoConversion tests that fields without mapconv tag are not converted.
func TestMapconv_NoConversion(t *testing.T) {
	type Src struct {
		Name string
	}
	type Dst struct {
		Name string
	}

	src := Src{Name: "John Doe"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "John Doe" {
		t.Errorf("expected Name = 'John Doe', got %q", dst.Name)
	}
}

// TestMapconv_AllFieldsTogether tests converting multiple fields at once.
func TestMapconv_AllFieldsTogether(t *testing.T) {
	type Src struct {
		Age    string `mapconv:"int"`
		Score  string `mapconv:"float64"`
		Active string `mapconv:"bool"`
		Count  string `mapconv:"int64"`
		Amount string `mapconv:"uint"`
		Name   string
	}
	type Dst struct {
		Age    int
		Score  float64
		Active bool
		Count  int64
		Amount uint
		Name   string
	}

	src := Src{
		Age:    "30",
		Score:  "88.5",
		Active: "true",
		Count:  "1000",
		Amount: "500",
		Name:   "Alice",
	}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Age != 30 {
		t.Errorf("expected Age = 30, got %d", dst.Age)
	}
	if dst.Score != 88.5 {
		t.Errorf("expected Score = 88.5, got %f", dst.Score)
	}
	if dst.Active != true {
		t.Errorf("expected Active = true, got %v", dst.Active)
	}
	if dst.Count != 1000 {
		t.Errorf("expected Count = 1000, got %d", dst.Count)
	}
	if dst.Amount != 500 {
		t.Errorf("expected Amount = 500, got %d", dst.Amount)
	}
	if dst.Name != "Alice" {
		t.Errorf("expected Name = 'Alice', got %q", dst.Name)
	}
}

// TestMapconv_InvalidInt tests error handling for invalid int conversion.
func TestMapconv_InvalidInt(t *testing.T) {
	type Src struct {
		Age string `mapconv:"int"`
	}
	type Dst struct {
		Age int
	}

	src := Src{Age: "not-a-number"}
	var dst Dst

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error for invalid int conversion, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.FieldPath != "Age" {
		t.Errorf("expected FieldPath = 'Age', got %q", mappingErr.FieldPath)
	}
}

// TestMapconv_InvalidFloat tests error handling for invalid float conversion.
func TestMapconv_InvalidFloat(t *testing.T) {
	type Src struct {
		Score string `mapconv:"float64"`
	}
	type Dst struct {
		Score float64
	}

	src := Src{Score: "not-a-float"}
	var dst Dst

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error for invalid float conversion, got nil")
	}
}

// TestMapconv_InvalidBool tests error handling for invalid bool conversion.
func TestMapconv_InvalidBool(t *testing.T) {
	type Src struct {
		Active string `mapconv:"bool"`
	}
	type Dst struct {
		Active bool
	}

	src := Src{Active: "maybe"}
	var dst Dst

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error for invalid bool conversion, got nil")
	}
}

// TestMapconv_NegativeInt tests converting negative string to int.
func TestMapconv_NegativeInt(t *testing.T) {
	type Src struct {
		Age string `mapconv:"int"`
	}
	type Dst struct {
		Age int
	}

	src := Src{Age: "-10"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Age != -10 {
		t.Errorf("expected Age = -10, got %d", dst.Age)
	}
}

// TestMapconv_NegativeFloat tests converting negative string to float.
func TestMapconv_NegativeFloat(t *testing.T) {
	type Src struct {
		Score string `mapconv:"float64"`
	}
	type Dst struct {
		Score float64
	}

	src := Src{Score: "-95.5"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Score != -95.5 {
		t.Errorf("expected Score = -95.5, got %f", dst.Score)
	}
}

// TestMapconv_EmptyString tests converting empty string (should error).
func TestMapconv_EmptyString(t *testing.T) {
	type Src struct {
		Age string `mapconv:"int"`
	}
	type Dst struct {
		Age int
	}

	src := Src{Age: ""}
	var dst Dst

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error for empty string conversion, got nil")
	}
}

// TestMapconv_BoolVariants tests various bool string representations.
func TestMapconv_BoolVariants(t *testing.T) {
	type Src struct {
		Active string `mapconv:"bool"`
	}
	type Dst struct {
		Active bool
	}

	testCases := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{"t", true},
		{"f", false},
		{"T", true},
		{"F", false},
		{"TRUE", true},
		{"FALSE", false},
	}

	for _, tc := range testCases {
		src := Src{Active: tc.input}
		var dst Dst

		if err := Map(&dst, src); err != nil {
			t.Errorf("unexpected error for input %q: %v", tc.input, err)
			continue
		}

		if dst.Active != tc.expected {
			t.Errorf("input %q: expected Active = %v, got %v", tc.input, tc.expected, dst.Active)
		}
	}
}

// TestMapconv_WithMapTag tests mapconv working together with map tag.
func TestMapconv_WithMapTag(t *testing.T) {
	type Src struct {
		UserAge string `map:"Age" mapconv:"int"`
	}
	type Dst struct {
		Age int
	}

	src := Src{UserAge: "35"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Age != 35 {
		t.Errorf("expected Age = 35, got %d", dst.Age)
	}
}

// TestMapconv_AllIntTypes tests all integer type conversions.
func TestMapconv_AllIntTypes(t *testing.T) {
	type Src struct {
		Int8Val  string `mapconv:"int8"`
		Int16Val string `mapconv:"int16"`
		Int32Val string `mapconv:"int32"`
	}
	type Dst struct {
		Int8Val  int8
		Int16Val int16
		Int32Val int32
	}

	src := Src{
		Int8Val:  "127",
		Int16Val: "32767",
		Int32Val: "2147483647",
	}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Int8Val != 127 {
		t.Errorf("expected Int8Val = 127, got %d", dst.Int8Val)
	}
	if dst.Int16Val != 32767 {
		t.Errorf("expected Int16Val = 32767, got %d", dst.Int16Val)
	}
	if dst.Int32Val != 2147483647 {
		t.Errorf("expected Int32Val = 2147483647, got %d", dst.Int32Val)
	}
}

// TestMapconv_AllUintTypes tests all unsigned integer type conversions.
func TestMapconv_AllUintTypes(t *testing.T) {
	type Src struct {
		Uint8Val  string `mapconv:"uint8"`
		Uint16Val string `mapconv:"uint16"`
		Uint32Val string `mapconv:"uint32"`
		Uint64Val string `mapconv:"uint64"`
	}
	type Dst struct {
		Uint8Val  uint8
		Uint16Val uint16
		Uint32Val uint32
		Uint64Val uint64
	}

	src := Src{
		Uint8Val:  "255",
		Uint16Val: "65535",
		Uint32Val: "4294967295",
		Uint64Val: "18446744073709551615",
	}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Uint8Val != 255 {
		t.Errorf("expected Uint8Val = 255, got %d", dst.Uint8Val)
	}
	if dst.Uint16Val != 65535 {
		t.Errorf("expected Uint16Val = 65535, got %d", dst.Uint16Val)
	}
	if dst.Uint32Val != 4294967295 {
		t.Errorf("expected Uint32Val = 4294967295, got %d", dst.Uint32Val)
	}
	if dst.Uint64Val != 18446744073709551615 {
		t.Errorf("expected Uint64Val = 18446744073709551615, got %d", dst.Uint64Val)
	}
}

// TestMapconv_Float32 tests float32 conversion.
func TestMapconv_Float32(t *testing.T) {
	type Src struct {
		Value string `mapconv:"float32"`
	}
	type Dst struct {
		Value float32
	}

	src := Src{Value: "3.14"}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Use approximate comparison for float32.
	if dst.Value < 3.13 || dst.Value > 3.15 {
		t.Errorf("expected Value ≈ 3.14, got %f", dst.Value)
	}
}

// TestMapconv_UnsupportedType tests error for unsupported target type.
func TestMapconv_UnsupportedType(t *testing.T) {
	type Src struct {
		Value string `mapconv:"complex128"`
	}
	type Dst struct {
		Value complex128
	}

	src := Src{Value: "1+2i"}
	var dst Dst

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error for unsupported target type, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.Reason == "" {
		t.Error("expected non-empty Reason in error")
	}
}
