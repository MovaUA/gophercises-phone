package phone_test

import (
	"testing"

	"github.com/movaua/gophercises-phone/pkg/phone"
)

func TestNorm(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{input: "1234567890", want: "1234567890"},
		{input: "123 456 7891", want: "1234567891"},
		{input: "(123) 456 7892", want: "1234567892"},
		{input: "(123) 456-7893", want: "1234567893"},
		{input: "123-456-7894", want: "1234567894"},
		{input: "123-456-7890", want: "1234567890"},
		{input: "(123)456-7892", want: "1234567892"},
		{input: "there is no phone here", want: ""},
		{input: "", want: ""},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			if got := phone.Norm(tc.input); got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

// always store the result of benchmark to a package level variable
// so the compiler cannot eliminate the Benchmark itself.
var benchmarkResult string

func BenchmarkNorn(b *testing.B) {
	benchCases := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456-7893",
		"there is no phone",
		"",
	}
	var r string
	for _, input := range benchCases {
		input := input
		b.Run(input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// records the result to prevent the compiler eliminating the function call
				r = phone.Norm(input)
			}
		})
	}

	// always store the result of benchmark to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	benchmarkResult = r
}
