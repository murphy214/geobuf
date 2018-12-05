package splitcombine

/* 

Now supports mapping of entire functions on to an arbitary tile size.

The tile mapping is at first a little counter intuitive, and its inputs may seem odd.

For exaple the bounds are used to determine the maximum amount of tiles possible to map to. 

OS defaults to a maximum of 1024 files so we can't go over this. 

Therefore, we abstract out a mappings recursively in a way which limits the maximum number of open files.

The cost of this is of course high levels of reads and writes.

However, this i/o after the initial mapping is trivial by writing a custom read to only get the tileid from a featuer byte array.

Effectively leaving making much more equivalent to simple disk i/o rather than a structure serialization.

This will be pretty signicant..

*/

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
	"math"
	"github.com/murphy214/pbf"
)

type TileMap func(feature *geojson.Feature) map[m.TileID]*geojson.Feature


// getting size of the estimated square grid of tiles
func GetSizeGrid(bb m.Extrema,zoom int) int {
	ne := []float64{bb.E,bb.N}
	sw := []float64{bb.W,bb.S}
	tile1 := m.Tile(ne[0],ne[1],zoom)
	tile2 := m.Tile(sw[0],sw[1],zoom)
	deltax := math.Abs(float64(tile1.X - tile2.X))
	deltay := math.Abs(float64(tile1.Y - tile2.Y))
	fmt.Println(deltax*deltay)
	return int(deltax * deltay)
}

// configuration structure for tile mapping
type TileConfig struct {
	InPlace bool // create new file or replace existing
	Bounds m.Extrema // expected bounds of the mapping
	Zoom int 
	OutputFileName string // defaults to "new_mapped.geobuf"
	Concurrent int // defaults at 1k
}


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
		return 0
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		// Could not obtain stat, handle error
		return 0
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

