package geobuf

import (
	"bufio"
	"github.com/murphy214/geobuf/geobuf_raw"
	"os"
	//"io"
	"bytes"
	"encoding/gob"
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/pbf"
	"github.com/murphy214/protoscan"
	"github.com/paulmach/go.geojson"
)

// protobuf scanner implementation
type Reader struct {
	FileBool bool
	Reader   *protoscan.ProtobufScanner
	Filename string
	File     *os.File
	Buf      []byte
	MetaData MetaData
}

// sub file contained within a geobuf
type SubFile struct {
	Positions      [2]int
	NumberFeatures int
	Size           int
}

// struct for handling metadata
type MetaData struct {
	FileSize       int
	NumberFeatures int
	Files          map[string]*SubFile
	Bounds         m.Extrema
}

// lints metadata
func (metadata *MetaData) LintMetaData(pos int) {
	for _, v := range metadata.Files {
		v.Positions = [2]int{v.Positions[0] + pos, v.Positions[1] + pos}
		v.Size = v.Positions[1] - v.Positions[0]
	}
}

// creates a reader for a byte array
func ReaderBuf(bytevals []byte) *Reader {
	buffer := bytes.NewReader(bytevals)
	buf := &Reader{Reader: protoscan.NewProtobufScanner(buffer), Buf: bytevals, FileBool: false}
	buf.CheckMetaData()
	return buf
}

// creates a reader for file
func ReaderFile(filename string) *Reader {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(file)

	buf := &Reader{
		Reader:   protoscan.NewProtobufScanner(reader),
		Filename: filename,
		FileBool: true,
		File:     file,
	}
	buf.CheckMetaData()
	return buf
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
func (reader *Reader) BytesIndicies() ([]byte, [2]int) {
	return reader.Reader.ProtobufIndicies()
}

// alias for the Protobuf() method
// again more expressive for our use case
func (reader *Reader) Feature() *geojson.Feature {
	return geobuf_raw.ReadFeature(reader.Bytes())
}

// alias for the Protobuf() method
// again more expressive for our use case
func (reader *Reader) FeatureIndicies() (*geojson.Feature, [2]int) {
	bytevals, indicies := reader.BytesIndicies()
	return geobuf_raw.ReadFeature(bytevals), indicies
}

// reads a single feature form bytes
func ReadFeature(bytevals []byte) *geojson.Feature {
	return geobuf_raw.ReadFeature(bytevals)
}

// reads a feature
func ReadKeys(bytevals []byte) []string {
	pbfval := pbf.PBF{Pbf: bytevals, Length: len(bytevals)}
	keys := []string{}
	key, val := pbfval.ReadKey()
	if key == 1 && val == 0 {
		pbfval.ReadVarint()
		key, val = pbfval.ReadKey()
	}
	for key == 2 && val == 2 {
		// starting properties shit here

		size := pbfval.ReadVarint()
		endpos := pbfval.Pos + size
		//pbfval.ReadKey()
		pbfval.Pos += 1
		keys = append(keys, pbfval.ReadString())

		pbfval.Pos = endpos
		key, val = pbfval.ReadKey()
	}

	return keys
}

// reads a feature
func ReadBoundingBox(bytevals []byte) []float64 {
	pos := len(bytevals) - 1
	alloc := make([]byte, 32)
	allocpos := 31
	boolval := true
	for boolval {
		alloc[allocpos] = bytevals[pos]
		if bytevals[pos] == 42 {
			boolval = false
		}
		pos--
		allocpos--
	}

	bb := make([]float64, 4)
	pbfval := pbf.NewPBF(alloc[allocpos+3:])
	bb[0] = float64(pbfval.ReadSVarintPower())
	bb[1] = float64(pbfval.ReadSVarintPower())
	bb[2] = float64(pbfval.ReadSVarintPower())
	bb[3] = float64(pbfval.ReadSVarintPower())
	return bb
}

func (reader *Reader) ReadAll() []*geojson.Feature {
	feats := []*geojson.Feature{}
	for reader.Next() {
		feats = append(feats, reader.Feature())
	}
	return feats
}

// resets a reader to be read again
func (reader *Reader) Reset() {
	if reader.FileBool {
		file, err := os.Open(reader.Filename)
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

// reads an indicies ready to append
func (reader *Reader) ReadIndAppend(inds [2]int) []byte {
	inds[0] = inds[0] - len(EncodeVarint(uint64(inds[1]-inds[0]))) - 1
	bytevals := make([]byte, inds[1]-inds[0])
	reader.File.ReadAt(bytevals, int64(inds[0]))
	return bytevals
}

// read feature from an indicies
func (reader *Reader) ReadIndFeature(inds [2]int) *geojson.Feature {
	bytevals := make([]byte, inds[1]-inds[0])
	reader.File.ReadAt(bytevals, int64(inds[0]))
	return ReadFeature(bytevals)
}

// this functions types into the underlying protoscan implementation
// and reconfigures the protoscan to start reading a certain position
func (reader *Reader) Seek(pos int) {
	if reader.FileBool {
		reader.File.Seek(int64(pos), 0)
		myreader := bufio.NewReader(reader.File)
		reader.Reader = protoscan.NewProtobufScanner(myreader)
		reader.Reader.TotalPosition = pos
	} else {
		buffer := bytes.NewReader(reader.Buf)
		buffer.Seek(int64(pos), 0)
		reader.Reader = protoscan.NewProtobufScanner(buffer)
		reader.Reader.TotalPosition = pos
	}
}

// reads the metadata from a raw bytes set
func WriteMetaData(meta MetaData) string {
	bb := bytes.NewBuffer([]byte{})
	dec := gob.NewEncoder(bb)
	err := dec.Encode(meta)
	if err != nil {
		fmt.Println(err)
	}
	return string(bb.Bytes())
}

// reads the metadata from a raw bytes set
func ReadMetaData(bytevals []byte) MetaData {
	dec := gob.NewDecoder(bytes.NewBuffer(bytevals))
	var q MetaData
	err := dec.Decode(&q)
	if err != nil {
		fmt.Println(err)
	}
	return q
}

// checks for metadata
func (reader *Reader) CheckMetaData() {
	reader.Next()
	feature := reader.Feature()
	// if the metadata feature exists
	if len(feature.Geometry.Point) == 5 {
		bytevals := []byte(feature.Properties["metadata"].(string))
		reader.MetaData = ReadMetaData(bytevals)
		reader.MetaData.LintMetaData(reader.Reader.TotalPosition)
	} else {
		reader.Reset()
	}

}
