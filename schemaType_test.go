/* schemaType_test.go
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
	_ "strings"
	"testing"
)

var NumType *SchemaType

// TestNewSchemaType will test NewSchemaType
func TestNewSchemaType(t *testing.T) {
	BootstrapSetKV("instance/icon-tasklist/{8bbd5acc-0dae-11eb-ad1d-e0d55e200f1c}")

	var parseErr error
	if NumType, parseErr = NewSchemaType("uint32 1000"); parseErr != nil {
		t.Errorf("Failed to parse our uint32 schema type string: %s", parseErr)
	}

	if NumType.Type != "uint32" { // Not reported as uint32
		t.Errorf("Expected uint32, got %v instead.", NumType.Type)
	}

	if NumType.UintVal != 1000 { // Not our expected num
		t.Errorf("Expected value of 1000, got %v instead.", NumType.UintVal)
	}
}

// TestSchemaTypeDuplicate will test SchemaType's Duplicate
func TestSchemaTypeDuplicate(t *testing.T) {
	dup := NumType.Duplicate()

	if dup.Type != NumType.Type {
		t.Errorf("Failed to duplicate SchemaType, mismatching type. Got %s, expected %s", dup.Type, NumType.Type)
	}

	if dup.UintVal != NumType.UintVal {
		t.Errorf("Failed to duplicate SchemaType, mismatching value. Got %v, expected %v", dup.UintVal, NumType.UintVal)
	}
}

// TestMatches will test if the schema types match
func TestMatches(t *testing.T) {
	dup := NumType.Duplicate()

	if !NumType.Matches(dup) { // Dup does not match
		t.Error("NumType and dup do not match when they should.")
	}
}

// TestSchemaTypeString will test the SchemaType's String
func TestSchemaTypeString(t *testing.T) {
	str := NumType.String()

	if str != "uint32 1000" {
		t.Errorf("Expected NumType to be uint32 1000, got %s instead.", str)
	}
}
