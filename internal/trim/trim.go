package trim

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/lamber92/go-cover/internal/metadata"
	"github.com/lamber92/go-cover/internal/utils"
)

// Do 按传入的保留规则对 覆盖率信息 进行修剪
func Do(sourcePath string, reserveRulesPath string) (out utils.Packages, err error) {
	file, err := os.OpenFile(sourcePath, os.O_RDONLY, 0666)
	if err != nil {
		err = fmt.Errorf("failed to open coverage.json. path: %s, err: %v", sourcePath, err)
		return
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("failed to read coverage.json. err: %v", err)
		return
	}
	packages, err := utils.UnmarshalJson(buffer)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal coverage.json")
	}

	info, err := utils.LoadReservedInfo(reserveRulesPath)
	if err != nil {
		return
	}
	out, err = TrimPackages(packages, info.Rules)
	return
}

// TrimPackages 按传入的保留规则对 覆盖率信息 进行修剪
func TrimPackages(packages utils.Packages, rules metadata.ReservedRules) (out utils.Packages, err error) {
	out, err = trimPackages(packages, rules)
	return
}

// trimPackages 裁剪包
func trimPackages(source utils.Packages, reserveRules metadata.ReservedRules) (out utils.Packages, err error) {
	prefix, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("failed to get pwd. err: %v", err)
		return
	}
	prefix = prefix + string(filepath.Separator)

	out = make(utils.Packages, 0)
	// 规则：只保保留规则中的文件及行号相关代码
	for _, pkg := range source {
		newPkg := &metadata.Package{
			Name:      pkg.Name,
			Functions: make([]*metadata.Function, 0),
		}

		for _, function := range pkg.Functions {
			newFunction := &metadata.Function{
				Name:       function.Name,
				File:       function.File,
				Start:      function.Start,
				End:        function.End,
				Statements: make([]*metadata.Statement, 0),
				NewLineSet: make(map[int]struct{}),
			}

			// 遍历方法需要保留的条件
			for file, rule := range reserveRules {
				// 是否包含该文件名
				if !strings.Contains(function.File, utils.FixPathSeparator(prefix+file)) {
					continue
				}
				// 是否包含该行(求两个闭区间是否有交集)
				if !judgeTwoAreaOverlap(function.StartLine, function.EndLine, rule.StartLine, rule.EndLine) {
					continue
				}

				// 遍历语句需要保留的条件
				for _, stmt := range function.Statements {
					// 累加每个Statements对象中的行号查看是否命中需要保留的行号
					hitLines := make([]int, 0) // 用于记录命中的行号组
					for i := stmt.StartLine; i <= stmt.EndLine; i++ {
						if _, exist := rule.LinesSet[i]; exist {
							hitLines = append(hitLines, i)
						}
					}

					// 如果有命中的行号，记录到保留组
					if len(hitLines) > 0 {
						newFunction.Statements = append(newFunction.Statements, &metadata.Statement{
							Start:     stmt.Start,
							End:       stmt.End,
							StartLine: stmt.StartLine,
							EndLine:   stmt.EndLine,
							Reached:   stmt.Reached,
						})
						// 记录新行号
						for _, v := range hitLines {
							newFunction.NewLineSet[v] = struct{}{}
						}
					}
				}
			}

			if len(newFunction.Statements) > 0 {
				newPkg.Functions = append(newPkg.Functions, newFunction)
			}
		}

		if len(newPkg.Functions) > 0 {
			out = append(out, newPkg)
		}
	}

	return
}

// judgeTwoAreaOverlap 判断两个区间是否重叠(左闭右闭区间)
func judgeTwoAreaOverlap(start1, end1, start2, end2 int) bool {
	return start2 <= end1 && end2 >= start1
}
