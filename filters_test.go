package xim

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestAddFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.Add("label1", "abc dあいbCh", "sample")
	filter.Add("label2", "abc debch iJあdeN", "sample")

	built := filter.MustBuild()
	assertBuiltFilter(t, built, map[string]bool{
		"label1 abc dあいbCh":       true,
		"label1 sample":           true,
		"label2 abc debch iJあdeN": true,
		"label2 sample":           true,
	})
}

func TestAddBigramsFilter(t *testing.T) {
	t.Run("indexes == single character", func(t *testing.T) {
		filter := NewFilters(nil)
		filter.AddBigrams("label1", "a")
		filter.AddBigrams("label2", "b")

		expected := make(map[string]bool)
		expected["label1 "+"a"] = true
		expected["label2 "+"b"] = true

		built := filter.MustBuild()
		assertBuiltIndex(t, built, expected)
	})

	t.Run("indexes == 2 or more characters", func(t *testing.T) {
		filter := NewFilters(nil)
		filter.AddBigrams("label1", "abc dあいbCh")
		filter.AddBigrams("label2", "abc debch iJあdeN")

		expected := make(map[string]bool)
		for _, s := range Bigrams("abc dあいbCh") {
			expected["label1 "+s] = true
		}
		for _, s := range Bigrams("abc debch iJあdeN") {
			expected["label2 "+s] = true
		}

		built := filter.MustBuild()
		assertBuiltIndex(t, built, expected)
	})
}

func TestAddBiunigramsFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.AddBiunigrams("label1", "abc dあいbCh")
	filter.AddBiunigrams("label2", "abc debch iJあdeN")

	// Use Bigrams for Filters
	expected := make(map[string]bool)
	for _, s := range Bigrams("abc dあいbCh") {
		expected["label1 "+s] = true
	}
	for _, s := range Bigrams("abc debch iJあdeN") {
		expected["label2 "+s] = true
	}

	built := filter.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestAddPrefixFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.AddPrefix("label1", "abc dあいbCh")
	filter.AddPrefix("label2", "abc debch iJあdeN")

	built := filter.MustBuild()
	assertBuiltFilter(t, built, map[string]bool{
		"label1 abc dあいbCh":       true,
		"label2 abc debch iJあdeN": true,
	})
}

func TestAddSuffixFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.AddSuffix("label1", "abc dあいbCh")
	filter.AddSuffix("label2", "abc debch iJあdeN")

	built := filter.MustBuild()
	assertBuiltFilter(t, built, map[string]bool{
		"label1 abc dあいbCh":       true,
		"label2 abc debch iJあdeN": true,
	})
}

func TestAddSomethingFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.AddSomething("label1", []string{"abc dあいbCh", "abc debch iJあdeN"})
	filter.AddSomething("label2", 123)
	now := time.Now()
	filter.AddSomething("label3", now)

	built := filter.MustBuild()
	assertBuiltFilter(t, built, map[string]bool{
		"label1 abc dあいbCh":                      true,
		"label1 abc debch iJあdeN":                true,
		"label2 123":                             true,
		fmt.Sprintf("label3 %d", now.UnixNano()): true,
	})
}

func TestAddAllFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.Add("label1", "abc dあいbCh", "sample")
	filter.AddBigrams("label2", "abc dあいbCh")
	filter.AddBiunigrams("label3", "abc dあいbCh")
	filter.AddPrefix("label4", "abc dあいbCh")
	filter.AddSomething("label5", []string{"abc dあいbCh", "AbcdeF"})

	expected := make(map[string]bool)

	// Add
	expected["label1 abc dあいbCh"] = true
	expected["label1 sample"] = true

	// AddBigrams
	for _, s := range Bigrams("abc dあいbCh") {
		expected["label2 "+s] = true
	}

	// AddBiunigrams
	for _, s := range Bigrams("abc dあいbCh") {
		expected["label3 "+s] = true
	}

	// AddPrefix
	expected["label4 abc dあいbCh"] = true

	// AddSomething
	expected["label5 abc dあいbCh"] = true
	expected["label5 AbcdeF"] = true

	built := filter.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestBuildFilter(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		filter := NewFilters(nil)
		for i := 0; i < MaxIndexesSize; i++ {
			filter.Add(fmt.Sprintf("label%d", i), "abc")
		}

		built, err := filter.Build()
		if err != nil {
			t.Errorf("error = %s, wants = nil", err)
		}
		if built == nil {
			t.Error("built = nil, wants != nil")
		}
	})

	t.Run("Config.CompositeIdxLabels > MaxCompositeIndexLabels", func(t *testing.T) {
		labels := make([]string, MaxCompositeIndexLabels+1)
		for i := 0; i < len(labels); i++ {
			labels[i] = string(rune('a' + i))
		}

		filter := NewFilters(&Config{CompositeIdxLabels: labels})
		if _, err := filter.Build(); err == nil {
			t.Error("error = nil, wants != nil")
		}
	})

	t.Run("Number of Build() results > MaxIndexesSize", func(t *testing.T) {
		filter := NewFilters(nil)
		for i := 0; i < MaxIndexesSize+1; i++ {
			filter.Add(fmt.Sprintf("label%d", i), "abc")
		}

		if _, err := filter.Build(); err == nil {
			t.Error("error = nil, wants != nil")
		}
	})
}

