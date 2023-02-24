package report

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type statusHolder struct {
	Status CheckStatus `json:"status"`
}

func TestCheckStatusMarshalJSON(t *testing.T) {
	tests := []struct {
		input CheckStatus
		want  string
	}{
		{
			input: StatusError,
			want:  `{"status":"Error"}`,
		},
		{
			input: StatusFailed,
			want:  `{"status":"Failed"}`,
		},
		{
			input: StatusPassed,
			want:  `{"status":"Passed"}`,
		},
		{
			input: StatusSkipped,
			want:  `{"status":"Skipped"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			got, err := json.Marshal(&statusHolder{tt.input})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}
