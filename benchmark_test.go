package mapper

import (
	"testing"
)

// =============================================================================
// Benchmark Struct Definitions
// =============================================================================
//
// These structs are designed to represent realistic DTO ↔ domain model mapping
// scenarios commonly found in production Go applications. The field counts and
// types are chosen to approximate real-world payloads while remaining manageable
// for analysis.

// --- Flat Structs (Simple Field Mapping) ---

// BenchSrcFlat represents a typical API response DTO with primitive fields.
// 10 fields is representative of simple DTOs in production systems.
type BenchSrcFlat struct {
	ID          int64
	Name        string
	Email       string
	Age         int
	Score       float64
	Active      bool
	Country     string
	Department  string
	CreatedYear int
	UpdatedYear int
}

// BenchDstFlat is the domain model counterpart to BenchSrcFlat.
type BenchDstFlat struct {
	ID          int64
	Name        string
	Email       string
	Age         int
	Score       float64
	Active      bool
	Country     string
	Department  string
	CreatedYear int
	UpdatedYear int
}

// --- Nested Structs (2-3 levels deep) ---

// BenchSrcAddress represents an embedded address in a larger struct.
type BenchSrcAddress struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// BenchDstAddress is the domain model counterpart.
type BenchDstAddress struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// BenchSrcCompany represents a nested company object.
type BenchSrcCompany struct {
	Name    string
	Address BenchSrcAddress
	TaxID   string
}

// BenchDstCompany is the domain model counterpart.
type BenchDstCompany struct {
	Name    string
	Address BenchDstAddress
	TaxID   string
}

// BenchSrcNested represents a user DTO with nested address and company.
type BenchSrcNested struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Address   BenchSrcAddress
	Company   BenchSrcCompany
	Tags      []string
}

// BenchDstNested is the domain model counterpart.
type BenchDstNested struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Address   BenchDstAddress
	Company   BenchDstCompany
	Tags      []string
}

// --- Slice/Collection Structs ---

// BenchSrcItem represents an item in an order/cart.
type BenchSrcItem struct {
	SKU      string
	Name     string
	Quantity int
	Price    float64
	Weight   float64
}

// BenchDstItem is the domain model counterpart.
type BenchDstItem struct {
	SKU      string
	Name     string
	Quantity int
	Price    float64
	Weight   float64
}

// BenchSrcOrder represents an order with a slice of items.
type BenchSrcOrder struct {
	OrderID     string
	CustomerID  int64
	Items       []BenchSrcItem
	TotalAmount float64
	Status      string
}

// BenchDstOrder is the domain model counterpart.
type BenchDstOrder struct {
	OrderID     string
	CustomerID  int64
	Items       []BenchDstItem
	TotalAmount float64
	Status      string
}

// --- Pointer Fields ---

// BenchSrcOptional represents a struct with optional (pointer) fields.
// This is common in PATCH-style APIs where fields may or may not be present.
type BenchSrcOptional struct {
	ID       int64
	Name     *string
	Email    *string
	Age      *int
	Score    *float64
	Active   *bool
	Address  *BenchSrcAddress
	Metadata map[string]string
}

// BenchDstOptional is the domain model counterpart.
type BenchDstOptional struct {
	ID       int64
	Name     *string
	Email    *string
	Age      *int
	Score    *float64
	Active   *bool
	Address  *BenchDstAddress
	Metadata map[string]string
}

// --- Deep Nesting (for stress testing) ---

// BenchSrcLevel3 is the innermost level.
type BenchSrcLevel3 struct {
	Value   int
	Message string
}

// BenchDstLevel3 is the domain counterpart.
type BenchDstLevel3 struct {
	Value   int
	Message string
}

// BenchSrcLevel2 contains Level3.
type BenchSrcLevel2 struct {
	Name  string
	Inner BenchSrcLevel3
}

// BenchDstLevel2 is the domain counterpart.
type BenchDstLevel2 struct {
	Name  string
	Inner BenchDstLevel3
}

