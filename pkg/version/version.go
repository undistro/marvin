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

package version

var (
	version = "dev"
	major   = ""
	minor   = ""
	commit  = ""
	date    = ""
)

type Info struct {
	Version string `json:"version"`
	Major   string `json:"major"`
	Minor   string `json:"minor"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func (i Info) String() string {
	return i.Version
}

func Get() Info {
	return Info{
		Version: version,
		Major:   major,
		Minor:   minor,
		Commit:  commit,
		Date:    date,
	}
}
