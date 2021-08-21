/* schemaType.go
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
	"strconv"
	"strings"
)

// NewSchemaType will attempt to convert the provided key/val into a SchemaType
func NewSchemaType(rawVal string) (sT *SchemaType, parseErr error) {
	sT = &SchemaType{Val: rawVal} // Ensure we always set the raw value

	if (rawVal == "false") || (rawVal == "true") { // Is a boolean
		sT.Type = "bool" // Define as a boolean

		if rawVal == "false" { // Is false
			sT.BoolVal = false
		} else { // Is true
			sT.BoolVal = true
		}
	} else if strings.HasPrefix(rawVal, "uint32") { // If the string starts with uint32
		rawVal = strings.Replace(rawVal, "uint32", "", -1)
		rawVal = strings.TrimSpace(rawVal) // Remove any whitespace

		var i uint64
		i, parseErr = strconv.ParseUint(rawVal, 10, 32)

		if parseErr == nil { // Didn't fail to parse our string as a uint
			sT.Type = "uint32"
			sT.UintVal = uint32(i) // Convert the uint64 to uint32 and set
		} else { // Failed to parse
			return
		}
	} else if floaty, floatParseErr := strconv.ParseFloat(rawVal, 64); floatParseErr == nil { // Attempt to convert to a float64 before int
		sT.FloatVal = floaty
		sT.FloatHadTrailingZero = strings.HasSuffix(rawVal, ".0") // This is useful for double->Go float64 and Go float64->double conversion
		sT.Type = "float64"
	} else if inty, intParseErr := strconv.ParseInt(rawVal, 10, 32); intParseErr == nil { // Attempt to convert to an int64 that is convertable to an int32
		sT.IntVal = int32(inty)
		sT.Type = "int32"
	} else { // Treat as a string, Val already set
		sT.Type = "string" // Define as a string
	}

	return
}

// Duplicate will duplicate this SchemaType
func (sT *SchemaType) Duplicate() *SchemaType {
	newSt := SchemaType{
		Type:                 sT.Type,
		BoolVal:              sT.BoolVal,
		FloatHadTrailingZero: sT.FloatHadTrailingZero,
		FloatVal:             sT.FloatVal,
		IntVal:               sT.IntVal,
		UintVal:              sT.UintVal,
		Val:                  sT.Val,
	}

	return &newSt
}

// Matches will check if the provided SchemaType matches this one
func (sT *SchemaType) Matches(oST *SchemaType) (matches bool) {
	if sT.Type != oST.Type { // Types don't match
		return
	}

	switch sT.Type {
	case "bool":
		return sT.BoolVal == oST.BoolVal
	case "uint32":
		return sT.UintVal == oST.UintVal
	case "int32":
		return sT.IntVal == oST.IntVal
	case "float64":
		return sT.FloatVal == oST.FloatVal
	default:
		return sT.Val == oST.Val
	}
}

// String will convert our SchemaType back to a string
// Note this only converts the value itself and not the key
func (sT *SchemaType) String() string {
	switch sT.Type {
	case "bool":
		return strconv.FormatBool(sT.BoolVal)
	case "uint32":
		return "uint32 " + strconv.FormatUint(uint64(sT.UintVal), 10)
	case "int32":
		return strconv.FormatInt(int64(sT.IntVal), 10)
	case "float64":
		floatString := strconv.FormatFloat(sT.FloatVal, 'G', -1, 64) // Use G for max digits, no trailing zeroes

		if !strings.Contains(floatString, ".") && sT.FloatHadTrailingZero { // Has no decimal and had one when we created the type
			floatString += ".0" // Add the .0 back
		}

		return floatString
	default: // Fall back (string, array string, etc)
		return sT.Val // Add our value directly
	}
}
