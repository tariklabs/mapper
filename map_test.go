package mapper

import (
	"testing"
)

type SourceWithMap struct {
	Labels   map[string]string
	Counts   map[string]int
	Metadata map[string]any
}

type DestWithMap struct {
	Labels   map[string]string
	Counts   map[string]int
	Metadata map[string]any
}

// TestMap_SameType_StringString tests mapping map[string]string -> map[string]string.
func TestMap_SameType_StringString(t *testing.T) {
	src := SourceWithMap{
		Labels: map[string]string{"env": "prod", "region": "us-east"},
	}
	var dst DestWithMap

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Labels) != len(src.Labels) {
		t.Errorf("expected %d labels, got %d", len(src.Labels), len(dst.Labels))
	}

	for k, v := range src.Labels {
		if dst.Labels[k] != v {
			t.Errorf("expected Labels[%q] = %q, got %q", k, v, dst.Labels[k])
		}
	}
}

// TestMap_SameType_StringInt tests mapping map[string]int -> map[string]int.
func TestMap_SameType_StringInt(t *testing.T) {
	src := SourceWithMap{
		Counts: map[string]int{"apples": 10, "oranges": 20, "bananas": 15},
	}
	var dst DestWithMap

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Counts) != len(src.Counts) {
		t.Errorf("expected %d counts, got %d", len(src.Counts), len(dst.Counts))
	}

	for k, v := range src.Counts {
		if dst.Counts[k] != v {
			t.Errorf("expected Counts[%q] = %d, got %d", k, v, dst.Counts[k])
		}
	}
}

// TestMap_NilMap tests that nil maps remain nil after mapping.
func TestMap_NilMap(t *testing.T) {
	src := SourceWithMap{
		Labels: nil,
	}
	var dst DestWithMap

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Labels != nil {
		t.Errorf("expected nil map, got %v", dst.Labels)
	}
}

// TestMap_EmptyMap tests that empty maps remain empty (not nil) after mapping.
func TestMap_EmptyMap(t *testing.T) {
	src := SourceWithMap{
		Labels: map[string]string{},
	}
	var dst DestWithMap

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Labels == nil {
		t.Error("expected empty map, got nil")
	}

	if len(dst.Labels) != 0 {
		t.Errorf("expected empty map, got %d elements", len(dst.Labels))
	}
}

// TestMap_DeepCopy_NoSharedReference verifies that modifying source doesn't affect destination.
func TestMap_DeepCopy_NoSharedReference(t *testing.T) {
	src := SourceWithMap{
		Labels: map[string]string{"key": "original"},
	}
	var dst DestWithMap

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Modify source after mapping.
	src.Labels["key"] = "modified"
	src.Labels["new"] = "added"

	// Destination should be unaffected.
	if dst.Labels["key"] != "original" {
		t.Errorf("expected Labels[key] = 'original', got %q (source modification affected destination)", dst.Labels["key"])
	}

	if _, exists := dst.Labels["new"]; exists {
		t.Error("new key should not exist in destination (source modification affected destination)")
	}
}

// Different type map tests (value type conversion)
type SourceWithInt32Map struct {
	Values map[string]int32
}

type DestWithInt64Map struct {
	Values map[string]int64
}

// TestMap_DifferentValueType_Int32ToInt64 tests converting map[string]int32 -> map[string]int64.
func TestMap_DifferentValueType_Int32ToInt64(t *testing.T) {
	src := SourceWithInt32Map{
		Values: map[string]int32{"a": 100, "b": 200, "c": 300},
	}
	var dst DestWithInt64Map

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Values) != len(src.Values) {
		t.Errorf("expected %d values, got %d", len(src.Values), len(dst.Values))
	}

	for k, v := range src.Values {
		if dst.Values[k] != int64(v) {
			t.Errorf("expected Values[%q] = %d, got %d", k, v, dst.Values[k])
		}
	}
}

type SourceWithFloat32Map struct {
	Scores map[string]float32
}

type DestWithFloat64Map struct {
	Scores map[string]float64
}

// TestMap_DifferentValueType_Float32ToFloat64 tests converting map[string]float32 -> map[string]float64.
func TestMap_DifferentValueType_Float32ToFloat64(t *testing.T) {
	src := SourceWithFloat32Map{
		Scores: map[string]float32{"math": 95.5, "english": 88.0},
	}
	var dst DestWithFloat64Map

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Scores) != len(src.Scores) {
		t.Errorf("expected %d scores, got %d", len(src.Scores), len(dst.Scores))
	}

	for k, v := range src.Scores {
		// Use approximate comparison for floats.
		expected := float64(v)
		if dst.Scores[k] < expected-0.01 || dst.Scores[k] > expected+0.01 {
			t.Errorf("expected Scores[%q] ≈ %f, got %f", k, expected, dst.Scores[k])
		}
	}
}

