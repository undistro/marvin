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
