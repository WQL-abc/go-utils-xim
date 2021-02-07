package xim

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	labels := make([]string, MaxCompositeIndexLabels+1)
	for i := 0; i < len(labels); i++ {
		labels[i] = string(rune('a' + i))
	}

	t.Run("len(CompositeIdxLabels)<=MaxCompositeIndexLabels", func(tr *testing.T) {
		conf := &Config{CompositeIdxLabels: labels[:MaxCompositeIndexLabels]}
		if _, err := ValidateConfig(conf); err != nil {
			tr.Errorf("expected: error = nil, but was: [%v]\n", err)
		}
		if validated, _ := ValidateConfig(conf); validated != conf {
			tr.Errorf("validated, _ := ValidateConfig(conf) expected: validated = conf, but was: validated = %#v\n", validated)
		}
	})

	t.Run("len(CompositeIdxLabels)>MaxCompositeIndexLabels", func(tr *testing.T) {
		conf := &Config{CompositeIdxLabels: labels[:MaxCompositeIndexLabels+1]}
		if _, err := ValidateConfig(conf); err == nil {
			tr.Error("CompositeIdxLabels > MaxCompositeIndexLabels expected: err != nil, but was: err = nil\n")
		}
	})

	t.Run("ValidateConfig(DefaultConfig)", func(tr *testing.T) {
		if _, err := ValidateConfig(DefaultConfig); err != nil {
			tr.Errorf("expected: error = nil, but was: error = [%v]\n", err)
		}
	})
}

func TestMustValidateConfig(t *testing.T) {
	labels := make([]string, MaxCompositeIndexLabels+1)
	for i := 0; i < len(labels); i++ {
		labels[i] = string(rune('a' + i))
	}

	t.Run("CompositeIdxLabels<=MaxCompositeIndexLabels", func(tr *testing.T) {
		defer func() {
			if rec := recover(); rec != nil {
				tr.Errorf("expected: not panic, was: panic = [%v]\n", rec)
			}
		}()

		conf := &Config{CompositeIdxLabels: labels[:MaxCompositeIndexLabels]}
		MustValidateConfig(conf)
	})

	t.Run("CompositeIdxLabels>MaxCompositeIndexLabels", func(tr *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				tr.Error("expected: panic, was: not panic\n")
			}
		}()

		conf := &Config{CompositeIdxLabels: labels[:MaxCompositeIndexLabels+1]}
		MustValidateConfig(conf)
	})
}

func TestAddIndexAndFilter(t *testing.T) {
	idx := NewIndexes(nil)
	idx.Add("label1", "abc dあいbCh", "sample")

	filter := NewFilters(nil)
	filter.Add("label1", "abc dあいbCh", "sample")

	builtIndexes := idx.MustBuild()
	builtFilters := filter.MustBuild()

	// All the contents of filter are present in index
	for builtFilter := range builtFilters {
		if !contains(t, builtIndexes, builtFilter) {
			t.Errorf("filter: %s not contains", builtFilter)
		}
	}
}

func TestAddBigramsIndexAndFilter(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddBigrams("label1", "abc dあいbCh")
	builtIndexes := idx.MustBuild()

	filter := NewFilters(nil)
	filter.AddBigrams("label1", "dあいb") // mid match of idx
	builtFilters := filter.MustBuild()

	// All the contents of filter are present in index
	for builtFilter := range builtFilters {
		if !contains(t, builtIndexes, builtFilter) {
			t.Errorf("filter: %s not contains", builtFilter)
		}
	}
}

func TestAddBiunigramsIndexAndFilter(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddBiunigrams("label1", "abc dあいbCh")
	builtIndexes := idx.MustBuild()

	filter := NewFilters(nil)
	filter.AddBiunigrams("label1", "dあいb") // mid match of idx
	builtFilters := filter.MustBuild()

	// All the contents of filter are present in index
	for builtFilter := range builtFilters {
		if !contains(t, builtIndexes, builtFilter) {
			t.Errorf("filter: %s not contains", builtFilter)
		}
	}
}

func TestInFilterIndexAndFilter(t *testing.T) {
	inBuilder := NewInBuilder()
	status1 := inBuilder.NewBit()
	status2 := inBuilder.NewBit()
	status3 := inBuilder.NewBit()

	idx := NewIndexes(nil)
	idx.Add("label1", inBuilder.Indexes(status1)...)
	idx.Add("label2", inBuilder.Indexes(status1, status2, status3)...)
	builtIndexes := idx.MustBuild()

	filter := NewFilters(nil)
	filter.Add("label1", inBuilder.Filter(status1, status2, status3))
	filter.Add("label2", inBuilder.Filter(status1))
	builtFilters := filter.MustBuild()

	// All the contents of filter are present in index
	for builtFilter := range builtFilters {
		if !contains(t, builtIndexes, builtFilter) {
			t.Errorf("filter: %s not contains", builtFilter)
		}
	}
}

func contains(t *testing.T, m map[string]bool, target string) bool {
	t.Helper()
	if _, ok := m[target]; ok {
		return true
	}
	return false
}
