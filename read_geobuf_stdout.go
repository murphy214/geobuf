package geobuf

import (
	"github.com/murphy214/geobuf_raw"
	"fmt"
	//"encoding/csv"
	"github.com/paulmach/go.geojson"
	""
	"io"
	"os"
	"strings"
	//"time"
	//"io/ioutil"
)

var toptags = []string{"building", "source", "highway", "addr:housenumber", "addr:street", "name", "addr:city", "addr:postcode", "natural", "addr:country", "source:date", "landuse", "surface", "power", "waterway", "start_date", "tiger:cfcc", "tiger:county", "amenity", "oneway", "tiger:reviewed", "wall", "created_by", "building:levels", "ref", "ref:bag", "maxspeed", "height", "barrier", "service", "tiger:name_base", "lanes", "attribution", "access", "tiger:name_type", "source:addr", "addr:place", "type", "ele", "layer", "tracktype", "place", "tiger:tlid", "tiger:source", "leisure", "tiger:upload_uuid", "foot", "railway", "bicycle", "operator", "tiger:zip_left", "addr:suburb", "yh:WIDTH", "tiger:zip_right", "bridge", "tiger:separated", "addr:conscriptionnumber", "addr:state", "shop", "addr:city:simc", "note", "lacounty:bld_id", "lacounty:ain", "ref:ruian:building", "source_ref", "lit", "yh:STRUCTURE", "yh:TYPE", "building:units", "name:en", "addr:province", "building:ruian:type", "yh:TOTYUMONO", "yh:WIDTH_RANK", "man_made", "osak:identifier", "osak:municipality_no", "osak:revision", "osak:street_no", "is_in", "ref:ruian:addr", "leaf_type", "addr:interpolation", "NHD:FCode", "NHD:ComID", "public_transport", "NHD:ReachCode", "intermittent", "roof:shape", "boundary", "tourism", "crossing", "tunnel", "building:flats", "addr:street:sym_ul", "NHD:RESOLUTION", "width", "gauge", "water", "entrance", "import", "website", "admin_level", "sport", "nhd:reach_code", "electrified", "NHD:way_id", "NHD:FType", "footway", "nhd:com_id", "tiger:name_direction_prefix", "wheelchair", "source:geometry", "sidewalk", "voltage", "fixme", "source:maxspeed", "smoothness", "description", "network", "opening_hours", "gnis:feature_id", "phone", "building:material", "tiger:name_base_1", "wikidata", "nycdoitt:bin", "nhd:fdate", "parking", "bus", "gnis:fcode", "religion", "emergency", "wikipedia", "leaf_cycle", "gnis:ftype", "ref:linz:address_id", "frequency", "motor_vehicle", "species", "name:ru", "source:name", "area", "is_in:state", "horse", "historic", "usage", "restriction", "raba:id", "name_1", "alt_name", "is_in:country", "gnis:created", "material", "LINZ:source_version", "addr:streetnumber", "is_in:state_code", "chicago:building_id", "osak:street_name", "cycleway", "denotation", "roof:material", "gnis:county_id", "wetland", "gnis:state_id", "fire_hydrant:type", "osak:municipality_name", "LINZ:layer", "osak:house_no", "LINZ:dataset", "addr:full", "addr:district", "NHD:FDate", "shelter", "NHD:FTYPE", "roof:colour", "postal_code", "note:ja", "building:use", "osak:subdivision", "source:conscriptionnumber", "cuisine", "addr:street:name", "route", "addr:street:type", "building:part", "it:fvg:ctrn:code", "it:fvg:ctrn:revision", "ref:ruian", "junction", "denomination", "hgv", "source:position", "noexit", "KSJ2:curve_id", "information"}

// BoundingBox implementation as per https://tools.ietf.org/html/rfc7946
// BoundingBox syntax: "bbox": [west, south, east, north]
// BoundingBox defaults "bbox": [-180.0, -90.0, 180.0, 90.0]
func BoundingBox_Points(pts [][]float64) []float64 {
	// setting opposite default values
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	for _, pt := range pts {
		x, y := pt[0], pt[1]
		// can only be one condition
		// using else if reduces one comparison
		if x < west {
			west = x
		}
		if x > east {
			east = x
		}

		if y < south {
			south = y
		}
		if y > north {
			north = y
		}
	}
	return []float64{west, south, east, north}
}

func Push_Two_BoundingBoxs(bb1 []float64, bb2 []float64) []float64 {
	// setting opposite default values
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	// setting bb1 and bb2
	west1, south1, east1, north1 := bb1[0], bb1[1], bb1[2], bb1[3]
	west2, south2, east2, north2 := bb2[0], bb2[1], bb2[2], bb2[3]

	// handling west values: min
	if west1 < west2 {
		west = west1
	} else {
		west = west2
	}

	// handling south values: min
	if south1 < south2 {
		south = south1
	} else {
		south = south2
	}

	// handling east values: max
	if east1 > east2 {
		east = east1
	} else {
		east = east2
	}

	// handling north values: max
	if north1 > north2 {
		north = north1
	} else {
		north = north2
	}

	return []float64{west, south, east, north}
}

