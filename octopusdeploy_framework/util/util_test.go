package util

import "testing"

func TestStringSlicesEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected bool
	}{
		{"both empty", []string{}, []string{}, true},
		{"both nil", nil, nil, true},
		{"same order", []string{"a", "b"}, []string{"a", "b"}, true},
		{"different order", []string{"a", "b"}, []string{"b", "a"}, true},
		{"different length", []string{"a"}, []string{"a", "b"}, false},
		{"different values", []string{"a", "b"}, []string{"a", "c"}, false},
		{"same duplicate counts", []string{"a", "a", "b"}, []string{"b", "a", "a"}, true},
		{"different duplicate counts", []string{"a", "a", "b"}, []string{"a", "b", "b"}, false},
		{"duplicate against distinct", []string{"a", "a"}, []string{"a", "b"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := StringSlicesEqual(test.a, test.b); actual != test.expected {
				t.Errorf("StringSlicesEqual(%v, %v) = %v, expected %v", test.a, test.b, actual, test.expected)
			}

			if actual := StringSlicesEqual(test.b, test.a); actual != test.expected {
				t.Errorf("StringSlicesEqual(%v, %v) = %v, expected %v", test.b, test.a, actual, test.expected)
			}
		})
	}
}
