package geobuf_new

import (
	"./geobuf_raw"
	geo "./geobuf_raw/geobuf"
	"os"
	"bufio"
	//"io"
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/paulmach/go.geojson"
    "github.com/murphy214/protoscan"
)

// protobuf scanner implementation
type Reader struct {
	FileBool bool
	Reader *protoscan.ProtobufScanner
	Filename string
	Feature_Bytes []byte
}

// creates a reader for a byte array
func Reader_Bytes(bytevals []byte) *Reader {
	return &Reader{Reader:protoscan.NewProtobufScanner(bytes.NewReader(bytevals)),FileBool:false}
}

// creates a reader for file
func Reader_File(filename string) *Reader {
	file,err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	return &Reader{
			Reader:protoscan.NewProtobufScanner(bufio.NewReader(file)),
			Filename:filename,
			FileBool:true,
		}
}

// alias for the Scan method on reader
// next is a little more expressive
func (reader *Reader) Next() bool {
	return reader.Reader.Scan()
}

// alias for the Protobuf() method 
// again more expressive for our use case
func (reader *Reader) Bytes() []byte {
	return reader.Reader.Protobuf()
}

// alias for the Protobuf() method 
// again more expressive for our use case
func (reader *Reader) Feature() *geojson.Feature {
	feature := &geo.Feature{}
	err := proto.Unmarshal(reader.Bytes(),feature)
	if err != nil {
		fmt.Println(err)
	}
	return geobuf_raw.Read_Feature(feature)
}


