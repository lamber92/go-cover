package utils

import (
	"encoding/json"
	"io"

	"github.com/lamber92/go-cover/internal/metadata"
)

func MarshalJson(w io.Writer, packages []*metadata.Package) error {
	return json.NewEncoder(w).Encode(struct{ Packages []*metadata.Package }{packages})
}

func UnmarshalJson(data []byte) (packages []*metadata.Package, err error) {
	result := &struct{ Packages []*metadata.Package }{}
	err = json.Unmarshal(data, result)
	if err == nil {
		packages = result.Packages
	}
	return
}
