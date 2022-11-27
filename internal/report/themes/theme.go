package themes

import (
	"fmt"

	"github.com/lamber92/go-cover/internal/report/types"
)

var themes = []types.Beautifier{
	defaultTheme{},
	defaultDiffTheme{},
}

// currTheme 用于渲染的主题。
var (
	currTheme     types.Beautifier = defaultTheme{}
	currDiffTheme types.Beautifier = defaultDiffTheme{}
)

// List 返回所有可用的主题。
func List() []types.Beautifier {
	return themes
}

// Get 获取一个主题的名字。如果没有找到则返回 nil。
func Get(name string) types.Beautifier {
	for _, t := range themes {
		if t.Name() == name {
			return t
		}
	}
	return nil
}

// Use 采用将用于呈现的主题的名称。
// 返回未知主题的错误。
func Use(name string) error {
	p := Get(name)
	if p == nil {
		return fmt.Errorf("unknown theme %q", name)
	}
	currTheme = p
	return nil
}

// Current 返回用于呈现 HTML 的主题。
func Current() types.Beautifier {
	return currTheme
}

// CurrentDiff 返回用于呈现增量 HTML 的主题。
func CurrentDiff() types.Beautifier {
	return currDiffTheme
}
