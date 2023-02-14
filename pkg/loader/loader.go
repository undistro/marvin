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

var supportedExt = map[string]bool{
	".yaml": true,
	".yml":  true,
	".json": true,
}

func LoadChecks(root string) (map[string]checks.Check, error) {
	c, _, err := load(root)
	return c, err
}

func LoadChecksAndTests(root string) (map[string]checks.Check, map[string][]checks.Test, error) {
	return load(root)
}

func load(root string) (map[string]checks.Check, map[string][]checks.Test, error) {
	tests := make(map[string][]checks.Test)
	check := make(map[string]checks.Check)

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
		check[k] = c
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return check, tests, nil
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
