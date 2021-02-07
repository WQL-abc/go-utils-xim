package xim

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/xerrors"
)

// Filters - filters builder for extra indexes.
type Filters struct {
	m    indexesMap // key=label, value=index set
	conf *Config
}

// NewFilters - creates and initializes a new Filters.
func NewFilters(conf *Config) *Filters {
	if conf == nil {
		conf = DefaultConfig
	}
	return &Filters{
		m:    make(indexesMap),
		conf: conf,
	}
}

func (filters *Filters) add(label string, indexes ...string) {
	for _, idx := range indexes {
		if filters.conf.IgnoreCase {
			idx = strings.ToLower(idx)
		}

		if _, ok := filters.m[label]; !ok {
			filters.m[label] = make(map[string]struct{})
		}

		filters.m[label][idx] = struct{}{}
	}
}

// Add - adds new filters with a label.
func (filters *Filters) Add(label string, indexes ...string) *Filters {
	filters.add(label, indexes...)
	return filters
}

// AddBigrams - adds new bigram filters with a label.
func (filters *Filters) AddBigrams(label string, s string) *Filters {
	return filters.AddBiunigrams(label, s)
}

// AddBiunigrams - adds new biunigram filters with a label.
func (filters *Filters) AddBiunigrams(label string, s string) *Filters {
	if runeLen := utf8.RuneCountInString(s); runeLen == 1 {
		filters.Add(label, s)
	} else if runeLen > 1 {
		filters.Add(label, Bigrams(s)...)
	}
	return filters
}

// AddPrefix - adds a new prefix filters with a label.
func (filters *Filters) AddPrefix(label string, s string) *Filters {
	// don't need to split prefixes on filters
	return filters.Add(label, s)
}

// AddSuffix - adds a new suffix filters with a label.
func (filters *Filters) AddSuffix(label string, s string) *Filters {
	// don't need to split suffixes on filters
	return filters.Add(label, s)
}

// AddSomething - adds new filter with a label.
// The indexes can be a slice or a string convertible value.
func (filters *Filters) AddSomething(label string, indexes interface{}) *Filters {
	addSomething(filters, label, indexes)
	return filters
}

// Build - builds filters to save.
func (filters *Filters) Build() (map[string]bool, error) {
	built := buildIndexes(filters.m, filters.conf.CompositeIdxLabels)

	if len(filters.conf.CompositeIdxLabels) > 1 {
		cis, err := createCompositeIndexes(filters.conf.CompositeIdxLabels, filters.m, true)
		if err != nil {
			return nil, err
		}
		for s, b := range cis {
			built[s] = b
		}
	}

	if filters.conf.SaveNoFiltersIndex && len(built) == 0 {
		built[IndexNoFilters] = true
	}

	if len(built) > MaxIndexesSize {
		return nil, xerrors.Errorf("index size exceeds %d", MaxIndexesSize)
	}

	return built, nil
}

// MustBuild - builds filters to save and panics with error.
func (filters Filters) MustBuild() map[string]bool {
	built, err := filters.Build()
	if err != nil {
		panic(err)
	}
	return built
}
