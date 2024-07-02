// version_test.go
package params

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	expectedVersion := fmt.Sprintf("%d.%d.%d-%s", VersionMajor, VersionMinor, VersionPatch, VersionMeta)
	if Version != expectedVersion {
		t.Errorf("Expected version %s, but got %s", expectedVersion, Version)
	}
}

func TestVersionFormat(t *testing.T) {
	expectedVersion := "1.0.0-unstable"
	if Version != expectedVersion {
		t.Errorf("Expected version %s, but got %s", expectedVersion, Version)
	}
}

func TestVersionComponents(t *testing.T) {
	if VersionMajor != 1 {
		t.Errorf("Expected VersionMajor to be 1, but got %d", VersionMajor)
	}
	if VersionMinor != 0 {
		t.Errorf("Expected VersionMinor to be 0, but got %d", VersionMinor)
	}
	if VersionPatch != 0 {
		t.Errorf("Expected VersionPatch to be 0, but got %d", VersionPatch)
	}
	if VersionMeta != "unstable" {
		t.Errorf("Expected VersionMeta to be 'unstable', but got %s", VersionMeta)
	}
}
