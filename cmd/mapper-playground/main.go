package main

import (
	"fmt"

	"github.com/tariklabs/mapper"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Mapper Playground - Feature Demo                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	// Basic Features
	example1_BasicMapping()
	example2_TagMapping()
	example3_StringConversion()

	// Complex Types
	example4_SliceMapping()
	example5_NestedStructs()
	example6_MapMapping()

	// Advanced Features
	example7_DeepCopyVerification()
	example8_IgnoreZeroSource()
	example9_StrictMode()
	example10_PointerHandling()
	example11_CustomTagName()
	example12_MaxDepth()

	// Real-World Example
	example13_AllFeaturesCombined()

	fmt.Println("\n✅ All examples completed successfully!")
}

// ═══════════════════════════════════════════════════════════════════════════════
// BASIC FEATURES
// ═══════════════════════════════════════════════════════════════════════════════

// Example 1: Basic Field Mapping (by name)
func example1_BasicMapping() {
	fmt.Println("\n━━━ Example 1: Basic Field Mapping ━━━")

	type Source struct {
		Name  string
		Age   int
		Email string
	}

	type Destination struct {
		Name  string
		Age   int
		Email string
	}

	src := Source{Name: "Alice", Age: 30, Email: "alice@example.com"}
	var dst Destination

	if err := mapper.Map(&dst, src); err != nil {
		panic(err)
	}

	fmt.Printf("Source:      %+v\n", src)
	fmt.Printf("Destination: %+v\n", dst)
}

// Example 2: Tag-Based Field Mapping
func example2_TagMapping() {
	fmt.Println("\n━━━ Example 2: Tag-Based Mapping (map tag) ━━━")

	type APIRequest struct {
		FullName    string `map:"Name"`
		YearsOld    int    `map:"Age"`
		ContactInfo string `map:"Email"`
	}

	type User struct {
		Name  string
		Age   int
		Email string
	}

	src := APIRequest{FullName: "Bob", YearsOld: 25, ContactInfo: "bob@example.com"}
	var dst User

	if err := mapper.Map(&dst, src); err != nil {
		panic(err)
	}

	fmt.Printf("Source (API): %+v\n", src)
	fmt.Printf("Destination:  %+v\n", dst)
	fmt.Println("  → Fields mapped via `map` tag aliases")
}

// Example 3: String Conversion (mapconv tag)
func example3_StringConversion() {
	fmt.Println("\n━━━ Example 3: String Conversion (mapconv tag) ━━━")

	type FormInput struct {
		Age      string `mapconv:"int"`
		Price    string `mapconv:"float64"`
		IsActive string `mapconv:"bool"`
		Count    string `mapconv:"int64"`
		Quantity string `mapconv:"uint"`
	}

	type Product struct {
		Age      int
		Price    float64
		IsActive bool
		Count    int64
		Quantity uint
	}

	src := FormInput{
		Age:      "42",
		Price:    "19.99",
		IsActive: "true",
		Count:    "1000000",
		Quantity: "50",
	}
	var dst Product

	if err := mapper.Map(&dst, src); err != nil {
		panic(err)
	}

	fmt.Printf("Source (strings): %+v\n", src)
	fmt.Printf("Destination:      %+v\n", dst)
	fmt.Printf("  → Age:      %T = %v\n", dst.Age, dst.Age)
	fmt.Printf("  → Price:    %T = %v\n", dst.Price, dst.Price)
	fmt.Printf("  → IsActive: %T = %v\n", dst.IsActive, dst.IsActive)
}

// ═══════════════════════════════════════════════════════════════════════════════
// COMPLEX TYPES
// ═══════════════════════════════════════════════════════════════════════════════

// Example 4: Slice Mapping with Type Conversion
func example4_SliceMapping() {
	fmt.Println("\n━━━ Example 4: Slice Mapping ━━━")

	type Source struct {
		Tags    []string
		Scores  []int32
		Empty   []string
		NilList []int
	}

	type Destination struct {
		Tags    []string
		Scores  []int64 // Different type: int32 → int64
		Empty   []string
		NilList []int
	}

	src := Source{
		Tags:    []string{"go", "mapper", "reflection"},
		Scores:  []int32{95, 87, 92},
		Empty:   []string{},
		NilList: nil,
	}
	var dst Destination

	if err := mapper.Map(&dst, src); err != nil {
		panic(err)
	}

	fmt.Printf("Source:      %+v\n", src)
	fmt.Printf("Destination: %+v\n", dst)
	fmt.Printf("  → Empty slice is nil: %v (expected: false)\n", dst.Empty == nil)
	fmt.Printf("  → Nil slice is nil:   %v (expected: true)\n", dst.NilList == nil)
	fmt.Printf("  → Scores converted:   %T → %T\n", src.Scores, dst.Scores)
}

