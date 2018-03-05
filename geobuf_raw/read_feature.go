package geobuf_raw

import (
	"github.com/paulmach/go.geojson"
)



// reads a feature
func ReadFeature(bytevals []byte) *geojson.Feature {
	pbf := PBF{Pbf:bytevals,Length:len(bytevals)}
	var geomtype string
	feature := &geojson.Feature{Properties:map[string]interface{}{}}

	key,val := pbf.ReadKey()
	if key == 1 && val == 0 {
		feature.ID = pbf.ReadVarint()
		key,val = pbf.ReadKey()
	}
	for key == 2 && val == 2 {
		// starting properties shit here

		size := pbf.ReadVarint()
		endpos := pbf.Pos + size
		//pbf.ReadKey()
		pbf.Pos += 1
		keyvalue := pbf.ReadString()

		pbf.Pos += 1
		pbf.ReadVarint()
		newkey,_ := pbf.ReadKey()
		switch newkey {
		case 1:
			feature.Properties[keyvalue] = pbf.ReadString()			
		case 2:
			feature.Properties[keyvalue] = pbf.ReadFloat()			
		case 3:
			feature.Properties[keyvalue] = pbf.ReadDouble()			
		case 4:
			feature.Properties[keyvalue] = pbf.ReadInt64()			
		case 5:
			feature.Properties[keyvalue] = pbf.ReadUInt64()			
		case 6:
			feature.Properties[keyvalue] = pbf.ReadUInt64()			
		case 7:
			feature.Properties[keyvalue] = pbf.ReadBool()			
		}
		pbf.Pos = endpos
		key,val = pbf.ReadKey()
	}
	if key == 3 && val == 0 {
		switch int(pbf.Pbf[pbf.Pos]) {
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
		pbf.Pos += 1
		key,val = pbf.ReadKey()
	}
	if key == 4 && val == 2 {
		size := pbf.ReadVarint()
		endpos := pbf.Pos + size

		switch geomtype {
		case "Point":
			feature.Geometry = geojson.NewPointGeometry(pbf.ReadPoint(endpos))
		case "LineString":
			feature.Geometry = geojson.NewLineStringGeometry(pbf.ReadLine(0,endpos))
		case "Polygon":
			feature.Geometry = geojson.NewPolygonGeometry(pbf.ReadPolygon(endpos))
		case "MultiPoint":
			feature.Geometry = geojson.NewMultiPointGeometry(pbf.ReadLine(0,endpos)...)
		case "MultiLineString":
			feature.Geometry = geojson.NewMultiLineStringGeometry(pbf.ReadPolygon(endpos)...)			
		case "MultiPolygon":
			feature.Geometry = geojson.NewMultiPolygonGeometry(pbf.ReadMultiPolygon(endpos)...)			

		}
		key,val = pbf.ReadKey()

	}
	if key == 5 && val == 2 {
		feature.BoundingBox = pbf.ReadBoundingBox()
	}
	return feature
}