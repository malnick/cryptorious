package config

import (
	"testing"
)

var tc = Config{}

func TestSetDefaults(t *testing.T) {
	err := tc.setDefaults()

	if err != nil {
		t.Error("Exepcted no errors setting default config, got", err.Error())
	}
}
