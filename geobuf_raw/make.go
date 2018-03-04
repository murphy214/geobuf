package geobuf_raw

import (
	"github.com/paulmach/go.geojson"
	"math"
	geo "./geobuf"
	"reflect"
)

var powerfactor = math.Pow(10.0,7.0)

// converts a single pt
func Convert_Pt(pt []float64) []int64 {
	newpt := make([]int64,2)
	newpt[0] = int64(pt[0] * math.Pow(10.0,7.0))
	newpt[1] = int64(pt[1] * math.Pow(10.0,7.0))	
	return newpt
}

// param encoding
func paramEnc(value int64) uint64 {
	return uint64((value << 1) ^ (value >> 31))
}

// makes a given point
func Make_Point(pt []float64) []uint64 {
	point := Convert_Pt(pt)
	return []uint64{paramEnc(point[0]),paramEnc(point[1])}
}

// makes a line
func Make_Line(line [][]float64) ([]uint64,[]int64) {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	//oldpt := Convert_Pt(line[0])
	newline := make([]uint64,len(line)*2)
	deltapt := make([]int64,2)
	pt := make([]int64,2)
	oldpt := make([]int64,2)


	for i,point := range line {
		x,y := point[0],point[1]
		if x < west {
			west = x
		} else if x > east {
			east = x
		}

		if y < south {
			south = y
		} else if y > north {
			north = y
		}

		pt = Convert_Pt(point)
		if i == 0 {
			newline[0] = paramEnc(pt[0])
			newline[1] = paramEnc(pt[1])
		} else {
			deltapt = []int64{pt[0] - oldpt[0],pt[1] - oldpt[1]}
			newline[i*2] = paramEnc(deltapt[0])
			newline[i*2+1] = paramEnc(deltapt[1])
		}
		oldpt = pt
	}
	return newline,[]int64{int64(west * powerfactor),
	int64(south * powerfactor), 
	int64(east * powerfactor),
	int64(north * powerfactor),}
}

// creates a polygon 
func Make_Polygon(polygon [][][]float64) ([]uint64,[]int64) {
	geometry := []uint64{}
	bb := []int64{}
	for i,cont := range polygon {
		geometry = append(geometry,uint64(len(cont) * 2))

		tmpgeom,tmpbb := Make_Line(cont)
		geometry = append(geometry,tmpgeom...)
		if i == 0 {
			bb = tmpbb
		}
	}
	return geometry,bb
}

// creates a multi polygon array 
func Make_MultiPolygon(multipolygon [][][][]float64) ([]uint64,[]int64) {
	geometry := []uint64{}
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	west,south,east,north = west * powerfactor,south * powerfactor,east * powerfactor,north * powerfactor
	bb := []int64{int64(west),int64(south),int64(east),int64(north)}

	for _,polygon := range multipolygon {
		geometry = append(geometry,uint64(len(polygon)))
		tempgeom,tempbb := Make_Polygon(polygon)
		geometry = append(geometry,tempgeom...)
		if bb[0] > tempbb[0] {
			bb[0] = tempbb[0]
		}
		if bb[1] > tempbb[1] {
			bb[1] = tempbb[1]
		}
		if bb[2] < tempbb[2] {
			bb[2] = tempbb[2]
		}
		if bb[3] < tempbb[3] {
			bb[3] = tempbb[3]
		}
	}
	return geometry,bb
}	


// syntax wise geometry will only have basic types
// flat for liens no header
// flaot for points
// header for rings of polygon
func Make_Geom(geom *geojson.Geometry) ([]uint64,geo.GeomType,[]int64) {
	var geomtype geo.GeomType

	switch geom.Type {
	case "Point":
		newpt := Make_Point(geom.Point)
		newpt2 := Convert_Pt(geom.Point)
		return newpt,geo.GeomType_POINT,[]int64{newpt2[0],newpt2[1],newpt2[0],newpt2[1]}
	case "LineString":
		geometry,bb :=  Make_Line(geom.LineString)
		return geometry,geo.GeomType_LINESTRING,bb
	case "Polygon":
		geometry,bb :=  Make_Polygon(geom.Polygon)
		return geometry,geo.GeomType_POLYGON,bb
	case "MultiPoint":
		// multipoint code here
		geometry,bb := Make_Line(geom.MultiPoint)
		return geometry,geo.GeomType_MULTIPOINT,bb
	case "MultiLineString":
		// multi line string code here
		geometry,bb := Make_Polygon(geom.MultiLineString)
		return geometry,geo.GeomType_MULTILINESTRING,bb		
	case "MultiPolygon":
		// multi polygon code here
		geometry,bb := Make_MultiPolygon(geom.MultiPolygon)
		return geometry,geo.GeomType_MULTIPOLYGON,bb
	}

	return []uint64{},geomtype,[]int64{}
} 

// reflects a tile value back and stuff
func Reflect_Value(v interface{}) *geo.Value {
	tv := new(geo.Value)
	//fmt.Print(v)
	vv := reflect.ValueOf(v)
	kd := vv.Kind()

	switch kd {
	case reflect.String:
		val := string(vv.String())
		tv.StringValue = val

	case reflect.Float32:
		val := float32(vv.Float())
		tv.FloatValue = val
	case reflect.Float64:
		// do double here
		val := float64(vv.Float())
		tv.DoubleValue = val
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		val := int64(vv.Int())
		tv.IntValue = val
	case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
		val := uint64(vv.Uint())
		tv.UintValue = val
	case reflect.Bool:
		val := vv.Bool()
		tv.BoolValue = val
	}

	return tv
}

func Make_Properties(properties map[string]interface{}) map[string]*geo.Value {
	newmap := map[string]*geo.Value{}
	for k,v := range properties {
		newmap[k] = Reflect_Value(v)
	}
	return newmap
}


func Make_Feature(feat *geojson.Feature) *geo.Feature {
	// getting geometry and geom type
	geom,geom_type,bb := Make_Geom(feat.Geometry)

	id,boolval := feat.ID.(uint64)


	feature := &geo.Feature{
				Geometry:geom,
				Type:geom_type,
				Properties:Make_Properties(feat.Properties),
				BoundingBox:bb,
	}
	if boolval {
		feature.Id = id
	}
	return feature
}