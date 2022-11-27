package metadata

import (
	"fmt"
	"strings"
)

type BranchesInfo struct {
	TargetBranchName  string
	CurrentBranchName string
	StartHashID       string
	EndHashID         string
}

// FormatBranchesInfo 格式化分支信息
func (info *BranchesInfo) FormatBranchesInfo() string {
	return fmt.Sprintf("%s,%s:%s,%s", info.TargetBranchName, info.CurrentBranchName, info.StartHashID, info.EndHashID)
}

// ParseBranchesInfo 解析分支信息
func ParseBranchesInfo(source string) (*BranchesInfo, error) {
	tmp := strings.Split(source, ":")
	if len(tmp) != 2 {
		return nil, fmt.Errorf("invalid commit_ids info. [%s]", source)
	}
	s1 := strings.Split(tmp[0], ",")
	if len(s1) != 2 {
		return nil, fmt.Errorf("invalid branches info. [%s]", tmp[0])
	}
	s2 := strings.Split(tmp[1], ",")
	if len(s2) != 2 {
		return nil, fmt.Errorf("invalid commit_ids info. [%s]", tmp[1])
	}
	return &BranchesInfo{
		TargetBranchName:  s1[0],
		CurrentBranchName: s1[1],
		StartHashID:       s2[0],
		EndHashID:         s2[1],
	}, nil
}
