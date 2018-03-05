package geobuf

import (
	"github.com/murphy214/geobuf/geobuf_raw"
	"os"
	"bufio"
	//"io"
	"bytes"
	"fmt"
	"github.com/paulmach/go.geojson"
    "github.com/murphy214/protoscan"
)

// protobuf scanner implementation
type Reader struct {
	FileBool bool
	Reader *protoscan.ProtobufScanner
	Filename string
	File *os.File
	Buf []byte
}

// creates a reader for a byte array
func ReaderBuf(bytevals []byte) *Reader {
	buffer := bytes.NewReader(bytevals)
	return &Reader{Reader:protoscan.NewProtobufScanner(buffer),Buf:bytevals,FileBool:false}
}

// creates a reader for file
func ReaderFile(filename string) *Reader {
	file,err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(file)
	return &Reader{
			Reader:protoscan.NewProtobufScanner(reader),
			Filename:filename,
			FileBool:true,
			File:file,
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
	return geobuf_raw.ReadFeature(reader.Bytes())
}


func ReadFeature(bytevals []byte) *geojson.Feature {
	return geobuf_raw.ReadFeature(bytevals)
}

func (reader *Reader) ReadAll() []*geojson.Feature {
	feats := []*geojson.Feature{}
	for reader.Next() {
		feats = append(feats,reader.Feature())
	}
	return feats
}

// resets a reader to be read again
func (reader *Reader) Reset() {
	if reader.FileBool {
		file,err := os.Open(reader.Filename)
		if err != nil {
			fmt.Println(err)
		}
		read := bufio.NewReader(file)
		reader.Reader = protoscan.NewProtobufScanner(read)
	} else {
		buffer := bytes.NewReader(reader.Buf)
		reader.Reader = protoscan.NewProtobufScanner(buffer)
	}
}

