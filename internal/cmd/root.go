package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-cover",
	Short: "go-cover is a converter for go coverage profile --> html report",
	Long: `go-cover is a converter for go coverage profile --> html report
	the command is: convert ${coverage.profile} --report ${report-mode}`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
