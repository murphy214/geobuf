package geobuf_raw

import (
	//"github.com/murphy214/geobuf_new"
	"github.com/murphy214/geobuf/geobuf_raw"
	"github.com/golang/protobuf/proto"
	//"fmt"
	"testing"
	"github.com/paulmach/go.geojson"

)

var bytes2 = []byte{24, 45, 68, 84, 251, 33, 9, 64,0}
var pbf = &PBF{Pbf:bytes2,Length:8}


var bytes3 = []byte{24, 45, 68, 84, 251, 33, 9, 64}
var pbf2 = &PBF{Pbf:bytes3,Length:8}


var feature_s = `{"id":1000001,"type":"Feature","bbox":[-83.647031,33.698307,-83.275933,33.9659119],"geometry":{"type":"MultiPolygon","coordinates":[[[[-83.537385,33.9659119],[-83.5084519,33.931233],[-83.4155119,33.918541],[-83.275933,33.847977],[-83.306619,33.811444],[-83.28034,33.7617739],[-83.29145,33.7343149],[-83.406189,33.698307],[-83.479523,33.802265],[-83.505928,33.81776],[-83.533165,33.820923],[-83.647031,33.9061979],[-83.537385,33.9659119]]],[[[-83.537385,33.9659119],[-83.5084519,33.931233],[-83.4155119,33.918541],[-83.275933,33.847977],[-83.306619,33.811444],[-83.28034,33.7617739],[-83.29145,33.7343149],[-83.406189,33.698307],[-83.479523,33.802265],[-83.505928,33.81776],[-83.533165,33.820923],[-83.647031,33.9061979],[-83.537385,33.9659119]]],[[[-83.537385,33.9659119],[-83.5084519,33.931233],[-83.4155119,33.918541],[-83.275933,33.847977],[-83.306619,33.811444],[-83.28034,33.7617739],[-83.29145,33.7343149],[-83.406189,33.698307],[-83.479523,33.802265],[-83.505928,33.81776],[-83.533165,33.820923],[-83.647031,33.9061979],[-83.537385,33.9659119]]]]},"properties":{"AREA":"13219","COLORKEY":"#03E174","area":"13219","index":1109}}`
var feature,_ = geojson.UnmarshalFeature([]byte(feature_s))
var feature_geo = geobuf_raw.MakeFeature(feature)
var bytevals,_ = proto.Marshal(feature_geo)


func Benchmark_Write_Old(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		f := geobuf_raw.MakeFeature(feature)
		proto.Marshal(f)
	}
}


func Benchmark_Write_New(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		WriteFeature(feature)
		//fmt.Println(string(bytevals))
	}
}

func Benchmark_Read_Old(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		f := &geo.Feature{}
		proto.Unmarshal(bytevals,f)
		geobuf_raw.ReadFeature(f)
	}
}


func Benchmark_Read_New(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		ReadFeature(bytevals)
		//fmt.Println(string(bytevals))
	}
}

