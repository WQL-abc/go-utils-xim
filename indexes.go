package xim

import (
	"strings"

	"golang.org/x/xerrors"
)

// Indexes - extra indexes for firestore query.
type Indexes struct {
	m    indexesMap // key=label, value=indexes
	conf *Config
}

// NewIndexes - creates and initializes a new Indexes.
func NewIndexes(conf *Config) *Indexes {
	if conf == nil {
		conf = DefaultConfig
	}
	return &Indexes{
		m:    make(indexesMap),
		conf: conf,
	}
}

func (idxs *Indexes) add(label string, indexes ...string) {
	for _, idx := range indexes {
		if idxs.conf.IgnoreCase {
			idx = strings.ToLower(idx)
		}

		if _, ok := idxs.m[label]; !ok {
			idxs.m[label] = make(map[string]struct{})
		}

		idxs.m[label][idx] = struct{}{}
	}
}

// Add - adds new indexes with a label.
func (idxs *Indexes) Add(label string, indexes ...string) *Indexes {
	idxs.add(label, indexes...)
	return idxs
}

// AddBigrams - adds new bigram indexes with a label.
func (idxs *Indexes) AddBigrams(label string, s string) *Indexes {
	return idxs.Add(label, Bigrams(s)...)
}

// AddBiunigrams - adds new biunigram indexes with a label.
func (idxs *Indexes) AddBiunigrams(label string, s string) *Indexes {
	return idxs.Add(label, Biunigrams(s)...)
}

// AddPrefixes - adds new prefix indexes with a label.
func (idxs *Indexes) AddPrefixes(label string, s string) *Indexes {
	return idxs.Add(label, Prefixes(s)...)
}

// AddSuffixes - adds new prefix indexes with a label.
func (idxs *Indexes) AddSuffixes(label string, s string) *Indexes {
	return idxs.Add(label, Suffixes(s)...)
}

// AddSomething - adds new indexes with a label.
// The indexes can be a slice or a string convertible value.
func (idxs *Indexes) AddSomething(label string, indexes interface{}) *Indexes {
	addSomething(idxs, label, indexes)
	return idxs
}

// Build - builds indexes to save.
func (idxs Indexes) Build() (map[string]bool, error) {
	built := buildIndexes(idxs.m, nil)

	if len(idxs.conf.CompositeIdxLabels) > 1 {
		cis, err := createCompositeIndexes(idxs.conf.CompositeIdxLabels, idxs.m, false)
		if err != nil {
			return nil, err
		}
		for s, b := range cis {
			built[s] = b
		}
	}

	if idxs.conf.SaveNoFiltersIndex {
		built[IndexNoFilters] = true
	}

	if len(built) > MaxIndexesSize {
		return nil, xerrors.Errorf("index size exceeds %d", MaxIndexesSize)
	}

	return built, nil
}

// MustBuild - builds indexes to save and panics with error.
func (idxs Indexes) MustBuild() map[string]bool {
	built, err := idxs.Build()
	if err != nil {
		panic(err)
	}
	return built
}
