package geobuf_raw

import (
	"math"
	"github.com/murphy214/pbf"
	"github.com/paulmach/go.geojson"
)

// gets the dimension size of the geometry if possible
func getdimsize(geom *geojson.Geometry) int {
	switch geom.Type {
	case "Point":
		return len(geom.Point)
	case "MultiPoint":
		if len(geom.MultiPoint)>0 {
			return len(geom.MultiPoint[0])
		}
		return 2 
	case "LineString":
		if len(geom.LineString)>0 {
			return len(geom.LineString[0])
		}
		return 2 
	case "Polygon":
		if len(geom.Polygon)>0&&len(geom.Polygon[0])>0 {
			return len(geom.Polygon[0][0])
		}
		return 2 
	case "MultiLineString":
		if len(geom.MultiLineString)>0&&len(geom.MultiLineString[0])>0 {
			return len(geom.MultiLineString[0][0])
		}
		return 2 
	case "MultiPolygon":
		if len(geom.MultiPolygon)>0&&len(geom.MultiPolygon[0])>0&&len(geom.MultiPolygon[0][0])>0 {
			return len(geom.MultiPolygon[0][0][0])
		}
		return 2 
	default:
		return 2
	} 
}

// returns the geom_type and dim_size
func geomcode_details(x int) (int,int) {
	if x <= 6 { 
		return x,2 
	} else {
		dim_size := x - ((x >> 4) << 4) 
		geom_type := x >> 4
		return geom_type,dim_size 
	}
}

// makes the geometry code with dim_size embedded
func makegeomcode(geom_type,dim_size int) byte {
	return byte((geom_type << 4) + dim_size)
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func readpoint(mpbf pbf.PBF,endpos int,dim_size int) []float64 {
	for mpbf.Pos < endpos {
		pt := make([]float64,dim_size)
		for j := 0;j < dim_size; j++ {
			pt[j] += round(mpbf.ReadSVarintPower(),0.5,7)
		}
		return pt
	}
	return []float64{}
}

func readpolygon(mpbf pbf.PBF,endpos int,dim_size int) [][][]float64 {
	polygon := [][][]float64{}
	for mpbf.Pos < endpos {
		num := mpbf.ReadVarint()
		polygon = append(polygon, readline(mpbf,num, endpos,dim_size))
	}
	return polygon
}

func readmultipolygon(mpbf pbf.PBF,endpos int,dim_size int) [][][][]float64 {
	multipolygon := [][][][]float64{}
	for mpbf.Pos < endpos {
		num_rings := mpbf.ReadVarint()
		polygon := make([][][]float64, num_rings)
		for i := 0; i < num_rings; i++ {
			num := mpbf.ReadVarint()
			polygon[i] = readline(mpbf,num, endpos,dim_size)
		}
		multipolygon = append(multipolygon, polygon)
	}
	return multipolygon
}

func readline(mpbf pbf.PBF,num int, endpos int,dim_size int) [][]float64 {
	pt := make([]float64,dim_size)
	if num == 0 {
		for startpos := mpbf.Pos; startpos < endpos; startpos++ {
			if mpbf.Pbf[startpos] <= 127 {
				num += 1
			}
		}
		newlist := make([][]float64, num/dim_size)
		for i := 0; i < num/dim_size; i++ {
			rdpt := make([]float64,dim_size)
			for j := 0;j < dim_size; j++ {
				pt[j] += mpbf.ReadSVarintPower()
				rdpt[j] = round(pt[j],0.5,7)
			}
			newlist[i] = rdpt
		}
		return newlist
	} else {
		newlist := make([][]float64, num/dim_size)
		for i := 0; i < num/dim_size; i++ {
			rdpt := make([]float64,dim_size)
			for j := 0;j < dim_size; j++ {
				pt[j] += mpbf.ReadSVarintPower()
				rdpt[j] = round(pt[j],0.5,7)
			}
			newlist[i] = rdpt
		}
		return newlist
	}
	return [][]float64{}
}

func readboundingbox(mpbf pbf.PBF) []float64 {
	bb := make([]float64, 4)
	mpbf.ReadVarint()
	bb[0] = float64(mpbf.ReadSVarintPower())
	bb[1] = float64(mpbf.ReadSVarintPower())
	bb[2] = float64(mpbf.ReadSVarintPower())
	bb[3] = float64(mpbf.ReadSVarintPower())
	return bb
}

// ############### STARTING geom writing
// converts a single pt
func ConvertPt(pt []float64,dim_size int) []int64 {
	newpt := make([]int64, dim_size)
	for pos := range pt {
		newpt[pos] = int64(pt[pos] * math.Pow(10.0, 7.0))
	}
	return newpt
}

// param encoding
func paramEnc(value int64) uint64 {
	return uint64((value << 1) ^ (value >> 31))
}


func writepointbs(pt []float64,dim_size int) []byte {
	point := ConvertPt(pt,dim_size)
	return append([]byte{34}, WritePackedUint64([]uint64{paramEnc(point[0]), paramEnc(point[1])})...)
}

// makes a line to int64
func writeline(line [][]float64,dim_size int) ([]uint64, []int64) {
	//geometry := []uint64{}
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	//oldpt := Convert_Pt(line[0])
	newline := make([]uint64, len(line)*dim_size)
	// deltapt := make([]int64, dim_size)
	pt := make([]int64, dim_size)
	oldpt := make([]int64, dim_size)

	for i, point := range line {
		x, y := point[0], point[1]
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

		pt = ConvertPt(point,dim_size)
		if i == 0 {
			for j := 0; j < dim_size; j++ {
				newline[j] = paramEnc(pt[j])
			}
		} else {
			for j := 0; j < dim_size; j++ {
				newline[i*dim_size+j] = paramEnc(pt[j] - oldpt[j])
			}
		}
		oldpt = pt
	}

	return newline, []int64{int64(west * powerfactor),
		int64(south * powerfactor),
		int64(east * powerfactor),
		int64(north * powerfactor)}
}

// writes a line to bytes 
func writelinebs(line [][]float64,dim_size int) ([]byte, []int64) {
	newline,bbs := writeline(line,dim_size)
	return append([]byte{34}, WritePackedUint64(newline)...),bbs
}


func writepolygon(polygon [][][]float64,dim_size int) ([]uint64, []int64) {
	geometry := []uint64{}
	bb := []int64{}
	for i, cont := range polygon {
		geometry = append(geometry, uint64(len(cont)*2))

		tmpgeom, tmpbb := writeline(cont,dim_size)
		geometry = append(geometry, tmpgeom...)
		if i == 0 {
			bb = tmpbb
		}
	}
	return geometry, bb
}

func writepolygonbs(polygon [][][]float64,dim_size int) ([]byte, []int64) {
	newline,bbs := writepolygon(polygon,dim_size)
	return append([]byte{34}, WritePackedUint64(newline)...),bbs
}

// creates a multi polygon array
func writemultipolygonbs(multipolygon [][][][]float64,dim_size int) ([]byte, []int64) {
	geometryb := []byte{34}
	geometry := []uint64{}
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	west, south, east, north = west*powerfactor, south*powerfactor, east*powerfactor, north*powerfactor
	bb := []int64{int64(west), int64(south), int64(east), int64(north)}

	for _, polygon := range multipolygon {
		geometry = append(geometry, uint64(len(polygon)))
		tempgeom, tempbb := writepolygon(polygon,dim_size)
		geometry = append(geometry, tempgeom...)
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
	geometryb = append(geometryb, WritePackedUint64(geometry)...)
	return geometryb, bb
}
