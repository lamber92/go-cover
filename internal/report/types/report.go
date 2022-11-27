package types

import (
	"go/token"
	"html"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lamber92/go-cover/internal/metadata"
)

// ProjectURL 是项目在 GitHub 上的站点。
const ProjectURL = "https://github.com/lamber92/go-cover"

const (
	hitPrefix  = "    "
	missPrefix = "MISS"
)

// ReportPackageList 是报表的包列表。
type ReportPackageList []ReportPackage

// ReportPackage 保存有关 Go 包的数据、它的函数和一些统计信息。
type ReportPackage struct {
	Pkg               *metadata.Package
	Functions         ReportFunctionList
	TotalStatements   int
	ReachedStatements int
}

// PercentageReached 计算包测试达到的语句的百分比。
func (rp *ReportPackage) PercentageReached() float64 {
	var rv float64
	if rp.TotalStatements > 0 {
		rv = float64(rp.ReachedStatements) / float64(rp.TotalStatements) * 100
	}
	return rv
}

// ReportFunction 是一个带有一些附加统计信息的元数据函数。
type ReportFunction struct {
	*metadata.Function
	StatementsReached int
}

// FunctionLine 保存代码行、它在源文件中的行号以及测试是否到达它。
type FunctionLine struct {
	Code       string
	LineNumber int
	Missed     bool
	NewCode    bool
}

// CoveragePercent 是函数的代码覆盖率百分比。如果函数没有语句，则返回 100。
func (f ReportFunction) CoveragePercent() float64 {
	reached := f.StatementsReached
	var stmtPercent float64 = 0
	if len(f.Statements) > 0 {
		stmtPercent = float64(reached) / float64(len(f.Statements)) * 100
	} else if len(f.Statements) == 0 {
		stmtPercent = 100
	}
	return stmtPercent
}

// ShortFileName 返回函数文件名的基本路径。为了方便在主题的HTML模板中使用而提供。
func (f ReportFunction) ShortFileName() string {
	return filepath.Base(f.File)
}

// Lines 返回有关所有函数代码行的信息。
func (f ReportFunction) Lines() []FunctionLine {
	type annotator struct {
		fset  *token.FileSet
		files map[string]*token.File
	}
	a := &annotator{}
	a.fset = token.NewFileSet()
	a.files = make(map[string]*token.File)

	// 加载行信息文件。可能矫枉过正，也许只是在这里计算偏移量的线。
	setContent := false
	file := a.files[f.File]
	if file == nil {
		info, err := os.Stat(f.File)
		if err != nil {
			panic(err)
		}
		file = a.fset.AddFile(f.File, a.fset.Base(), int(info.Size()))
		setContent = true
	}

	data, err := ioutil.ReadFile(f.File)
	if err != nil {
		panic(err)
	}

	if setContent {
		// 这会处理内容并记录行号信息。
		file.SetLinesForContent(data)
	}

	statements := f.Statements[:]
	lineno := file.Line(file.Pos(f.Start))
	lines := strings.Split(string(data)[f.Start:f.End], "\n")
	fls := make([]FunctionLine, len(lines))

	for i, line := range lines {
		lineno := lineno + i
		statementFound := false
		hit := false
		for j := 0; j < len(statements); j++ {
			start := file.Line(file.Pos(statements[j].Start))
			if start == lineno {
				statementFound = true
				if !hit && statements[j].Reached > 0 {
					hit = true
				}
				statements = append(statements[:j], statements[j+1:]...)
			}
		}
		hitmiss := hitPrefix
		newCode := false
		if statementFound && !hit {
			hitmiss = missPrefix
		}
		// 判断是否是代码
		if len(f.NewLineSet) > 0 {
			_, newCode = f.NewLineSet[lineno]
		}
		fls[i] = FunctionLine{
			Missed:     hitmiss == missPrefix,
			NewCode:    newCode,
			LineNumber: lineno,
			Code:       html.EscapeString(strings.Replace(line, "\t", "        ", -1)),
		}
	}
	return fls
}

// ReportFunctionList is a list of functions for a report.
type ReportFunctionList []ReportFunction

func (l ReportFunctionList) Len() int {
	return len(l)
}

// Less
// TODO: make sort method configurable?
func (l ReportFunctionList) Less(i, j int) bool {
	var left, right float64
	if len(l[i].Statements) > 0 {
		left = float64(l[i].StatementsReached) / float64(len(l[i].Statements))
	}
	if len(l[j].Statements) > 0 {
		right = float64(l[j].StatementsReached) / float64(len(l[j].Statements))
	}
	if left < right {
		return true
	}
	return left == right && len(l[i].Statements) < len(l[j].Statements)
}

func (l ReportFunctionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// BranchesInfo 提交信息
type BranchesInfo struct {
	TargetBranchName  string
	CurrentBranchName string
	StartHashID       string
	EndHashID         string
}
