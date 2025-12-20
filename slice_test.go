package mapper

import (
	"testing"
)

// Test structs for slice mapping.

type SourceWithSlice struct {
	Tags    []string
	Numbers []int
	Scores  []float64
	Flags   []bool
}

type DestWithSlice struct {
	Tags    []string
	Numbers []int
	Scores  []float64
	Flags   []bool
}

type SourceWithInt32Slice struct {
	Values []int32
}

type DestWithInt64Slice struct {
	Values []int64
}

type SourceWithFloat32Slice struct {
	Values []float32
}

type DestWithFloat64Slice struct {
	Values []float64
}

// TestSlice_SameType_String tests mapping []string -> []string.
func TestSlice_SameType_String(t *testing.T) {
	src := SourceWithSlice{
		Tags: []string{"go", "mapper", "reflection"},
	}
	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Tags) != len(src.Tags) {
		t.Errorf("expected %d tags, got %d", len(src.Tags), len(dst.Tags))
	}

	for i, tag := range src.Tags {
		if dst.Tags[i] != tag {
			t.Errorf("expected tag[%d] = %q, got %q", i, tag, dst.Tags[i])
		}
	}
}

// TestSlice_SameType_Int tests mapping []int -> []int.
func TestSlice_SameType_Int(t *testing.T) {
	src := SourceWithSlice{
		Numbers: []int{1, 2, 3, 4, 5},
	}
	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Numbers) != len(src.Numbers) {
		t.Errorf("expected %d numbers, got %d", len(src.Numbers), len(dst.Numbers))
	}

	for i, num := range src.Numbers {
		if dst.Numbers[i] != num {
			t.Errorf("expected number[%d] = %d, got %d", i, num, dst.Numbers[i])
		}
	}
}

// TestSlice_SameType_Float64 tests mapping []float64 -> []float64.
func TestSlice_SameType_Float64(t *testing.T) {
	src := SourceWithSlice{
		Scores: []float64{1.1, 2.2, 3.3},
	}
	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Scores) != len(src.Scores) {
		t.Errorf("expected %d scores, got %d", len(src.Scores), len(dst.Scores))
	}

	for i, score := range src.Scores {
		if dst.Scores[i] != score {
			t.Errorf("expected score[%d] = %f, got %f", i, score, dst.Scores[i])
		}
	}
}

// TestSlice_SameType_Bool tests mapping []bool -> []bool.
func TestSlice_SameType_Bool(t *testing.T) {
	src := SourceWithSlice{
		Flags: []bool{true, false, true},
	}
	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Flags) != len(src.Flags) {
		t.Errorf("expected %d flags, got %d", len(src.Flags), len(dst.Flags))
	}

	for i, flag := range src.Flags {
		if dst.Flags[i] != flag {
			t.Errorf("expected flag[%d] = %v, got %v", i, flag, dst.Flags[i])
		}
	}
}

// TestSlice_NilSlice tests that nil slice in source results in nil slice in destination.
func TestSlice_NilSlice(t *testing.T) {
	src := SourceWithSlice{
		Tags: nil, // explicitly nil
	}
	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Tags != nil {
		t.Errorf("expected nil slice, got %v", dst.Tags)
	}
}

// TestSlice_EmptySlice tests that empty slice in source results in empty slice (not nil) in destination.
func TestSlice_EmptySlice(t *testing.T) {
	src := SourceWithSlice{
		Tags: []string{}, // empty, not nil
	}
	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Tags == nil {
		t.Errorf("expected empty slice, got nil")
	}

	if len(dst.Tags) != 0 {
		t.Errorf("expected 0 elements, got %d", len(dst.Tags))
	}
}

// TestSlice_DeepCopy tests that modifying source slice after mapping doesn't affect destination.
func TestSlice_DeepCopy(t *testing.T) {
	src := SourceWithSlice{
		Tags: []string{"original"},
	}
	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Modify source after mapping.
	src.Tags[0] = "modified"

	// Destination should still have original value.
	if dst.Tags[0] != "original" {
		t.Errorf("expected 'original', got %q - destination was affected by source modification", dst.Tags[0])
	}
}