// Example 5: Nested Struct Mapping
func example5_NestedStructs() {
	fmt.Println("\n━━━ Example 5: Nested Struct Mapping ━━━")

	type SrcAddress struct {
		Street string
		City   string
		Zip    string
	}

	type SrcCompany struct {
		Name    string
		Address SrcAddress
	}

	type SrcEmployee struct {
		Name    string
		Age     int
		Company SrcCompany
	}

	type DstAddress struct {
		Street string
		City   string
		Zip    string
	}

	type DstCompany struct {
		Name    string
		Address DstAddress
	}

	type DstEmployee struct {
		Name    string
		Age     int
		Company DstCompany
	}

	src := SrcEmployee{
		Name: "Charlie",
		Age:  35,
		Company: SrcCompany{
			Name: "TechCorp",
			Address: SrcAddress{
				Street: "123 Innovation Way",
				City:   "San Francisco",
				Zip:    "94105",
			},
		},
	}
	var dst DstEmployee

	if err := mapper.Map(&dst, src); err != nil {
		panic(err)
	}

	fmt.Printf("Source:\n")
	fmt.Printf("  Name: %s, Age: %d\n", src.Name, src.Age)
	fmt.Printf("  Company: %s\n", src.Company.Name)
	fmt.Printf("  Address: %s, %s %s\n", src.Company.Address.Street, src.Company.Address.City, src.Company.Address.Zip)

	fmt.Printf("\nDestination (3 levels deep):\n")
	fmt.Printf("  Name: %s, Age: %d\n", dst.Name, dst.Age)
	fmt.Printf("  Company: %s\n", dst.Company.Name)
	fmt.Printf("  Address: %s, %s %s\n", dst.Company.Address.Street, dst.Company.Address.City, dst.Company.Address.Zip)
}

// Example 6: Map Mapping
func example6_MapMapping() {
	fmt.Println("\n━━━ Example 6: Map Mapping ━━━")

	type Source struct {
		Labels   map[string]string
		Counts   map[string]int32
		Metadata map[string]any
		Config   map[string]map[string]string
	}

	type Destination struct {
		Labels   map[string]string
		Counts   map[string]int64 // Different value type: int32 → int64
		Metadata map[string]any
		Config   map[string]map[string]string
	}

	src := Source{
		Labels: map[string]string{"env": "production", "region": "us-east"},
		Counts: map[string]int32{"errors": 0, "warnings": 5, "requests": 1000},
		Metadata: map[string]any{
			"version": "2.0",
			"enabled": true,
		},
		Config: map[string]map[string]string{
			"database": {"host": "localhost", "port": "5432"},
			"cache":    {"host": "redis", "port": "6379"},
		},
	}
	var dst Destination

	if err := mapper.Map(&dst, src); err != nil {
		panic(err)
	}

	fmt.Printf("Source:\n")
	fmt.Printf("  Labels:   %v\n", src.Labels)
	fmt.Printf("  Counts:   %v (%T values)\n", src.Counts, src.Counts["errors"])
	fmt.Printf("  Metadata: %v\n", src.Metadata)
	fmt.Printf("  Config:   %v\n", src.Config)

	fmt.Printf("\nDestination:\n")
	fmt.Printf("  Labels:   %v\n", dst.Labels)
	fmt.Printf("  Counts:   %v (%T values)\n", dst.Counts, dst.Counts["errors"])
	fmt.Printf("  Metadata: %v\n", dst.Metadata)
	fmt.Printf("  Config:   %v\n", dst.Config)
}

// ═══════════════════════════════════════════════════════════════════════════════
// ADVANCED FEATURES
// ═══════════════════════════════════════════════════════════════════════════════

// Example 7: Deep Copy Verification
func example7_DeepCopyVerification() {
	fmt.Println("\n━━━ Example 7: Deep Copy (No Shared References) ━━━")

	type Data struct {
		Items  []string
		Config map[string]string
	}

	src := Data{
		Items:  []string{"original", "values"},
		Config: map[string]string{"key": "original"},
	}
	var dst Data

	if err := mapper.Map(&dst, src); err != nil {
		panic(err)
	}

	fmt.Printf("Before modification:\n")
	fmt.Printf("  Source.Items:  %v\n", src.Items)
	fmt.Printf("  Dest.Items:    %v\n", dst.Items)
	fmt.Printf("  Source.Config: %v\n", src.Config)
	fmt.Printf("  Dest.Config:   %v\n", dst.Config)

	// Modify source
	src.Items[0] = "MODIFIED"
	src.Config["key"] = "MODIFIED"

	fmt.Printf("\nAfter modifying source:\n")
	fmt.Printf("  Source.Items:  %v\n", src.Items)
	fmt.Printf("  Dest.Items:    %v (unaffected!)\n", dst.Items)
	fmt.Printf("  Source.Config: %v\n", src.Config)
	fmt.Printf("  Dest.Config:   %v (unaffected!)\n", dst.Config)
}

