package geobuf_raw

import (
	"github.com/paulmach/go.geojson"
	"math"
	geo "github.com/murphy214/geobuf_new/geobuf_raw/geobuf"
)

// decodes a given delta
func DecodeDelta(nume uint64) float64 {
	num := int(nume)
	if num%2 == 1 {
		return float64((num + 1) / -2) / powerfactor
	} else {
		return float64(num / 2) / powerfactor
	}
}

// reads a point
func ReadPoint(geom []uint64) []float64 {
	return []float64{DecodeDelta(geom[0]),DecodeDelta(geom[1])}
}

// reads a line
func ReadLine(line []uint64) [][]float64 {
	newline := make([][]float64,len(line)/2)
	pt := []float64{0.0,0.0}
	for i := 0; i < len(newline); i++ {
		deltapt := ReadPoint(line[i*2:i*2+2])
		newpt := []float64{pt[0] + deltapt[0],pt[1] + deltapt[1]}
		newline[i] = newpt
		pt = newpt

	}
	return newline
}

func ReadPolygon(polygon []uint64) [][][]float64 {
	pos := 0
	newpolygon := [][][]float64{}
	for pos < len(polygon) {
		size := int(polygon[pos])
		pos += 1
		line := ReadLine(polygon[pos:pos+size])
		line[len(line)-1] = line[0]
		newpolygon = append(newpolygon,line)
		pos += size
	}
	return newpolygon
}


func ReadMultiPolygon(multipolygon []uint64) [][][][]float64 {
	pos := 0
	newmultipolygon := [][][][]float64{}
	for pos < len(multipolygon) {
		ringsize := int(multipolygon[pos])
		pos += 1
		currentring := 0
		startpos := pos
		for currentring < ringsize {
			size := int(multipolygon[pos])
			pos += 1
			pos += size
			currentring += 1
		}
		newmultipolygon = append(newmultipolygon,ReadPolygon(multipolygon[startpos:pos]))
	}
	return newmultipolygon
}


// decodes a geometry
func ReadGeometry(geom []uint64,geomtype geo.GeomType) *geojson.Geometry {
	switch geomtype {
	case geo.GeomType_POINT:
		return geojson.NewPointGeometry(ReadPoint(geom))
	case geo.GeomType_LINESTRING:	
		return geojson.NewLineStringGeometry(ReadLine(geom))
	case geo.GeomType_POLYGON:
		return geojson.NewPolygonGeometry(ReadPolygon(geom))
	case geo.GeomType_MULTIPOINT:
		return geojson.NewMultiPointGeometry(ReadLine(geom)...)
	case geo.GeomType_MULTILINESTRING:
		return geojson.NewMultiLineStringGeometry(ReadPolygon(geom)...)
	case geo.GeomType_MULTIPOLYGON:
		return geojson.NewMultiPolygonGeometry(ReadMultiPolygon(geom)...)
	}

	return &geojson.Geometry{}
}

func GetValue(vall *geo.Value) interface{} {
	val := *vall
	if val.StringValue != "" {
		return val.StringValue
	}
	if val.FloatValue != 0.0 {
		return val.FloatValue
	}
	if val.DoubleValue != 0 {
		return val.DoubleValue
	}
	if val.IntValue != 0 {
		return val.IntValue
	}
	if val.UintValue != 0 {
		return val.UintValue
	}
	if val.SintValue != 0 {
		return val.SintValue
	}
	if val.BoolValue != false {
		return val.BoolValue
	}
	return val.BoolValue
}


func GetID(id interface{}) (interface{},bool) {
	if id == nil {
		return "",false
	} 
	return id,true
}

func ReadFeature(feat *geo.Feature) *geojson.Feature {
	feature := &geojson.Feature{Properties:map[string]interface{}{}}
	feature.Geometry = ReadGeometry(feat.Geometry,feat.Type)
	
	// adding id
	if feat.Id != 0 {
		feature.ID = feat.Id
	}

	// getting properties
	for k,v := range feat.Properties {
		feature.Properties[k] = GetValue(v)
	}
	val := math.Pow(10.0,7.0)
	if len(feat.BoundingBox) > 0 {
		feature.BoundingBox = []float64{float64(feat.BoundingBox[0]) / val,
			float64(feat.BoundingBox[1]) / val,
			float64(feat.BoundingBox[2]) / val,
			float64(feat.BoundingBox[3]) / val,}
		//feature.BoundingBox = []float64{}
		//for i := range 
	}


	return feature

}