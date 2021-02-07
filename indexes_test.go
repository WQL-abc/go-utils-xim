package xim

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestAddIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.Add("label1", "abc dあいbCh", "sample")
	idx.Add("label2", "abc debch iJあdeN", "sample")

	built := idx.MustBuild()
	assertBuiltIndex(t, built, map[string]bool{
		"label1 abc dあいbCh":       true,
		"label1 sample":           true,
		"label2 abc debch iJあdeN": true,
		"label2 sample":           true,
	})
}

func TestAddBigrams(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddBigrams("label1", "abc dあいbCh")
	idx.AddBigrams("label2", "abc debch iJあdeN")

	expected := make(map[string]bool)
	for _, s := range Bigrams("abc dあいbCh") {
		expected["label1 "+s] = true
	}
	for _, s := range Bigrams("abc debch iJあdeN") {
		expected["label2 "+s] = true
	}

	built := idx.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestAddBiunigramsIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddBiunigrams("label1", "abc dあいbCh")
	idx.AddBiunigrams("label2", "abc debch iJあdeN")

	expected := make(map[string]bool)
	for _, s := range Biunigrams("abc dあいbCh") {
		expected["label1 "+s] = true
	}
	for _, s := range Biunigrams("abc debch iJあdeN") {
		expected["label2 "+s] = true
	}

	built := idx.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestAddPrefixesIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddPrefixes("label1", "abc dあいbCh")
	idx.AddPrefixes("label2", "abc debch iJあdeN")

	expected := make(map[string]bool)
	for _, s := range Prefixes("abc dあいbCh") {
		expected["label1 "+s] = true
	}
	for _, s := range Prefixes("abc debch iJあdeN") {
		expected["label2 "+s] = true
	}

	built := idx.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestAddSuffixesIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddSuffixes("label1", "abc dあいbCh")
	idx.AddSuffixes("label2", "abc debch iJあdeN")

	expected := make(map[string]bool)
	for _, s := range Suffixes("abc dあいbCh") {
		expected["label1 "+s] = true
	}
	for _, s := range Suffixes("abc debch iJあdeN") {
		expected["label2 "+s] = true
	}

	built := idx.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestAddSomethingIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddSomething("label1", []string{"abc dあいbCh", "abc debch iJあdeN"})
	idx.AddSomething("label2", 123)
	now := time.Now()
	idx.AddSomething("label3", now)

	built := idx.MustBuild()
	assertBuiltIndex(t, built, map[string]bool{
		"label1 abc dあいbCh":                      true,
		"label1 abc debch iJあdeN":                true,
		"label2 123":                             true,
		fmt.Sprintf("label3 %d", now.UnixNano()): true,
	})
}

func TestAddAllIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.Add("label1", "abc dあいbCh", "sample")
	idx.AddBigrams("label2", "abc dあいbCh")
	idx.AddBiunigrams("label3", "abc dあいbCh")
	idx.AddPrefixes("label4", "abc dあいbCh")
	idx.AddSomething("label5", []string{"abc dあいbCh", "AbcdeF"})

	expected := make(map[string]bool)

	// Add
	expected["label1 abc dあいbCh"] = true
	expected["label1 sample"] = true

	// AddBigrams
	for _, s := range Bigrams("abc dあいbCh") {
		expected["label2 "+s] = true
	}

	// AddBiunigrams
	for _, s := range Biunigrams("abc dあいbCh") {
		expected["label3 "+s] = true
	}

	// AddPrefixes
	for _, s := range Prefixes("abc dあいbCh") {
		expected["label4 "+s] = true
	}

	// AddSomething
	expected["label5 abc dあいbCh"] = true
	expected["label5 AbcdeF"] = true

	built := idx.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestBuildIndex(t *testing.T) {
	t.Run("Success", func(tr *testing.T) {
		idx := NewIndexes(nil)
		for i := 0; i < MaxIndexesSize; i++ {
			idx.Add(fmt.Sprintf("label%d", i), "abc")
		}

		built, err := idx.Build()
		if err != nil {
			tr.Errorf("error = %s, wants = nil", err)
		}
		if built == nil {
			tr.Error("built = nil, wants != nil")
		}
	})

	t.Run("Config.CompositeIdxLabels > MaxCompositeIndexLabels", func(tr *testing.T) {
		labels := make([]string, MaxCompositeIndexLabels+1)
		for i := 0; i < len(labels); i++ {
			labels[i] = string(rune('a' + i))
		}

		idx := NewIndexes(&Config{CompositeIdxLabels: labels})
		if _, err := idx.Build(); err == nil {
			tr.Error("error = nil, wants != nil")
		}
	})

	t.Run("Number of Build() results > MaxIndexesSize", func(tr *testing.T) {
		idx := NewIndexes(nil)
		for i := 0; i < MaxIndexesSize+1; i++ {
			idx.Add(fmt.Sprintf("label%d", i), "abc")
		}

		if _, err := idx.Build(); err == nil {
			tr.Error("error = nil, wants != nil")
		}
	})
}