// BenchSrcLevel1 contains Level2.
type BenchSrcLevel1 struct {
	ID    int64
	Inner BenchSrcLevel2
}

// BenchDstLevel1 is the domain counterpart.
type BenchDstLevel1 struct {
	ID    int64
	Inner BenchDstLevel2
}

// BenchSrcDeep is the outermost struct with 4 levels of nesting.
type BenchSrcDeep struct {
	Name  string
	Inner BenchSrcLevel1
}

// BenchDstDeep is the domain counterpart.
type BenchDstDeep struct {
	Name  string
	Inner BenchDstLevel1
}

// =============================================================================
// Test Data Fixtures
// =============================================================================
//
// Pre-allocated fixtures ensure consistent, repeatable benchmark inputs.
// Using package-level variables avoids allocation during benchmark iterations.

var (
	benchSrcFlat = BenchSrcFlat{
		ID:          12345,
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		Age:         35,
		Score:       87.5,
		Active:      true,
		Country:     "United States",
		Department:  "Engineering",
		CreatedYear: 2020,
		UpdatedYear: 2024,
	}

	benchSrcNested = BenchSrcNested{
		ID:        67890,
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@company.com",
		Address: BenchSrcAddress{
			Street:     "123 Main Street",
			City:       "San Francisco",
			State:      "CA",
			PostalCode: "94102",
			Country:    "USA",
		},
		Company: BenchSrcCompany{
			Name: "Tech Corp",
			Address: BenchSrcAddress{
				Street:     "456 Corporate Blvd",
				City:       "Palo Alto",
				State:      "CA",
				PostalCode: "94301",
				Country:    "USA",
			},
			TaxID: "12-3456789",
		},
		Tags: []string{"developer", "senior", "backend"},
	}

	benchSrcDeep = BenchSrcDeep{
		Name: "DeepRoot",
		Inner: BenchSrcLevel1{
			ID: 111,
			Inner: BenchSrcLevel2{
				Name: "Level2",
				Inner: BenchSrcLevel3{
					Value:   42,
					Message: "innermost",
				},
			},
		},
	}
)

// buildBenchSrcOrder creates an order with the specified number of items.
// Item count is parameterized to test scaling behavior.
func buildBenchSrcOrder(itemCount int) BenchSrcOrder {
	items := make([]BenchSrcItem, itemCount)
	for i := 0; i < itemCount; i++ {
		items[i] = BenchSrcItem{
			SKU:      "SKU-" + string(rune('A'+i%26)),
			Name:     "Product Item",
			Quantity: (i % 5) + 1,
			Price:    float64(10 + i%100),
			Weight:   float64(i%10) + 0.5,
		}
	}
	return BenchSrcOrder{
		OrderID:     "ORD-2024-00001",
		CustomerID:  12345,
		Items:       items,
		TotalAmount: 999.99,
		Status:      "pending",
	}
}

// buildBenchSrcOptional creates an optional struct with all fields populated.
func buildBenchSrcOptional() BenchSrcOptional {
	name := "Optional User"
	email := "optional@example.com"
	age := 28
	score := 92.5
	active := true
	return BenchSrcOptional{
		ID:     99999,
		Name:   &name,
		Email:  &email,
		Age:    &age,
		Score:  &score,
		Active: &active,
		Address: &BenchSrcAddress{
			Street:     "789 Optional Ave",
			City:       "Seattle",
			State:      "WA",
			PostalCode: "98101",
			Country:    "USA",
		},
		Metadata: map[string]string{
			"source":   "api",
			"version":  "v2",
			"platform": "web",
		},
	}
}

// =============================================================================
// Benchmarks: Flat Struct Mapping
// =============================================================================

