package cmd

import (
	"log"
	"os"

	"github.com/jinzhu/copier"
	"github.com/lamber92/go-cover/internal/convert"
	"github.com/lamber92/go-cover/internal/diff"
	"github.com/lamber92/go-cover/internal/metadata"
	"github.com/lamber92/go-cover/internal/report"
	"github.com/lamber92/go-cover/internal/trim"
	"github.com/lamber92/go-cover/internal/utils"
	"github.com/spf13/cobra"
)

const (
	outputModeAll      = "all"       // 都要
	outputModeOnlyFull = "full-only" // 只要全量报告
	outputModeOnlyDiff = "diff-only" // 只要增量报告
	outputModeOnlyJson = "json-only" // 只要全量的json(只输出到stdout)
)

var (
	outputMode string
	css        string
	difference string
)

var covertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert ${coverage.profile}",
	Long:  "convert ${coverage.profile}",
	Run: func(cmd *cobra.Command, args []string) {
		runConvert(args)
	},
}

func init() {
	covertCmd.Flags().StringVarP(&outputMode, "output-mode", "o", outputModeAll, "Options: 'full-only' or 'diff-only'; Default: 'all'")
	covertCmd.Flags().StringVarP(&css, "css-format", "f", "", "The file-path witch record customized report themes within CSS-format")
	covertCmd.Flags().StringVarP(&difference, "diff", "d", "", "The file-path witch record code difference information")
	covertCmd.Flags().StringVarP(&currentBranch, "current-branch", "c", "", "The current branch under test")
	covertCmd.Flags().StringVarP(&targetBranch, "target-branch", "t", defaultTargetBranch, "The branch that was compared to find the difference")
	covertCmd.Flags().StringVarP(&hashIdsRangeParam, "hash-ids-range", "i", "", "The range of hash-ids that need to be reserved. format: 'start-hash-id,end-hash-id'")

	rootCmd.AddCommand(covertCmd)
}

func runConvert(args []string) {
	if len(args) == 0 {
		log.Fatalln("Expected at least one coverage profile.")
		return
	}

	packages, err := convert.Do(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	switch outputMode {
	case outputModeOnlyJson:
		// 如果是只要json, 完成直接退出
		if err = utils.MarshalJson(os.Stdout, packages); err != nil {
			log.Fatalf("Failed to generate json. err: %v\n", err)
		}
	case outputModeOnlyFull:
		buildFullReport(packages)
	case outputModeOnlyDiff:
		buildDiffReport(packages)
	case outputModeAll:
		buildFullReport(packages)
		buildDiffReport(packages)
	default:
		log.Fatalf("Unsupported output mode. [%s]", outputMode)
	}

	return
}

func buildFullReport(packages utils.Packages) {
	newPkg := make(utils.Packages, 0)
	if err := copier.CopyWithOption(&newPkg, &packages, copier.Option{DeepCopy: true}); err != nil {
		log.Fatalf("Handle packages data failed. err: %v\n", err)
	}

	currentBranch, err := utils.GetCurrentBranch()
	if err != nil {
		log.Fatalf(err.Error())
	}
	param := &report.GenerateHTMLParam{
		Packages:     newPkg,
		CSS:          css,
		Dir:          ".",
		FileName:     utils.FullHTML,
		BranchesInfo: &metadata.BranchesInfo{CurrentBranchName: currentBranch},
	}
	if err = report.GenerateHTML(param); err != nil {
		log.Fatalf("Failed to generate full-coverage-report. err: %v\n", err)
		return
	}
	log.Println("Generate full-coverage-report success.")
}

func buildDiffReport(packages utils.Packages) {
	var (
		diffPackages utils.Packages
		branchesInfo *metadata.BranchesInfo
		err          error
	)

	if len(difference) > 0 {
		info, err := utils.LoadReservedInfo(difference)
		if err != nil {
			return
		}
		branchesInfo = info.Branches
		diffPackages, err = trim.TrimPackages(packages, info.Rules)
		if err != nil {
			log.Fatalf("Failed to trim diff-coverage. err: %v\n", err)
			return
		}
	} else {
		hashIdsRange, err := parseHashIdsRange()
		if err != nil {
			log.Fatalln(err)
		}
		diffMgr, err := diff.Do(currentBranch, targetBranch, hashIdsRange)
		if err != nil {
			log.Fatalln(err)
		}
		branchesInfo = &metadata.BranchesInfo{
			TargetBranchName:  diffMgr.TargetBranch,
			CurrentBranchName: diffMgr.CurrentBranch,
			StartHashID:       diffMgr.CommitHashIdRange[0],
			EndHashID:         diffMgr.CommitHashIdRange[1],
		}
		diffPackages, err = trim.TrimPackages(packages, diffMgr.ConvToReservedRules())
		if err != nil {
			log.Fatalf("Failed to trim diff-coverage. err: %v\n", err)
			return
		}
	}

	param := &report.GenerateHTMLParam{
		Packages:     diffPackages,
		CSS:          css,
		Dir:          ".",
		FileName:     utils.DiffHTML,
		BranchesInfo: branchesInfo,
	}
	if err = report.GenerateHTML(param); err != nil {
		log.Fatalf("Failed to generate diff-coverage-report. err: %v\n", err)
		return
	}
	log.Println("Generate diff-coverage-report success.")
}
