package convert

import (
	"go/build"

	"github.com/lamber92/go-cover/internal/metadata"
	"github.com/lamber92/go-cover/internal/utils"
	"golang.org/x/tools/cover"
)

type packagesCache map[string]*build.Package

func Do(filename string) (ps utils.Packages, err error) {
	profiles, err := cover.ParseProfiles(filename)
	if err != nil {
		return
	}
	var (
		packages = make(packagesCache)
		conv     = converter{packages: make(map[string]*metadata.Package)}
	)
	for _, p := range profiles {
		if err = conv.convertProfile(packages, p); err != nil {
			return
		}
	}
	for _, pkg := range conv.packages {
		ps.AppendPackage(pkg)
	}
	return
}
