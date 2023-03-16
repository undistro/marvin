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

package types

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
