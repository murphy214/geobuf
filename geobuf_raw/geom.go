package geobuf_raw

import (
	"math"

	"github.com/murphy214/pbf"
)

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

func readpoint(mpbf pbf.PBF,endpos int) []float64 {
	for mpbf.Pos < endpos {
		x := mpbf.ReadSVarintPower()
		y := mpbf.ReadSVarintPower()
		return []float64{round(x, .5, 7), round(y, .5, 7)}
	}
	return []float64{}
}

func readpolygon(mpbf pbf.PBF,endpos int) [][][]float64 {
	polygon := [][][]float64{}
	for mpbf.Pos < endpos {
		num := mpbf.ReadVarint()
		polygon = append(polygon, readline(mpbf,num, endpos))
	}
	return polygon
}

func readmultipolygon(mpbf pbf.PBF,endpos int) [][][][]float64 {
	multipolygon := [][][][]float64{}
	for mpbf.Pos < endpos {
		num_rings := mpbf.ReadVarint()
		polygon := make([][][]float64, num_rings)
		for i := 0; i < num_rings; i++ {
			num := mpbf.ReadVarint()
			polygon[i] = readline(mpbf,num, endpos)
		}
		multipolygon = append(multipolygon, polygon)
	}
	return multipolygon
}

func readline(mpbf pbf.PBF,num int, endpos int) [][]float64 {
	var x, y float64
	if num == 0 {

		for startpos := mpbf.Pos; startpos < endpos; startpos++ {
			if mpbf.Pbf[startpos] <= 127 {
				num += 1
			}
		}
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += mpbf.ReadSVarintPower()
			y += mpbf.ReadSVarintPower()
			newlist[i] = []float64{round(x, .5, 7), round(y, .5, 7)}
		}
		return newlist
	} else {
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += mpbf.ReadSVarintPower()
			y += mpbf.ReadSVarintPower()

			newlist[i] = []float64{round(x, .5, 7), round(y, .5, 7)}

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


// converts a single pt
func ConvertPt(pt []float64) []int64 {
	newpt := make([]int64, 2)
	newpt[0] = int64(pt[0] * math.Pow(10.0, 7.0))
	newpt[1] = int64(pt[1] * math.Pow(10.0, 7.0))
	return newpt
}

// param encoding
func paramEnc(value int64) uint64 {
	return uint64((value << 1) ^ (value >> 31))
}


func writepointbs(pt []float64) []byte {
	point := ConvertPt(pt)
	return append([]byte{34}, WritePackedUint64([]uint64{paramEnc(point[0]), paramEnc(point[1])})...)
}

// makes a line to int64
func writeline(line [][]float64) ([]uint64, []int64) {
	//geometry := []uint64{}
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	//oldpt := Convert_Pt(line[0])
	newline := make([]uint64, len(line)*2)
	deltapt := make([]int64, 2)
	pt := make([]int64, 2)
	oldpt := make([]int64, 2)

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

		pt = ConvertPt(point)
		if i == 0 {
			newline[0] = paramEnc(pt[0])
			newline[1] = paramEnc(pt[1])
		} else {
			deltapt = []int64{pt[0] - oldpt[0], pt[1] - oldpt[1]}
			newline[i*2] = paramEnc(deltapt[0])
			newline[i*2+1] = paramEnc(deltapt[1])
		}
		oldpt = pt
	}

	return newline, []int64{int64(west * powerfactor),
		int64(south * powerfactor),
		int64(east * powerfactor),
		int64(north * powerfactor)}
}

// writes a line to bytes 
func writelinebs(line [][]float64) ([]byte, []int64) {
	newline,bbs := writeline(line)
	return append([]byte{34}, WritePackedUint64(newline)...),bbs
}


func writepolygon(polygon [][][]float64) ([]uint64, []int64) {
	geometry := []uint64{}
	bb := []int64{}
	for i, cont := range polygon {
		geometry = append(geometry, uint64(len(cont)*2))

		tmpgeom, tmpbb := writeline(cont)
		geometry = append(geometry, tmpgeom...)
		if i == 0 {
			bb = tmpbb
		}
	}
	//geometryb = append(geometryb,WritePackedUint64(geometry)...)
	return geometry, bb
}

func writepolygonbs(polygon [][][]float64) ([]byte, []int64) {
	newline,bbs := writepolygon(polygon)
	return append([]byte{34}, WritePackedUint64(newline)...),bbs
}

// creates a multi polygon array
func writemultipolygonbs(multipolygon [][][][]float64) ([]byte, []int64) {
	geometryb := []byte{34}
	geometry := []uint64{}
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	west, south, east, north = west*powerfactor, south*powerfactor, east*powerfactor, north*powerfactor
	bb := []int64{int64(west), int64(south), int64(east), int64(north)}

	for _, polygon := range multipolygon {
		geometry = append(geometry, uint64(len(polygon)))
		tempgeom, tempbb := writepolygon(polygon)
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
