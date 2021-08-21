/* utils.go
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
	"os/exec"
	"strings"
)

// DconfDump will attempt to dump the contents of the provided path
// If no path is provided, / (root) is used
func DconfDump(path string) (content []byte, dumpErr error) {
	if path == "" {
		path = "/"
	}

	if _, dumpErr = exec.LookPath("dconf"); dumpErr != nil { // Failed to look up dconf
		dumpErr = ErrNoDconfInPath
		return
	}

	dconfCmd := exec.Command("dconf", "dump", path)
	content, dumpErr = dconfCmd.Output() // Run and output its stdout to dconfOutput
	return
}

// RemoveFromStringArr will remove the specified string from our array
func RemoveFromStringArr(arr []string, removeString string) []string {
	newList := []string{} // Create a new array of items to retain
	arrLen := len(arr)

	if arrLen != 0 {
		for index := 0; index < arrLen; index++ { // For each item in this array
			item := arr[index]

			if item != removeString { // If this item does not match our removeString
				newList = append(newList, item) // Append
			}
		}
	}

	return newList
}

// TrimSectionSlashes will trim the provided section string of any prefixed and suffixed forward slash
func TrimSectionSlashes(section string) string {
	s := strings.TrimPrefix(section, "/") // Ensure our source doesn't start with /
	return strings.TrimSuffix(s, "/")     // Ensure our source doesn't end with /
}
