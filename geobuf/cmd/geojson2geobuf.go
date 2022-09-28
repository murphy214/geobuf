package cmd

import (

	g "github.com/murphy214/geobuf"

	//"github.com/murphy214/lrs_backend"
	"fmt"

	"github.com/spf13/cobra"
)


func init() {
	rootCmd.AddCommand(geojson2geobufCmd)
	// rootCmd.PersistentFlags().StringVarP(&filename, "infilename", "f", "", "The input filename to be read from.")
	// rootCmd.PersistentFlags().StringVarP(&outfilename, "outfilename", "o", "", "The output filename to be written to.")
	// viper.BindPFlag("outfilename", rootCmd.PersistentFlags().Lookup("outfilename"))

	//rootCmd.PersistentFlags().IntVarP(&resolution, "limit", "l", 1000, "limit of blocks to be open")
	//viper.BindPFlag("limit", rootCmd.PersistentFlags().Lookup("limit"))
}

var geojson2geobufCmd = &cobra.Command{
	Use:   "json2buf",
	Short: "Converts a geojson to geobuf",
	Long:  `

Usage:
geobuf json2buf -f a.geojson -o a.geobuf 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(filename) == 0 {
			fmt.Println("Input filename not given!")
		} else if len(filename) > 0 && len(outfilename) > 0 {
			fmt.Printf("Converting geojson: %s to file %s\n",filename,outfilename)
			g.ConvertGeojson(filename,outfilename)
		} else if len(filename) > 0 {
			fmt.Println("Output filename not given!")
		}
		//lrs.RenderTableComplete()
	},
}
