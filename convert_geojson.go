package geobuf

import (
	"fmt"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"os"
	"sync"
	//"log"
	//"io"
	"strings"
	"time"
)

// structure used for converting geojson
type Geojson_File struct {
	Features []*geojson.Feature
	Count    int
	File     *os.File
	Pos      int64
	Feat_Pos int
}

// creates a geojosn
func NewGeojson(filename string) Geojson_File {

	file, _ := os.Open(filename)

	bytevals := make([]byte, 100)

	file.ReadAt(bytevals, int64(0))
	boolval := false
	var startpos int
	for ii, i := range string(bytevals) {
		if string(i) == "[" && boolval == false {
			startpos = ii
			boolval = true
		}
	}

	return Geojson_File{File: file, Pos: int64(startpos + 1)}
}

// reads a chunk of a geojson file
func (geojsonfile *Geojson_File) ReadChunk(size int) []string {
	var bytevals []byte
	if size > int(geojsonfile.Pos)+10000000 {
		bytevals = make([]byte, 10000000)
	} else {
		bytevals = make([]byte, size-int(geojsonfile.Pos))
	}

	geojsonfile.File.ReadAt(bytevals, geojsonfile.Pos)
	debt := 0
	//fmt.Println(string(bytevals)[:10])
	newlist := []int{}
	boolval := false
	//fmt.Println("\n",string(bytevals[0:2]),"\n")
	//old_debt := 10
	for i, run := range string(bytevals) {
		//fmt.Println(string(run))
		if "{" == string(run) {
			//fmt.Println("hre")
			boolval = true
			if debt == 0 {
				newlist = append(newlist, i)
			}
			debt += 1
		} else if "}" == string(run) && boolval == true {
			debt -= 1
			if debt == 0 {
				newlist = append(newlist, i)
			}
		}
		//fmt.Println(debt)
		//old_debt = debt
		//string(bytevals)
	}
	boolval = false
	row := []int{}
	geojsons := []string{}
	//fmt.Println(newlist)
	for _, i := range newlist {
		row = append(row, i)
		if boolval == false {
			boolval = true
		} else if boolval == true {
			//fmt.Println(row)
			vals := string(bytevals[row[0]:row[1]])

			geojsons = append(geojsons, vals+"}")

			row = []int{}
			boolval = false
		}

	}
	var newpos int64
	if len(newlist) > 0 {
		newpos = geojsonfile.Pos + int64(newlist[len(newlist)-1])
	} else {
		newpos = int64(size)
	}
	geojsonfile.Pos = newpos
	//fmt.Println(len(geojsons))
	return geojsons
}

// adds featuers
func AddFeatures(geobuf *Writer, feats []string, count int, s time.Time) int {
	var wg sync.WaitGroup
	for _, i := range feats {
		wg.Add(1)
		go func(i string) {
			//fmt.Println(i)
			//fmt.Println(i+"}")
			if len(i) > 0 {
				feat, err := geojson.UnmarshalFeature([]byte(i))
				//fmt.Println(i,feat)
				if err != nil {
					fmt.Println(err, feat)
				} else {
					if feat.Geometry != nil {
						//fmt.Println(ReadFeature(geobuf_raw.WriteFeature(feat)).Geometry)
						geobuf.WriteFeature(feat)
					} else {
						fmt.Println(feat)
					}
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	count += len(feats)
	fmt.Printf("\r%d features created from raw geojson string in %s", count, time.Now().Sub(s))

	return count
}

func GetFilesize(filename string) int {
	fi, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
	}

	// get the Size
	size := fi.Size()
	return int(size)
}

// function used for converting geojson to geobuf
func ConvertGeojson(infile string, outfile string) {
	s := time.Now()
	size := GetFilesize(infile)

	geobuf := WriterFileNew(outfile)
	geojsonfile := NewGeojson(infile)
	count := 0
	feats := []string{"d"}
	//fmt.Println(feats)
	for len(feats) > 0 {
		feats = geojsonfile.ReadChunk(size)
		count = AddFeatures(geobuf, feats, count, s)
	}
}

// function used for converting geojson to geobuf
func ConvertGeobuf(infile string, outfile string) {

	geobuf := ReaderFile(infile)
	//geojsonfile := NewGeojson(outfile)
	//fc.FeatureCollection{}

	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
	}
	file.WriteString(`{"type": "FeatureCollection", "features": [`)

	for geobuf.Next() {
		feature := geobuf.Feature()
		s, _ := feature.MarshalJSON()
		if geobuf.Next() {
			file.Write(append(s, 44))

		} else {
			file.Write(s)

		}
	}
	file.WriteString("]}")
}

type MapFunc func(feature *geojson.Feature) *geojson.Feature

// function used for converting geojson to geobuf
func MapGeobuf(infile string, newfile string, mapfunc MapFunc) {
	geobuf := ReaderFile(infile)
	geobuf2 := WriterFileNew(newfile)
	for geobuf.Next() {
		feature := geobuf.Feature()
		feature = mapfunc(feature)
		geobuf2.WriteFeature(feature)
	}
}

func BenchmarkRead(filename_geojson string, filename_geobuf string) {
	// swapping filenames if needed
	if strings.Contains(filename_geojson, ".geobuf") {
		dummy := filename_geojson
		filename_geojson = filename_geobuf
		filename_geobuf = dummy
	}

	// timing geojson read
	s := time.Now()
	bytevals, err := ioutil.ReadFile(filename_geojson)
	if err != nil {
		fmt.Println(err)
	}
	_, err = geojson.UnmarshalFeatureCollection(bytevals)
	if err != nil {
		fmt.Println(err)
	}
	end_geojson := time.Now().Sub(s)

	s = time.Now()
	geobuf := ReaderFile(filename_geobuf)
	for geobuf.Next() {
		geobuf.Feature()
	}
	end_geobuf := time.Now().Sub(s)

	fmt.Printf("Time to Read Geojson File: %s\nTime to Read Geobuf File: %s\n", end_geojson, end_geobuf)
}

func BenchmarkWrite(filename_geojson string, filename_geobuf string) {
	// swapping filenames if needed
	if strings.Contains(filename_geojson, ".geobuf") {
		dummy := filename_geojson
		filename_geojson = filename_geobuf
		filename_geobuf = dummy
	}

	// getting feature colllection
	bytevals, err := ioutil.ReadFile(filename_geojson)
	if err != nil {
		fmt.Println(err)
	}
	fc, err := geojson.UnmarshalFeatureCollection(bytevals)
	if err != nil {
		fmt.Println(err)
	}

	// timing geojson read
	s := time.Now()
	_, err = fc.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	end_geojson := time.Now().Sub(s)

	s = time.Now()
	geobuf := WriterBufNew()
	for _, feature := range fc.Features {
		geobuf.WriteFeature(feature)
	}
	geobuf.Bytes()
	end_geobuf := time.Now().Sub(s)

	fmt.Printf("Time to Write Geojson File: %s\nTime to Write Geobuf File: %s\n", end_geojson, end_geobuf)
}
