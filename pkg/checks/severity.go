// Copyright 2023 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package checks

import (
	"fmt"
	"strings"
)

type Severity int

const (
	SeverityUnknown Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityLow:
		return "Low"
	case SeverityMedium:
		return "Medium"
	case SeverityHigh:
		return "High"
	case SeverityCritical:
		return "Critical"
	default:
		return ""
	}
}

func ParseSeverity(s string) Severity {
	switch s {
	case "Low":
		return SeverityLow
	case "Medium":
		return SeverityMedium
	case "High":
		return SeverityHigh
	case "Critical":
		return SeverityCritical
	default:
		return SeverityUnknown
	}
}

func (s *Severity) UnmarshalJSON(b []byte) error {
	*s = ParseSeverity(strings.Trim(string(b), "\""))
	return nil
}

func (s Severity) MarshalJSON() ([]byte, error) {
	v := fmt.Sprintf(`"%s"`, s.String())
	return []byte(v), nil
}
