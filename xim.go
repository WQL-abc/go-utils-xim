package xim

import (
	"bytes"
	"fmt"

	"golang.org/x/xerrors"
)

const (
	IndexNoFilters          = "__NF__" // index to be used for no-filters.
	MaxIndexesSize          = 512      // maximum size of indexes.
	MaxCompositeIndexLabels = 8        // maximum number of labels for composite index.
)

const combinationIndexSeparator = ";"

// Config - describe extra indexes configuration.
type Config struct {
	CompositeIdxLabels []string // label list which defines composite indexes to improve the search performance
	IgnoreCase         bool     // defines whether to ignore case on search
	SaveNoFiltersIndex bool     // defines whether to save IndexNoFilters index.
}

// DefaultConfig - default configuration.
var DefaultConfig = new(Config)

// ValidateConfig - validates Config fields.
func ValidateConfig(conf *Config) (*Config, error) {
	if len(conf.CompositeIdxLabels) > MaxCompositeIndexLabels {
		return nil, xerrors.Errorf("CompositeIdxLabels size exceeds %d", MaxCompositeIndexLabels)
	}
	return conf, nil
}

// MustValidateConfig - validates fields and panics if it's invalid.
func MustValidateConfig(conf *Config) *Config {
	conf, err := ValidateConfig(conf)
	if err != nil {
		panic(err)
	}
	return conf
}

// common indexes map
// key=label, value=index set
type indexesMap map[string]map[string]struct{}

// buildIndexes - builds indexes from m.
// m is map[label]tokens.
func buildIndexes(m indexesMap, labelsToExclude []string) map[string]bool {
	idxSet := make(map[string]bool)

	excludeSet := make(map[string]struct{})
	for _, l := range labelsToExclude {
		excludeSet[l] = struct{}{}
	}

	for label, tokens := range m {
		if _, ok := excludeSet[label]; ok {
			continue
		}
		for t := range tokens {
			idxSet[fmt.Sprintf("%s %s", label, t)] = true
		}
	}

	return idxSet
}

func appendCombinationIndex(indexes, index string) string {
	buf := bytes.NewBufferString(indexes)
	if buf.Len() > 0 {
		buf.WriteString(combinationIndexSeparator)
	}
	buf.WriteString(index)
	return buf.String()
}

// createCompositeIndexes - creates composite indexes of labels from m.
// It reduces zig-zag merge join latency.
// forFilters is used for Filters.
func createCompositeIndexes(labels []string, m indexesMap, forFilters bool) (map[string]bool, error) {
	if len(labels) > MaxCompositeIndexLabels {
		return nil, xerrors.Errorf("CompositeIdxLabels size exceeds %d", MaxCompositeIndexLabels)
	}

	indexes := make(map[string]bool, 64)

	f := func(combination uint8, index string, someNew bool) {
		if forFilters && !someNew {
			return
		}
		indexes[fmt.Sprintf("%d %s", combination, index)] = true
	}

	// used indexes sets for filters
	usedIndexes := make(indexesMap)

	// generate combination indexes with bit operation
	// mapping each labels to each bits.

	// construct recursive funcs at first.
	// reverse loop for labels so that the first label will be right-end bit.
	var combinationForFilter uint8
	for i := len(labels) - 1; i >= 0; i-- {
		i := i
		prevF := f
		idxLabel := labels[i]

		if len(m[idxLabel]) > 0 {
			combinationForFilter |= 1 << uint(i)
		}

		f = func(combination uint8, index string, someNew bool) {
			if combination&(1<<uint(i)) == 0 {
				// no process bit for the combination.
				prevF(combination, index, someNew)
				return
			}
			// check process bit for the combination.
			tokens := make([]string, 0, len(m[idxLabel]))
			for token := range m[idxLabel] {
				tokens = append(tokens, token)
			}
			for _, token := range tokens {
				combinationIndex := appendCombinationIndex(index, token)

				// check if the token is already used for filters
				if _, ok := usedIndexes[idxLabel]; !ok {
					usedIndexes[idxLabel] = make(map[string]struct{})
				}
				if _, ok := usedIndexes[idxLabel][token]; !ok {
					usedIndexes[idxLabel][token] = struct{}{}
					someNew = true
				}

				prevF(combination, combinationIndex, someNew) // recursive call
			}
		}
	}

	// now generate indexes.
	if forFilters {
		f(combinationForFilter, "", false)
	} else {
		for i := 3; i < (1 << uint(len(labels))); i++ {
			if (i & (i - 1)) == 0 {
				// do not save single index
				continue
			}
			f(uint8(i), "", false)
		}
	}

	return indexes, nil
}
