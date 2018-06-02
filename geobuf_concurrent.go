package geobuf

import (
	//"fmt"
	"github.com/murphy214/geobuf/geobuf_raw"
	"github.com/paulmach/go.geojson"
	"sync"
	//"time"
)

// the data structure for mapping a function against a concurrent reader
type Reader_Concurrent_Function struct {
	Count      int
	BoolVal    bool
	StartBool  bool
	Channel    chan interface{}
	Reader     *Reader
	Pos        int
	TotalCount int
}

func (reader *Reader_Concurrent_Function) Next() bool {
	if reader.TotalCount == reader.Pos+1 && !reader.StartBool && !reader.BoolVal {
		return false
	}

	return ((reader.StartBool || reader.Count != 0) || reader.BoolVal) || (reader.TotalCount > reader.Pos && !reader.StartBool)
}

func (reader *Reader_Concurrent_Function) Increment() bool {
	reader.BoolVal = reader.Reader.Next()
	reader.TotalCount++
	reader.Count++

	return reader.BoolVal
}

func (reader *Reader_Concurrent_Function) Defer() {
	reader.Pos++
	reader.Count--
}

// a masking function
func Mask(val interface{}) interface{} {
	return val
}

// this is the mapped function that will be returned
type MapFunction func(feature *geojson.Feature) interface{}

// this function instantiates the concurrent reader
func NewGeobufReaderFunction(geobuf *Reader, myfunc MapFunction) *Reader_Concurrent_Function {
	newreader := &Reader_Concurrent_Function{Reader: geobuf, Channel: make(chan interface{}), BoolVal: true, StartBool: true}
	var wg sync.WaitGroup
	go func() {
		newreader.StartBool = false

		for newreader.Increment() {
			wg.Add(1)

			newreader.BoolVal = !newreader.Reader.Reader.EndBool

			bytevals := newreader.Reader.Bytes()
			go func(bytevals []byte) {
				defer wg.Done()
				feature := geobuf_raw.ReadFeature(bytevals)
				newreader.Channel <- myfunc(feature)
			}(bytevals)
			if newreader.Count > 100 {
				wg.Wait()
			}

		}
		wg.Wait()

	}()

	return newreader
}

// this function will return the next available value from a function
func (reader *Reader_Concurrent_Function) Value() interface{} {
	reader.Defer()
	f := <-reader.Channel

	return f
}

// the data structure for a geojson feature concurrent reader
type Reader_Concurrent struct {
	Count      int
	BoolVal    bool
	StartBool  bool
	Channel    chan *geojson.Feature
	Reader     *Reader
	Pos        int
	TotalCount int
}

func (reader *Reader_Concurrent) Next() bool {
	if reader.TotalCount == reader.Pos+1 && !reader.StartBool && !reader.BoolVal {
		return false
	}
	return ((reader.StartBool || reader.Count != 0) || reader.BoolVal) || (reader.TotalCount > reader.Pos && !reader.StartBool)
}

func (reader *Reader_Concurrent) SubFileNext() bool {
	if reader.TotalCount == reader.Pos+1 && !reader.StartBool && !reader.BoolVal {
		return false
	}
	return ((reader.StartBool || reader.Count != 0) || reader.BoolVal) || (reader.TotalCount > reader.Pos && !reader.StartBool) && reader.Reader.Reader.TotalPosition < reader.Reader.SubFileEnd
}

func (reader *Reader_Concurrent) Increment() bool {
	reader.BoolVal = reader.Reader.Next()
	reader.TotalCount++
	reader.Count++

	return reader.BoolVal
}

func (reader *Reader_Concurrent) Defer() {
	reader.Pos++
	reader.Count--

}

// this function instantiates the concurrent reader
func NewGeobufReaderConcurrent(geobuf *Reader) *Reader_Concurrent {
	newreader := &Reader_Concurrent{Reader: geobuf, Channel: make(chan *geojson.Feature), BoolVal: true, StartBool: true}
	var wg sync.WaitGroup
	go func() {
		newreader.StartBool = false

		for newreader.Increment() {
			wg.Add(1)

			newreader.BoolVal = !newreader.Reader.Reader.EndBool

			bytevals := newreader.Reader.Bytes()
			go func(bytevals []byte) {
				defer wg.Done()
				newreader.Channel <- geobuf_raw.ReadFeature(bytevals)
			}(bytevals)
			if newreader.Count > 100 {
				wg.Wait()
			}

		}
		wg.Wait()

	}()

	return newreader
}

// this function will return the next available value from a function
func (reader *Reader_Concurrent) Feature() *geojson.Feature {
	reader.Defer()
	f := <-reader.Channel
	return f
}
