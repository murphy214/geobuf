package geobuf_new

import (
	"./geobuf_raw"
	"os"
	"bufio"
	//"io"
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/paulmach/go.geojson"
)

const maxVarintBytes = 10 // maximum Length of a varint

func EncodeVarint(x uint64) []byte {
	var buf [maxVarintBytes]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

// the writer struct
type Writer struct {
	Filename string
	Writer *bufio.Writer
	FileBool bool
	Buffer *bytes.Buffer
	File *os.File
}

// creates a writer struct
func Writer_File_New(filename string) Writer {
	file,err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	return Writer{Filename:filename,Writer:bufio.NewWriter(file),FileBool:true,File:file}
}

// creates a writer struct
func Writer_File(filename string) Writer {
	file,err := os.OpenFile(filename,os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	return Writer{Filename:filename,Writer:bufio.NewWriter(file),FileBool:true,File:file}
}

// creates a writer buffer new
func Writer_Buf_New() Writer {
    var b bytes.Buffer
    return Writer{Writer:bufio.NewWriter(&b),Buffer:&b,FileBool:false}
}	

// creates a writer buffer 
func Writer_Buf(bytevals []byte) Writer {
	buffer := bytes.NewBuffer(bytevals)
    return Writer{Writer:bufio.NewWriter(buffer),Buffer:buffer,FileBool:false}
}	

// writing feature
func (writer *Writer) Write_Feature(feature *geojson.Feature) {
	bytevals,err := proto.Marshal(geobuf_raw.Make_Feature(feature))
	if err != nil {
		fmt.Println(err)
	}
	
	// writing the appended bytevals to the writer

	bytevals = append(
						append(
							[]byte{10},EncodeVarint(uint64(len(bytevals)))...
						),
					bytevals...)
	if writer.FileBool {
		writer.File.Write(bytevals)
	} else {
		writer.Writer.Write(bytevals)
    	if err = writer.Writer.Flush(); err != nil {
        	fmt.Println(err)
    	}
    }
	//writer.Writer.Reset()
}



