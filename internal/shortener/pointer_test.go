package shortener

import "testing"

type counter struct {
	value int
}

func (c counter) incrementValue() {
	c.value++
}

func (c *counter) incrementPointer() {
	c.value++
}

func TestValueReceiverDoesNotChangeOriginal(t *testing.T) {
	example := counter{
		value: 10,
	}
	example.incrementValue()
	if example.value != 10 {
		t.Errorf("expected 10, got %d", example.value)
	}
}

func TestPointerReceiverChangesOriginal(t *testing.T) {
	example := counter{
		value: 10,
	}
	example.incrementPointer()
	if example.value != 11 {
		t.Errorf("expected 11, got %d", example.value)
	}
}
