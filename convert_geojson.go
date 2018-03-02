package geobuf_new


import (
	//"io/ioutil"
	"github.com/paulmach/go.geojson"
	"fmt"
	"os"
	"sync"
	//"log"
	//"io"
	//"strings"
)

// structure used for converting geojson
type Geojson_File struct {
	Features []*geojson.Feature
	Count int
	File *os.File
	Pos int64
	Feat_Pos int
}

// creates a geojosn
func New_Geojson(filename string) Geojson_File {

	file,_ := os.Open(filename)

	bytevals := make([]byte,100)

	file.ReadAt(bytevals,int64(0))
	boolval := false
	var startpos int
	for ii,i := range string(bytevals) {
		if string(i) == "[" && boolval == false {
			startpos = ii
			boolval = true
		}
	}

	return Geojson_File{File:file,Pos:int64(startpos+1)}
}

// reads a chunk of a geojson file
func (geojsonfile *Geojson_File) Read_Chunk(size int) []string {
	var bytevals []byte
	if size > int(geojsonfile.Pos) + 10000000 {
		bytevals = make([]byte,10000000)
	} else {
		bytevals = make([]byte,size - int(geojsonfile.Pos))
	}



	geojsonfile.File.ReadAt(bytevals,geojsonfile.Pos)
	debt := 0
	//fmt.Println(string(bytevals)[:10])
	newlist := []int{}
	boolval := false
	//fmt.Println("\n",string(bytevals[0:2]),"\n")
	//old_debt := 10
	for i,run := range string(bytevals) {
		//fmt.Println(string(run))
		if "{" == string(run) {
			//fmt.Println("hre")
			boolval = true
			if debt == 0 {
				newlist = append(newlist,i)
			}
			debt += 1
		} else if "}" == string(run) && boolval == true {
			debt -= 1
			if debt == 0 {
				newlist = append(newlist,i)
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
	for _,i := range newlist {
		row = append(row,i)
		if boolval == false {
			boolval = true
		} else if boolval == true {
			//fmt.Println(row)
			vals := string(bytevals[row[0]:row[1]])

			geojsons = append(geojsons,vals + "}")

			
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

// splits feature into smaller groups
func Split_Feature(i *geojson.Feature) []*geojson.Feature {
	if  i.Geometry == nil {
		return []*geojson.Feature{}
	} else if i.Geometry.Type == "MultiLineString" {
		props := i.Properties
		newfeats := []*geojson.Feature{}
		for _,newline := range i.Geometry.MultiLineString {
			newfeats = append(newfeats,&geojson.Feature{Geometry:&geojson.Geometry{LineString:newline,Type:"LineString"},Properties:props,ID:i.ID})
		}

		return newfeats
	} else if i.Geometry.Type == "MultiPolygon" {
		props := i.Properties

		newfeats := []*geojson.Feature{}
		for _,newline := range i.Geometry.MultiPolygon {
			newfeats = append(newfeats,&geojson.Feature{Geometry:&geojson.Geometry{Polygon:newline,Type:"Polygon"},Properties:props,ID:i.ID})
		}
		return newfeats
	} else if i.Geometry.Type == "MultiPoint" {
		props := i.Properties


		newfeats := []*geojson.Feature{}
		for _,newline := range i.Geometry.MultiPoint {
			newfeats = append(newfeats,&geojson.Feature{Geometry:&geojson.Geometry{Point:newline,Type:"Point"},Properties:props,ID:i.ID})
		}
		return newfeats
	} else {

		return []*geojson.Feature{i}
	}
	return []*geojson.Feature{}
}



// adds featuers
func Add_Features(geobuf Writer,feats []string,count int) int {
	var wg sync.WaitGroup
	for _,i := range feats {
		wg.Add(1)
		go func(i string) {
			//fmt.Println(i)
			//fmt.Println(i+"}")
			if len(i) > 0 {
				feat,err := geojson.UnmarshalFeature([]byte(i))
				//fmt.Println(i,feat)
				if err != nil {
					fmt.Println(err)
				}
			//fmt.Println(feat)
				feats2 := Split_Feature(feat)
				for _,feat2 := range feats2 {
					feat2.ID = count
			//bytevals = append(bytevals,Write_Feature(feat)...)
					geobuf.Write_Feature(feat2)
					count += 1
				}
			}
			fmt.Printf("\r[%d/%d] Creating geojson_buf from raw geojson string",count,len(feats))
			wg.Done()
		}(i)
	}
	wg.Wait()
	return count

}


func Get_Filesize(filename string) int {
	fi, err := os.Stat(filename);
	if err != nil {
		fmt.Println(err)
	}

	// get the Size
	size := fi.Size()
	return int(size)
}

// function used for converting geojson to geobuf
func Convert_Geojson(infile string,outfile string) {
	size :=Get_Filesize(infile)

	geobuf := Writer_File_New(outfile)
	geojsonfile := New_Geojson(infile)
	count := 0
	feats := []string{"d"}
	//fmt.Println(feats)
	for len(feats) > 0 {
		feats = geojsonfile.Read_Chunk(size)

		count = Add_Features(geobuf,feats,count)
	}
	//geobuf.Writer.Flush()
}