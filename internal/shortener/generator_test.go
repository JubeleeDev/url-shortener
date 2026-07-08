package shortener

import "testing"

func TestGenerateCodeReturnsRequestedLength(t *testing.T) {
	code, err := GenerateCode(5)

	if err != nil {
		t.Fatalf("expected err == nil, got %v", err)
	}

	if len(code) != 5 {
		t.Errorf("expected code length 5, got %d", len(code))
	}

}

func TestGenerateCodeReturnsErrorForInvalidLength(t *testing.T) {
	cases := []struct {
		name   string
		length int
	}{
		{name: "zero", length: 0},
		{name: "negative", length: -3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GenerateCode(tc.length)

			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
