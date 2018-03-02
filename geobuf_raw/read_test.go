package geobuf_raw

import (
	"testing"
	"github.com/golang/protobuf/proto"
	geo "./geobuf"
	"fmt"
	"github.com/paulmach/go.geojson"
)


var new_feat = Make_Feature(feat)
var bytevals,_ = proto.Marshal(new_feat)
var bytevals2,_ = feat.MarshalJSON()


func Benchmark_Read_Feature(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		Read_Feature(new_feat)
	}
}


func Benchmark_Read_Feature_Old(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_,err := geojson.UnmarshalFeature(bytevals2)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Benchmark_Read_Feature_New(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		feature := &geo.Feature{}
		err := proto.Unmarshal(bytevals,feature)
		if err != nil {
			fmt.Println(err)
		}
		Read_Feature(feature)
	}
}