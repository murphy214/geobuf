package cmd

import (
	"os"
	"os/exec"

	g "github.com/murphy214/geobuf"

	//"github.com/murphy214/lrs_backend"
	"fmt"

	"github.com/spf13/cobra"
)

var outfilename string
var filename string

func init() {
	rootCmd.AddCommand(geobuf2geojsonCmd)
	rootCmd.PersistentFlags().StringVarP(&filename, "infilename", "f", "", "The input filename to be read from.")
	rootCmd.PersistentFlags().StringVarP(&outfilename, "outfilename", "o", "", "The output filename to be written to.")
	// viper.BindPFlag("outfilename", rootCmd.PersistentFlags().Lookup("outfilename"))

	//rootCmd.PersistentFlags().IntVarP(&resolution, "limit", "l", 1000, "limit of blocks to be open")
	//viper.BindPFlag("limit", rootCmd.PersistentFlags().Lookup("limit"))
}

var geobuf2geojsonCmd = &cobra.Command{
	Use:   "buf2json",
	Short: "Converts a geobuf to geojson",
	Long:  `

Usage:
geobuf buf2json -f a.geobuf -o a.geojson 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(filename) == 0 {
			fmt.Println("Input filename not given!")
		} else if len(filename) > 0 && len(outfilename) > 0 {
			fmt.Println(filename)
			g.ConvertGeobuf(filename,outfilename)
		} else if len(filename) > 0 {
			defer os.Remove("tmp.geojson")
			fmt.Println(filename)
			g.ConvertGeobuf(filename,"tmp.geojson")
			bs,_ := exec.Command("cat","tmp.geojson").Output()
			fmt.Println(string(bs))
		}
		//lrs.RenderTableComplete()
	},
}
