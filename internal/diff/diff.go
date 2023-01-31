package diff

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"

	"github.com/lamber92/go-cover/internal/metadata"
	"github.com/lamber92/go-cover/internal/utils"
	"github.com/spf13/cast"
)

type diff struct {
	CurrentBranch     string
	TargetBranch      string
	CommitHashIdRange [2]string

	commitHashIdSet   map[string]struct{}
	filePathM2LineNos map[string]map[string]struct{}
}

func NewDiffManager(currentBranch, targetBranch string) (d *diff, err error) {
	d = &diff{
		CommitHashIdRange: [2]string{},
		commitHashIdSet:   make(map[string]struct{}),
		filePathM2LineNos: make(map[string]map[string]struct{}),
	}

	if len(currentBranch) == 0 {
		currentBranch, err = utils.GetCurrentBranch()
		if err != nil {
			return
		}
	}
	d.CurrentBranch = d.fixBranch(currentBranch)
	d.TargetBranch = d.fixBranch(targetBranch)
	return
}

// fixBranch 修复分支名，确保是远程分支名称
func (*diff) fixBranch(branch string) string {
	if strings.HasPrefix(branch, "origin/") {
		return branch
	}
	return "origin/" + branch
}

func Do(currentBranch, targetBranch string, hashIdsRange []string) (d *diff, err error) {
	if d, err = NewDiffManager(currentBranch, targetBranch); err != nil {
		return
	}
	// 判断是否有范围限制
	if len(hashIdsRange) > 0 {
		if err = d.listDiffCommitHashIdsWithLimit(hashIdsRange); err != nil {
			return
		}
	} else {
		if err = d.listDiffCommitHashIds(); err != nil {
			return
		}
	}
	if err = d.listCommitModifyFiles(); err != nil {
		return
	}
	if err = d.listCommitModifyLineNos(); err != nil {
		return
	}
	return
}

func (d *diff) listDiffCommitHashIds() error {
	cmd := exec.Command("git", "log", fmt.Sprintf("%s..%s", d.TargetBranch, d.CurrentBranch), "--oneline")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get commit hash information for differences between branches. err: %v", err)
	}

	commitHashIds := make([]string, 0)

	br := bufio.NewReader(bytes.NewBuffer(output))
	for {
		buff, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}

		tmp := bytes.Split(buff, []byte(" "))
		if len(tmp) > 0 {
			hashId := string(tmp[0])
			commitHashIds = append(commitHashIds, hashId)
			d.commitHashIdSet[hashId] = struct{}{}
		}
	}

	if len(commitHashIds) > 0 {
		d.CommitHashIdRange[0] = commitHashIds[0]
		d.CommitHashIdRange[1] = commitHashIds[len(commitHashIds)-1]
	}

	return nil
}

func (d *diff) listDiffCommitHashIdsWithLimit(hashIdsRange []string) error {
	cmd := exec.Command("git", "log", fmt.Sprintf("%s..%s", d.TargetBranch, d.CurrentBranch), "--oneline")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get commit hash information for differences between branches. err: %v", err)
	}

	var (
		tmpHashIds  = make([]string, 0)
		startHashId = d.cutHashId(hashIdsRange[0]) // 裁剪hash_id长度
		endHashId   = d.cutHashId(hashIdsRange[1]) // 裁剪hash_id长度
		hitStartId  = false                        // 表示起始id是否命中
		hitEndId    = false                        // 表示终止id是否命中
	)

	br := bufio.NewReader(bytes.NewBuffer(output))
	for {
		buff, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}

		tmp := bytes.Split(buff, []byte(" "))
		if len(tmp) > 0 {
			tmpHashIds = append(tmpHashIds, string(tmp[0]))
		}
	}

	// 过滤出真正需要的hash-id集合
	for _, hashId := range tmpHashIds {
		if hashId == startHashId {
			d.CommitHashIdRange[0] = hashId
			hitStartId = true
		}

		// 命中起始id后才开始记录
		if hitStartId {
			d.commitHashIdSet[hashId] = struct{}{}
		}

		// 命中终止id，直接退出
		if hashId == endHashId {
			d.CommitHashIdRange[1] = hashId
			hitEndId = true
			break
		}
	}

	// 检查起始id和终止id是否都被命中
	if !hitStartId || !hitEndId {
		return fmt.Errorf("hash-ids range is not in current branch commit. range-ids: [%+v], full-ids: [%+v]",
			hashIdsRange, tmpHashIds)
	}

	return nil
}

