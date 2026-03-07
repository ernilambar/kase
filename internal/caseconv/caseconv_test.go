package caseconv

import (
	"testing"
)

func TestWords(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"Hello World", []string{"Hello", "World"}},
		{"hello world", []string{"hello", "world"}},
		{"hello_world", []string{"hello", "world"}},
		{"hello-world", []string{"hello", "world"}},
		{"HelloWorld", []string{"Hello", "World"}},
		{"helloWorld", []string{"hello", "World"}},
		{"  foo   bar  ", []string{"foo", "bar"}},
		{"", nil},
		{"   ", nil},
	}
	for _, tt := range tests {
		got := Words(tt.in, false)
		if !sliceEqual(got, tt.want) {
			t.Errorf("Words(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestToKebab(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Hello World", "hello-world"},
		{"hello world", "hello-world"},
		{"helloWorld", "hello-world"},
		{"HelloWorld", "hello-world"},
	}
	for _, tt := range tests {
		got := ToKebab(tt.in, false)
		if got != tt.want {
			t.Errorf("ToKebab(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestToSnake(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Hello World", "hello_world"},
		{"hello-world", "hello_world"},
	}
	for _, tt := range tests {
		got := ToSnake(tt.in, false)
		if got != tt.want {
			t.Errorf("ToSnake(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestToCamel(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"hello world", "helloWorld"},
		{"Hello World", "helloWorld"},
		{"hello_world", "helloWorld"},
	}
	for _, tt := range tests {
		got := ToCamel(tt.in, false)
		if got != tt.want {
			t.Errorf("ToCamel(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestToPascal(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"hello world", "HelloWorld"},
		{"Hello World", "HelloWorld"},
		{"hello_world", "HelloWorld"},
	}
	for _, tt := range tests {
		got := ToPascal(tt.in, false)
		if got != tt.want {
			t.Errorf("ToPascal(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestToTitle(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"hello world", "Hello World"},
		{"hello_world", "Hello World"},
		{"Hello World", "Hello World"},
	}
	for _, tt := range tests {
		got := ToTitle(tt.in, false)
		if got != tt.want {
			t.Errorf("ToTitle(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEmptyInputConversion(t *testing.T) {
	empty := ""
	if ToKebab(empty, false) != "" || ToSnake(empty, false) != "" || ToCamel(empty, false) != "" || ToPascal(empty, false) != "" || ToTitle(empty, false) != "" {
		t.Error("converters should return empty string for empty input")
	}
}

func TestUnicode(t *testing.T) {
	// Default: accented chars normalized to ASCII
	got := ToKebab("café au_lait", false)
	want := "cafe-au-lait"
	if got != want {
		t.Errorf("ToKebab(Unicode, false) = %q, want %q", got, want)
	}
}

func TestRawPreservesAccents(t *testing.T) {
	// --raw: accented chars preserved
	got := ToKebab("café au_lait", true)
	want := "café-au-lait"
	if got != want {
		t.Errorf("ToKebab(Unicode, true/raw) = %q, want %q", got, want)
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"hello_world", "snake"},
		{"hello-world", "kebab"},
		{"HelloWorld", "pascal"},
		{"helloWorld", "camel"},
		{"Hello World", "title"},
		{"hello world", "title"},
		{"hello", "camel"},
		{"Hello", "pascal"},
		{"", "unknown"},
		{"   ", "unknown"},
		{"foo_bar-baz", "unknown"}, // both _ and -
	}
	for _, tt := range tests {
		got := Detect(tt.in)
		if got != tt.want {
			t.Errorf("Detect(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