// TestSlice_DeepCopy_Append tests that appending to source slice doesn't affect destination.
func TestSlice_DeepCopy_Append(t *testing.T) {
	src := SourceWithSlice{
		Numbers: []int{1, 2, 3},
	}
	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	originalLen := len(dst.Numbers)

	// Append to source after mapping.
	src.Numbers = append(src.Numbers, 4, 5, 6)

	// Destination length should be unchanged.
	if len(dst.Numbers) != originalLen {
		t.Errorf("expected length %d, got %d - destination was affected by source append", originalLen, len(dst.Numbers))
	}
}

// TestSlice_DifferentType_Int32ToInt64 tests converting []int32 -> []int64.
func TestSlice_DifferentType_Int32ToInt64(t *testing.T) {
	src := SourceWithInt32Slice{
		Values: []int32{1, 2, 3, 4, 5},
	}
	var dst DestWithInt64Slice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Values) != len(src.Values) {
		t.Errorf("expected %d values, got %d", len(src.Values), len(dst.Values))
	}

	for i, val := range src.Values {
		if dst.Values[i] != int64(val) {
			t.Errorf("expected value[%d] = %d, got %d", i, val, dst.Values[i])
		}
	}
}

// TestSlice_DifferentType_Float32ToFloat64 tests converting []float32 -> []float64.
func TestSlice_DifferentType_Float32ToFloat64(t *testing.T) {
	src := SourceWithFloat32Slice{
		Values: []float32{1.1, 2.2, 3.3},
	}
	var dst DestWithFloat64Slice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Values) != len(src.Values) {
		t.Errorf("expected %d values, got %d", len(src.Values), len(dst.Values))
	}

	for i, val := range src.Values {
		// Use approximate comparison for float conversion.
		expected := float64(val)
		if dst.Values[i] != expected {
			t.Errorf("expected value[%d] = %f, got %f", i, expected, dst.Values[i])
		}
	}
}

// TestSlice_DifferentType_NilSlice tests nil handling with different types.
func TestSlice_DifferentType_NilSlice(t *testing.T) {
	src := SourceWithInt32Slice{
		Values: nil,
	}
	var dst DestWithInt64Slice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Values != nil {
		t.Errorf("expected nil slice, got %v", dst.Values)
	}
}

// TestSlice_DifferentType_EmptySlice tests empty slice handling with different types.
func TestSlice_DifferentType_EmptySlice(t *testing.T) {
	src := SourceWithInt32Slice{
		Values: []int32{},
	}
	var dst DestWithInt64Slice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Values == nil {
		t.Errorf("expected empty slice, got nil")
	}

	if len(dst.Values) != 0 {
		t.Errorf("expected 0 elements, got %d", len(dst.Values))
	}
}

// TestSlice_LargeSlice tests mapping a large slice for performance sanity.
func TestSlice_LargeSlice(t *testing.T) {
	const size = 10000
	src := SourceWithSlice{
		Numbers: make([]int, size),
	}
	for i := 0; i < size; i++ {
		src.Numbers[i] = i
	}

	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Numbers) != size {
		t.Errorf("expected %d numbers, got %d", size, len(dst.Numbers))
	}

	// Check first and last elements.
	if dst.Numbers[0] != 0 {
		t.Errorf("expected first element = 0, got %d", dst.Numbers[0])
	}
	if dst.Numbers[size-1] != size-1 {
		t.Errorf("expected last element = %d, got %d", size-1, dst.Numbers[size-1])
	}
}

