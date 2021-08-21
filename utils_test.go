/* utils_test.go
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

// TestRemoveFromStringArr will test RemoveFromStringArr
func TestRemoveFromStringArr(t *testing.T) {
	list := []string{
		"abc",
		"def",
		"ghi",
	}

	newList := RemoveFromStringArr(list, "def")

	if len(newList) != 2 {
		t.Errorf("Failed to remove \"def\" from our list. Provided %v and got %v", list, newList)
	}
}

// TestTrimSectionSlashes will test TrimSectionSlashes
func TestTrimSectionSlashes(t *testing.T) {
	if TrimSectionSlashes("/com/solus-project/") != "com/solus-project" {
		t.Error("Failed to properly trim the prefixed and suffixed / from \"/com/solus-project/\"")
	}
}
