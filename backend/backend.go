package backend

import (
	m "github.com/murphy214/mercantile"
	g "github.com/murphy214/geobuf"
	//"github.com/murphy214/paulmach"
	"bytes"
	"fmt"
	"encoding/json"
	"strconv"
	"net/url"
	"github.com/murphy214/geobuf/geobuf_raw"
)

// the backend is an abstraction structure for doign things about the wjhoel blcok
// this really isnt the best for like bulk processnig just data visualization and analysis ideally
type Backend struct {
	Geobuf *g.Reader
	Limit int // defaults to 1 million
	TileBool bool // whether or not the data is tiled
	Zoom int // the zoom the data is tiled at 	
}

// This structure parses a map context and returns a filter
type MapContext struct {
	Bounds m.Extrema
	Zoom float64
	Filter *FeatureFilter
}

// adds the underlying bounds to the mapcontext filter
func (mapcontext *MapContext) AddBounds() *MapContext {
	newfilt := &FeatureFilter{
		Bounds:mapcontext.Bounds,
		BoundsBool:true,
		Key:"$geometry",
		Operator:Intersects,
	}
	mapcontext.Filter = &FeatureFilter{
		Operator:AndOperator,
		Filters:[]*FeatureFilter{newfilt,mapcontext.Filter},
	}
	return mapcontext
}

// parses a given string query from a url 
func ParseQuery(stringval string) *MapContext {
	val := &MapContext{}
	mv, err := url.ParseQuery(stringval)
	if err != nil {
		fmt.Println(err)
	}

	// getting the filter string
	filterstr := mv["filter"][0]
	var mm []interface{}
	err = json.Unmarshal([]byte(filterstr),&mm)
	if err != nil {
		fmt.Println(err)
	}
	val.Filter = ParseAll(mm)
	
	// getting the bounds string
	bds_str := mv["bounds"][0]
	var vv []float64
	err = json.Unmarshal([]byte(bds_str), &vv)
	if err != nil {
		fmt.Println(err)
	}
	if len(vv) == 4 {
		w,s,e,n := vv[0],vv[1],vv[2],vv[3]
		val.Bounds = m.Extrema{N:n,S:s,E:e,W:w}
	}	
	
	// getting the zoom string
	zoomstr := mv["zoom"][0]
	zoom,err := strconv.ParseFloat(zoomstr,64)
	if err != nil {
		fmt.Println(err)
	}
	val.Zoom = zoom

	return val.AddBounds()
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