// Different key type tests
type SourceWithIntKeyMap struct {
	Data map[int]string
}

type DestWithInt64KeyMap struct {
	Data map[int64]string
}

// TestMap_DifferentKeyType_IntToInt64 tests converting map[int]string -> map[int64]string.
func TestMap_DifferentKeyType_IntToInt64(t *testing.T) {
	src := SourceWithIntKeyMap{
		Data: map[int]string{1: "one", 2: "two", 3: "three"},
	}
	var dst DestWithInt64KeyMap

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Data) != len(src.Data) {
		t.Errorf("expected %d entries, got %d", len(src.Data), len(dst.Data))
	}

	for k, v := range src.Data {
		if dst.Data[int64(k)] != v {
			t.Errorf("expected Data[%d] = %q, got %q", k, v, dst.Data[int64(k)])
		}
	}
}

// Nested map tests
type SourceWithNestedMap struct {
	Config map[string]map[string]string
}

type DestWithNestedMap struct {
	Config map[string]map[string]string
}

// TestMap_NestedMaps tests mapping nested maps.
func TestMap_NestedMaps(t *testing.T) {
	src := SourceWithNestedMap{
		Config: map[string]map[string]string{
			"database": {"host": "localhost", "port": "5432"},
			"cache":    {"host": "redis", "port": "6379"},
		},
	}
	var dst DestWithNestedMap

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Config) != len(src.Config) {
		t.Errorf("expected %d config entries, got %d", len(src.Config), len(dst.Config))
	}

	for outerKey, innerMap := range src.Config {
		dstInner, exists := dst.Config[outerKey]
		if !exists {
			t.Errorf("expected Config[%q] to exist", outerKey)
			continue
		}
		for innerKey, v := range innerMap {
			if dstInner[innerKey] != v {
				t.Errorf("expected Config[%q][%q] = %q, got %q", outerKey, innerKey, v, dstInner[innerKey])
			}
		}
	}
}

// TestMap_NestedMaps_DeepCopy verifies nested maps are deep copied.
func TestMap_NestedMaps_DeepCopy(t *testing.T) {
	src := SourceWithNestedMap{
		Config: map[string]map[string]string{
			"database": {"host": "localhost"},
		},
	}
	var dst DestWithNestedMap

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src.Config["database"]["host"] = "modified"
	src.Config["database"]["new"] = "added"

	if dst.Config["database"]["host"] != "localhost" {
		t.Errorf("nested map modification affected destination: got %q", dst.Config["database"]["host"])
	}

	if _, exists := dst.Config["database"]["new"]; exists {
		t.Error("new nested key should not exist in destination")
	}
}

// Map with slice values tests
type SourceWithMapOfSlices struct {
	Groups map[string][]string
}

type DestWithMapOfSlices struct {
	Groups map[string][]string
}

// TestMap_WithSliceValues tests mapping maps with slice values.
func TestMap_WithSliceValues(t *testing.T) {
	src := SourceWithMapOfSlices{
		Groups: map[string][]string{
			"admins": {"alice", "bob"},
			"users":  {"charlie", "diana", "eve"},
		},
	}
	var dst DestWithMapOfSlices

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Groups) != len(src.Groups) {
		t.Errorf("expected %d groups, got %d", len(src.Groups), len(dst.Groups))
	}

	for k, v := range src.Groups {
		if len(dst.Groups[k]) != len(v) {
			t.Errorf("expected Groups[%q] length = %d, got %d", k, len(v), len(dst.Groups[k]))
		}
		for i, member := range v {
			if dst.Groups[k][i] != member {
				t.Errorf("expected Groups[%q][%d] = %q, got %q", k, i, member, dst.Groups[k][i])
			}
		}
	}
}

// TestMap_WithSliceValues_DeepCopy verifies slice values in maps are deep copied.
func TestMap_WithSliceValues_DeepCopy(t *testing.T) {
	src := SourceWithMapOfSlices{
		Groups: map[string][]string{
			"admins": {"alice", "bob"},
		},
	}
	var dst DestWithMapOfSlices

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src.Groups["admins"][0] = "modified"
	src.Groups["admins"] = append(src.Groups["admins"], "new")

	if dst.Groups["admins"][0] != "alice" {
		t.Errorf("slice modification affected destination: got %q", dst.Groups["admins"][0])
	}

	if len(dst.Groups["admins"]) != 2 {
		t.Errorf("destination slice length changed: got %d", len(dst.Groups["admins"]))
	}
}

