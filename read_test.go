package geobuf_new

import (
	"testing"
	"github.com/golang/protobuf/proto"
	geo "./geobuf"
	"fmt"
)


var new_feat = Make_Feature(feat)
var bytevals,_ = proto.Marshal(new_feat)


func Benchmark_Read_Feature(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		Read_Feature(new_feat)
	}
}


func Benchmark_Read_Feature_Raw(b *testing.B) {
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