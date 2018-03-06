package geobuf_raw

import (
	//"github.com/murphy214/geobuf_new"
	//"github.com/murphy214/geobuf/geobuf_raw"
	//"github.com/golang/protobuf/proto"
	//"fmt"
	"testing"
	//"github.com/murphy214/pbf"
	"github.com/paulmach/go.geojson"
	"math"
)



var precision = math.Pow(10.0,-7.0)

var feature_s = `{"id":1000001,"type":"Feature","bbox":[-83.647031,33.698307,-83.275933,33.9659119],"geometry":{"type":"MultiPolygon","coordinates":[[[[-83.537385,33.9659119],[-83.5084519,33.931233],[-83.4155119,33.918541],[-83.275933,33.847977],[-83.306619,33.811444],[-83.28034,33.7617739],[-83.29145,33.7343149],[-83.406189,33.698307],[-83.479523,33.802265],[-83.505928,33.81776],[-83.533165,33.820923],[-83.647031,33.9061979],[-83.537385,33.9659119]]],[[[-83.537385,33.9659119],[-83.5084519,33.931233],[-83.4155119,33.918541],[-83.275933,33.847977],[-83.306619,33.811444],[-83.28034,33.7617739],[-83.29145,33.7343149],[-83.406189,33.698307],[-83.479523,33.802265],[-83.505928,33.81776],[-83.533165,33.820923],[-83.647031,33.9061979],[-83.537385,33.9659119]]],[[[-83.537385,33.9659119],[-83.5084519,33.931233],[-83.4155119,33.918541],[-83.275933,33.847977],[-83.306619,33.811444],[-83.28034,33.7617739],[-83.29145,33.7343149],[-83.406189,33.698307],[-83.479523,33.802265],[-83.505928,33.81776],[-83.533165,33.820923],[-83.647031,33.9061979],[-83.537385,33.9659119]]]]},"properties":{"AREA":"13219","COLORKEY":"#03E174","area":"13219","index":1109}}`
var feature,_ = geojson.UnmarshalFeature([]byte(feature_s))
var bytevals = WriteFeature(feature)


var polygon_s = `{"geometry": {"type": "Polygon", "coordinates": [[[-7.734374999999999, 25.799891182088334], [10.8984375, -34.016241889667015], [45.703125, 17.644022027872726], [-5.9765625, 26.43122806450644], [-7.734374999999999, 25.799891182088334]]]}, "type": "Feature", "properties": {}}`
var multipolygon_s = `{"type":"Feature","properties":{},"geometry":{"type":"MultiPolygon","coordinates":[[[[-71.71875,51.17934297928927],[-36.2109375,-49.15296965617039],[30.585937499999996,0.3515602939922709],[29.179687499999996,59.17592824927136],[-38.3203125,70.72897946208789],[-71.71875,51.17934297928927]]],[[[33.3984375,74.68325030051861],[75.234375,16.29905101458183],[76.2890625,64.77412531292873],[32.6953125,75.23066741281573],[33.3984375,74.68325030051861]]]]}}`
var linestring_s = `{"geometry": {"type": "LineString", "coordinates": [[10.8984375, 56.17002298293205], [16.5234375, -2.108898659243126], [59.4140625, 42.032974332441405], [61.17187499999999, 42.293564192170095]]}, "type": "Feature", "properties": {}}`	
var multilinestring_s = `{"geometry": {"type": "MultiLineString", "coordinates": [[[-48.1640625, 47.754097979680026], [-9.140625, 4.214943141390651], [15.468749999999998, -9.102096738726443]], [[10.8984375, 56.17002298293205], [16.5234375, -2.108898659243126], [59.4140625, 42.032974332441405], [61.17187499999999, 42.293564192170095]]]}, "type": "Feature", "properties": {}}`
var point_s = `{"geometry": {"type": "Point", "coordinates": [-48.1640625, 47.754097979680026]}, "type": "Feature", "properties": {}}`
var multipoint_s = `{"geometry": {"type": "MultiPoint", "coordinates": [[-48.1640625, 47.754097979680026], [-9.140625, 4.214943141390651]]}, "type": "Feature", "properties": {}}`



