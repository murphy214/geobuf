package geobuf_raw

import (
	pbff "github.com/murphy214/pbf"
	"github.com/paulmach/go.geojson"
)

// reads a feature
func ReadFeature(bytevals []byte) *geojson.Feature {
	pbfval := pbff.PBF{Pbf: bytevals, Length: len(bytevals)}
	var geomtype string
	feature := &geojson.Feature{Properties: map[string]interface{}{}}

	key, val := pbfval.ReadKey()
	if key == 1 && val == 0 {
		feature.ID = pbfval.ReadVarint()
		key, val = pbfval.ReadKey()
	}
	for key == 2 && val == 2 {
		// starting properties shit here

		size := pbfval.ReadVarint()
		endpos := pbfval.Pos + size
		//pbfval.ReadKey()
		pbfval.Pos += 1
		keyvalue := pbfval.ReadString()

		pbfval.Pos += 1
		pbfval.Varint()
		newkey, _ := pbfval.ReadKey()
		switch newkey {
		case 1:
			feature.Properties[keyvalue] = pbfval.ReadString()
		case 2:
			feature.Properties[keyvalue] = pbfval.ReadFloat()
		case 3:
			feature.Properties[keyvalue] = pbfval.ReadDouble()
		case 4:
			feature.Properties[keyvalue] = pbfval.ReadInt64()
		case 5:
			feature.Properties[keyvalue] = pbfval.ReadUInt64()
		case 6:
			feature.Properties[keyvalue] = pbfval.ReadUInt64()
		case 7:
			feature.Properties[keyvalue] = pbfval.ReadBool()
		}
		pbfval.Pos = endpos
		key, val = pbfval.ReadKey()
	}
	if key == 3 && val == 0 {
		switch int(pbfval.Pbf[pbfval.Pos]) {
		case 1:
			geomtype = "Point"
		case 2:
			geomtype = "LineString"
		case 3:
			geomtype = "Polygon"
		case 4:
			geomtype = "MultiPoint"
		case 5:
			geomtype = "MultiLineString"
		case 6:
			geomtype = "MultiPolygon"
		}
		pbfval.Pos += 1
		key, val = pbfval.ReadKey()
	}
	if key == 4 && val == 2 {
		size := pbfval.ReadVarint()
		endpos := pbfval.Pos + size

		switch geomtype {
		case "Point":
			feature.Geometry = geojson.NewPointGeometry(ReadPoint(pbfval,endpos))
		case "LineString":
			feature.Geometry = geojson.NewLineStringGeometry(ReadLine(pbfval,0, endpos))
		case "Polygon":
			feature.Geometry = geojson.NewPolygonGeometry(ReadPolygon(pbfval,endpos))
		case "MultiPoint":
			feature.Geometry = geojson.NewMultiPointGeometry(ReadLine(pbfval,0, endpos)...)
		case "MultiLineString":
			feature.Geometry = geojson.NewMultiLineStringGeometry(ReadPolygon(pbfval,endpos)...)
		case "MultiPolygon":
			feature.Geometry = geojson.NewMultiPolygonGeometry(ReadMultiPolygon(pbfval,endpos)...)

		}
		key, val = pbfval.ReadKey()

	}
	if key == 5 && val == 2 {
		feature.BoundingBox = ReadBoundingBox(pbfval)
	}
	return feature
}


// reads a feature
func ReadBB(bytevals []byte) []float64 {
	pbfval := pbff.PBF{Pbf: bytevals, Length: len(bytevals)}

	key, val := pbfval.ReadKey()
	if key == 1 && val == 0 {
		pbfval.ReadVarint()
		key, val = pbfval.ReadKey()
	}
	for key == 2 && val == 2 {
		// starting properties shit here

		size := pbfval.ReadVarint()
		endpos := pbfval.Pos + size
		pbfval.Pos = endpos
		key, val = pbfval.ReadKey()
	}
	if key == 3 && val == 0 {
		pbfval.Pos += 1
		key, val = pbfval.ReadKey()
	}
	if key == 4 && val == 2 {
		size := pbfval.ReadVarint()
		endpos := pbfval.Pos + size
		pbfval.Pos = endpos
		key, val = pbfval.ReadKey()

	}
	if key == 5 && val == 2 {
		return ReadBoundingBox(pbfval)
	}
	return []float64{}
}



// TO-DO: needs to be abstracted out into appropriate geobuf package
// geobuf functions i still have in here
func ReadPoint(pbf pbff.PBF,endpos int) []float64 {
	for pbf.Pos < endpos {
		x := pbf.ReadSVarintPower()
		y := pbf.ReadSVarintPower()
		return []float64{pbff.Round(x, .5, 7), pbff.Round(y, .5, 7)}
	}
	return []float64{}
}

// TO-DO: needs to be abstracted out into appropriate geobuf package
// reads a line
func ReadLine(pbf pbff.PBF,num int, endpos int) [][]float64 {
	var x, y float64
	if num == 0 {

		for startpos := pbf.Pos; startpos < endpos; startpos++ {
			if pbf.Pbf[startpos] <= 127 {
				num += 1
			}
		}
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += pbf.ReadSVarintPower()
			y += pbf.ReadSVarintPower()
			newlist[i] = []float64{pbff.Round(x, .5, 7), pbff.Round(y, .5, 7)}
		}

		return newlist
	} else {
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += pbf.ReadSVarintPower()
			y += pbf.ReadSVarintPower()

			newlist[i] = []float64{pbff.Round(x, .5, 7), pbff.Round(y, .5, 7)}

		}
		return newlist
	}
	return [][]float64{}
}

// TO-DO: needs to be abstracted out into appropriate geobuf package
func ReadPolygon(pbf pbff.PBF,endpos int) [][][]float64 {
	polygon := [][][]float64{}
	for pbf.Pos < endpos {
		num := pbf.ReadVarint()
		polygon = append(polygon, ReadLine(pbf,num, endpos))
	}
	return polygon
}

// TO-DO: needs to be abstracted out into appropriate geobuf package
func ReadMultiPolygon(pbf pbff.PBF,endpos int) [][][][]float64 {
	multipolygon := [][][][]float64{}
	for pbf.Pos < endpos {
		num_rings := pbf.ReadVarint()
		polygon := make([][][]float64, num_rings)
		for i := 0; i < num_rings; i++ {
			num := pbf.ReadVarint()
			polygon[i] = ReadLine(pbf,num, endpos)
		}
		multipolygon = append(multipolygon, polygon)
	}
	return multipolygon
}

// TO-DO: needs to be abstracted out into appropriate geobuf package
func ReadBoundingBox(pbf pbff.PBF) []float64 {
	bb := make([]float64, 4)
	pbf.ReadVarint()
	bb[0] = float64(pbf.ReadSVarintPower())
	bb[1] = float64(pbf.ReadSVarintPower())
	bb[2] = float64(pbf.ReadSVarintPower())
	bb[3] = float64(pbf.ReadSVarintPower())
	return bb
}