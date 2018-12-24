package backend

import (
	g "github.com/murphy214/geobuf"
	//"github.com/murphy214/paulmach"
	"bytes"
	"fmt"
	"github.com/murphy214/geobuf/geobuf_raw"
)

// the backend is an abstraction structure for doign things about the wjhoel blcok
// this really isnt the best for like bulk processnig just data visualization and analysis ideally
type Backend struct {
	Geobuf *g.Reader
	Limit int // defaults to 1 million
}


// creates a new backend
func NewBackend(reader *g.Reader) *Backend {
	return &Backend{Geobuf:reader,Limit:1000000}
}

func addtobuf(newlist [][]byte,mybuf *bytes.Buffer) {
	c := make(chan string,len(newlist))
	for _,bs := range newlist {
		go func(bs []byte,c chan string) {
			bs,_ = geobuf_raw.ReadFeature(bs).MarshalJSON()
			c <- string(bs)
		}(bs,c)
	}
	for range newlist {
		out := <- c
		mybuf.WriteString(out+",")
	}
}

// converts a geobuf object to a raw string 
func convertgeobuf(buf *g.Reader) string {
	mybuf := bytes.NewBuffer([]byte(`{"type": "FeatureCollection", "features": [`))

	buf.Next()
	if buf.Next() {
		feature := buf.Feature()
		bs,err := feature.MarshalJSON()
		if err != nil {
			fmt.Println(err)
		}
		mybuf.WriteString(string(bs)+",")
	}
	increment := 10000
	newlist := [][]byte{}
	for buf.Next() {
		newlist = append(newlist,buf.Bytes())
		if len(newlist) == increment {
			/*
			c := make(chan []byte,len(newlist))
			for _,bs := range newlist {
				go func(bs []byte,c chan []byte) {
					bs,_ := geobuf_raw.ReadFeature(bs).MarshalJSON()
					c <- bs
				}(bs,c)
			}
			for range newlist {
				out := <- c
				mybuf.WriteString(out+",")
			}
			*/

			addtobuf(newlist, mybuf)
			newlist = [][]byte{}
		}
	}

	if len(newlist) != 0 {	
		lastindex := len(newlist)-1
		addtobuf(newlist[:lastindex], mybuf)

		// getting last feature
		lastfeaturebs := newlist[lastindex]
		lastfeature := geobuf_raw.ReadFeature(lastfeaturebs)
		lastfeaturebs,_ = lastfeature.MarshalJSON()
		mybuf.WriteString(string(lastfeaturebs)+"]}")
	}
	return mybuf.String()
}



// filters a backend entirely
func (backend *Backend) FilterBackend(featurefilter *FeatureFilter) *g.Reader {
	bufw := g.WriterBufNew()
	backend.Geobuf.Reset()

	// iterating through each feature applygint the filter
	i := 0
	for backend.Geobuf.Next() {
		feature := backend.Geobuf.Feature()
		if featurefilter.Filter(feature) && i < backend.Limit {
			bufw.WriteFeature(feature)
			i++
		}
	}

	return bufw.Reader()
}

// filters the backend and returns a geojson string
func (backend *Backend) FilterBackendGeoJSON(featurefilter *FeatureFilter) string {
	return convertgeobuf(backend.FilterBackend(featurefilter))
}