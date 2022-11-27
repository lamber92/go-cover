package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	Separator  = string(filepath.Separator)
	TimeFormat = "2006_01_02_15_04_05"
	FullHTML   = "full.html"
	DiffHTML   = "diff.html"
)

// CreateFile 创建文件
func CreateFile(dir, fileName string) (file *os.File, err error) {
	path := dir + Separator + fileName
	if exists(path) {
		// 这里直接覆写原文件
		if file, err = os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0666); err != nil {
			err = fmt.Errorf("fail to open file. err: %v", err)
			return
		}
	} else {
		if err = os.MkdirAll(dir, 0766); err != nil {
			return
		}
		if file, err = os.Create(path); err != nil {
			err = fmt.Errorf("fail to create file. err: %v", err)
			return
		}
	}
	return
}

// exists 判断路径是否存在
func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// FixPathSeparator 修复不同平台上文件路径/符差异的问题
func FixPathSeparator(source string) string {
	if runtime.GOOS == "windows" {
		source = strings.ReplaceAll(source, "/", Separator)
	} else {
		source = strings.ReplaceAll(source, "\\", Separator)
	}
	return source
}