// BenchmarkMap_Flat measures the baseline performance for simple struct mapping.
// This represents the best-case scenario with no nesting or complex types.
func BenchmarkMap_Flat(b *testing.B) {
	src := benchSrcFlat
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstFlat
		if err := Map(&dst, src); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMap_Flat_Pointer tests mapping when source is passed as a pointer.
// This is a common pattern in HTTP handlers.
func BenchmarkMap_Flat_Pointer(b *testing.B) {
	src := &benchSrcFlat
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstFlat
		if err := Map(&dst, src); err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// Benchmarks: Nested Struct Mapping
// =============================================================================

// BenchmarkMap_Nested measures performance with 2-3 levels of struct nesting.
// This is representative of typical domain models with embedded value objects.
func BenchmarkMap_Nested(b *testing.B) {
	src := benchSrcNested
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstNested
		if err := Map(&dst, src); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMap_DeepNested measures performance with 4 levels of nesting.
// This stress-tests the recursive traversal path.
func BenchmarkMap_DeepNested(b *testing.B) {
	src := benchSrcDeep
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstDeep
		if err := Map(&dst, src); err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// Benchmarks: Slice Mapping
// =============================================================================

// BenchmarkMap_Slice_10Items measures slice mapping with a small collection.
func BenchmarkMap_Slice_10Items(b *testing.B) {
	src := buildBenchSrcOrder(10)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstOrder
		if err := Map(&dst, src); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMap_Slice_100Items measures slice mapping with a medium collection.
// 100 items is typical for paginated API responses.
func BenchmarkMap_Slice_100Items(b *testing.B) {
	src := buildBenchSrcOrder(100)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstOrder
		if err := Map(&dst, src); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMap_Slice_1000Items measures slice mapping with a large collection.
// This tests scaling behavior for bulk operations.
func BenchmarkMap_Slice_1000Items(b *testing.B) {
	src := buildBenchSrcOrder(1000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstOrder
		if err := Map(&dst, src); err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// Benchmarks: Pointer/Optional Fields
// =============================================================================

// BenchmarkMap_Optional_AllPopulated measures mapping when all optional fields
// have values. This tests the pointer dereferencing and allocation path.
func BenchmarkMap_Optional_AllPopulated(b *testing.B) {
	src := buildBenchSrcOptional()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstOptional
		if err := Map(&dst, src); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMap_Optional_SomeNil measures mapping when some optional fields
// are nil. This tests the nil-check fast path.
func BenchmarkMap_Optional_SomeNil(b *testing.B) {
	name := "Partial User"
	src := BenchSrcOptional{
		ID:       11111,
		Name:     &name,
		Email:    nil, // nil
		Age:      nil, // nil
		Score:    nil, // nil
		Active:   nil, // nil
		Address:  nil, // nil
		Metadata: nil, // nil
	}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstOptional
		if err := Map(&dst, src); err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// Benchmarks: Mapper Options
// =============================================================================

// BenchmarkMapWithOptions_IgnoreZeroSource measures the overhead of the
// ignoreZeroSource option, which requires zero-value checks on each field.
func BenchmarkMapWithOptions_IgnoreZeroSource(b *testing.B) {
	src := benchSrcFlat
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstFlat
		if err := MapWithOptions(&dst, src, WithIgnoreZeroSource()); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMapWithOptions_CustomTag measures the overhead of using a custom
// tag name instead of the default "map" tag.
func BenchmarkMapWithOptions_CustomTag(b *testing.B) {
	src := benchSrcFlat
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstFlat
		if err := MapWithOptions(&dst, src, WithTagName("json")); err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// Benchmarks: Cache Behavior (Warm vs Cold)
// =============================================================================

// BenchmarkMap_CacheWarm measures performance when struct metadata is already
// cached. This is the common case after the first mapping of a type pair.
func BenchmarkMap_CacheWarm(b *testing.B) {
	// Warm the cache before benchmarking
	var warmup BenchDstFlat
	_ = Map(&warmup, benchSrcFlat)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstFlat
		if err := Map(&dst, benchSrcFlat); err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// Baseline: Manual Mapping
// =============================================================================
//
// These functions provide a performance baseline by implementing manual
// field-by-field assignment. This represents the theoretical minimum overhead
// for struct mapping, against which the reflection-based mapper can be compared.

// manualMapFlat performs direct field assignment without reflection.
func manualMapFlat(dst *BenchDstFlat, src *BenchSrcFlat) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Email = src.Email
	dst.Age = src.Age
	dst.Score = src.Score
	dst.Active = src.Active
	dst.Country = src.Country
	dst.Department = src.Department
	dst.CreatedYear = src.CreatedYear
	dst.UpdatedYear = src.UpdatedYear
}

// manualMapNested performs direct field assignment for nested structs.
func manualMapNested(dst *BenchDstNested, src *BenchSrcNested) {
	dst.ID = src.ID
	dst.FirstName = src.FirstName
	dst.LastName = src.LastName
	dst.Email = src.Email

	// Nested address
	dst.Address.Street = src.Address.Street
	dst.Address.City = src.Address.City
	dst.Address.State = src.Address.State
	dst.Address.PostalCode = src.Address.PostalCode
	dst.Address.Country = src.Address.Country

	// Nested company
	dst.Company.Name = src.Company.Name
	dst.Company.TaxID = src.Company.TaxID
	dst.Company.Address.Street = src.Company.Address.Street
	dst.Company.Address.City = src.Company.Address.City
	dst.Company.Address.State = src.Company.Address.State
	dst.Company.Address.PostalCode = src.Company.Address.PostalCode
	dst.Company.Address.Country = src.Company.Address.Country

	// Slice (deep copy)
	if src.Tags != nil {
		dst.Tags = make([]string, len(src.Tags))
		copy(dst.Tags, src.Tags)
	}
}

// manualMapOrder performs direct field assignment for slices of structs.
func manualMapOrder(dst *BenchDstOrder, src *BenchSrcOrder) {
	dst.OrderID = src.OrderID
	dst.CustomerID = src.CustomerID
	dst.TotalAmount = src.TotalAmount
	dst.Status = src.Status

	if src.Items != nil {
		dst.Items = make([]BenchDstItem, len(src.Items))
		for i := range src.Items {
			dst.Items[i].SKU = src.Items[i].SKU
			dst.Items[i].Name = src.Items[i].Name
			dst.Items[i].Quantity = src.Items[i].Quantity
			dst.Items[i].Price = src.Items[i].Price
			dst.Items[i].Weight = src.Items[i].Weight
		}
	}
}

// BenchmarkBaseline_ManualFlat provides the baseline for flat struct mapping.
func BenchmarkBaseline_ManualFlat(b *testing.B) {
	src := benchSrcFlat
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstFlat
		manualMapFlat(&dst, &src)
	}
}

// BenchmarkBaseline_ManualNested provides the baseline for nested struct mapping.
func BenchmarkBaseline_ManualNested(b *testing.B) {
	src := benchSrcNested
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstNested
		manualMapNested(&dst, &src)
	}
}

// BenchmarkBaseline_ManualSlice_100Items provides the baseline for slice mapping.
func BenchmarkBaseline_ManualSlice_100Items(b *testing.B) {
	src := buildBenchSrcOrder(100)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var dst BenchDstOrder
		manualMapOrder(&dst, &src)
	}
}

// =============================================================================
// Parallel Benchmarks
// =============================================================================
//
// These benchmarks test thread-safety and contention on the metadata cache
// under concurrent load. The mapper uses sync.RWMutex for cache access,
// so these tests reveal any lock contention issues.

// BenchmarkMap_Parallel_Flat tests concurrent flat struct mapping.
func BenchmarkMap_Parallel_Flat(b *testing.B) {
	src := benchSrcFlat
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var dst BenchDstFlat
			if err := Map(&dst, src); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMap_Parallel_Nested tests concurrent nested struct mapping.
func BenchmarkMap_Parallel_Nested(b *testing.B) {
	src := benchSrcNested
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var dst BenchDstNested
			if err := Map(&dst, src); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMap_Parallel_Slice tests concurrent slice mapping.
func BenchmarkMap_Parallel_Slice(b *testing.B) {
	src := buildBenchSrcOrder(100)
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var dst BenchDstOrder
			if err := Map(&dst, src); err != nil {
				b.Fatal(err)
			}
		}
	})
}
