/* schema_test.go
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
	"testing"
)

var ExampleContent []byte

var (
	TestSectionID = "applets/{7a6ad2d4-a770-11eb-a9a4-0242de3bcd68}"
)

// TestParseSchemaLine will test ParseSchemaLine
func TestParseSchemaLine(t *testing.T) {
	var key string
	var sT *SchemaType

	key, sT = ParseSchemaLine("dark-theme=true")

	if key != "dark-theme" {
		t.Error("Failed to get correct key, expected: dark-theme")
	}

	if sT.Type != "bool" { // Failed to parse as a boolean
		t.Errorf("dark-theme was not parsed as expected bool, got %v instead.", sT.Type)
	}

	if sT.BoolVal != true { // Not correct boolean value
		t.Errorf("dark-theme does not have expected value of true, got %v instead.", sT.BoolVal)
	}

	key, sT = ParseSchemaLine("layout='solus-fortitude'")

	if key != "layout" { // Not layout
		t.Error("Failed to get correct key, expected: layout")
	}

	if sT.Type != "string" { // Not a string
		t.Errorf("layout was not parsed as expected string, got %v instead.", sT.Type)
	}

	if sT.Val != "'solus-fortitude'" { // Not 'solus-foritude'
		t.Errorf("layout does not have expected value of 'solus-fortitude', got %v instead.", sT.Val)
	}

	key, sT = ParseSchemaLine("speed=-0.63571428571428568")

	if key != "speed" { // Failed to get the schema key
		t.Error("Failed to get correct key, expected: speed")
	}

	if sT.Type != "float64" { // Did not get expected type
		t.Errorf("size was not parsed as expected float64, got %v instead.", sT.Type)
	}

	if sT.FloatVal != -0.63571428571428568 { // Did not get correct value
		t.Errorf("size does not have expected value of -0.63571428571428568, got %v instead.", sT.FloatVal)
	}

	key, sT = ParseSchemaLine("size=uint32 39")

	if key != "size" { // Failed to get the schema key
		t.Error("Failed to get correct key, expected: size")
	}

	if sT.Type != "uint32" { // Did not get expected type
		t.Errorf("size was not parsed as expected uint32, got %v instead.", sT.Type)
	}

	if sT.UintVal != 39 { // Did not get correct value
		t.Errorf("size does not have expected value of 39, got %v instead.", sT.UintVal)
	}
}

// TestGetSection will test GetSection
func TestGetSection(t *testing.T) {
	var getErr error
	if TestSchemaKV, getErr = TestSchema.GetSection(TestSectionID); getErr != nil { // Failed to get the section
		t.Fatalf("Failed to get section called: applets/{7a6ad2d4-a770-11eb-a9a4-0242de3bcd68}")
	}
}

// TestAddSection will test AddSection
func TestAddSection(t *testing.T) {
	var addErr error
	if addErr = TestSchema.AddSection(TestSectionID, TestSchemaKV); addErr == nil { // If we successfully added a duplicate section
		t.Error("Added section applets/{7a6ad2d4-a770-11eb-a9a4-0242de3bcd68}, which already existed. That should not be allowed.")
	}

	if addErr = TestSchema.AddSection("applet/look-at-me-i-am-special", TestSchemaKV); addErr != nil { // Failed to add this TestSchemaKV
		t.Errorf("Failed to add section applet/look-at-me-i-am-special: %v", addErr)
	}
}

// TestDeleteSection will test DeleteSections
func TestDeleteSections(t *testing.T) {
	testID := "applet/look-at-me-i-am-special"
	TestSchema.DeleteSections(testID) // Delete this section

	if TestSchema.HasSection(testID) { // Still have it
		t.Errorf("Failed to delete section: %s", testID)
	}
}

// TestDeleteSectionsWithPrefix will test DeleteSectionsWithPrefix
func TestDeleteSectionsWithPrefix(t *testing.T) {
	testID := "instance/budgie-menu/{8bbab560-0dae-11eb-ad1d-e0d55e200f1c}"
	TestSchema.DeleteSectionsWithPrefix("instance") // Delete everything with instance

	if TestSchema.HasSection(testID) { // Still have this instance test ID
		t.Errorf("Failed to delete sections starting with instance, as %s still exists", testID)
	}
}

// TestHasSection will test HasSection
func TestHasSection(t *testing.T) {
	if !TestSchema.HasSection(TestSectionID) { // Failed to get a valid section
		t.Errorf("Failed to get valid section by key: %s", TestSectionID)
	}
}

// TestMigrateSectionsWithName will test MigrateSectionsWithName
func TestMigrateSectionsWithName(t *testing.T) {
	oldPanelKey := "panels/{8bba3bf8-0dae-11eb-ad1d-e0d55e200f1c}"
	newPanelKey := "moved-panels/{8bba3bf8-0dae-11eb-ad1d-e0d55e200f1c}"

	TestSchema.MigrateSectionsWithName("panels/", "moved-panels/", false)

	if TestSchema.HasSection(oldPanelKey) { // Old panel key still exists
		t.Errorf("Key %s still exists. Should have been moved to %s", oldPanelKey, newPanelKey)
	}

	if !TestSchema.HasSection(newPanelKey) { // New panel key does not exist
		t.Errorf("New key %s does not exist after migration", newPanelKey)
	}
}
