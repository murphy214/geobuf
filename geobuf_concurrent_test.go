package geobuf

import (
	//"fmt"
	ld "github.com/murphy214/ld-geojson"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"os"
	"testing"
)

// creating the two test files
func I() int {
	ld.Convert_FeatureCollection("test_data/wv.geojson", "test_data/wv_ld.geojson")
	return 0
}

var _ = I()

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

// benchmarks line delimitted geojson read
func BenchmarkLineDelimitedGeojsonRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		ldjson := ld.Read_LD_Geojson("test_data/wv_ld.geojson")
		for ldjson.Next() {
			ldjson.Feature()
		}
	}
	os.Remove("test_data/wv_ld.geojson")
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

// benchmarks a geobuf geojson read
func BenchmarkGeobufConcurrentRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		buf := ReaderFile("test_data/wv.geobuf")
		buff := NewGeobufReaderConcurrent(buf)
		for buff.Next() {
			buff.Feature()
		}
	}
}
