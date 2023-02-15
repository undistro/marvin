package loader

import (
	"io/fs"
	"log"

	"github.com/undistro/marvin/internal/builtins"
)

var Builtins ChecksMap

func init() {
	check, _, walkFn := walkDir(builtins.EmbbedChecksFS.ReadFile)
	err := fs.WalkDir(builtins.EmbbedChecksFS, ".", walkFn)
	if err != nil {
		log.Fatal(err)
	}
	Builtins = check
}
