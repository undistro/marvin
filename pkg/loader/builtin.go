package loader

import (
	"io/fs"
	"log"

	"github.com/undistro/marvin/internal/builtins"
)

var Builtins ChecksMap

func init() {
	checks, _, walkFn := walkDir(builtins.EmbbedChecksFS.ReadFile, "builtin:")
	err := fs.WalkDir(builtins.EmbbedChecksFS, ".", walkFn)
	if err != nil {
		log.Fatal(err)
	}
	Builtins = checks
}
