package xim

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

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

// Add - adds new indexes with a label.
func (idxs *Indexes) Add(label string, indexes ...string) {
	for _, idx := range indexes {
		if idxs.conf.IgnoreCase {
			idx = strings.ToLower(idx)
		}

		if _, ok := idxs.m[label]; !ok {
			idxs.m[label] = make(map[string]struct{})
		}

		idxs.m[label][idx] = struct{}{}
	}
	return
}

// AddBigrams - adds new bigram indexes with a label.
func (idxs *Indexes) AddBigrams(label string, s string) *Indexes {
	idxs.Add(label, Bigrams(s)...)
	return idxs
}

// AddBiunigrams - adds new biunigram indexes with a label.
func (idxs *Indexes) AddBiunigrams(label string, s string) *Indexes {
	idxs.Add(label, Biunigrams(s)...)
	return idxs
}

// AddPrefixes - adds new prefix indexes with a label.
func (idxs *Indexes) AddPrefixes(label string, s string) *Indexes {
	idxs.Add(label, Prefixes(s)...)
	return idxs
}

// AddSuffixes - adds new prefix indexes with a label.
func (idxs *Indexes) AddSuffixes(label string, s string) *Indexes {
	idxs.Add(label, Suffixes(s)...)
	return idxs
}

// AddSomething - adds new indexes with a label.
// The indexes can be a slice or a string convertible value.
func (idxs *Indexes) AddSomething(label string, indexes interface{}) *Indexes {
	v := reflect.Indirect(reflect.ValueOf(indexes))

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			index := v.Index(i)
			if !index.CanInterface() {
				continue
			}
			idxs.Add(label, fmt.Sprintf("%v", index.Interface()))
		}
	case reflect.Struct:
		if v.Type() == timeType {
			unix := v.Interface().(time.Time).UnixNano()
			idxs.Add(label, strconv.FormatInt(unix, 10))
			break
		}
		fallthrough
	default:
		if v.CanInterface() {
			idxs.Add(label, fmt.Sprintf("%v", v.Interface()))
		}
	}

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
