package types

import "text/template"

// Beautifier 定义了一个用于呈现 HTML 覆盖率统计信息的主题。
type Beautifier interface {
	// Name is the name of the theme.
	Name() string
	// Description is a single line comment about the theme.
	Description() string
	// Template is the content that will be rendered.
	Template() *template.Template
	Data() *TemplateData
}

// TemplateData 具有用于渲染的 HTML 模板所需的所有字段。
type TemplateData struct {
	// CSS is the stylesheet content that will be embedded in the HTML page.
	CSS string
	// When is the date time of report generation.
	When string
	// Overview holds data used for an additional header in case of multiple Go packages
	// have been analysed. Can be used for a high level summary. Is nil if the report has
	// only one package.
	Overview *ReportPackage
	// Packages is the list of all Go packages analysed.
	Packages ReportPackageList
	// ProjectURL is the project's site on GitHub.
	ProjectURL string
	// BranchesInfo
	BranchesInfo *BranchesInfo //
}
