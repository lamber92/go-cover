package metadata

type Rule struct {
	StartLine int              // 需要保留的起始行号
	EndLine   int              // 需要保留的结束行号
	LinesSet  map[int]struct{} // 需要保留的行号集合
}

type ReservedRules map[string]*Rule

func (m *ReservedRules) Add(k string, v *Rule) {
	(*m)[k] = v
}

func (m *ReservedRules) Delete(k string) {
	delete(*m, k)
}

type ReservedInfo struct {
	Branches *BranchesInfo // 分支信息
	Rules    ReservedRules // map[go文件路径]*Rule
}
