/* structs.go
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

// Modification defines a desired change to a SchemaType
type Modification struct {
	// ReplaceValues is an array of RegExp (regular expression) for search to the replacement value
	// This is expected to be a length of two, the first being the value we are looking for and the second is the value we are replacing it with
	// The first value can be an exact string or regex
	ReplaceValues []string `toml:"replaceValue"`

	// Value is the raw value we are applying as the value for the modification
	Value string `toml:"value"`
}

// Schema is a map of paths to key values
type Schema struct {
	Order []string             // Our fixed order
	Map   map[string]*SchemaKV // Our Map of Sections (like com/solus-project/budgie-desktop/instance/icon-tasklist)
	Path  string               // Path for the Schema
}

// SchemaKV is a map of keys to our SchemaType
type SchemaKV struct {
	Order []string               // Our fixed order
	Keys  map[string]*SchemaType // Our Map of Keys in each Section
}

// SchemaType is our defined type
// This type will have a defined Type (e.g. "bool") as Type and the designated type set
// This allows us to perform less type checking and reflection during marshal and unmarshalling
type SchemaType struct {
	Type string

	BoolVal              bool
	FloatHadTrailingZero bool
	FloatVal             float64
	IntVal               int32
	UintVal              uint32
	Val                  string
}