// TestMap_WithPointerValues tests that pointer values are deep copied.
func TestMap_WithPointerValues(t *testing.T) {
	type Value struct {
		Name string
	}
	type Src struct {
		Data map[string]*Value
	}
	type Dst struct {
		Data map[string]*Value
	}

	original := &Value{Name: "original"}
	src := Src{
		Data: map[string]*Value{"key": original},
	}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Data["key"] == nil {
		t.Fatal("expected pointer value to be copied")
	}

	if dst.Data["key"].Name != "original" {
		t.Errorf("expected Name = 'original', got %q", dst.Data["key"].Name)
	}

	src.Data["key"].Name = "modified"

	if dst.Data["key"].Name != "original" {
		t.Errorf("pointer modification affected destination: got %q (expected 'original')", dst.Data["key"].Name)
	}
}

// TestMap_WithPointerValues_Nil tests nil pointer values in maps.
func TestMap_WithPointerValues_Nil(t *testing.T) {
	type Value struct {
		Name string
	}
	type Src struct {
		Data map[string]*Value
	}
	type Dst struct {
		Data map[string]*Value
	}

	src := Src{
		Data: map[string]*Value{"key": nil},
	}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Data["key"] != nil {
		t.Errorf("expected nil pointer, got %v", dst.Data["key"])
	}
}

// Additional same-type tests
type SourceWithIntKeyStringValue struct {
	Data map[int]string
}

type DestWithIntKeyStringValue struct {
	Data map[int]string
}

// TestMap_SameType_IntKeyStringValue tests mapping map[int]string -> map[int]string.
func TestMap_SameType_IntKeyStringValue(t *testing.T) {
	src := SourceWithIntKeyStringValue{
		Data: map[int]string{1: "one", 2: "two", 3: "three"},
	}
	var dst DestWithIntKeyStringValue

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Data) != len(src.Data) {
		t.Errorf("expected %d entries, got %d", len(src.Data), len(dst.Data))
	}

	for k, v := range src.Data {
		if dst.Data[k] != v {
			t.Errorf("expected Data[%d] = %q, got %q", k, v, dst.Data[k])
		}
	}
}

// TestMap_AllKeysPresent verifies all keys from source are present in destination.
func TestMap_AllKeysPresent(t *testing.T) {
	src := SourceWithMap{
		Labels: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
			"key4": "value4",
			"key5": "value5",
		},
	}
	var dst DestWithMap

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for k := range src.Labels {
		if _, exists := dst.Labels[k]; !exists {
			t.Errorf("key %q from source not found in destination", k)
		}
	}
}

// Both key and value conversion test
type SourceWithInt32KeyInt32Value struct {
	Data map[int32]int32
}

type DestWithInt64KeyInt64Value struct {
	Data map[int64]int64
}

// TestMap_ConvertBothKeyAndValue tests converting both key and value types.
func TestMap_ConvertBothKeyAndValue(t *testing.T) {
	src := SourceWithInt32KeyInt32Value{
		Data: map[int32]int32{1: 100, 2: 200, 3: 300},
	}
	var dst DestWithInt64KeyInt64Value

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Data) != len(src.Data) {
		t.Errorf("expected %d entries, got %d", len(src.Data), len(dst.Data))
	}

	for k, v := range src.Data {
		dstVal, exists := dst.Data[int64(k)]
		if !exists {
			t.Errorf("key %d not found in destination", k)
			continue
		}
		if dstVal != int64(v) {
			t.Errorf("expected Data[%d] = %d, got %d", k, v, dstVal)
		}
	}
}

// TestMap_DifferentType_IntToFloat tests converting map[string]int -> map[string]float64.
func TestMap_DifferentType_IntToFloat(t *testing.T) {
	type Src struct {
		Values map[string]int
	}
	type Dst struct {
		Values map[string]float64
	}

	src := Src{
		Values: map[string]int{"a": 10, "b": 20},
	}
	var dst Dst

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Values["a"] != 10.0 {
		t.Errorf("expected Values[a] = 10.0, got %f", dst.Values["a"])
	}
	if dst.Values["b"] != 20.0 {
		t.Errorf("expected Values[b] = 20.0, got %f", dst.Values["b"])
	}
}

// TestMap_ErrorMessage_IncompatibleKeys verifies error message clearly states key incompatibility.
func TestMap_ErrorMessage_IncompatibleKeys(t *testing.T) {
	type Src struct {
		Data map[string]int
	}
	type Dst struct {
		Data map[int]int // string -> int not convertible
	}

	src := Src{Data: map[string]int{"key": 1}}
	var dst Dst

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if !contains(mappingErr.Reason, "key") {
		t.Errorf("error message should mention 'key', got: %q", mappingErr.Reason)
	}
}

