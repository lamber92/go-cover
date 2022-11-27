package cmd

import (
	"log"
	"os"

	"github.com/lamber92/go-cover/internal/trim"
	"github.com/lamber92/go-cover/internal/utils"
	"github.com/spf13/cobra"
)

var trimCmd = &cobra.Command{
	Use:   "trim",
	Short: "trim ${coverage.json}",
	Long:  "trim ${coverage.json}",
	Run: func(cmd *cobra.Command, args []string) {
		runTrim(args)
	},
}

func init() {
	trimCmd.Flags().StringVarP(&difference, "diff", "d", "", "the file-path witch record code difference information")

	rootCmd.AddCommand(trimCmd)
}

func runTrim(args []string) {
	if len(args) == 0 {
		log.Fatalln("Expected at least one coverage json.")
		return
	}
	if len(difference) == 0 {
		log.Fatalln("difference path is empty. skip generating diff-coverage-report.")
		return
	}
	packages, err := trim.Do(args[0], difference)
	if err != nil {
		log.Fatalln(err)
		return
	}
	if err = utils.MarshalJson(os.Stdout, packages); err != nil {
		log.Fatalln(err)
		return
	}
	return
}
