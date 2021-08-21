/* errors.go
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

// This file is our centralized location for defining errors

import (
	"errors"
)

var (
	// ErrKeyAlreadyExists is an error we return when we already have a key in a schema key-value store. Mostly useful for validating during section adding.
	ErrKeyAlreadyExists = errors.New("key already exists in schemakv")

	// ErrKeyNotExists is an error we return when we do not have a key in a schema
	ErrKeyNotExists = errors.New("key does not exist")

	// ErrModCannotDoReplace is an error we return when we cannot do a replacement of a value
	ErrModCannotDoReplace = errors.New("cannot perform replace modification, schematype is not of ArrayAsString or string")

	// ErrModNoReplaceValueOrValue is an error we return if we cannot perform a modification without a value
	ErrModNoReplaceValueOrValue = errors.New("cannot perform modification, no replacevalue or value specified")

	// ErrNoContentProvided is an error we return when no content is provided when attempt to import
	ErrNoContentProvided error = errors.New("no content provided as a byte slice")

	// ErrNoDconfInPath is an error we return if we could not find dconf in the path during a dconf operation
	ErrNoDconfInPath error = errors.New("no dconf found in path")

	// ErrSectionDoesNotExist is an error we return if a section requested does not exist
	ErrSectionDoesNotExist = errors.New("section does not exist")

	// ErrSectionExists is an error we return if a section exists
	ErrSectionExists = errors.New("section exists")
)
