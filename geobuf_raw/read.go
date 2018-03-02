package geobuf_raw

import (
	"github.com/paulmach/go.geojson"
	"math"
	geo "./geobuf"
)

// decodes a given delta
func DecodeDelta(nume uint64) float64 {
	num := int(nume)
	if num%2 == 1 {
		return float64((num + 1) / -2) / math.Pow(10.0,7.0)
	} else {
		return float64(num / 2) / math.Pow(10.0,7.0)
	}
}

// reads a point
func Read_Point(geom []uint64) []float64 {
	return []float64{DecodeDelta(geom[0]),DecodeDelta(geom[1])}
}

// reads a line
func Read_Line(line []uint64) [][]float64 {
	newline := make([][]float64,len(line)/2)
	pt := []float64{0.0,0.0}
	for i := 0; i < len(newline); i++ {
		deltapt := Read_Point(line[i*2:i*2+2])
		newpt := []float64{pt[0] + deltapt[0],pt[1] + deltapt[1]}
		newline[i] = newpt
		pt = newpt

	}
	return newline
}

func Read_Polygon(polygon []uint64) [][][]float64 {
	pos := 0
	newpolygon := [][][]float64{}
	for pos < len(polygon) {
		size := int(polygon[pos])
		pos += 1
		newpolygon = append(newpolygon,Read_Line(polygon[pos:pos+size]))
		pos += size
	}
	return newpolygon
}


// decodes a geometry
func Read_Geometry(geom []uint64,geomtype geo.GeomType) *geojson.Geometry {
	if geomtype == geo.GeomType_POINT {
		return geojson.NewPointGeometry(Read_Point(geom))
	} else if geomtype == geo.GeomType_LINESTRING {
		return geojson.NewLineStringGeometry(Read_Line(geom))
	} else if geomtype == geo.GeomType_POLYGON {
		return geojson.NewPolygonGeometry(Read_Polygon(geom))
	}
	return &geojson.Geometry{}
}

func Get_Value(vall *geo.Value) interface{} {
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


func Get_ID(id interface{}) (interface{},bool) {
	if id == nil {
		return "",false
	} 
	return id,true
}

func Read_Feature(feat *geo.Feature) *geojson.Feature {
	feature := &geojson.Feature{Properties:map[string]interface{}{}}
	feature.Geometry = Read_Geometry(feat.Geometry,feat.Type)
	
	// adding id
	if feat.Id != 0 {
		feature.ID = feat.Id
	}

	// getting properties
	for k,v := range feat.Properties {
		feature.Properties[k] = Get_Value(v)
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