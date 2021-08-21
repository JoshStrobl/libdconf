/* main_test.go
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
	"fmt"
	"os"
	"testing"
)

var TestSchema *Schema
var TestSchemaKV *SchemaKV

// Basically for setting global variables
func Bootstrap() {
	content, exampleReadErr := os.ReadFile("examples/com__solus-project__budgie-panel")

	if exampleReadErr != nil { // Failed to read our required example content
		fmt.Printf("Failed to read required example content: %s\n", exampleReadErr)
		os.Exit(1)
	}

	TestSchema, _ = NewSchema("/com/solus-project/budgie-panel/", content) // Attempt to read our content
}

// BootstrapSetKV will just bootstrap the global key value
func BootstrapSetKV(sectionID string) {
	TestSchemaKV, _ = TestSchema.GetSection(sectionID) // Set our testing schema
}

func TestMain(m *testing.M) {
	Bootstrap()
	os.Exit(m.Run())
}
