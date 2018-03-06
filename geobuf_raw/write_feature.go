package geobuf_raw

import (
	"math"
	"github.com/paulmach/go.geojson"
	"reflect"
	"github.com/murphy214/pbf"
)


// converts a single pt
func ConvertPt(pt []float64) []int64 {
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
func MakePoint(pt []float64) []byte {
	point := ConvertPt(pt)
	geometryb := []byte{34}
	return append(geometryb,WritePackedUint64([]uint64{paramEnc(point[0]),paramEnc(point[1])})...)
}

// makes a line
func MakeLine(line [][]float64) ([]byte,[]int64) {
	geometryb := []byte{34}
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

		pt = ConvertPt(point)
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
	geometryb = append(geometryb,WritePackedUint64(newline)...)

	return geometryb,[]int64{int64(west * powerfactor),
	int64(south * powerfactor), 
	int64(east * powerfactor),
	int64(north * powerfactor),}
}

// makes a line
func MakeLine2(line [][]float64) ([]uint64,[]int64) {
	//geometry := []uint64{}
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

		pt = ConvertPt(point)
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
func MakePolygon(polygon [][][]float64) ([]byte,[]int64) {
	geometryb := []byte{34}
	geometry := []uint64{}
	bb := []int64{}
	for i,cont := range polygon {
		geometry = append(geometry,uint64(len(cont) * 2))

		tmpgeom,tmpbb := MakeLine2(cont)
		geometry = append(geometry,tmpgeom...)
		if i == 0 {
			bb = tmpbb
		}
	}
	geometryb = append(geometryb,WritePackedUint64(geometry)...)
	return geometryb,bb
}


// creates a polygon 
func MakePolygon2(polygon [][][]float64) ([]uint64,[]int64) {
	geometry := []uint64{}
	bb := []int64{}
	for i,cont := range polygon {
		geometry = append(geometry,uint64(len(cont) * 2))

		tmpgeom,tmpbb := MakeLine2(cont)
		geometry = append(geometry,tmpgeom...)
		if i == 0 {
			bb = tmpbb
		}
	}
	//geometryb = append(geometryb,WritePackedUint64(geometry)...)
	return geometry,bb
}

// creates a multi polygon array 
func MakeMultiPolygon(multipolygon [][][][]float64) ([]byte,[]int64) {
	geometryb := []byte{34}
	geometry := []uint64{}
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	west,south,east,north = west * powerfactor,south * powerfactor,east * powerfactor,north * powerfactor
	bb := []int64{int64(west),int64(south),int64(east),int64(north)}

	for _,polygon := range multipolygon {
		geometry = append(geometry,uint64(len(polygon)))
		tempgeom,tempbb := MakePolygon2(polygon)
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
	geometryb = append(geometryb,WritePackedUint64(geometry)...)
	return geometryb,bb
}	

func MakeKeyValue(key string,value interface{}) []byte {
	array1 := []byte{18}
	array3 := []byte{10}
	array4 := pbf.EncodeVarint(uint64(len(key)))
	array5 := []byte(key)
	array6 := WriteValue(value)
	array2 := pbf.EncodeVarint(uint64(len(array3)+len(array4)+len(array5)+len(array6)))
	return AppendAll(array1,array2,array3,array4,array5,array6)

}

// writes a feature
func WriteFeature(feat  *geojson.Feature) []byte {
	newbytes := []byte{8}

	if feat.ID != nil {
		vv := reflect.ValueOf(feat.ID)
		kd := vv.Kind()
		switch kd {
		case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
			newbytes = append(newbytes,pbf.EncodeVarint(uint64(vv.Int()))...)
		}
	} else {
		newbytes = []byte{}
	}

	// writing each key value property
	for k,v := range feat.Properties {
		newbytes = append(newbytes,MakeKeyValue(k,v)...)
	}
	if feat.Geometry != nil {
		switch feat.Geometry.Type {
		case "Point":
			// 
			newbytes = append(newbytes,[]byte{24,1}...)
			geomb := MakePoint(feat.Geometry.Point)
			newbytes = append(newbytes,geomb...)
		case "LineString":
			// 
			newbytes = append(newbytes,[]byte{24,2}...)
			geomb,_ := MakeLine(feat.Geometry.LineString)
			newbytes = append(newbytes,geomb...)

		case "Polygon":
			// here
			newbytes = append(newbytes,[]byte{24,3}...)
			geomb,_ := MakePolygon(feat.Geometry.Polygon)
			newbytes = append(newbytes,geomb...)
		case "MultiPoint":
			// here
			newbytes = append(newbytes,[]byte{24,4}...)
			geomb,_ := MakeLine(feat.Geometry.MultiPoint)
			newbytes = append(newbytes,geomb...)
		case "MultiLineString":
			newbytes = append(newbytes,[]byte{24,5}...)
			geomb,_ := MakePolygon(feat.Geometry.MultiLineString)
			newbytes = append(newbytes,geomb...)
		case "MultiPolygon":
			newbytes = append(newbytes,[]byte{24,6}...)
			geomb,_ := MakeMultiPolygon(feat.Geometry.MultiPolygon)
			newbytes = append(newbytes,geomb...) 
		}
	}
	return newbytes
}