// adds a set of bytes representing a features
func (splitter *Splitter) AddBytes(key string,bs []byte) {
	splitter.NumberFeatures++
	buf, boolval := splitter.SplitMap[key]
	if !boolval {
		splitter.SplitMap[key] = g.WriterFileNew(key + ".geobuf")
		buf = splitter.SplitMap[key]
	}
	buf.Write(bs)
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

	// closing all underlying writers
	for _,bufw := range splitter.SplitMap {
		bufw.Close()
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
	//fmt.Printf("Combined the %d Sub Files\n", len(filenames))

	// removing all the intermediate files
	for _, i := range filenames {
		os.RemoveAll(i)
	}
	//fmt.Println("Removed All Sub Files")
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


// this is a corner case that continues to come up
// if splitting a file and the sub file needs to be split again
// each subfile contains its own subfile that needs to be handled
// combine all the sub files into one big file 
// then dump meta data into another
func CombineFileSubFiles(filenames []string) {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	bb := m.Extrema{N: north, S: south, E: east, W: west}
	metadata := g.MetaData{
		Bounds:         bb,
		Files:          map[string]*g.SubFile{},
		NumberFeatures: 0,
	}
	currentpos := 0
	file,_ := os.Create("tmp.geobuf")
	for _,filename := range filenames {
		buf := g.ReaderFile(filename)
		filemap := buf.MetaData.Files
		// iterating through each file in flemap
		for k,v := range filemap {
			oldpos := currentpos
			filesize := v.Positions[1]-v.Positions[0]
			currentpos+=filesize
			rawbytes := make([]byte,int64(filesize))
			buf.File.ReadAt(rawbytes,int64(v.Positions[0]))
			subfile := &g.SubFile{
				Positions: [2]int{oldpos, currentpos},
				Size: currentpos - oldpos,
			}
			metadata.Files[k] = subfile
			file.Write(rawbytes)
		}
		buf.Close()
		os.Remove(filename)
	}
	// creating the feature that inddicates this contains metadata
	feature := geojson.NewPointFeature([]float64{0, 0})
	feature.Properties = map[string]interface{}{"metadata": g.WriteMetaData(metadata)}

	// creating new buffer that will soon replace the old one
	buf := g.WriterFileNew("meta.geobuf")
	buf.WriteFeature(feature)
	mycmd := "cat meta.geobuf tmp.geobuf >> tmp2.geobuf"

	// runnign the command string combining all the files into one
	cmd := exec.Command("bash", "-c", mycmd)
	cmd.Run()

	os.Remove("meta.geobuf")
	os.Remove("tmp.geobuf")
	os.Rename("tmp2.geobuf","tmp.geobuf")
}

// structure for finding overlapping values
func Overlapping_1D(box1min float64,box1max float64,box2min float64,box2max float64) bool {
	if box1max >= box2min && box2max >= box1min {
		return true
	} else {
		return false
	}
	return false
}

// returns a boolval for whether or not the bb intersects
func Intersect(bdsref m.Extrema,bds m.Extrema) bool {
	if Overlapping_1D(bdsref.W,bdsref.E,bds.W,bds.E) && Overlapping_1D(bdsref.S,bdsref.N,bds.S,bds.N) {
		return true
	} else {
		return false
	}
	return false
}

// this function is intended to be a light weight 
// read of feature bytes with minimal logic / alloc to 
// read the single feature property we need, "TILEID"
func LazyFeatureTileID(bs []byte) m.TileID {
	pbfval := pbf.PBF{Pbf: bs, Length: len(bs)}

	// id logic
	key, val := pbfval.ReadKey()
	if key == 1 && val == 0 {
		pbfval.Varint()
		key, val = pbfval.ReadKey()
	}

	// property logic / for loop
	// here we take noticeable deviation from prior implementations
	// the changes generally are from type serialization to byte traversal
	for key == 2 && val == 2 {
		// starting properties shit here

		size := pbfval.ReadVarint()
		endpos := pbfval.Pos + size
		//pbfval.ReadKey()
		pbfval.Pos += 1
		keyvalue := pbfval.ReadString()
		if keyvalue == "TILEID" {
			pbfval.Pos += 1
			pbfval.Varint()
			pbfval.ReadKey()
			return m.TileFromString(pbfval.ReadString())
		} else {
			pbfval.Pos += 1
			pbfval.Pos += pbfval.ReadVarint()
		}
		pbfval.Pos = endpos
		key, val = pbfval.ReadKey()
	}
	return m.TileID{}
}

// processes concurrent features
func concurrentfeatures(features []*geojson.Feature,splitter *Splitter,bds m.Extrema,zoom int,currentzoom int,mapfunc TileMap) {
	c := make(chan map[m.TileID]*geojson.Feature,len(features))
	for _,feature := range features {
		tmpbb := feature.BoundingBox
		newbb := m.Extrema{W:tmpbb[0],S:tmpbb[1],E:tmpbb[2],N:tmpbb[3]}
		if Intersect(newbb,bds) {
			go func(feature *geojson.Feature,zoom int,c chan map[m.TileID]*geojson.Feature) {
				c <- mapfunc(feature)
			}(feature,zoom,c)
		} else {
			c <- map[m.TileID]*geojson.Feature{}
		}
	}
	for range features {
		out := <-c
		for tile,feature := range out {
			// permeating the mapped tile to the current zoom level were mapping
			feature.Properties["TILEID"] = m.TilestrFile(tile)
			for int(tile.Z) != currentzoom {
				tile = m.Parent(tile)
			}
			splitter.AddFeature(m.TilestrFile(tile),feature)			
		}
	}
}  


type LoggingInitialMap struct {
	StartTime time.Time
}

// logs the number of features and number of files
func (logger *LoggingInitialMap) Log(numberoffeautures int,featurescreated,numberoffiles int,percent_read float64) {
	timepassed := time.Now().Sub(logger.StartTime).Seconds()
	fmt.Printf(
		"\rInitial Map| %.1f%% | Features: %dk | Features Created: %dk | Number of Files: %d | Features / s: %d",
		percent_read*100,
		numberoffeautures/1000,
		featurescreated/1000,
		numberoffiles,
		int(float64(numberoffeautures)/timepassed),
	)
}


var StartZoom = 0
var EndZoom = 0


// a powerful function that maps an entire geobuf into a submapping that can be navigated
// through geobufs seek api 
// one can effectively view this as a sort with about tiles 
func MapGeobuf(filename string,mapfunc TileMap, tileconfig *TileConfig) {
	// setting default outfilename
	// collisons with other names are possible but unlikely
	if len(tileconfig.OutputFileName) == 0 {
		tileconfig.OutputFileName = "new_mapped.geobuf"
	}

	if tileconfig.Concurrent == 0 {
		tileconfig.Concurrent = 1000
	}


	logger := LoggingInitialMap{time.Now()}

	// determining the largest size zoom we can start at
	size := GetSizeGrid(tileconfig.Bounds,tileconfig.Zoom)
	currentzoom := tileconfig.Zoom
	for size > 750 {
		size = size / 4 
		currentzoom--
	}	

	StartZoom = currentzoom
	EndZoom = tileconfig.Zoom
	
	// getting size
	filesize := GetSize(filename)


	// creating tile buffer
	buf := g.ReaderFile(filename)
	
	// creating new splitter 
	// this abstracts away an underlying mapping process out of memory
	splitter := NewSplitter(buf)

	// iterating through each feature
	i := 0
	features := []*geojson.Feature{}
	for buf.Next() {
		feature := buf.Feature()
		features = append(features,feature)


		/*
		THIS ADDS CONCURRENCY LEAVING HERE IF I DECIDE TO CHANGE SHIT
		tmpbb := feature.BoundingBox
		newbb := m.Extrema{W:tmpbb[0],S:tmpbb[1],E:tmpbb[2],N:tmpbb[3]}

		if Intersect(newbb,tileconfig.Bounds) {
			for tile,v := range mapfunc(feature) {
				// permeating the mapped tile to the current zoom level were mapping
				feature.Properties["TILEID"] = m.TilestrFile(tile)
				for int(tile.Z) != currentzoom {
					tile = m.Parent(tile)
				}
				splitter.AddFeature(m.TilestrFile(tile),v)			
			}
		}
		*/

		if len(features) == tileconfig.Concurrent {
			concurrentfeatures(features, splitter, tileconfig.Bounds, tileconfig.Zoom,currentzoom, mapfunc)
			features = []*geojson.Feature{}
		}
		if i%tileconfig.Concurrent==0 {
			percent_read := float64(buf.Reader.TotalPosition)/float64(filesize)
			logger.Log(i,splitter.NumberFeatures,len(splitter.SplitMap),percent_read)
		}
		i++
	}


	concurrentfeatures(features, splitter, tileconfig.Bounds, tileconfig.Zoom,currentzoom, mapfunc)

	// combining all the files
	splitter.Combine()
	os.Rename("tmp.geobuf","new.geobuf")

	// this code steps through a zoom level with steps as large as possible (4) 256 max open files
	// for every one of these jumps we require an read / write of the entire file
	// luckily this is almost pure i/o 
	fmt.Println("Starting in-place split combine on mapped features.")
	for currentzoom != tileconfig.Zoom {
		// determing the step we will take
		var delta int
		if tileconfig.Zoom - currentzoom < 4 {
			delta = tileconfig.Zoom - currentzoom
		} else {
			delta = 4
		}
		currentzoom+=delta
		fmt.Printf("Top level, Current zoom changed: %d, delta: %d\n",currentzoom,delta)
		// opening current geobuf iteration & getting filemap
		newbuf := g.ReaderFile("new.geobuf")
		filemap := newbuf.MetaData.Files
		newlist := []string{}
		i := 0
		for k := range filemap {
			// putting read cursor at the start of the tile k block
			newbuf.SubFileSeek(k)
			fmt.Printf("\rRemapping subfile %v to zoom %d subfiles [%d/%d]",k,currentzoom,i,len(filemap))
			// creating the subsplitter
			// and passing through the current split to the next zoom
			subsplitter := NewSplitter(newbuf)
			for newbuf.SubFileNext() {
				bs := newbuf.Bytes()
				tile := LazyFeatureTileID(bs)
				for int(tile.Z) != currentzoom {
					tile = m.Parent(tile)
				}
				subsplitter.AddBytes(m.TilestrFile(tile),bs)
			}
			
			// combining and cleaning up
			subsplitter.Combine()
			os.Rename("tmp.geobuf",k+".geobuf")
			newlist = append(newlist,k+".geobuf")
			i++
		}

		// cleaning up each iteration
		CombineFileSubFiles(newlist)
		os.Remove("new.geobuf")
		os.Rename("tmp.geobuf","new.geobuf")
	}

	// final clean up process
	if tileconfig.InPlace {
		os.Remove(filename)
		os.Rename("new.geobuf",filename)
	} else {
		os.Rename("new.geobuf",tileconfig.OutputFileName)
	}	
}