// Example 8: Ignore Zero Values (Patch Semantics)
func example8_IgnoreZeroSource() {
	fmt.Println("\n━━━ Example 8: Patch Semantics (WithIgnoreZeroSource) ━━━")

	type PatchRequest struct {
		Name  string
		Email string
		Age   int
	}

	type User struct {
		Name  string
		Email string
		Age   int
	}

	// Existing user in database
	existing := User{Name: "Charlie", Email: "charlie@old.com", Age: 28}

	// Partial update - only Name provided, others are zero values
	patch := PatchRequest{Name: "Charles", Email: "", Age: 0}

	fmt.Printf("Existing user: %+v\n", existing)
	fmt.Printf("Patch request: %+v\n", patch)

	if err := mapper.MapWithOptions(&existing, patch, mapper.WithIgnoreZeroSource()); err != nil {
		panic(err)
	}

	fmt.Printf("After patch:   %+v\n", existing)
	fmt.Println("  → Name updated, Email and Age preserved (zero values ignored)")
}

// Example 9: Strict Mode
func example9_StrictMode() {
	fmt.Println("\n━━━ Example 9: Strict Mode (WithStrictMode) ━━━")

	type Source struct {
		Name string
		// Missing 'Age' and 'Email' fields
	}

	type Destination struct {
		Name  string
		Age   int
		Email string
	}

	src := Source{Name: "Diana"}
	var dst Destination

	// Without strict mode (default) - no error, missing fields ignored
	err := mapper.Map(&dst, src)
	fmt.Printf("Without strict mode: err=%v, dst=%+v\n", err, dst)

	// With strict mode - error because 'Age' and 'Email' have no source
	dst = Destination{} // reset
	err = mapper.MapWithOptions(&dst, src, mapper.WithStrictMode())
	if err != nil {
		fmt.Printf("With strict mode:    err=%q\n", err.Error())
	}
	fmt.Println("  → Strict mode ensures all destination fields have a source")
}

// Example 10: Pointer Handling
func example10_PointerHandling() {
	fmt.Println("\n━━━ Example 10: Pointer Handling ━━━")

	// Value to Pointer
	type ValueSource struct {
		Name  string
		Value int
	}

	type PointerDest struct {
		Name  *string
		Value *int
	}

	src1 := ValueSource{Name: "Eve", Value: 100}
	var dst1 PointerDest

	if err := mapper.Map(&dst1, src1); err != nil {
		panic(err)
	}

	fmt.Printf("Value → Pointer:\n")
	fmt.Printf("  Source:      %+v\n", src1)
	fmt.Printf("  Destination: Name=%q, Value=%d\n", *dst1.Name, *dst1.Value)

	// Pointer to Value
	type PointerSource struct {
		Name  *string
		Value *int
	}

	type ValueDest struct {
		Name  string
		Value int
	}

	name := "Frank"
	value := 200
	src2 := PointerSource{Name: &name, Value: &value}
	var dst2 ValueDest

	if err := mapper.Map(&dst2, src2); err != nil {
		panic(err)
	}

	fmt.Printf("\nPointer → Value:\n")
	fmt.Printf("  Source:      Name=%q, Value=%d\n", *src2.Name, *src2.Value)
	fmt.Printf("  Destination: %+v\n", dst2)
}

// Example 11: Custom Tag Name
func example11_CustomTagName() {
	fmt.Println("\n━━━ Example 11: Custom Tag Name (WithTagName) ━━━")

	// Use 'json' tag instead of default 'map' tag
	type APIResponse struct {
		UserName  string `json:"name"`
		UserEmail string `json:"email"`
		UserAge   int    `json:"age"`
	}

	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	src := APIResponse{UserName: "Grace", UserEmail: "grace@example.com", UserAge: 32}
	var dst User

	// Use "json" tag for field mapping
	if err := mapper.MapWithOptions(&dst, src, mapper.WithTagName("json")); err != nil {
		panic(err)
	}

	fmt.Printf("Source (json tags): %+v\n", src)
	fmt.Printf("Destination:        %+v\n", dst)
	fmt.Println("  → Mapped using `json` tag instead of default `map` tag")
}

