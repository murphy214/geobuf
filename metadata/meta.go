package metadata

import (
	"fmt"
	g "github.com/murphy214/geobuf"
	raw "github.com/murphy214/geobuf/geobuf_raw"
	"github.com/paulmach/go.geojson"
	"os"
	"time"
)

type Meta struct {
	Type          string // geometry type
	Verticies     int    // geometry verticies
	Properties    int    // number properties
	SizeJSON      int    // size in bytes
	TimeReadJSON  int    // time in nanoseconds for read
	TimeWriteJSON int    // time in nanosecons for write
	SizeBUF       int    // size in bytes
	TimeReadBUF   int    // time in bytes for read
	TimeWriteBUF  int    // time in bytes for write
}

func (meta *Meta) MakeString() string {
	return fmt.Sprintf("%s,%d,%d,%d,%d,%d,%d,%d,%d\n", meta.Type, meta.Verticies, meta.Properties,
		meta.SizeJSON, meta.TimeReadJSON, meta.TimeWriteJSON, meta.SizeBUF, meta.TimeReadBUF, meta.TimeWriteBUF,
	)
}

//
type MetaCSV struct {
	File     *os.File
	FileName string
}

// creates a meta csv
func NewMetaCSV(filename string) *MetaCSV {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	metacsv := &MetaCSV{File: file, FileName: filename}
	metacsv.File.WriteString("type,verticies,number_properties,size_json,time_read_json,time_write_json,size_buf,time_read_buf,time_write_buf\n")
	return metacsv
}

// adds a single feature to the meta csv
func (metacsv *MetaCSV) AddMeta(feature *geojson.Feature) {
	featurestring := MakeMeta(feature).MakeString()
	metacsv.File.WriteString(featurestring)
}

// creates a metadata structure
func MakeMeta(feature *geojson.Feature) *Meta {
	meta := &Meta{}
	meta.Type = string(feature.Geometry.Type)
	var total int
	switch meta.Type {
	case "Point":
		meta.Verticies = 1
	case "MultiPoint":
		meta.Verticies = len(feature.Geometry.MultiPoint)
	case "LineString":
		meta.Verticies = len(feature.Geometry.LineString)
	case "MultiLineString":
		for _, line := range feature.Geometry.MultiLineString {
			total += len(line)
		}
		meta.Verticies = total
	case "Polygon":
		for _, line := range feature.Geometry.Polygon {
			total += len(line)
		}
		meta.Verticies = total
	case "MultiPolygon":
		for _, polygon := range feature.Geometry.MultiPolygon {
			for _, line := range polygon {
				total += len(line)
			}
		}
		meta.Verticies = total
	}
	meta.Properties = len(feature.Properties)
	s := time.Now()
	bytevals, err := feature.MarshalJSON()
	meta.TimeWriteJSON = int(time.Now().Sub(s).Nanoseconds())
	if err != nil {
		fmt.Println(err)
	}
	meta.SizeJSON = len(bytevals)
	s = time.Now()
	_, err = geojson.UnmarshalFeature(bytevals)
	meta.TimeReadJSON = int(time.Now().Sub(s).Nanoseconds())
	if err != nil {
		fmt.Println(err)
	}
	s = time.Now()
	bytevals = raw.WriteFeature(feature)
	meta.TimeWriteBUF = int(time.Now().Sub(s).Nanoseconds())
	meta.SizeBUF = len(bytevals)
	s = time.Now()
	feature = raw.ReadFeature(bytevals)
	meta.TimeReadBUF = int(time.Now().Sub(s).Nanoseconds())
	return meta
}

// creates a meta data csv from a geobuf reader with fields
// type,verticies,number_properties,size_json,time_read_json,time_write_json,size_buf,time_read_buf,time_write_buf
func CreateMetaCSV(buf *g.Reader, outfilecsv string) {
	outcsv := NewMetaCSV(outfilecsv)
	i := 0
	for buf.Next() {
		outcsv.AddMeta(buf.Feature())
		i++
		if i%1000 == 0 {
			fmt.Printf("\rTotal Number of Meta Data Features Completed %d", i)
		}
	}

}
