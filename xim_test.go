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
