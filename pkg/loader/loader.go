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
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/undistro/marvin/pkg/checks"
)

type (
	ChecksMap    map[string]checks.Check
	TestsMap     map[string][]checks.Test
	readFileFunc func(string) ([]byte, error)
)

func (cm ChecksMap) toList() []checks.Check {
	if cm == nil {
		return nil
	}
	list := make([]checks.Check, 0, len(cm))
	for _, c := range cm {
		list = append(list, c)
	}
	return list
}

var supportedExt = map[string]bool{
	".yaml": true,
	".yml":  true,
	".json": true,
}

func LoadChecks(root string) ([]checks.Check, error) {
	c, _, err := load(root)
	return c.toList(), err
}

func LoadChecksAndTests(root string) (ChecksMap, TestsMap, error) {
	return load(root)
}

func load(root string) (ChecksMap, TestsMap, error) {
	check, tests, walkFn := walkDir(os.ReadFile, false)
	err := filepath.WalkDir(root, walkFn)
	if err != nil {
		return nil, nil, err
	}
	return check, tests, nil
}

func walkDir(readFileFn readFileFunc, builtin bool) (ChecksMap, TestsMap, fs.WalkDirFunc) {
	tests := make(TestsMap)
	check := make(ChecksMap)
	return check, tests, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		ext := filepath.Ext(path)
		if !supportedExt[ext] {
			return nil // unsupported file
		}
		bs, err := readFileFn(path)
		if err != nil {
			return err
		}

		testSuffix := "_test" + ext
		isTest := strings.HasSuffix(path, testSuffix)
		if isTest {
			t, err := parseTests(ext, bs)
			if err != nil {
				return err
			}
			k := strings.TrimSuffix(path, testSuffix)
			tests[k] = t
			return nil
		}
		c, err := parseCheck(ext, bs)
		if err != nil {
			return err
		}
		c.Builtin = builtin
		c.Path = path
		k := strings.TrimSuffix(path, ext)
		if builtin {
			k = "builtin:" + k
		}
		check[k] = c
		return nil
	}
}

func parseCheck(ext string, bs []byte) (checks.Check, error) {
	obj := checks.Check{}
	return parse(ext, bs, obj)
}

func parseTests(ext string, bs []byte) ([]checks.Test, error) {
	var obj []checks.Test
	return parse(ext, bs, obj)
}

func parse[T any](ext string, bs []byte, obj T) (T, error) {
	var err error
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(bs, &obj)
	case ".json":
		err = json.Unmarshal(bs, &obj)
	default:
		return obj, fmt.Errorf("unsupported file extension: %s", ext)
	}
	return obj, err
}
