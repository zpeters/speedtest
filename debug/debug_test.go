package debug

import (
	"testing"
)

// FIXME: need a better test for var existence
func TestDebug(t *testing.T) {
	if (DEBUG == false) || (DEBUG == true) {
		// exists
	} else {
		t.Errorf("DEBUG missing")
	}

	if (QUIET == false) || (QUIET == true) {
		// exists
	} else {
		t.Errorf("QUIET missing")
	}

	if (REPORT == false) || (REPORT == true) {
		//  exists
	} else {
		t.Errorf("REPORT missing")
	}
}

func TestNull(t *testing.T) {
}
