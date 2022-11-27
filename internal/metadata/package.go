package metadata

import "fmt"

type Package struct {
	// 名称是包的规范路径。
	Name string `json:"Name,omitempty"`

	// Functions 是使用此包注册的函数列表。
	Functions []*Function `json:"Functions,omitempty"`
}

// Accumulate 会将提供的 Package 中的覆盖率信息累积到此 Package 中。
func (p *Package) Accumulate(p2 *Package) error {
	if p.Name != p2.Name {
		return fmt.Errorf("names do not match: %q != %q", p.Name, p2.Name)
	}
	if len(p.Functions) != len(p2.Functions) {
		return fmt.Errorf("function counts do not match: %d != %d", len(p.Functions), len(p2.Functions))
	}
	for i, f := range p.Functions {
		err := f.Accumulate(p2.Functions[i])
		if err != nil {
			return err
		}
	}
	return nil
}
