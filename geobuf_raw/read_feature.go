package geobuf_raw

import (
	"github.com/murphy214/pbf"
	"github.com/paulmach/go.geojson"
)

// reads a feature
func ReadFeature(bytevals []byte) *geojson.Feature {
	pbfval := pbf.PBF{Pbf: bytevals, Length: len(bytevals)}
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
			feature.Geometry = geojson.NewPointGeometry(readpoint(pbfval,endpos))
		case "LineString":
			feature.Geometry = geojson.NewLineStringGeometry(readline(pbfval,0, endpos))
		case "Polygon":
			feature.Geometry = geojson.NewPolygonGeometry(readpolygon(pbfval,endpos))
		case "MultiPoint":
			feature.Geometry = geojson.NewMultiPointGeometry(readline(pbfval,0, endpos)...)
		case "MultiLineString":
			feature.Geometry = geojson.NewMultiLineStringGeometry(readpolygon(pbfval,endpos)...)
		case "MultiPolygon":
			feature.Geometry = geojson.NewMultiPolygonGeometry(readmultipolygon(pbfval,endpos)...)

		}
		key, val = pbfval.ReadKey()

	}
	if key == 5 && val == 2 {
		feature.BoundingBox = readboundingbox(pbfval)
	}
	return feature
}


// reads a feature
func ReadBB(bytevals []byte) []float64 {
	pbfval := pbf.PBF{Pbf: bytevals, Length: len(bytevals)}

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
		return readboundingbox(pbfval)
	}
	return []float64{}
}