// TestMap_ErrorMessage_IncompatibleValues verifies error message clearly states value incompatibility.
func TestMap_ErrorMessage_IncompatibleValues(t *testing.T) {
	type Src struct {
		Data map[string]string
	}
	type Dst struct {
		Data map[string]int // string -> int not convertible
	}

	src := Src{Data: map[string]string{"key": "not-a-number"}}
	var dst Dst

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if !contains(mappingErr.Reason, "value") {
		t.Errorf("error message should mention 'value', got: %q", mappingErr.Reason)
	}
}

// TestMap_DifferentType_NilMap tests nil handling with different types.
func TestMap_DifferentType_NilMap(t *testing.T) {
	src := SourceWithInt32Map{
		Values: nil,
	}
	var dst DestWithInt64Map

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Values != nil {
		t.Errorf("expected nil map, got %v", dst.Values)
	}
}

// TestMap_DifferentType_EmptyMap tests empty map handling with different types.
func TestMap_DifferentType_EmptyMap(t *testing.T) {
	src := SourceWithInt32Map{
		Values: map[string]int32{},
	}
	var dst DestWithInt64Map

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Values == nil {
		t.Error("expected empty map, got nil")
	}

	if len(dst.Values) != 0 {
		t.Errorf("expected empty map, got %d elements", len(dst.Values))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Error cases
type SourceWithIncompatibleMapKey struct {
	Data map[string]string
}

type DestWithIncompatibleMapKey struct {
	Data map[int]string // string key cannot convert to int
}

// TestMap_IncompatibleKeyType tests error handling for incompatible key types.
func TestMap_IncompatibleKeyType(t *testing.T) {
	src := SourceWithIncompatibleMapKey{
		Data: map[string]string{"key": "value"},
	}
	var dst DestWithIncompatibleMapKey

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error for incompatible key types, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.FieldPath != "Data" {
		t.Errorf("expected FieldPath = 'Data', got %q", mappingErr.FieldPath)
	}
}

type SourceWithIncompatibleMapValue struct {
	Data map[string]string
}

type DestWithIncompatibleMapValue struct {
	Data map[string]struct{ Name string } // string value cannot convert to struct
}

// TestMap_IncompatibleValueType tests error handling for incompatible value types.
func TestMap_IncompatibleValueType(t *testing.T) {
	src := SourceWithIncompatibleMapValue{
		Data: map[string]string{"key": "value"},
	}
	var dst DestWithIncompatibleMapValue

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error for incompatible value types, got nil")
	}

	mappingErr, ok := err.(*MappingError)
	if !ok {
		t.Fatalf("expected *MappingError, got %T", err)
	}

	if mappingErr.FieldPath != "Data" {
		t.Errorf("expected FieldPath = 'Data', got %q", mappingErr.FieldPath)
	}
}

// Multiple maps in same struct
type SourceMultipleMaps struct {
	Labels   map[string]string
	Counts   map[string]int
	Metadata map[string]any
}

type DestMultipleMaps struct {
	Labels   map[string]string
	Counts   map[string]int
	Metadata map[string]any
}

// TestMap_MultipleMaps tests mapping multiple maps in the same struct.
func TestMap_MultipleMaps(t *testing.T) {
	src := SourceMultipleMaps{
		Labels:   map[string]string{"env": "prod"},
		Counts:   map[string]int{"errors": 0, "warnings": 5},
		Metadata: map[string]any{"version": "1.0", "count": 42},
	}
	var dst DestMultipleMaps

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Labels["env"] != "prod" {
		t.Errorf("expected Labels[env] = 'prod', got %q", dst.Labels["env"])
	}

	if dst.Counts["warnings"] != 5 {
		t.Errorf("expected Counts[warnings] = 5, got %d", dst.Counts["warnings"])
	}

	if dst.Metadata["version"] != "1.0" {
		t.Errorf("expected Metadata[version] = '1.0', got %v", dst.Metadata["version"])
	}
}

// =============================================================================
// Map with tag mapping
// =============================================================================

type SourceMapWithTag struct {
	UserLabels map[string]string `map:"Labels"`
}

type DestMapWithTag struct {
	Labels map[string]string
}

// TestMap_WithTagMapping tests map field mapping with tag.
func TestMap_WithTagMapping(t *testing.T) {
	src := SourceMapWithTag{
		UserLabels: map[string]string{"role": "admin"},
	}
	var dst DestMapWithTag

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Labels["role"] != "admin" {
		t.Errorf("expected Labels[role] = 'admin', got %q", dst.Labels["role"])
	}
}