var polygon,err = geojson.UnmarshalFeature([]byte(polygon_s))
var multipolygon,_ = geojson.UnmarshalFeature([]byte(multipolygon_s))
var linestring,_ = geojson.UnmarshalFeature([]byte(linestring_s))
var multilinestring,_ = geojson.UnmarshalFeature([]byte(multilinestring_s))
var point,_ = geojson.UnmarshalFeature([]byte(point_s))
var multipoint,_ = geojson.UnmarshalFeature([]byte(multipoint_s))

func Benchmark_Write_New(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		WriteFeature(feature)
		//fmt.Println(string(bytevals))
	}
}



func Benchmark_Read_New(b *testing.B) {
    b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		ReadFeature(bytevals)
		//fmt.Println(string(bytevals))
	}
}





func TestReadWritePolygon(t *testing.T) {
	bytevals2 := WriteFeature(polygon)
	newfeature := ReadFeature(bytevals2)

	if newfeature.Geometry.Type == newfeature.Geometry.Type {
		//fmt.Println(newfeature.Geometry.Polygon,polygon.Geometry.Polygon)

		for i := range newfeature.Geometry.Polygon {
			for j := range newfeature.Geometry.Polygon[i] {
				//fmt.Println(newfeature.Geometry.Polygon[i][j],polygon.Geometry.Polygon[i][j])
				deltax := math.Abs(newfeature.Geometry.Polygon[i][j][0] - polygon.Geometry.Polygon[i][j][0])
				deltay := math.Abs(newfeature.Geometry.Polygon[i][j][1] - polygon.Geometry.Polygon[i][j][1])
				if deltax > precision || deltay > precision {
					t.Error("Test_Read_Write_Polygon Geometry signicantly off %v %v",newfeature.Geometry,polygon.Geometry)
				}
			}
		}
	} else {
		t.Errorf("Test_Read_Write_Polygon %v %v",newfeature,polygon)
	}


}

func TestReadWriteMultiPolygon(t *testing.T) {
	bytevals2 := WriteFeature(multipolygon)
	newfeature := ReadFeature(bytevals2)

	if newfeature.Geometry.Type == newfeature.Geometry.Type {
		//fmt.Println(newfeature.Geometry.Polygon,polygon.Geometry.Polygon)

		for i := range newfeature.Geometry.MultiPolygon {
			for j := range newfeature.Geometry.MultiPolygon[i] {
				for k := range newfeature.Geometry.MultiPolygon[i][j] {
					//fmt.Println(newfeature.Geometry.Polygon[i][j],polygon.Geometry.Polygon[i][j])
					deltax := math.Abs(newfeature.Geometry.MultiPolygon[i][j][k][0] - multipolygon.Geometry.MultiPolygon[i][j][k][0])
					deltay := math.Abs(newfeature.Geometry.MultiPolygon[i][j][k][1] - multipolygon.Geometry.MultiPolygon[i][j][k][1])
					if deltax > precision || deltay > precision {
						t.Error("TestReadWriteMultiPolygon Geometry signicantly off %v %v",newfeature.Geometry,multipolygon.Geometry)
					}
				}
			}
		}
	} else {
		t.Errorf("TestReadWriteMultiPolygon %v %v",newfeature,multipolygon,multipolygon)
	}


}




