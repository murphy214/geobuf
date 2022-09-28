package cmd

import (
	"io/ioutil"
	"strconv"
	"strings"

	g "github.com/murphy214/geobuf"
	geojson "github.com/paulmach/go.geojson"

	//"github.com/murphy214/lrs_backend"
	"fmt"

	"github.com/spf13/cobra"
)

func parseinds(inds []string) []int {
	myinds := make([]int,len(inds))

	if len(inds) > 2 {
		return []int{}
	} else {
		for pos,ind := range inds {
			myint,err := strconv.ParseInt(ind,10,64)
			if err != nil {
				return []int{}
			}
			myinds[pos] = int(myint)
		} 
	}
	return myinds
}

var inds []string
func init() {
	rootCmd.AddCommand(printFeatureInd)
	rootCmd.PersistentFlags().StringArrayVarP(&inds, "inds", "i", []string{}, "The index positions of the feature collection")

	// rootCmd.PersistentFlags().StringVarP(&filename, "infilename", "f", "", "The input filename to be read from.")
	// rootCmd.PersistentFlags().StringVarP(&outfilename, "outfilename", "o", "", "The output filename to be written to.")
	// viper.BindPFlag("outfilename", rootCmd.PersistentFlags().Lookup("outfilename"))

	//rootCmd.PersistentFlags().IntVarP(&resolution, "limit", "l", 1000, "limit of blocks to be open")
	//viper.BindPFlag("limit", rootCmd.PersistentFlags().Lookup("limit"))
}

var printFeatureInd = &cobra.Command{
	Use:   "feature",
	Short: "Prints a feature in a geojson file",
	Long:  `

Usage:
geobuf feature -i 0 -f a.geojson 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(filename) == 0 || len(inds) == 0 {
			fmt.Println("Input filename or indexes not given!")
		} else if len(filename) > 0 && len(inds) > 0 {
			myinds := parseinds(inds)
			if len(myinds) == 1 {
				myind := myinds[0]
				if strings.HasSuffix(filename,"buf") {
					buf := g.ReaderFile(filename)
					ind := 0
					for buf.Next() {
						if ind == myind {
							feat := buf.Feature()
							bs,_ := feat.MarshalJSON()
							fmt.Println(string(bs))
						}
						buf.Reader.Protobuf()
						ind++
					}
				} else if strings.HasSuffix(filename,"json") {
					bs,_ := ioutil.ReadFile(filename)
					fc,err := geojson.UnmarshalFeatureCollection(bs)
					if err != nil {
						fmt.Println(err)
					}
					if len(fc.Features) > myind {
						feat := fc.Features[myind]
						bs,_ := feat.MarshalJSON()
						fmt.Println(string(bs))
					} 
				}
			} else if len(inds) == 2 {
				myind1,myind2 := myinds[0],myinds[1]
				if strings.HasSuffix(filename,"buf") {
					buf := g.ReaderFile(filename)
					ind := 0
					for buf.Next() {
						if myind1 <= ind && myind2 >= ind {
							feat := buf.Feature()
							bs,_ := feat.MarshalJSON()
							fmt.Println(string(bs))
						}
						buf.Reader.Protobuf()
						ind++
					}
				} else if strings.HasSuffix(filename,"json") {
					bs,_ := ioutil.ReadFile(filename)
					fc,err := geojson.UnmarshalFeatureCollection(bs)
					if err != nil {
						fmt.Println(err)
					}
					if len(fc.Features) > myind2 {
						feats := fc.Features[myind1:myind2]
						for _,feat := range feats {
							bs,_ := feat.MarshalJSON()
							fmt.Println(string(bs))
						}
					} 
				}
			}

		} else if len(filename) > 0 {
			fmt.Println("Output filename not given!")
		}
		//lrs.RenderTableComplete()
	},
}
