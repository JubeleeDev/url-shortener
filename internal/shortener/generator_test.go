package shortener

import "testing"

func TestGenerateCodeReturnsRequestedLength(t *testing.T) {
	code, err := GenerateCode(5)

	if err != nil {
		t.Fatalf("expected err == nil, got %v", err)
	}

	if len(code) != 5 {
		t.Errorf("expected code length = 5, got %d", len(code))
	}

}

func TestGenerateCodeReturnsErrorForZeroLength(t *testing.T) {
	_, err := GenerateCode(0)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGenerateCodeReturnsErrorForNegativeLength(t *testing.T) {
	_, err := GenerateCode(-3)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
