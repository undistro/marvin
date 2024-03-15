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

package loader

import (
	"io/fs"
	"log"

	"github.com/undistro/marvin/internal/builtins"
	"github.com/undistro/marvin/pkg/types"
)

// Builtins represents the builtins checks
var Builtins []types.Check

func init() {
	c, _, walkFn := walkDir(builtins.EmbedChecksFS.ReadFile, true)
	err := fs.WalkDir(builtins.EmbedChecksFS, ".", walkFn)
	if err != nil {
		log.Fatal(err)
	}
	Builtins = c.toList()
}
