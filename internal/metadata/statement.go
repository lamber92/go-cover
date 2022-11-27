package metadata

import (
	"fmt"
)

type Statement struct {
	// Start 是语句的起始偏移量。
	Start int `json:"Start,omitempty"`

	// End 是语句的结束偏移量。
	End int `json:"End,omitempty"`

	// StartLine 是函数的起始行号
	StartLine int `json:"StartLine,omitempty"`

	// EndLine 是函数的结束行号
	EndLine int `json:"EndLine,omitempty"`

	// Reached 是语句到被执行的次数。
	Reached int64 `json:"Reached,omitempty"`
}

// Accumulate 会将提供的 Statement 中的覆盖率信息累积到此 Statement 中。
func (s *Statement) Accumulate(s2 *Statement) error {
	if s.Start != s2.Start || s.End != s2.End {
		return fmt.Errorf("source ranges do not match: %d-%d != %d-%d", s.Start, s.End, s2.Start, s2.End)
	}
	s.Reached += s2.Reached
	return nil
}