func (d *diff) listCommitModifyFiles() error {
	for hashId := range d.commitHashIdSet {
		cmd := exec.Command("git", "show", hashId, "--name-only")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to get modify files information from commit hash id. err: %v", err)
		}
		br := bufio.NewReader(bytes.NewBuffer(output))
		for {
			buff, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}

			line := string(buff)
			if len(line) == 0 ||
				line[0] == ' ' ||
				!strings.Contains(line, ".go") {
				continue
			}

			d.filePathM2LineNos[line] = make(map[string]struct{}, 0)
		}
	}

	return nil
}

func (d *diff) listCommitModifyLineNos() error {
	for path, lineNos := range d.filePathM2LineNos {
		// https://git-scm.com/docs/git-blame
		cmd := exec.Command("git", "blame", path, "-w", "-s", "--show-name")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to get modify codes information from file[%s]. err: %v", path, err)
		}

		br := bufio.NewReader(bytes.NewBuffer(output))
		for {
			buff, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			if len(buff) == 0 {
				continue
			}

			line := string(buff) // 格式: hash-id file-path   line) codes......
			index := strings.Index(line, ")")
			if index == -1 {
				continue
			}
			line = line[:index] // 格式: hash-id file-path   line
			tmp := strings.Split(line, " ")
			hashId := d.cutHashId(tmp[0])
			lineNo := tmp[len(tmp)-1]

			if _, ok := d.commitHashIdSet[hashId]; ok {
				lineNos[lineNo] = struct{}{}
			}
		}

		d.filePathM2LineNos[path] = lineNos
	}

	return nil
}

func (d *diff) cutHashId(source string) string {
	if len(source) > 7 {
		return source[:7]
	}
	return source
}

// ConvToOutputFormat 转换为输出格式
func (d *diff) ConvToOutputFormat() []string {
	sorted := d.convToSortedMap()
	branchInfo := &metadata.BranchesInfo{
		TargetBranchName:  d.TargetBranch,
		CurrentBranchName: d.CurrentBranch,
		StartHashID:       d.CommitHashIdRange[0],
		EndHashID:         d.CommitHashIdRange[1],
	}

	result := make([]string, 0, len(sorted)+1)
	result = append(result, branchInfo.FormatBranchesInfo())
	for path, lineNos := range sorted {
		result = append(result, fmt.Sprintf("%s %s", path, strings.Join(lineNos, ",")))
	}

	return result
}

func (d *diff) ConvToReservedRules() metadata.ReservedRules {
	sorted := d.convToSortedMap()
	result := make(metadata.ReservedRules)
	for path, lineNos := range sorted {
		if len(lineNos) == 0 {
			continue
		}
		set := make(map[int]struct{})
		for _, v := range lineNos {
			set[cast.ToInt(v)] = struct{}{}
		}
		result[path] = &metadata.Rule{
			StartLine: cast.ToInt(lineNos[0]),
			EndLine:   cast.ToInt(lineNos[len(lineNos)-1]),
			LinesSet:  set,
		}
	}
	return result
}

func (d *diff) convToSortedMap() map[string][]string {
	sorted := make(map[string][]string)
	for filePath, lineNos := range d.filePathM2LineNos {
		lines := make([]uint, 0, len(lineNos))
		for no := range lineNos {
			lines = append(lines, cast.ToUint(no))
		}
		sort.Slice(lines, func(i, j int) bool {
			return lines[i] < lines[j]
		})

		strLines := make([]string, 0, len(lines))
		for _, v := range lines {
			strLines = append(strLines, cast.ToString(v))
		}
		sorted[filePath] = strLines
	}
	return sorted
}
