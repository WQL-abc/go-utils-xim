package xim

import (
	"fmt"
	"sort"
	"testing"
)

func TestString(t *testing.T) {
	b := bigram{'a', 'あ'}
	if b.String() != "aあ" {
		t.Errorf("exected:%s, but was: %s\n", "aあ", b.String())
	}
}

func TestToUnigrams(t *testing.T) {
	result := toUnigrams("abc dあいbCh")
	if len(result) != 8 {
		t.Errorf("len(result) exected:%d, but was: %d\n", 8, len(result))
	}

	expected := []rune{
		'a',
		'b',
		'c',
		'C',
		'd',
		'あ',
		'い',
		'h',
	}

	for _, r := range expected {
		r := r // escape: Using the variable on range scope `r` in loop literal
		t.Run(fmt.Sprintf("%s", string(r)), func(tr *testing.T) {
			if _, ok := result[r]; !ok {
				tr.Errorf("%s: Unigram notfound. '%s'", t.Name(), string(r))
			}
		})
	}
}

func TestToBigrams(t *testing.T) {
	result := toBigrams("abc debch iJあdeN")
	if len(result) != 9 {
		t.Errorf("len(result) exected:%d, but was: %d\n", 9, len(result))
	}

	expected := []bigram{
		{'a', 'b'},
		{'b', 'c'},
		{'d', 'e'},
		{'e', 'b'},
		{'c', 'h'},
		{'i', 'J'},
		{'J', 'あ'},
		{'あ', 'd'},
		{'e', 'N'},
	}

	for _, b := range expected {
		b := b // escape: Using the variable on range scope `b` in loop literal
		t.Run(fmt.Sprintf("%s,%s", string(b.a), string(b.b)), func(tr *testing.T) {
			if _, ok := result[b]; !ok {
				tr.Errorf("%s: Bigram notfound. %v", t.Name(), b)
			}
		})
	}
}

func TestBigrams(t *testing.T) {
	result := Bigrams("abc dあいbCh")
	if len(result) != 7 {
		t.Errorf("len(result) exected:%d, but was: %d\n", 6, len(result))
	}

	sort.Strings(result)

	expected := []string{
		"Ch",
		"ab",
		"bC",
		"bc",
		"dあ",
		"あい",
		"いb",
	}

	for i := range result {
		i := i // escape: Using the variable on range scope `i` in loop literal
		t.Run(fmt.Sprintf("result[%d]", i), func(tr *testing.T) {
			if result[i] != expected[i] {
				tr.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", t.Name(), result[i], expected[i])
			}
		})
	}
}

func TestBiunigrams(t *testing.T) {
	result := Biunigrams("abc dあいbCh")
	if len(result) != 15 {
		t.Errorf("len(result) expected: %d, but was: %d\n", 13, len(result))
	}

	sort.Strings(result)

	expected := []string{
		"C",
		"Ch",
		"a",
		"ab",
		"b",
		"bC",
		"bc",
		"c",
		"d",
		"dあ",
		"h",
		"あ",
		"あい",
		"い",
		"いb",
	}

	for i := range result {
		i := i // escape: Using the variable on range scope `i` in loop literal
		t.Run(fmt.Sprintf("result[%d]", i), func(tr *testing.T) {
			if result[i] != expected[i] {
				tr.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", t.Name(), result[i], expected[i])
			}
		})
	}
}

func TestPrefixes(t *testing.T) {
	result := Prefixes("abc dあいbCh")
	if len(result) != 9 {
		t.Errorf("len(result) exected:%d, but was: %d\n", 9, len(result))
	}

	sort.Strings(result)

	expected := []string{
		"a",
		"ab",
		"abc",
		"d",
		"dあ",
		"dあい",
		"dあいb",
		"dあいbC",
		"dあいbCh",
	}

	for i := range result {
		i := i // escape: Using the variable on range scope `i` in loop literal
		t.Run(fmt.Sprintf("result[%d]", i), func(tr *testing.T) {
			if result[i] != expected[i] {
				tr.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", t.Name(), result[i], expected[i])
			}
		})
	}
}

func TestSuffixes(t *testing.T) {
	result := Suffixes("abc dあいbCh")
	if len(result) != 9 {
		t.Errorf("len(result) exected:%d, but was: %d\n", 9, len(result))
	}

	sort.Strings(result)

	expected := []string{
		"Ch",
		"abc",
		"bCh",
		"bc",
		"c",
		"dあいbCh",
		"h",
		"あいbCh",
		"いbCh",
	}

	for i := range result {
		i := i // escape: Using the variable on range scope `i` in loop literal
		t.Run(fmt.Sprintf("result[%d]", i), func(tr *testing.T) {
			if result[i] != expected[i] {
				tr.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", t.Name(), result[i], expected[i])
			}
		})
	}
}