func TestReadWriteMultiLineString(t *testing.T) {
	bytevals2 := WriteFeature(multilinestring)
	newfeature := ReadFeature(bytevals2)
	if newfeature.Geometry.Type == newfeature.Geometry.Type {
		//fmt.Println(newfeature.Geometry.MultiLineString,multilinestring.Geometry.MultiLineString)

		for i := range newfeature.Geometry.MultiLineString {
			for j := range newfeature.Geometry.MultiLineString[i] {
				//fmt.Println(newfeature.Geometry.Polygon[i][j],polygon.Geometry.Polygon[i][j])
				deltax := math.Abs(newfeature.Geometry.MultiLineString[i][j][0] - multilinestring.Geometry.MultiLineString[i][j][0])
				deltay := math.Abs(newfeature.Geometry.MultiLineString[i][j][1] - multilinestring.Geometry.MultiLineString[i][j][1])
				if deltax > precision || deltay > precision {
					t.Error("TestReadWriteMultiLineString Geometry signicantly off %v %v",newfeature.Geometry,multilinestring.Geometry)
				}
			}
		}
	} else {
		t.Errorf("TestReadWriteMultiLineString %v %v",newfeature,multilinestring)
	}
}



func TestReadWriteLineString(t *testing.T) {
	bytevals2 := WriteFeature(linestring)
	newfeature := ReadFeature(bytevals2)
	if newfeature.Geometry.Type == newfeature.Geometry.Type {
		//fmt.Println(newfeature.Geometry.MultiLineString,multilinestring.Geometry.MultiLineString)

		for i := range newfeature.Geometry.LineString {
			//fmt.Println(newfeature.Geometry.Polygon[i][j],polygon.Geometry.Polygon[i][j])
			deltax := math.Abs(newfeature.Geometry.LineString[i][0] - linestring.Geometry.LineString[i][0])
			deltay := math.Abs(newfeature.Geometry.LineString[i][1] - linestring.Geometry.LineString[i][1])
			if deltax > precision || deltay > precision {
				t.Error("TestReadWriteLineString Geometry signicantly off %v %v",newfeature.Geometry,linestring.Geometry)
			}
		}
	} else {
		t.Errorf("TestReadWriteLineString %v %v",newfeature,linestring)
	}
}



func TestReadWriteMultiPoint(t *testing.T) {
	bytevals2 := WriteFeature(multipoint)
	newfeature := ReadFeature(bytevals2)
	if newfeature.Geometry.Type == newfeature.Geometry.Type {
		//fmt.Println(newfeature.Geometry.MultiLineString,multilinestring.Geometry.MultiLineString)

		for i := range newfeature.Geometry.MultiPoint {
			//fmt.Println(newfeature.Geometry.Polygon[i][j],polygon.Geometry.Polygon[i][j])
			deltax := math.Abs(newfeature.Geometry.MultiPoint[i][0] - multipoint.Geometry.MultiPoint[i][0])
			deltay := math.Abs(newfeature.Geometry.MultiPoint[i][1] - multipoint.Geometry.MultiPoint[i][1])
			if deltax > precision || deltay > precision {
				t.Error("TestReadWriteLineString Geometry signicantly off %v %v",newfeature.Geometry,multipoint.Geometry)
			}
		}
	} else {
		t.Errorf("TestReadWriteLineString %v %v",newfeature,multipoint)
	}
}



func TestReadWritePoint(t *testing.T) {
	bytevals2 := WriteFeature(point)
	newfeature := ReadFeature(bytevals2)
	if newfeature.Geometry.Type == newfeature.Geometry.Type {
		//fmt.Println(newfeature.Geometry.MultiLineString,multilinestring.Geometry.MultiLineString)

			//fmt.Println(newfeature.Geometry.Polygon[i][j],polygon.Geometry.Polygon[i][j])
		deltax := math.Abs(newfeature.Geometry.Point[0] - point.Geometry.Point[0])
		deltay := math.Abs(newfeature.Geometry.Point[1] - point.Geometry.Point[1])
		if deltax > precision || deltay > precision {
				t.Error("TestReadWriteLineString Geometry signicantly off %v %v",newfeature.Geometry,point.Geometry)
		}
	} else {
		t.Errorf("TestReadWriteLineString %v %v",newfeature,point)
	}
}


