package xim

import (
	"fmt"
	"reflect"
	"testing"
)

func assertBit(t *testing.T, title string, actual, expected Bit) {
	t.Helper()
	if actual != expected {
		t.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", title, actual, expected)
	}
}

func TestInBuilderBit(t *testing.T) {
	inBuilder := NewInBuilder()

	expected := []Bit{
		1,
		2,
		4,
		8,
		16,
		32,
		64,
		128,
		256,
	}

	for i := range expected {
		i := i // escape: Using the variable on range scope `i` in loop literal
		t.Run(fmt.Sprintf("%d", expected[i]), func(tr *testing.T) {
			assertBit(t, tr.Name(), inBuilder.NewBit(), expected[i])
		})
	}

	uintSize := 16

	for i := 9; i < uintSize-1; i++ {
		assertBit(t, t.Name(), inBuilder.NewBit(), Bit(1<<uint(i)))
	}

	assertBit(t, t.Name(), inBuilder.NewBit(), 1<<(uint(uintSize)-1))

	// overflow
	func() {
		defer func() {
			if rec := recover(); rec == nil {
				t.Errorf("panic expected")
			}
		}()

		inBuilder.NewBit()
	}()
}

func TestInBuilderIndexes(t *testing.T) {
	inBuilder := NewInBuilder()

	a := inBuilder.NewBit()
	b := inBuilder.NewBit()
	c := inBuilder.NewBit()
	d := inBuilder.NewBit()

	idxs := inBuilder.Indexes()

	if len(idxs) != 0 {
		t.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", t.Name(), len(idxs), 0)
	}

	idxs = inBuilder.Indexes(a, c)

	expected := []string{
		"1",
		"3",
		"4",
		"5",
		"6",
		"7",
		"9",
		"b",
		"c",
		"d",
		"e",
		"f",
	}
	if !reflect.DeepEqual(idxs, expected) {
		t.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", t.Name(), idxs, expected)
	}

	idxs = inBuilder.Indexes(b, d)

	expected = []string{
		"2",
		"3",
		"6",
		"7",
		"8",
		"9",
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
	}
	if !reflect.DeepEqual(idxs, expected) {
		t.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", t.Name(), idxs, expected)
	}
}

func TestInBuilderFilters(t *testing.T) {
	inBuilder := NewInBuilder()

	a := inBuilder.NewBit()
	b := inBuilder.NewBit()
	c := inBuilder.NewBit()
	d := inBuilder.NewBit()

	filter := inBuilder.Filter(a, c)
	expected := "5"

	if filter != expected {
		t.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", t.Name(), filter, expected)
	}

	filter = inBuilder.Filter(b, d)
	expected = "a"

	if filter != expected {
		t.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", t.Name(), filter, expected)
	}
}
