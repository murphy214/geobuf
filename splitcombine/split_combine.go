package splitcombine

import (
	"fmt"
	g "github.com/murphy214/geobuf"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/tile-cover"
	"github.com/paulmach/go.geojson"
	"os"
	"os/exec"
	"strings"
	"time"
)

// this structure allows for easy splitting of tiles
type Splitter struct {
	Bounds         m.Extrema
	Reader         *g.Reader
	NumberFeatures int
	SplitMap       map[string]*g.Writer
}

func NewSplitter(buf *g.Reader) *Splitter {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	bb := m.Extrema{N: north, S: south, E: east, W: west}
	return &Splitter{Bounds: bb, Reader: buf, SplitMap: map[string]*g.Writer{}}
}

// pushs two bounding box values
func PushTwoBoundingBoxs(bb1, bb2 m.Extrema) m.Extrema {
	// setting opposite default values
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	// setting bb1 and bb2
	west1, south1, east1, north1 := bb1.W, bb1.S, bb1.E, bb1.N
	west2, south2, east2, north2 := bb2.W, bb2.S, bb2.E, bb2.N

	// handling west values: min
	if west1 < west2 {
		west = west1
	} else {
		west = west2
	}

	// handling south values: min
	if south1 < south2 {
		south = south1
	} else {
		south = south2
	}

	// handling east values: max
	if east1 > east2 {
		east = east1
	} else {
		east = east2
	}

	// handling north values: max
	if north1 > north2 {
		north = north1
	} else {
		north = north2
	}

	return m.Extrema{N: north, S: south, E: east, W: west}
}

// gets the size of a given file
func GetSize(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		// Could not obtain stat, handle error
		fmt.Println(err)
	}

	return int(fi.Size())
}

// adds a single feature to the split map
func (splitter *Splitter) AddFeature(key string, feature *geojson.Feature) {
	if len(feature.BoundingBox) == 4 {
		bb := feature.BoundingBox
		west, south, east, north := bb[0], bb[1], bb[2], bb[3]
		bbnew := m.Extrema{N: north, S: south, E: east, W: west}
		splitter.Bounds = PushTwoBoundingBoxs(splitter.Bounds, bbnew)
	}

	splitter.NumberFeatures += 1
	buf, boolval := splitter.SplitMap[key]
	if !boolval {
		splitter.SplitMap[key] = g.WriterFileNew(key + ".geobuf")
		buf = splitter.SplitMap[key]
	}
	buf.WriteFeature(feature)
}

// maps a function that generates a key to an entire geobuf file
func (splitter *Splitter) MapToSubFiles(myfunc func(feature *geojson.Feature) []string) {
	i := 0
	s := time.Now()
	for splitter.Reader.Next() {
		feature := splitter.Reader.Feature()
		keys := myfunc(feature)
		for _, key := range keys {
			splitter.AddFeature(key, feature)
		}
		i++
		//fmt.Println(i)
		if i%1000 == 0 {
			fmt.Printf("\r%d Features Split in %s", i, time.Now().Sub(s))
		}
	}
	fmt.Println()

}

// combines all hte intermediate files and removes them
func (splitter *Splitter) Combine() {
	// creating metadata structure
	metadata := g.MetaData{
		Bounds:         splitter.Bounds,
		Files:          map[string]*g.SubFile{},
		NumberFeatures: splitter.NumberFeatures,
	}

	// iterating through each file and adding subfile metadata
	currentpos := 0
	filenames := []string{}
	for k := range splitter.SplitMap {
		oldpos := currentpos
		filesize := GetSize(k + ".geobuf")
		currentpos += filesize
		subfile := &g.SubFile{Positions: [2]int{oldpos, currentpos}, Size: currentpos - oldpos}
		metadata.Files[k] = subfile
		filenames = append(filenames, k+".geobuf")
	}

	// creating the feature that inddicates this contains metadata
	feature := geojson.NewPointFeature([]float64{0, 0})
	feature.Properties = map[string]interface{}{"metadata": g.WriteMetaData(metadata)}

	// creating new buffer that will soon replace the old one
	buf := g.WriterFileNew("tmp.geobuf")
	buf.WriteFeature(feature)

	// now preparing the command string
	mycmd := "cat " + strings.Join(filenames, " ") + " >> tmp.geobuf"

	// runnign the command string combining all the files into one
	cmd := exec.Command("bash", "-c", mycmd)
	cmd.Run()
	fmt.Printf("Combined the %d Sub Files\n", len(filenames))

	// removing all the intermediate files
	for _, i := range filenames {
		os.RemoveAll(i)
	}
	fmt.Println("Removed All Sub Files")
}

// wrapping all the methods up
func SplitCombineFile(buf *g.Reader, myfunc func(feature *geojson.Feature) []string) {
	splitter := NewSplitter(buf)
	splitter.MapToSubFiles(myfunc)
	splitter.Combine()
	os.Remove(buf.Filename)
	os.Rename("tmp.geobuf", buf.Filename)
}

// a function to split and combine tiles
func SplitCombineTiles(buf *g.Reader, zoom int) *g.Reader {
	// defining function
	myfunc := func(feature *geojson.Feature) []string {
		tiles := tilecover.TileCover(feature, zoom)
		mystringtiles := make([]string, len(tiles))
		for pos, tile := range tiles {
			mystringtiles[pos] = m.TilestrFile(tile)
		}
		return mystringtiles
	}
	// running operation
	SplitCombineFile(buf, myfunc)
	return g.ReaderFile(buf.Filename)
}
