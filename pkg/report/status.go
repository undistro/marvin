package report

import "fmt"

type CheckStatus int

const (
	StatusPassed CheckStatus = iota
	StatusSkipped
	StatusFailed
	StatusError
)

func (s CheckStatus) String() string {
	switch s {
	case StatusPassed:
		return "Passed"
	case StatusSkipped:
		return "Skipped"
	case StatusFailed:
		return "Failed"
	case StatusError:
		return "Error"
	default:
		return ""
	}
}

func (s CheckStatus) MarshalJSON() ([]byte, error) {
	v := fmt.Sprintf(`"%s"`, s.String())
	return []byte(v), nil
}
