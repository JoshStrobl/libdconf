/* schemaKV.go
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
	"regexp"
	"sort"
	"strings"
)

// AddKey will attempt to add the SchemaType for the provided key to the SchemaKV
// This will return an error if the key already exists
func (kv *SchemaKV) AddKey(key string, t *SchemaType) error {
	if kv.HasKey(key) { // Already exists
		return ErrKeyAlreadyExists
	}

	kv.Keys[key] = t
	kv.Order = append(kv.Order, key)
	sort.Strings(kv.Order)

	return nil
}

// DeleteKeys will delete all specified keys from a SchemaKV
func (kv *SchemaKV) DeleteKeys(keys ...string) {
	for _, key := range keys { // For each key
		kv.Order = RemoveFromStringArr(kv.Order, key) // Remove from the order
		delete(kv.Keys, key)                          // Delete the key
	}
}

// Duplicate will duplicate this SchemaKV into a new SchemaKV
func (kv *SchemaKV) Duplicate() *SchemaKV {
	newKv := SchemaKV{
		Order: []string{},
		Keys:  make(map[string]*SchemaType),
	}

	for kvKey, kvVal := range kv.Keys {
		newKv.Keys[kvKey] = kvVal.Duplicate() // Duplicate the SchemaType and assign it
	}

	copy(newKv.Order, kv.Order) // Make a copy of the old Kv order to the new one
	return &newKv
}

// GetVal will get the SchemaType value for the provided key, or return an error
func (kv *SchemaKV) GetVal(key string) (*SchemaType, error) {
	val, exists := kv.Keys[key]

	if !exists { // Key does not exist
		return nil, ErrKeyNotExists
	}

	return val, nil
}

// HasKey returns if we have this key
func (kv *SchemaKV) HasKey(key string) bool {
	_, exists := kv.Keys[key]
	return exists
}

// ModifyKey will attempt to modify a SchemaType specified by key, with the provided Modification
func (kv *SchemaKV) ModifyKey(key string, mod Modification) (modErr error) {
	if !kv.HasKey(key) { // If we don't have this key
		modErr = ErrKeyNotExists
		return
	}

	hasReplaceVal := len(mod.ReplaceValues) == 2
	hasValue := mod.Value != ""

	if !hasReplaceVal && !hasValue { // Have neither ReplaceValue or Value
		modErr = ErrModNoReplaceValueOrValue
		return
	}

	if mod.Value != "" { // If we have a value defined, so we're not doing something complex like string regex
		parsedSt, parseStErr := NewSchemaType(mod.Value) // Attempt to parse our provided value into a SchemaType

		if parseStErr != nil {
			modErr = parseStErr
			return
		}

		kv.Keys[key] = parsedSt // Just update our key with the new SchemaType
		return
	}

	existingSt, _ := kv.Keys[key] // Get the current SchemaType

	if existingSt.Type != "ArrayAsString" && existingSt.Type != "string" { // Not string manipulation
		if hasReplaceVal { // If we are trying to replace one string with another, and this isn't it
			modErr = ErrModCannotDoReplace
			return
		}

		return
	}

	finding := mod.ReplaceValues[0]
	replacement := mod.ReplaceValues[1]

	if strings.HasPrefix(finding, "re:") { // Is intended to be regex
		finding = strings.TrimPrefix(finding, "re:")

		reg, regCompileErr := regexp.Compile(finding) // Attempt to create our regexp.Regexp struct

		if regCompileErr != nil { // Not valid regexp
			modErr = regCompileErr
			return
		}

		existingSt.Val = reg.ReplaceAllString(existingSt.Val, replacement) // Replace all using regexp
	} else { // Not intended to be regex
		existingSt.Val = strings.ReplaceAll(existingSt.Val, finding, replacement) // Replace all instances
	}

	return
}

// MoveKey will attempt to move the source key to the destination.
// If the source does not exist or the destination already exists, returns an error
func (kv *SchemaKV) MoveKey(source string, dest string) error {
	sT, sourceExists := kv.Keys[source]

	if !sourceExists { // Source does not exist
		return ErrKeyNotExists
	}

	if _, exists := kv.Keys[dest]; exists { // Destination already exists
		return ErrKeyAlreadyExists
	}

	kv.AddKey(dest, sT)   // Add the new section
	kv.DeleteKeys(source) // Delete the source key

	return nil
}
