package metadata

import "fmt"

type Function struct {
	// Name 是函数的名称。
	// 如果函数有接收器，名称将采用 T.N 形式，其中 T 是类型，N 是名称。
	Name string `json:"Name,omitempty"`

	// File 是定义函数的文件的完整路径。
	File string `json:"File,omitempty"`

	// Start 是函数签名的起始偏移量。
	Start int `json:"Start,omitempty"`

	// End 是函数的结束偏移量。
	End int `json:"End,omitempty"`

	// StartLine 是函数的起始行号
	StartLine int `json:"StartLine,omitempty"`

	// EndLine 是函数的结束行号
	EndLine int `json:"EndLine,omitempty"`

	// Statements 是指使用此函数注册的语句。
	Statements []*Statement `json:"Statements,omitempty"`

	// NewLineSet 新代码行号集合。用于增量覆盖率。
	NewLineSet map[int]struct{} `json:"NewLineSet,omitempty"`
}

// Accumulate 会将提供的 Function 的覆盖率信息累积到此 Function 中。
func (f *Function) Accumulate(f2 *Function) error {
	if f.Name != f2.Name {
		return fmt.Errorf("names do not match: %q != %q", f.Name, f2.Name)
	}
	if f.File != f2.File {
		return fmt.Errorf("files do not match: %q != %q", f.File, f2.File)
	}
	if f.Start != f2.Start || f.End != f2.End {
		return fmt.Errorf("source ranges do not match: %d-%d != %d-%d", f.Start, f.End, f2.Start, f2.End)
	}
	if len(f.Statements) != len(f2.Statements) {
		return fmt.Errorf("number of statements do not match: %d != %d", len(f.Statements), len(f2.Statements))
	}
	for i, s := range f.Statements {
		err := s.Accumulate(f2.Statements[i])
		if err != nil {
			return err
		}
	}
	return nil
}
