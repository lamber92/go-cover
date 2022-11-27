package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lamber92/go-cover/internal/diff"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		runDiff()
	},
}

const defaultTargetBranch = "master"

var (
	currentBranch     string
	targetBranch      string
	hashIdsRangeParam string
)

func init() {
	diffCmd.Flags().StringVarP(&currentBranch, "current-branch", "c", "", "The current branch under test")
	diffCmd.Flags().StringVarP(&targetBranch, "target-branch", "t", defaultTargetBranch, "The branch that was compared to find the difference")
	diffCmd.Flags().StringVarP(&hashIdsRangeParam, "hash-ids-range", "i", "", "The range of hash-ids that need to be reserved. format: 'start-hash-id,end-hash-id'")

	rootCmd.AddCommand(diffCmd)
}

func runDiff() {
	hashIdsRange, err := parseHashIdsRange()
	if err != nil {
		log.Fatal(err)
	}
	diffMgr, err := diff.Do(currentBranch, targetBranch, hashIdsRange)
	if err != nil {
		log.Fatalln(err)
	}
	content := diffMgr.ConvToOutputFormat()

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for _, v := range content {
		if _, err = fmt.Fprintln(out, v); err != nil {
			log.Fatalf("failed to write diff-info to stdout. err: %v\n", err)
		}
	}
}

func parseHashIdsRange() ([]string, error) {
	hashIdsRange := make([]string, 0, 2)
	if len(hashIdsRangeParam) > 0 {
		tmp := strings.Split(hashIdsRangeParam, ",")
		if len(tmp) != 2 {
			return nil, fmt.Errorf("invalid hash-ids range fomart. [%s]\n", hashIdsRangeParam)
		}
		hashIdsRange = append(hashIdsRange, tmp[0], tmp[1])
	}
	return hashIdsRange, nil
}
