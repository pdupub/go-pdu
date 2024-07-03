// Copyright 2024 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

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
