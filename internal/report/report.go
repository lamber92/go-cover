package report

import (
	"fmt"
	"os"

	"github.com/lamber92/go-cover/internal/metadata"
	"github.com/lamber92/go-cover/internal/report/types"
	"github.com/lamber92/go-cover/internal/utils"
)

type GenerateHTMLParam struct {
	Packages     utils.Packages
	CSS          string
	Dir          string
	FileName     string
	BranchesInfo *metadata.BranchesInfo
}

// GenerateHTML 通过解析 go-convert/metadata 数据输出 HTML 报告。
// css 参数是自定义样式表的绝对路径。使用空字符串以使用可用的默认样式表。
func GenerateHTML(param *GenerateHTMLParam) error {
	// Custom stylesheet?
	stylesheet := ""
	if param.CSS != "" {
		if _, err := exists(param.CSS); err != nil {
			return fmt.Errorf("stylesheet(css) is not exists. err: %v", err)
		}
		stylesheet = param.CSS
	}

	reporter := newReport(param.Packages,
		stylesheet,
		&types.BranchesInfo{
			TargetBranchName:  param.BranchesInfo.TargetBranchName,
			CurrentBranchName: param.BranchesInfo.CurrentBranchName,
			StartHashID:       param.BranchesInfo.StartHashID,
			EndHashID:         param.BranchesInfo.EndHashID,
		})
	file, err := utils.CreateFile(param.Dir, param.FileName)
	if err != nil {
		return err
	}

	if param.FileName == utils.DiffHTML {
		if err = writeDiffReport(file, reporter); err != nil {
			return fmt.Errorf("generate HTML diff-report failed. err: %v", err)
		}
	} else {
		if err = writeFullReport(file, reporter); err != nil {
			return fmt.Errorf("generate HTML full-report failed. err: %v", err)
		}
	}

	return nil
}

type report struct {
	packages   utils.Packages
	stylesheet string // absolute path to CSS
	commit     *types.BranchesInfo
}

// newReport 创建一个新报表。
func newReport(ps utils.Packages, stylesheet string, commit *types.BranchesInfo) (r *report) {
	r = &report{
		packages:   ps,
		stylesheet: stylesheet,
		commit:     commit,
	}
	return
}

// clear 从报告中清除覆盖率信息。
func (r *report) clear() {
	r.packages = nil
}

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		return false, err
	}
	return true, nil
}
