package geobuf

import (
	"bufio"
	"github.com/murphy214/geobuf/geobuf_raw"
	"os"
	//"io"
	"bytes"
	"fmt"
	"io/ioutil"
	//"github.com/murphy214/protoscan"
	//"github.com/golang/protobuf/proto"
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
	Filename  string
	Writer    *bufio.Writer
	FileBool  bool
	Buffer    *bytes.Buffer
	File      *os.File
	Bytesvals []byte
}

var writersize = 64 * 4096

// creates a writer struct
func WriterFileNew(filename string) *Writer {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	return &Writer{Filename: filename, FileBool: true, File: file, Bytesvals: []byte{}}
}

// creates a writer struct
func WriterFile(filename string) *Writer {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	return &Writer{Filename: filename, FileBool: true, File: file}
}

// creates a writer buffer new
func WriterBufNew() *Writer {
	b := bytes.NewBuffer([]byte{})
	return &Writer{Writer: bufio.NewWriterSize(b, writersize), Buffer: b, FileBool: false}
}

// creates a writer buffer
func WriterBuf(bytevals []byte) *Writer {
	buffer := bytes.NewBuffer(bytevals)
	return &Writer{Writer: bufio.NewWriterSize(buffer, writersize), Buffer: buffer, FileBool: false}
}

// writing feature
func (writer *Writer) WriteFeature(feature *geojson.Feature) {
	bytevals := geobuf_raw.WriteFeature(feature)

	// writing the appended bytevals to the writer

	bytevals = append(
		append(
			[]byte{10}, EncodeVarint(uint64(len(bytevals)))...,
		),
		bytevals...)
	if writer.FileBool {
		writer.File.Write(bytevals)
	} else {
		writer.Writer.Write(bytevals)
	}

}

// writes a set of byte values representing a feature
// to the underlying writer
func (writer *Writer) Write(bytevals []byte) {
	bytevals = append(
		append(
			[]byte{10}, EncodeVarint(uint64(len(bytevals)))...,
		),
		bytevals...)
	if writer.FileBool {
		writer.File.Write(bytevals)
	} else {
		writer.Writer.Write(bytevals)
	}

}

// writes a set of byte values representing a feature
// to the underlying writer
func (writer *Writer) WriteRaw(bytevals []byte) {
	if writer.FileBool {
		writer.File.Write(bytevals)
	} else {
		writer.Writer.Write(bytevals)
	}

}

// adds a geobuf buffer value to an existing geobuf
func (writer *Writer) AddGeobuf(buf *Writer) {
	writer.Writer.Flush()
	if !buf.FileBool {
		buf.Writer.Flush()
		//buf.Writer = bufio.NewWriter(buf.Buffer)
		if writer.FileBool {
			writer.File.Write(buf.Buffer.Bytes())
		} else {
			writer.Writer.Write(buf.Buffer.Bytes())
		}
	}
}

// returns the bytes present in an underlying
// writer type buffer
func (writer *Writer) Bytes() []byte {
	writer.Writer.Flush()

	//
	if !writer.FileBool {
		return writer.Buffer.Bytes()

	} else {
		writer.File.Close()
		bytevals, _ := ioutil.ReadFile(writer.Filename)
		return bytevals
	}
	return []byte{}
}

// converts a writer into a reader
func (writer *Writer) Reader() *Reader {
	writer.Writer.Flush()

	if !writer.FileBool {
		newreader := ReaderBuf(writer.Bytes())
		return newreader
	} else {
		writer.File.Close()
		newreader := ReaderFile(writer.Filename)
		return newreader
	}
	return &Reader{}
}
