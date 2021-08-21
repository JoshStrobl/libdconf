/* schemaKV_test.go
 *
 * Copyright 2021 Joshua Strobl
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package libdconf

import (
	"strings"
	"testing"
)

// TestAddKey will test AddKey
func TestAddKey(t *testing.T) {
	BootstrapSetKV("instance/icon-tasklist/{8bbd5acc-0dae-11eb-ad1d-e0d55e200f1c}")

	key, showST := ParseSchemaLine("show-all-windows-on-click=true")

	if addErr := TestSchemaKV.AddKey(key, showST); addErr != nil {
		t.Errorf("Failed to add our show-all-windows-on-click schema type: %s", addErr)
	}
}

// TestDeleteKey will test DeleteKey
func TestDeleteKey(t *testing.T) {
	TestSchemaKV.DeleteKeys("show-all-windows-on-click") // Delete show-all-windows-on-click

	if TestSchemaKV.HasKey("show-all-windows-on-click") {
		t.Error("Failed to delete show-all-windows-on-click")
	}
}

// TestSchemaKVDuplicate will test SchemaKV's Duplicate
func TestSchemaKVDuplicate(t *testing.T) {
	pinnedKVDuplicate := TestSchemaKV.Duplicate() // Duplicate this key

	if !pinnedKVDuplicate.HasKey("only-pinned") { // Failed to duplicate properlty
		t.Errorf("Failed to properly duplicate KV, missing only-pinned.\n%v", pinnedKVDuplicate)
	}
}

// TestGetVal will test GetVal
func TestGetVal(t *testing.T) {
	val, getErr := TestSchemaKV.GetVal("only-pinned")

	if getErr != nil { // Failed to get the key
		t.Errorf("Failed to get only-pinned value: %s", getErr)
	}

	if val.Type != "bool" { // only-pinned is a bool so should be set as the val
		t.Errorf("only-pinned was expected to be bool, got %v instead.", val.Type)
	}
}

// TestHasKey will test HasKey
func TestHasKey(t *testing.T) {
	if !TestSchemaKV.HasKey("only-pinned") {
		t.Error("Failed to get only-pinned value, reported as not existing when it does.")
	}
}

// TestModifyKey will test ModifyKey
func TestModifyKey(t *testing.T) {
	celluloid := "io.github.celluloid_player.Celluloid.desktop"
	mod := Modification{
		ReplaceValues: []string{"io.github.GnomeMpv.desktop", celluloid},
		Value:         "",
	}

	if modErr := TestSchemaKV.ModifyKey("pinned-launchers", mod); modErr != nil { // Modify our pinned-launchers
		t.Errorf("Failed to modify pinned-launchers: %s", modErr)
	}

	val, _ := TestSchemaKV.GetVal("pinned-launchers") // Get pinned launchers SchemaType

	if !strings.Contains(val.Val, celluloid) {
		t.Errorf("Failed to properly perform replace for celluloid. SchemaType is: \n%v", val)
	}
}

// TestMoveKey tests MoveKey
func TestMoveKey(t *testing.T) {
	var moveErr error
	if moveErr = TestSchemaKV.MoveKey("pinned-launchers", "super-pinned-launchers"); moveErr != nil { // Attempt move
		t.Errorf("Failed to move pinned-launchers to super-pinned-launchers: %s", moveErr)
	}

	if _, moveErr = TestSchemaKV.GetVal("super-pinned-launchers"); moveErr != nil {
		t.Errorf("Failed to get super-pinned-launchers value: %s", moveErr)
	}
}