// this functions takes an array of bounding box objects and
// pushses them all out
func Expand_BoundingBoxs(bboxs [][]float64) []float64 {
	bbox := bboxs[0]
	for _, temp_bbox := range bboxs[1:] {
		bbox = Push_Two_BoundingBoxs(bbox, temp_bbox)
	}
	return bbox
}

// boudning box on a normal point geometry
// relatively useless
func BoundingBox_PointGeometry(pt []float64) []float64 {
	return []float64{pt[0], pt[1], pt[0], pt[1]}
}

// Returns BoundingBox for a MultiPoint
func BoundingBox_MultiPointGeometry(pts [][]float64) []float64 {
	return BoundingBox_Points(pts)
}

// Returns BoundingBox for a LineString
func BoundingBox_LineStringGeometry(line [][]float64) []float64 {
	return BoundingBox_Points(line)
}

// Returns BoundingBox for a MultiLineString
func BoundingBox_MultiLineStringGeometry(multiline [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, line := range multiline {
		bboxs = append(bboxs, BoundingBox_Points(line))
	}
	return Expand_BoundingBoxs(bboxs)
}

// Returns BoundingBox for a Polygon
func BoundingBox_PolygonGeometry(polygon [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, cont := range polygon {
		bboxs = append(bboxs, BoundingBox_Points(cont))
	}
	return Expand_BoundingBoxs(bboxs)
}

// Returns BoundingBox for a Polygon
func BoundingBox_MultiPolygonGeometry(multipolygon [][][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, polygon := range multipolygon {
		for _, cont := range polygon {
			bboxs = append(bboxs, BoundingBox_Points(cont))
		}
	}
	return Expand_BoundingBoxs(bboxs)
}

// Returns a BoundingBox for a geometry collection
func BoundingBox_GeometryCollection(gs []*geojson.Geometry) []float64 {
	bboxs := [][]float64{}
	for _, g := range gs {
		bboxs = append(bboxs, geobuf_raw.Get_BoundingBox(g))
	}
	return Expand_BoundingBoxs(bboxs)
}

// retrieves a boundingbox for a given geometry
func GetBoundingBox(g *geojson.Geometry) []float64 {
	switch g.Type {
	case "Point":
		return BoundingBox_PointGeometry(g.Point)
	case "MultiPoint":
		return BoundingBox_MultiPointGeometry(g.MultiPoint)
	case "LineString":
		return BoundingBox_LineStringGeometry(g.LineString)
	case "MultiLineString":
		return BoundingBox_MultiLineStringGeometry(g.MultiLineString)
	case "Polygon":
		return BoundingBox_PolygonGeometry(g.Polygon)
	case "MultiPolygon":
		return BoundingBox_MultiPolygonGeometry(g.MultiPolygon)

	}
	return []float64{}
}

func GetKeys(buf *Reader) ([]string, int) {
	keymap := map[string]string{}
	totalkeys := []string{}
	i := 0
	for buf.Next() {
		keys := ReadKeys(buf.Bytes())
		for _, key := range keys {
			_, boolval := keymap[key]
			if !boolval {
				keymap[key] = ""
				totalkeys = append(totalkeys, key)
			}
		}
		i++
	}
	totalkeys = append(totalkeys, []string{"Bounds", "Type", "Geometry"}...)
	buf.Reset()
	return totalkeys, i
}

func WriteRow(feature *geojson.Feature, keys []string) {
	bounds := GetBoundingBox(feature.Geometry)
	feature.Properties["Bounds"] = fmt.Sprintf("%f,%f,%f,%f", bounds[0], bounds[1], bounds[2], bounds[3])
	feature.Properties["Type"] = string(feature.Geometry.Type)
	s, _ := feature.Geometry.MarshalJSON()

	feature.Properties["Geometry"] = string(s)
	newrow := make([]string, len(keys))
	for pos, key := range keys {
		val, boolval := feature.Properties[key]
		//fmt.Println(fmt.Sprint(val), val, feature.Properties)

		if !boolval {
			val = ""
		}
		newrow[pos] = fmt.Sprint(val)
	}
	io.WriteString(os.Stdout, strings.Join(newrow, "|")+"\n")
}

func ReadGeobufCSV(filename string) {
	buf := ReaderFile(filename)
	//keys, _ := GetKeys(buf)
	keys := append(toptags[:50], []string{"Bounds", "Type", "Geometry"}...)
	io.WriteString(os.Stdout, strings.Join(keys, "|")+"\n")
	myfunc := func(feature *geojson.Feature) interface{} {
		WriteRow(feature, keys)
		return ""
	}

	for buf.Next() {
		myfunc(buf.Feature())
	}
}