// TestSlice_MultipleSliceFields tests mapping multiple slice fields at once.
func TestSlice_MultipleSliceFields(t *testing.T) {
	src := SourceWithSlice{
		Tags:    []string{"a", "b", "c"},
		Numbers: []int{1, 2, 3},
		Scores:  []float64{1.1, 2.2, 3.3},
		Flags:   []bool{true, false, true},
	}
	var dst DestWithSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(dst.Tags))
	}
	if len(dst.Numbers) != 3 {
		t.Errorf("expected 3 numbers, got %d", len(dst.Numbers))
	}
	if len(dst.Scores) != 3 {
		t.Errorf("expected 3 scores, got %d", len(dst.Scores))
	}
	if len(dst.Flags) != 3 {
		t.Errorf("expected 3 flags, got %d", len(dst.Flags))
	}
}

// Test structs for pointer to slice.

type SourceWithSlicePointer struct {
	Tags *[]string
}

type DestWithSlicePointer struct {
	Tags *[]string
}

type DestWithSliceValue struct {
	Tags []string
}

// TestSlice_PointerToSlice tests mapping *[]string -> *[]string.
func TestSlice_PointerToSlice(t *testing.T) {
	tags := []string{"go", "mapper"}
	src := SourceWithSlicePointer{
		Tags: &tags,
	}
	var dst DestWithSlicePointer

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Tags == nil {
		t.Fatal("expected non-nil pointer")
	}

	if len(*dst.Tags) != len(*src.Tags) {
		t.Errorf("expected %d tags, got %d", len(*src.Tags), len(*dst.Tags))
	}
}

// TestSlice_PointerToSlice_Nil tests mapping nil *[]string -> *[]string.
func TestSlice_PointerToSlice_Nil(t *testing.T) {
	src := SourceWithSlicePointer{
		Tags: nil,
	}
	var dst DestWithSlicePointer

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Tags != nil {
		t.Errorf("expected nil pointer, got %v", dst.Tags)
	}
}

// Test structs for incompatible slices.

type SourceWithStringSlice struct {
	Values []string
}

type DestWithIntSlice struct {
	Values []int
}

// TestSlice_IncompatibleTypes tests that incompatible slice types return an error.
func TestSlice_IncompatibleTypes(t *testing.T) {
	src := SourceWithStringSlice{
		Values: []string{"a", "b", "c"},
	}
	var dst DestWithIntSlice

	err := Map(&dst, src)
	if err == nil {
		t.Fatal("expected error for incompatible slice types, got nil")
	}

	// Verify it's a MappingError.
	if _, ok := err.(*MappingError); !ok {
		t.Errorf("expected *MappingError, got %T", err)
	}
}

// Test nested slices.

type SourceWithNestedSlice struct {
	Matrix [][]int
}

type DestWithNestedSlice struct {
	Matrix [][]int
}

// TestSlice_NestedSlice tests mapping [][]int -> [][]int.
func TestSlice_NestedSlice(t *testing.T) {
	src := SourceWithNestedSlice{
		Matrix: [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}
	var dst DestWithNestedSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Matrix) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(dst.Matrix))
	}

	for i, row := range src.Matrix {
		if len(dst.Matrix[i]) != len(row) {
			t.Errorf("row %d: expected %d elements, got %d", i, len(row), len(dst.Matrix[i]))
		}
		for j, val := range row {
			if dst.Matrix[i][j] != val {
				t.Errorf("matrix[%d][%d]: expected %d, got %d", i, j, val, dst.Matrix[i][j])
			}
		}
	}
}

// TestSlice_NestedSlice_DeepCopy tests that nested slices are also deep copied.
func TestSlice_NestedSlice_DeepCopy(t *testing.T) {
	src := SourceWithNestedSlice{
		Matrix: [][]int{
			{1, 2, 3},
		},
	}
	var dst DestWithNestedSlice

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Modify source after mapping.
	src.Matrix[0][0] = 999

	// Destination should still have original value.
	if dst.Matrix[0][0] != 1 {
		t.Errorf("expected 1, got %d - nested slice was not deep copied", dst.Matrix[0][0])
	}
}
