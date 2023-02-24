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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type severityHolder struct {
	Severity Severity `json:"severity"`
}

func TestSeverityMarshalJSON(t *testing.T) {
	tests := []struct {
		input Severity
		want  string
	}{
		{
			input: SeverityLow,
			want:  `{"severity":"Low"}`,
		},
		{
			input: SeverityMedium,
			want:  `{"severity":"Medium"}`,
		},
		{
			input: SeverityHigh,
			want:  `{"severity":"High"}`,
		},
		{
			input: SeverityCritical,
			want:  `{"severity":"Critical"}`,
		},
		{
			input: SeverityCritical,
			want:  `{"severity":"Critical"}`,
		},
		{
			input: SeverityUnknown,
			want:  `{"severity":""}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			input := severityHolder{tt.input}
			got, err := json.Marshal(&input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestSeverityUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input string
		want  Severity
	}{
		{
			input: `{"severity": "Low"}`,
			want:  SeverityLow,
		},
		{
			input: `{"severity": "Medium"}`,
			want:  SeverityMedium,
		},
		{
			input: `{"severity": "High"}`,
			want:  SeverityHigh,
		},
		{
			input: `{"severity": "Critical"}`,
			want:  SeverityCritical,
		},
		{
			input: `{"severity": "Critical"}`,
			want:  SeverityCritical,
		},
		{
			input: `{"severity": "foo"}`,
			want:  SeverityUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			var got severityHolder
			err := json.Unmarshal([]byte(tt.input), &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got.Severity)
		})
	}
}
