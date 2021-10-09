package geobuf

import (
	//"fmt"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"testing"
)


/*
All benchmarks are on a 178 mb geojson roads file of west virginia.byte
The differing filenames are just there different file representations.

*/

// benchmarks a normal geojson deserireadalization
func BenchmarkFeatureCollectionRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		bytevals, _ := ioutil.ReadFile("test_data/wv.geojson")
		_, _ = geojson.UnmarshalFeatureCollection(bytevals)
	}
}



// benchmarks a geobuf geojson read
func BenchmarkGeobufRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		buf := ReaderFile("test_data/wv.geobuf")
		for buf.Next() {
			buf.Feature()
		}
	}
}
