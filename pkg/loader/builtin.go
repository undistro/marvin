package loader

import (
	"io/fs"
	"log"

	"github.com/undistro/marvin/internal/builtins"
	"github.com/undistro/marvin/pkg/checks"
)

var Builtins []checks.Check

func init() {
	c, _, walkFn := walkDir(builtins.EmbbedChecksFS.ReadFile, true)
	err := fs.WalkDir(builtins.EmbbedChecksFS, ".", walkFn)
	if err != nil {
		log.Fatal(err)
	}
	Builtins = c.toList()
}
