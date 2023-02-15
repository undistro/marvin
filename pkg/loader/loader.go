package loader

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/undistro/marvin/pkg/scan"
)

var supportedExt = map[string]bool{
	".yaml": true,
	".yml":  true,
	".json": true,
}

func LoadChecks(root string) (map[string]scan.Check, error) {
	checks, _, err := load(root)
	return checks, err
}

func LoadChecksAndTests(root string) (map[string]scan.Check, map[string][]scan.Test, error) {
	return load(root)
}

func load(root string) (map[string]scan.Check, map[string][]scan.Test, error) {
	tests := make(map[string][]scan.Test)
	checks := make(map[string]scan.Check)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		ext := filepath.Ext(path)
		if !supportedExt[ext] {
			return nil // unsupported file
		}
		bs, err := os.ReadFile(path)
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
		k := strings.TrimSuffix(path, ext)
		checks[k] = c
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return checks, tests, nil
}

func parseCheck(ext string, bs []byte) (scan.Check, error) {
	obj := scan.Check{}
	return parse(ext, bs, obj)
}

func parseTests(ext string, bs []byte) ([]scan.Test, error) {
	var obj []scan.Test
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
