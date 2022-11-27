package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"

	"github.com/lamber92/go-cover/internal/metadata"
)

// Packages 表示一组 Package 结构
type Packages []*metadata.Package

// AppendPackage 方法可用于将包覆盖率结果合并到集合中
func (ps *Packages) AppendPackage(p *metadata.Package) {
	i := sort.Search(len(*ps), func(i int) bool {
		return (*ps)[i].Name >= p.Name
	})
	if i < len(*ps) && (*ps)[i].Name == p.Name {
		(*ps)[i].Accumulate(p)
	} else {
		head := (*ps)[:i]
		tail := append([]*metadata.Package{p}, (*ps)[i:]...)
		*ps = append(head, tail...)
	}
}

// ReadPackages 获取文件名列表并将其内容解析为 Packages 对象
// 特殊文件名“-”可用于指示标准输入
// 忽略重复的文件名
func ReadPackages(filenames []string) (ps Packages, err error) {
	copy_ := make([]string, len(filenames))
	copy(copy_, filenames)
	filenames = copy_
	sort.Strings(filenames)

	// Eliminate duplicates.
	unique := []string{filenames[0]}
	if len(filenames) > 1 {
		for _, f := range filenames[1:] {
			if f != unique[len(unique)-1] {
				unique = append(unique, f)
			}
		}
	}

	// 打开文件
	var files []*os.File
	for _, f := range filenames {
		if f == "-" {
			files = append(files, os.Stdin)
		} else {
			file, err := os.Open(f)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			files = append(files, os.Stdin)
		}
	}

	// 解析文件，积累包。
	for _, file := range files {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		result := &struct{ Packages []*metadata.Package }{}
		err = json.Unmarshal(data, result)
		if err != nil {
			return nil, err
		}
		for _, p := range result.Packages {
			ps.AppendPackage(p)
		}
	}
	return ps, nil
}
