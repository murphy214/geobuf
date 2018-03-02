package geobuf_new

import (
	"testing"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
)

func Benchmark_Read_FeatureCollection_Old(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		bytevals,_ := ioutil.ReadFile("test_data/county.geojson")
		geojson.UnmarshalFeatureCollection(bytevals)
	}
}


func Benchmark_Read_FeatureCollection_New(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		geobuf := Reader_File("test_data/county.geobuf")
		for geobuf.Next() {
			geobuf.Feature()
		}
	}
}