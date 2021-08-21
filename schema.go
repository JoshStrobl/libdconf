/* schema.go
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
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

var (
	// SectionRegexp is our regular expression for a section name
	SectionRegexp = regexp.MustCompile(`^\[([A-Za-z0-9-_\:\.\/{}]+)\]`)
)

// NewSchema will attempt to create a new Schema based on the root of the dconf directory
// content is optional. If no content is specified, we will attempt to do a dconf dump instead
// This will be done user running the command.
// If we fail to dump or parse the schema, we will return an error
func NewSchema(path string, content []byte) (schema *Schema, readErr error) {
	if content == nil || (len(content) == 0) { // content not specified or has no content
		readErr = ErrNoContentProvided
		return
	}

	schema = &Schema{
		Map:   make(map[string]*SchemaKV),
		Order: []string{},
		Path:  path,
	}

	cs := string(content[:]) // Convert our byte array to a string

	lines := strings.Split(cs, "\n") // Split on new line

	var currentSection string // Our current section, sub-folder, whatever you want to call it
	var currentKV *SchemaKV   // Key value store

	for _, line := range lines {
		if line == "" { // If this is an empty new line
			if currentSection != "" && currentKV != nil { // If we have a current section
				schema.Map[currentSection] = currentKV
			}

			currentSection = "" // Reset to an empty section name
			currentKV = nil     // Force as nil for future check
			continue
		}

		if currentKV == nil { // No valid map
			currentKV = &SchemaKV{
				Order: []string{},
				Keys:  make(map[string]*SchemaType),
			}
		}

		sectionMatches := SectionRegexp.FindAllString(line, 1) // Find all our matches, limit to 1 though

		if len(sectionMatches) == 1 { // If we have a match
			currentSection = sectionMatches[0]
			currentSection = strings.TrimPrefix(currentSection, "[")
			currentSection = strings.TrimSuffix(currentSection, "]")

			schema.Order = append(schema.Order, currentSection) // Add our section
		} else { // If we didn't find a section per our regex
			key, sT := ParseSchemaLine(line) // Attempt to parse our schema type

			if sT.Type != "" { // If we got something. omg yay
				currentKV.AddKey(key, sT) // Add the key
			}
		}
	}

	return
}

// ParseSchemaLine will parse our key=val line in an attempt to figure out its type
func ParseSchemaLine(line string) (key string, t *SchemaType) {
	keyValArr := strings.SplitN(line, "=", 2) // Split between key and value

	if len(keyValArr) != 2 { // Not a key=val
		return
	}

	key = keyValArr[0]     // Set key to first position in array
	rawVal := keyValArr[1] // Set our raw value

	parsedSt, parseErr := NewSchemaType(rawVal) // Attempt to parse our "raw" value to a SchemaType

	if parseErr == nil {
		t = parsedSt
	}

	return
}

// AddSection will attempt to add the provided SchemaKV as the provided section name
// This will return an error if the section already exists
func (schema *Schema) AddSection(section string, sT *SchemaKV) (addErr error) {
	if schema.HasSection(section) { // Section already exists
		addErr = ErrSectionExists
		return
	}

	schema.Map[section] = sT
	schema.Order = append(schema.Order, section) // Append to the existing order
	return
}

// DeleteSections will delete all sections specified should they match exactly
// If you want to match by prefix, use the DeleteSectionsWithPrefix func
func (schema *Schema) DeleteSections(sections ...string) {
	for _, section := range sections { // For each section
		section = TrimSectionSlashes(section)
		delete(schema.Map, section)                               // Delete the section
		schema.Order = RemoveFromStringArr(schema.Order, section) // Remove the section from the string array
	}
}

// DeleteSectionsWithPrefix will delete all sections specified and all sections which begin with the section prefixed
func (schema *Schema) DeleteSectionsWithPrefix(sections ...string) {
	for sectionInMap := range schema.Map { // For each section in the map
		for _, section := range sections { // For each of the sections
			section = TrimSectionSlashes(section)

			if strings.HasPrefix(sectionInMap, section) { // If this section in the map begins with the section specified
				schema.DeleteSections(sectionInMap)
				break // Break immediately since we don't need to check the other sections
			}
		}
	}
}

// GetSection will attempt to get the SchemaKV associated with the provided section
func (schema *Schema) GetSection(section string) (kv *SchemaKV, getErr error) {
	var exists bool
	if kv, exists = schema.Map[section]; !exists { // If the section does not exist
		kv = nil
		getErr = ErrSectionDoesNotExist
	}

	return
}

// HasSection will check if our Schema has the provided section
func (schema *Schema) HasSection(section string) (exists bool) {
	_, exists = schema.Map[section]
	return
}

// ImportIntoDconf will import this Schema into its path via dconf load
func (schema *Schema) ImportIntoDconf() (importErr error) {
	schemaContent := []byte(schema.String())

	if len(schemaContent) == 0 { // No content
		importErr = ErrNoContentProvided
		return
	}

	bytesReader := bytes.NewReader(schemaContent) // Create a bytes reader for the file contents

	var dconfPath string
	dconfPath, importErr = exec.LookPath("dconf") // Get the path to the dconf binary

	if importErr != nil { // Failed to look up the path
		importErr = ErrNoDconfInPath
	}

	dconfLoad := exec.Cmd{
		Path: dconfPath,
		Args: []string{
			"load",
			"/",
		},
		Stdin: bytesReader,
	}

	if importErr = dconfLoad.Start(); importErr != nil { // Failed to run the command
		importErr = fmt.Errorf("failed during execution of dconf load: %s", importErr)
		return
	}

	dconfLoad.Wait()
	return
}

// MigrateSectionsWithName will migrate sections with the source prefix specified, remapping them to have the destination prefix.
// If exact is set to true, we will only migrate the section if it is an exact match
func (schema *Schema) MigrateSectionsWithName(source string, dest string, exact bool) {
	source = TrimSectionSlashes(source)
	dest = TrimSectionSlashes(dest)

	for sectionKey, sectionKV := range schema.Map { // For each section in the map
		if !strings.HasPrefix(sectionKey, source) || // Doesn't have source in the name of the section
			(exact && sectionKey != source) { // exact specified and section key isn't the same as source
			continue // Skip it
		}

		newSection := strings.TrimPrefix(sectionKey, source) // Remove the source
		newSection = dest + newSection                       // Prepend the destination

		schema.AddSection(newSection, sectionKV.Duplicate()) // Add the new section
		schema.DeleteSections(sectionKey)                    // Delete the old section
	}
}

// String will convert our Schema back to a String
func (schema *Schema) String() (schemaString string) {
	lines := []string{}        // Set our lines that we'll use to ensure newlines and the like
	sort.Strings(schema.Order) // Sort our order

	for _, section := range schema.Order { // Use order so our sections are organized alphabetically
		kv := schema.Map[section] // Get our key/value

		if kv == nil || len(kv.Keys) == 0 { // No keys
			continue
		}

		sectionLabel := fmt.Sprintf("[%s]", section) // Ensure we re-add [ and ]

		lines = append(lines, sectionLabel) // Add our section

		sort.Strings(kv.Order) // Order our SchemaKV

		for _, orderedKey := range kv.Order { // For each of our ordered keys
			sT := kv.Keys[orderedKey]
			lines = append(lines, orderedKey+"="+sT.String()) // Add in alphabetical order
		}

		lines = append(lines, "", "") // Add explicit new line after all our key/values
	}

	schemaString = strings.Join(lines, "\n")                      // Join our strings, separated by newline
	schemaString = strings.ReplaceAll(schemaString, "\n\n", "\n") // Replace double newlines with just one

	return
}
