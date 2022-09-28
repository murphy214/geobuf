package cmd

import (
	"strings"

	g "github.com/murphy214/geobuf"

	//"github.com/murphy214/lrs_backend"
	"fmt"

	"github.com/spf13/cobra"
)



func init() {
	rootCmd.AddCommand(csvCmd)

	// rootCmd.PersistentFlags().StringVarP(&filename, "infilename", "f", "", "The input filename to be read from.")
	// rootCmd.PersistentFlags().StringVarP(&outfilename, "outfilename", "o", "", "The output filename to be written to.")
	// viper.BindPFlag("outfilename", rootCmd.PersistentFlags().Lookup("outfilename"))

	//rootCmd.PersistentFlags().IntVarP(&resolution, "limit", "l", 1000, "limit of blocks to be open")
	//viper.BindPFlag("limit", rootCmd.PersistentFlags().Lookup("limit"))
}

var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Writes a csv file from a geojson or geobuf file to std out.",
	Long:  `

Usage:
geobuf csv -f a.geojson 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(filename) == 0 {
			fmt.Println("Input filename or indexes not given!")
		} else if len(filename) > 0 {
			if strings.HasSuffix(filename,"buf") {
				g.ReadGeobufCSVNew(filename)
			} else if strings.HasSuffix(filename,"json") {
				g.ReadGeoJSONCSV(filename)
			}
		}
		//lrs.RenderTableComplete()
	},
}
