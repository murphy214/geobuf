package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)


func init() {
}

var rootCmd = &cobra.Command{
	Use:   "geobuf",
	Short: "Command line tool to interface with geobuf data",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

