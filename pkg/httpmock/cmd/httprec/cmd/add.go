package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/processone/dpk/pkg/httpmock"
)

var URI string

// addCmd represents the get command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a request to a scenario",
	Long: `Add command will add an HTTP request to a scenario.
It will create the scenario file if it does not exist already.`,
	Run: func(cmd *cobra.Command, args []string) {
		scenarioFile := args[0]
		recorder := httpmock.Recorder{Logger: log.New(os.Stderr, "", 0)}
		if err := recorder.Record(URI, scenarioFile); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Args = cobra.ExactArgs(1) // Scenario name
	addCmd.Flags().StringVarP(&URI, "url", "u", "", "URL endpoint for HTTP request to record")
	addCmd.Example = "  httprec add fixtures/scenario1 -url https://www.process-one.net/"
	if err := addCmd.MarkFlagRequired("url"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