func TestMustBuildFilter(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		defer func() {
			if rec := recover(); rec != nil {
				t.Errorf("expected: not panic, was: panic = [%v]\n", rec)
			}
		}()

		filter := NewFilters(nil)
		for i := 0; i < MaxIndexesSize; i++ {
			filter.Add(fmt.Sprintf("label%d", i), "abc")
		}

		built := filter.MustBuild()
		if built == nil {
			t.Error("built = nil, wants != nil")
		}
	})

	t.Run("Config.CompositeIdxLabels > MaxCompositeIndexLabels", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected: panic, was: not panic\n")
			}
		}()

		labels := make([]string, MaxCompositeIndexLabels+1)
		for i := 0; i < len(labels); i++ {
			labels[i] = string(rune('a' + i))
		}

		filter := NewFilters(&Config{CompositeIdxLabels: labels})
		filter.MustBuild()
	})

	t.Run("Number of MustBuild() results > MaxIndexesSize", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected: panic, was: not panic\n")
			}
		}()

		filter := NewFilters(nil)
		for i := 0; i < MaxIndexesSize+1; i++ {
			filter.Add(fmt.Sprintf("label%d", i), "abc")
		}

		filter.MustBuild()
	})
}

func TestFilterConfigCompositeIdxLabels(t *testing.T) {
	filter := NewFilters(&Config{CompositeIdxLabels: []string{"label1", "label2", "label3"}})
	filter.Add("label1", "a")
	filter.Add("label2", "b")
	filter.Add("label3", "c")

	built := filter.MustBuild()

	//   c b a
	//   -----
	//3  0 1 1
	//5  1 0 1
	//6  1 1 0
	//7  1 1 1
	assertBuiltIndex(t, built, map[string]bool{
		// "3 a;b",
		// "5 a;c",
		// "6 b;c",

		// now needs only indexes with all specified labels
		"7 a;b;c": true,
	})
}

func TestFilterConfigIgnoreCase(t *testing.T) {
	filter := NewFilters(&Config{IgnoreCase: true})
	filter.Add("label1", "abc dあいbCh", "saMPle")
	filter.AddBigrams("label2", "abc dあいbCh")
	filter.AddBiunigrams("label3", "abc dあいbCh")
	filter.AddPrefix("label4", "abc dあいbCh")
	filter.AddSomething("label5", []string{"abc dあいbCh", "AbcdeF"})

	expected := make(map[string]bool)

	// Add
	expected["label1 abc dあいbch"] = true
	expected["label1 sample"] = true

	// AddBigrams
	for _, s := range Bigrams("abc dあいbch") {
		expected["label2 "+s] = true
	}

	// AddBiunigrams
	for _, s := range Bigrams("abc dあいbch") {
		expected["label3 "+s] = true
	}

	// AddPrefix
	expected["label4 abc dあいbch"] = true

	// AddSomething
	expected["label5 abc dあいbch"] = true
	expected["label5 abcdef"] = true

	built := filter.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestFilterConfigSaveNoFiltersIndex(t *testing.T) {
	t.Run("No filter Add", func(t *testing.T) {
		filter := NewFilters(&Config{SaveNoFiltersIndex: true})

		built := filter.MustBuild()

		assertBuiltIndex(t, built, map[string]bool{
			IndexNoFilters: true,
		})
	})

	t.Run("One or more filters are Add", func(t *testing.T) {
		filter := NewFilters(&Config{SaveNoFiltersIndex: true})
		filter.Add("label1", "a")

		built := filter.MustBuild()

		assertBuiltIndex(t, built, map[string]bool{
			"label1 a": true,
		})
	})
}

func assertBuiltFilter(t *testing.T, actual, expected map[string]bool) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected, actual: `%v`, expected: `%v`", actual, expected)
	}
}