// Example 12: Max Depth Protection
func example12_MaxDepth() {
	fmt.Println("\n━━━ Example 12: Max Depth Protection (WithMaxDepth) ━━━")

	type Node struct {
		Value int
		Next  *Node
	}

	// Create a linked list: 1 -> 2 -> 3 -> 4 -> 5
	src := Node{Value: 1}
	current := &src
	for i := 2; i <= 5; i++ {
		current.Next = &Node{Value: i}
		current = current.Next
	}

	var dst Node

	// With default depth (64), this works fine
	if err := mapper.Map(&dst, src); err != nil {
		panic(err)
	}
	fmt.Printf("Default depth (64): Successfully mapped linked list\n")

	// Print the mapped list
	fmt.Printf("  Mapped values: ")
	for n := &dst; n != nil; n = n.Next {
		fmt.Printf("%d ", n.Value)
	}
	fmt.Println()

	// With restricted depth, deep structures fail
	dst = Node{} // reset
	err := mapper.MapWithOptions(&dst, src, mapper.WithMaxDepth(3))
	if err != nil {
		fmt.Printf("WithMaxDepth(3):    err=%q\n", err.Error())
	}
	fmt.Println("  → MaxDepth prevents stack overflow from circular references")
}

// ═══════════════════════════════════════════════════════════════════════════════
// REAL-WORLD EXAMPLE
// ═══════════════════════════════════════════════════════════════════════════════

// Example 13: All Features Combined
func example13_AllFeaturesCombined() {
	fmt.Println("\n━━━ Example 13: Real-World Scenario ━━━")
	fmt.Println("Converting an API request to a domain model with multiple features:")

	type AddressInput struct {
		Street string `map:"StreetAddress"`
		City   string
		Zip    string `map:"PostalCode"`
	}

	type CreateOrderRequest struct {
		CustomerName string            `map:"Customer"`
		TotalPrice   string            `map:"Total" mapconv:"float64"`
		ItemCount    string            `mapconv:"int"`
		IsPriority   string            `mapconv:"bool"`
		Tags         []string          `map:"Labels"`
		Metadata     map[string]string `map:"Info"`
		Shipping     AddressInput      `map:"ShippingAddress"`
	}

	type Address struct {
		StreetAddress string
		City          string
		PostalCode    string
	}

	type Order struct {
		Customer        string
		Total           float64
		ItemCount       int
		IsPriority      bool
		Labels          []string
		Info            map[string]string
		ShippingAddress Address
	}

	src := CreateOrderRequest{
		CustomerName: "ACME Corp",
		TotalPrice:   "1299.99",
		ItemCount:    "5",
		IsPriority:   "true",
		Tags:         []string{"urgent", "wholesale"},
		Metadata:     map[string]string{"region": "US", "priority": "high"},
		Shipping: AddressInput{
			Street: "456 Commerce Blvd",
			City:   "New York",
			Zip:    "10001",
		},
	}
	var dst Order

	if err := mapper.Map(&dst, src); err != nil {
		panic(err)
	}

	fmt.Printf("\nInput (API Request):\n")
	fmt.Printf("  CustomerName: %q (string)\n", src.CustomerName)
	fmt.Printf("  TotalPrice:   %q (string)\n", src.TotalPrice)
	fmt.Printf("  ItemCount:    %q (string)\n", src.ItemCount)
	fmt.Printf("  IsPriority:   %q (string)\n", src.IsPriority)
	fmt.Printf("  Tags:         %v\n", src.Tags)
	fmt.Printf("  Metadata:     %v\n", src.Metadata)
	fmt.Printf("  Shipping:     %+v\n", src.Shipping)

	fmt.Printf("\nOutput (Domain Model):\n")
	fmt.Printf("  Customer:        %q\n", dst.Customer)
	fmt.Printf("  Total:           %.2f (float64)\n", dst.Total)
	fmt.Printf("  ItemCount:       %d (int)\n", dst.ItemCount)
	fmt.Printf("  IsPriority:      %v (bool)\n", dst.IsPriority)
	fmt.Printf("  Labels:          %v\n", dst.Labels)
	fmt.Printf("  Info:            %v\n", dst.Info)
	fmt.Printf("  ShippingAddress: %+v\n", dst.ShippingAddress)

	fmt.Println("\nFeatures demonstrated:")
	fmt.Println("  ✓ Field aliasing via `map` tag")
	fmt.Println("  ✓ String-to-type conversion via `mapconv` tag")
	fmt.Println("  ✓ Nested struct mapping")
	fmt.Println("  ✓ Slice deep copy")
	fmt.Println("  ✓ Map deep copy")
}
