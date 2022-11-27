package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/lamber92/go-cover/internal/metadata"
	"github.com/lamber92/go-cover/internal/report/themes"
	"github.com/lamber92/go-cover/internal/report/types"
)

// writeFullReport 写全量报告
func writeFullReport(w io.Writer, r *report) error {
	return (&basicWriter{theme: themes.Current()}).Do(w, r)
}

// writeDiffReport 写增量报告
func writeDiffReport(w io.Writer, r *report) error {
	return (&basicWriter{theme: themes.CurrentDiff()}).Do(w, r)
}

type basicWriter struct {
	theme types.Beautifier
}

// Do 向给定的编写器打印一份覆盖率报告。
func (b *basicWriter) Do(w io.Writer, r *report) error {
	theme := b.theme
	data := theme.Data()

	css := data.CSS
	if len(r.stylesheet) > 0 {
		// Inline CSS.
		f, err := os.Open(r.stylesheet)
		if err != nil {
			return fmt.Errorf("print report. err: %v", err)
		}
		style, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("read style. err: %v", err)
		}
		css = string(style)
	}
	reportPackages := make(types.ReportPackageList, len(r.packages))
	for i, pkg := range r.packages {
		reportPackages[i] = b.buildReportPackage(pkg)
	}

	data.CSS = css
	data.Packages = reportPackages
	data.BranchesInfo = r.commit

	if len(reportPackages) > 1 {
		rv := types.ReportPackage{
			Pkg: &metadata.Package{Name: "Report Total"},
		}
		for _, rp := range reportPackages {
			rv.ReachedStatements += rp.ReachedStatements
			rv.TotalStatements += rp.TotalStatements
		}
		data.Overview = &rv
	}
	if err := theme.Template().Execute(w, data); err != nil {
		return fmt.Errorf("execute template. err: %v", err)
	}
	return nil
}

type reverse struct {
	sort.Interface
}

func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func (b *basicWriter) buildReportPackage(pkg *metadata.Package) types.ReportPackage {
	rv := types.ReportPackage{
		Pkg:       pkg,
		Functions: make(types.ReportFunctionList, len(pkg.Functions)),
	}
	for i, fn := range pkg.Functions {
		reached := 0
		for _, stmt := range fn.Statements {
			if stmt.Reached > 0 {
				reached++
			}
		}
		rv.Functions[i] = types.ReportFunction{Function: fn, StatementsReached: reached}
		rv.TotalStatements += len(fn.Statements)
		rv.ReachedStatements += reached
	}
	sort.Sort(reverse{rv.Functions})
	return rv
}