func TestMustBuildIndex(t *testing.T) {
	t.Run("Success", func(tr *testing.T) {
		defer func() {
			if rec := recover(); rec != nil {
				tr.Errorf("expected: not panic, was: panic = [%v]\n", rec)
			}
		}()

		idx := NewIndexes(nil)
		for i := 0; i < MaxIndexesSize; i++ {
			idx.Add(fmt.Sprintf("label%d", i), "abc")
		}

		built := idx.MustBuild()
		if built == nil {
			tr.Error("built = nil, wants != nil")
		}
	})

	t.Run("Config.CompositeIdxLabels > MaxCompositeIndexLabels", func(tr *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				tr.Error("expected: panic, was: not panic\n")
			}
		}()

		labels := make([]string, MaxCompositeIndexLabels+1)
		for i := 0; i < len(labels); i++ {
			labels[i] = string(rune('a' + i))
		}

		idx := NewIndexes(&Config{CompositeIdxLabels: labels})
		idx.MustBuild()
	})

	t.Run("Number of MustBuild() results > MaxIndexesSize", func(tr *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				tr.Error("expected: panic, was: not panic\n")
			}
		}()

		idx := NewIndexes(nil)
		for i := 0; i < MaxIndexesSize+1; i++ {
			idx.Add(fmt.Sprintf("label%d", i), "abc")
		}

		idx.MustBuild()
	})
}

func TestIndexConfigCompositeIdxLabels(t *testing.T) {
	idx := NewIndexes(&Config{CompositeIdxLabels: []string{"label1", "label2", "label3"}})
	idx.Add("label1", "a")
	idx.Add("label2", "b")
	idx.Add("label3", "c")
	idx.Add("label4", "d")

	built := idx.MustBuild()

	//   c b a
	//   -----
	//3  0 1 1
	//5  1 0 1
	//6  1 1 0
	//7  1 1 1
	assertBuiltIndex(t, built, map[string]bool{
		"label1 a": true,
		"label2 b": true,
		"label3 c": true,
		"label4 d": true,
		"3 a;b":    true,
		"5 a;c":    true,
		"6 b;c":    true,
		"7 a;b;c":  true,
	})
}

func TestIndexConfigIgnoreCase(t *testing.T) {
	idx := NewIndexes(&Config{IgnoreCase: true})

	idx.Add("label1", "abc dあいbCh", "saMPle")
	idx.AddBigrams("label2", "abc dあいbCh")
	idx.AddBiunigrams("label3", "abc dあいbCh")
	idx.AddPrefixes("label4", "abc dあいbCh")
	idx.AddSomething("label5", []string{"abc dあいbCh", "AbcdeF"})

	expected := make(map[string]bool)

	// Add
	expected["label1 abc dあいbch"] = true
	expected["label1 sample"] = true

	// AddBigrams
	for _, s := range Bigrams("abc dあいbch") {
		expected["label2 "+s] = true
	}

	// AddBiunigrams
	for _, s := range Biunigrams("abc dあいbch") {
		expected["label3 "+s] = true
	}

	// AddPrefixes
	for _, s := range Prefixes("abc dあいbch") {
		expected["label4 "+s] = true
	}

	// AddSomething
	expected["label5 abc dあいbch"] = true
	expected["label5 abcdef"] = true

	built := idx.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestIndexConfigSaveNoFiltersIndex(t *testing.T) {
	idx := NewIndexes(&Config{SaveNoFiltersIndex: true})
	idx.Add("label1", "a")

	built := idx.MustBuild()
	assertBuiltIndex(t, built, map[string]bool{
		"label1 a":     true,
		IndexNoFilters: true,
	})
}

func assertBuiltIndex(t *testing.T, actual, expected map[string]bool) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected, actual: `%v`, expected: `%v`", actual, expected)
	}
}
