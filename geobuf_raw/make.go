package geobuf_raw

import (
	"github.com/paulmach/go.geojson"
	"math"
	geo "./geobuf"
	"reflect"
)

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
	return newline,[]int64{int64(west * math.Pow(10.0,7.0)),
	int64(south * math.Pow(10.0,7.0)), 
	int64(east * math.Pow(10.0,7.0)),
	int64(north * math.Pow(10.0,7.0)),}
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


// syntax wise geometry will only have basic types
// flat for liens no header
// flaot for points
// header for rings of polygon
func Make_Geom(geom *geojson.Geometry) ([]uint64,geo.GeomType,[]int64) {
	var geomtype geo.GeomType
	if geom.Type == "Point" {
		newpt := Make_Point(geom.Point)
		newpt2 := Convert_Pt(geom.Point)

		return newpt,geo.GeomType_POINT,[]int64{newpt2[0],newpt2[1],newpt2[0],newpt2[1]}
	} else if geom.Type == "LineString" {
		geometry,bb :=  Make_Line(geom.LineString)
		return geometry,geo.GeomType_LINESTRING,bb
	} else if geom.Type == "Polygon" {
		geometry,bb :=  Make_Polygon(geom.Polygon)
		return geometry,geo.GeomType_POLYGON,bb
	}
	return []uint64{},geomtype,[]int64{}
} 

// reflects a tile value back and stuff
func Reflect_Value(v interface{}) *geo.Value {
	var tv *geo.Value
	//fmt.Print(v)
	vv := reflect.ValueOf(v)
	kd := vv.Kind()
	if (reflect.Float64 == kd) || (reflect.Float32 == kd) {
		//fmt.Print(v, "float", k)
		tv = Make_Tv_Float(float64(vv.Float()))
		//hash = Hash_Tv(tv)
	} else if (reflect.Int == kd) || (reflect.Int8 == kd) || (reflect.Int16 == kd) || (reflect.Int32 == kd) || (reflect.Int64 == kd) || (reflect.Uint8 == kd) || (reflect.Uint16 == kd) || (reflect.Uint32 == kd) || (reflect.Uint64 == kd) {
		//fmt.Print(v, "int", k)
		tv = Make_Tv_Int(int(vv.Int()))
		//hash = Hash_Tv(tv)
	} else if reflect.String == kd {
		//fmt.Print(v, "str", k)
		tv = Make_Tv_String(string(vv.String()))
		//hash = Hash_Tv(tv)

	} else {
		tv := new(geo.Value)
		t := ""
		tv.StringValue = t
	}
	return tv
}

// makes a tile_value string
func Make_Tv_String(stringval string) *geo.Value {
	tv := new(geo.Value)
	t := string(stringval)
	tv.StringValue = t
	return tv
}

// makes a tile value float
func Make_Tv_Float(val float64) *geo.Value {
	tv := new(geo.Value)
	t := float64(val)
	tv.DoubleValue = t
	return tv
}

// makes a tile value int
func Make_Tv_Int(val int) *geo.Value {
	tv := new(geo.Value)
	t := int64(val)
	tv.SintValue = t
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