package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBoolPointerFromType(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Bool
		expected *bool
		isNil    bool
	}{
		{
			name:     "null value returns nil",
			input:    types.BoolNull(),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "unknown value returns nil",
			input:    types.BoolUnknown(),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "true value returns pointer to true",
			input:    types.BoolValue(true),
			expected: boolPtr(true),
			isNil:    false,
		},
		{
			name:     "false value returns pointer to false",
			input:    types.BoolValue(false),
			expected: boolPtr(false),
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := boolPointerFromType(tt.input)
			if tt.isNil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("expected non-nil pointer, got nil")
				} else if *result != *tt.expected {
					t.Errorf("expected %v, got %v", *tt.expected, *result)
				}
			}
		})
	}
}

func TestStringPointerFromType(t *testing.T) {
	tests := []struct {
		name     string
		input    types.String
		expected *string
		isNil    bool
	}{
		{
			name:     "null value returns nil",
			input:    types.StringNull(),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "unknown value returns nil",
			input:    types.StringUnknown(),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "string value returns pointer",
			input:    types.StringValue("test"),
			expected: stringPtr("test"),
			isNil:    false,
		},
		{
			name:     "empty string value returns pointer to empty string",
			input:    types.StringValue(""),
			expected: stringPtr(""),
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringPointerFromType(tt.input)
			if tt.isNil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("expected non-nil pointer, got nil")
				} else if *result != *tt.expected {
					t.Errorf("expected %v, got %v", *tt.expected, *result)
				}
			}
		})
	}
}

func TestInt64PointerFromType(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Int64
		expected *int64
		isNil    bool
	}{
		{
			name:     "null value returns nil",
			input:    types.Int64Null(),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "unknown value returns nil",
			input:    types.Int64Unknown(),
			expected: nil,
			isNil:    true,
		},
		{
			name:     "positive value returns pointer",
			input:    types.Int64Value(42),
			expected: int64Ptr(42),
			isNil:    false,
		},
		{
			name:     "zero value returns pointer to zero",
			input:    types.Int64Value(0),
			expected: int64Ptr(0),
			isNil:    false,
		},
		{
			name:     "negative value returns pointer",
			input:    types.Int64Value(-1),
			expected: int64Ptr(-1),
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int64PointerFromType(tt.input)
			if tt.isNil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("expected non-nil pointer, got nil")
				} else if *result != *tt.expected {
					t.Errorf("expected %v, got %v", *tt.expected, *result)
				}
			}
		})
	}
}

func TestBoolValueOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Bool
		fallback bool
		expected bool
	}{
		{
			name:     "null value returns fallback",
			input:    types.BoolNull(),
			fallback: true,
			expected: true,
		},
		{
			name:     "unknown value returns fallback",
			input:    types.BoolUnknown(),
			fallback: false,
			expected: false,
		},
		{
			name:     "true value returns true",
			input:    types.BoolValue(true),
			fallback: false,
			expected: true,
		},
		{
			name:     "false value returns false",
			input:    types.BoolValue(false),
			fallback: true,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := boolValueOrDefault(tt.input, tt.fallback)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDiffFunction(t *testing.T) {
	tests := []struct {
		name        string
		state       []string
		plan        []string
		expectedAdd []string
		expectedRem []string
	}{
		{
			name:        "empty state and plan",
			state:       []string{},
			plan:        []string{},
			expectedAdd: []string{},
			expectedRem: []string{},
		},
		{
			name:        "add items to empty state",
			state:       []string{},
			plan:        []string{"a", "b", "c"},
			expectedAdd: []string{"a", "b", "c"},
			expectedRem: []string{},
		},
		{
			name:        "remove all items",
			state:       []string{"a", "b", "c"},
			plan:        []string{},
			expectedAdd: []string{},
			expectedRem: []string{"a", "b", "c"},
		},
		{
			name:        "no changes",
			state:       []string{"a", "b", "c"},
			plan:        []string{"a", "b", "c"},
			expectedAdd: []string{},
			expectedRem: []string{},
		},
		{
			name:        "add and remove items",
			state:       []string{"a", "b"},
			plan:        []string{"b", "c"},
			expectedAdd: []string{"c"},
			expectedRem: []string{"a"},
		},
		{
			name:        "complex scenario",
			state:       []string{"1", "2", "3", "4"},
			plan:        []string{"2", "4", "5", "6"},
			expectedAdd: []string{"5", "6"},
			expectedRem: []string{"1", "3"},
		},
		{
			name:        "large list performance test",
			state:       generateStringList(100),
			plan:        generateStringList(100),
			expectedAdd: []string{},
			expectedRem: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			add, rem := diffFunction(tt.state, tt.plan)

			if !stringSlicesEqual(add, tt.expectedAdd) {
				t.Errorf("add: expected %v, got %v", tt.expectedAdd, add)
			}

			if !stringSlicesEqual(rem, tt.expectedRem) {
				t.Errorf("remove: expected %v, got %v", tt.expectedRem, rem)
			}
		})
	}
}

// Helper functions for tests
func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[string]bool)
	for _, item := range a {
		aMap[item] = true
	}

	for _, item := range b {
		if !aMap[item] {
			return false
		}
	}

	return true
}

func generateStringList(size int) []string {
	result := make([]string, size)
	for i := 0; i < size; i++ {
		result[i] = string(rune('a' + (i % 26)))
	}
	return result
}
