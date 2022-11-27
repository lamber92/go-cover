package convert

import (
	"go/build"
	"os"

	"github.com/lamber92/go-cover/internal/metadata"
	"github.com/lamber92/go-cover/internal/utils"
	"golang.org/x/tools/cover"
)

type packagesCache map[string]*build.Package

func Convert(filename string) error {
	var (
		ps       utils.Packages
		packages = make(packagesCache)
	)

	converter := converter{
		packages: make(map[string]*metadata.Package),
	}
	profiles, err := cover.ParseProfiles(filename)
	if err != nil {
		return err
	}
	for _, p := range profiles {
		if err := converter.convertProfile(packages, p); err != nil {
			return err
		}
	}

	for _, pkg := range converter.packages {
		ps.AppendPackage(pkg)
	}

	if err := utils.MarshalJson(os.Stdout, ps); err != nil {
		return err
	}
	return nil
}
