package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/lamber92/go-cover/internal/metadata"
)

// LoadReservedInfo 加载保留规则
func LoadReservedInfo(path string) (results *metadata.ReservedInfo, err error) {
	// 文件格式如下：
	// target_branch_name,current_branch_name:start_commit_id,end_commit_id
	// xxx.go line1,line2,line3.....
	// yyy.go line1,line4,line9.....
	//
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		err = fmt.Errorf("fail to open diff-file. path: %s, err: %v", path, err)
		return
	}
	defer file.Close()

	results = &metadata.ReservedInfo{
		Branches: nil,
		Rules:    make(metadata.ReservedRules, 0),
	}

	// 按行读取规则文件内容
	reader := bufio.NewReader(file)
	firstLine := true
	for {
		buff, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		data := string(buff)

		// 第一行是分支信息
		if firstLine {
			if results.Branches, err = metadata.ParseBranchesInfo(data); err != nil {
				return nil, err
			}
			firstLine = false
			continue
		}

		// 拆分文件与行列表
		columns := strings.Split(data, " ")
		if len(columns) != 2 {
			log.Printf("invalid columns: %+v\n", columns)
			continue
		}
		// 拆分行列表
		lines := strings.Split(columns[1], ",")
		if len(lines) == 0 {
			log.Printf("invalid lines: %+v\n", lines)
			continue
		}

		// 构造最终结果
		var (
			newLines    = StringsToInts(lines)
			newLinesSet = make(map[int]struct{})
		)
		for _, v := range newLines {
			newLinesSet[v] = struct{}{}
		}
		results.Rules.Add(columns[0], &metadata.Rule{
			StartLine: newLines[0],
			EndLine:   newLines[len(newLines)-1],
			LinesSet:  newLinesSet,
		})
	}

	return
}
