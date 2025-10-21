package xray

import "testing"

func TestConfigEqualsObservatorySections(t *testing.T) {
	base := &Config{
		Observatory:      []byte(`{"enabled":false}`),
		BurstObservatory: []byte(`{"sample":1}`),
	}

	identical := *base
	if !base.Equals(&identical) {
		t.Fatalf("expected identical configs to be equal")
	}

	diffObservatory := *base
	diffObservatory.Observatory = []byte(`{"enabled":true}`)
	if base.Equals(&diffObservatory) {
		t.Fatalf("expected configs with different observatory to be unequal")
	}

	diffBurstObservatory := *base
	diffBurstObservatory.BurstObservatory = []byte(`{"sample":2}`)
	if base.Equals(&diffBurstObservatory) {
		t.Fatalf("expected configs with different burst observatory to be unequal")
	}
}
