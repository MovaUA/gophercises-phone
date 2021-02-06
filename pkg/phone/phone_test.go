package phone_test

import (
	"testing"

	"github.com/movaua/gophercises-phone/pkg/phone"
)

func TestNorm(t *testing.T) {
	tests := []struct {
		actual string
		want   string
	}{
		{actual: "1234567890", want: "1234567890"},
		{actual: "123 456 7891", want: "1234567891"},
		{actual: "(123) 456 7892", want: "1234567892"},
		{actual: "(123) 456-7893", want: "1234567893"},
		{actual: "123-456-7894", want: "1234567894"},
		{actual: "123-456-7890", want: "1234567890"},
		{actual: "1234567892", want: "1234567892"},
		{actual: "(123)456-7892", want: "1234567892"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.actual, func(t *testing.T) {
			if got := phone.Norm(tt.actual); got != tt.want {
				t.Fatalf("want %q, got %q", tt.want, got)
			}
		})
	}
}
