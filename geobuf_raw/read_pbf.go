package geobuf_raw

import (
	//"io/ioutil"
	//"fmt"
	"fmt"
	"math"
	//"vector-tile/2.1"
	//"github.com/golang/protobuf/proto"
)

var powerfactor = math.Pow(10.0,7.0)
var Powerfactor = math.Pow(10.0,7.0)

// 
type PBF struct {
	Pbf []byte
	Pos int
	Length int
}

const maxVarintBytes = 10 // maximum Length of a varint

// EncodeVarint returns the varint encoding of x.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
// Not used by the package itself, but helpful to clients
// wishing to use the same encoding.
func EncodeVarint(x uint64) []byte {
	var buf [maxVarintBytes]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}
 
// DecodeVarint reads a varint-encoded integer from the slice.
// It returns the integer and the number of bytes consumed, or
// zero if there is not enough.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func DecodeVarint2(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}

func DecodeVarint(buf []byte) (x uint64) { 
	i := 0
	if buf[i] < 0x80 {
		return uint64(buf[i])
	}

	var b uint64
	// we already checked the first byte
	x = uint64(buf[i]) - 0x80
	i++

	b = uint64(buf[i])
	i++
	x += b << 7
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 7

	b = uint64(buf[i])
	i++
	x += b << 14
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 14

	b = uint64(buf[i])
	i++
	x += b << 21
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 21

	b = uint64(buf[i])
	i++
	x += b << 28
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 28

	b = uint64(buf[i])
	i++
	x += b << 35
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 35

	b = uint64(buf[i])
	i++
	x += b << 42
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 42

	b = uint64(buf[i])
	i++
	x += b << 49
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 49

	b = uint64(buf[i])
	i++
	x += b << 56
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 56

	b = uint64(buf[i])
	i++
	x += b << 63
	if b&0x80 == 0 {
		goto done
	}
	// x -= 0x80 << 63 // Always zero.

	return 0

done:
	return x
}

// a much faster key integration (microseconds to nanoseconds)
// returns the value number and key number for a given byte
func Key(x byte) (byte, byte) {
	//fmt.Printf("%08b\n",x)
	val := x >> 3

	// if the x value has a value in the 8 place
	if int(x) >= 8 {
		x = x & 0x07

	} else {
		return val, x
	}
	// if the x value has a value in the 16 place
	if int(x) >= 16 {
		x = x & 0x0f

	} else {
		return val, x
	}

	if int(x) >= 32 {
		x = x & 0x1f

	} else {
		return val, x
	}

	if int(x) >= 64 {
		x = x & 0x3f

	} else {
		return val, x
	}

	if int(x) >= 128 {
		x = x & 0x7f

	} else {
		return val, x
	}

	return val, x

}


func ReadInt32(buf []byte) int32 {
	if len(buf) == 4 {
    	return int32(((int(buf[0])) | (int(buf[1]) << 8) | (int(buf[2]) << 16)) + (int(buf[3]) << 24))
	} else if len(buf) == 3 {
    	return int32(((int(buf[0])) | (int(buf[1]) << 8) | (int(buf[2]) << 16)))
	} else if len(buf) == 2 {
    	return int32(((int(buf[0])) | (int(buf[1]) << 8)))
	} else if len(buf) == 1 {
    	return int32(buf[0])
	}
	return int32(0)
}


func ReadUInt32(buf []byte) uint32 {
	if len(buf) == 4 {
    	return uint32(((int(buf[0])) | (int(buf[1]) << 8) | (int(buf[2]) << 16)) + (int(buf[3]) * 0x1000000))
	} else if len(buf) == 3 {
    	return uint32(((int(buf[0])) | (int(buf[1]) << 8) | (int(buf[2]) << 16)))
	} else if len(buf) == 2 {
    	return uint32(((int(buf[0])) | (int(buf[1]) << 8)))
	} else if len(buf) == 1 {
    	return uint32(buf[0])
	}
	return uint32(0)
}


// reads a uint64 from a list of bytes
func ReadUint64(bytes []byte) uint64 {
	return DecodeVarint(bytes)
}

// reads a uint64 from a list of bytes
func ReadInt64(bytes []byte) int64 {
	return int64(DecodeVarint(bytes))
}



func (pbf *PBF) ReadKey() (byte,byte) {
	
	var key,val byte
	if pbf.Pos > pbf.Length - 1 {
		key,val = 100,100
	} else {
		key,val = Key(pbf.Pbf[pbf.Pos])
		pbf.Pos += 1

	}
	return key,val	
	//pbf.Pos += 1
	//return Key(pbf.Pbf[pbf.Pos-1])
}


func (pbf *PBF) ReadVarint() int {
	if pbf.Pbf[pbf.Pos] <= 127 {
		pbf.Pos += 1
		return int(pbf.Pbf[pbf.Pos-1])
	}
	if pbf.Pos + 1 >= pbf.Length {
		if pbf.Pos + 1 == pbf.Length {
			pbf.Pos += 1
		}
		return 0
	}
	startPos := pbf.Pos 
	for pbf.Pbf[pbf.Pos] > 127 {
		pbf.Pos += 1
	}
	pbf.Pos += 1
	//if pbf.Pos - startPos == 1 {
	//	return int(pbf.Pbf[startPos])
	//}
	return int(DecodeVarint(pbf.Pbf[startPos:pbf.Pos]))
}

func (pbf *PBF) ReadSVarint() float64 {
	num := pbf.ReadVarint()
	if num%2 == 1 {
		return float64((num + 1) / -2)
	} else {
		return float64(num / 2)
	}
	return float64(0)
}

func (pbf *PBF) ReadSVarintPower() float64 {
	num := pbf.ReadVarint()
	if num%2 == 1 {
		return float64((num + 1) / -2) / powerfactor
	} else {
		return float64(num / 2) / powerfactor 
	}
	return float64(0)
}

// var int bytes
func (pbf *PBF) Varint() []byte {
	startPos := pbf.Pos 
	for pbf.Pbf[pbf.Pos] > 127 {
		pbf.Pos += 1
	}
	pbf.Pos += 1
	return pbf.Pbf[startPos:pbf.Pos]
}



func (pbf *PBF) ReadFixed32() uint32 {
	val := ReadUInt32(pbf.Pbf[pbf.Pos:pbf.Pos+4])

	pbf.Pos += 4
	return val
}

func (pbf *PBF) ReadUInt32() uint32 {
	return ReadUInt32(pbf.Varint())
}




func (pbf *PBF) ReadSFixed32() int32 {
	val := ReadInt32(pbf.Pbf[pbf.Pos:pbf.Pos+4])
	pbf.Pos += 4
	return val
}

func (pbf *PBF) ReadInt32() int32 {
	return ReadInt32(pbf.Varint())
}

// reads a uint64 from a list of bytes
func (pbf *PBF) ReadFixed64() uint64 {
	v := DecodeVarint(pbf.Pbf[pbf.Pos:pbf.Pos+8])
	pbf.Pos += 8
	return v
}

func (pbf *PBF) ReadUInt64() uint64 {
	return ReadUint64(pbf.Varint())
}

// reads a uint64 from a list of bytes
func (pbf *PBF) ReadSFixed64() int64 {
	v := DecodeVarint(pbf.Pbf[pbf.Pos:pbf.Pos+8])
	pbf.Pos += 8
	return int64(v)
}


func (pbf *PBF) ReadInt64() int64 {
	return ReadInt64(pbf.Varint())
}


func (pbf *PBF) ReadDouble() float64 {
	a := pbf.Pos
	pbf.Pos += 8
	return math.Float64frombits(uint64(pbf.Pbf[a]) | uint64(pbf.Pbf[a+1])<<8 | uint64(pbf.Pbf[a+2])<<16 | uint64(pbf.Pbf[a+3])<<24 | uint64(pbf.Pbf[a+4])<<32 | uint64(pbf.Pbf[a+5])<<40 | uint64(pbf.Pbf[a+6])<<48 | uint64(pbf.Pbf[a+7])<<56)
}

func (pbf *PBF) ReadFloat() float32 {
	a := pbf.Pos
	pbf.Pos += 4
	return math.Float32frombits(uint32(pbf.Pbf[a]) | uint32(pbf.Pbf[a+1])<<8 | uint32(pbf.Pbf[a+2])<<16 | uint32(pbf.Pbf[a+3])<<24)
}

func (pbf *PBF) ReadString() string {
	size := pbf.ReadVarint()
	stringval := string(pbf.Pbf[pbf.Pos:pbf.Pos+size])
	pbf.Pos += size
	return stringval
}

func (pbf *PBF) ReadBool() bool {

	if pbf.Pbf[pbf.Pos] == 1 {
		pbf.Pos += 1
		return true
	} else {
		pbf.Pos += 1
		return false
	}

	//pbf.Byte()
	/*
	//size := pbf.ReadVarint()
	buf := pbf.Pbf[pbf.Pos:pbf.Pos+1]
	pbf.Pos += 1

	if buf[0] == 1 {
		return true
	} else if buf[0] == 0 {
		return false
	}
	*/
	return false
}

func (pbf *PBF) ReadPacked() []uint32 {

	endpos := pbf.Pos + pbf.ReadVarint()
	//fmt.Println(pbf.Pbf[pbf.Pos])
	// potential uint32 values
	//fmt.Println(endpos)
	vals := make([]uint32,pbf.Length)
	currentpos := 0
	//fmt.Println(uint32(byte(32)))
	for pbf.Pos < endpos {
		startpos := pbf.Pos

		for pbf.Pbf[pbf.Pos] > 127 {
			pbf.Pos += 1
		}
		pbf.Pos += 1

		switch pbf.Pos - startpos {

		case 1:
			vals[currentpos] = uint32(pbf.Pbf[startpos])
			currentpos += 1
		//} else if startpos - startpos == 2 {
		case 2:
			vals[currentpos] = (uint32(pbf.Pbf[startpos])) | (uint32(pbf.Pbf[startpos+1]) << 8)
			currentpos += 1
		//} else if startpos - startpos == 3 {
		case 3:
			vals[currentpos] = (uint32(pbf.Pbf[startpos])) | (uint32(pbf.Pbf[startpos+1]) << 8) | (uint32(pbf.Pbf[startpos+2]) << 16)
			currentpos += 1
		//} else if startpos - startpos == 4 {
		case 4:
			vals[currentpos] = (uint32(pbf.Pbf[startpos])) | (uint32(pbf.Pbf[startpos+1]) << 8) | (uint32(pbf.Pbf[startpos+2]) << 16) + (uint32(pbf.Pbf[startpos+3]) * 0x1000000)
			currentpos += 1
		} 		
	}
	return vals[:currentpos]
}

func (pbf *PBF) ReadPoint(endpos int) []float64 {
	for pbf.Pos < endpos {
		x := pbf.ReadSVarintPower()
		y := pbf.ReadSVarintPower()
		return []float64{x,y}
	}
	return []float64{}
}

func (pbf *PBF) ReadLine(num int,endpos int) [][]float64 {
	var x,y float64
	if num == 0 {
		
		for startpos := pbf.Pos; startpos < endpos; startpos++  {
			if pbf.Pbf[startpos] <= 127 {
				num += 1
			} 
		}
		newlist := make([][]float64,num/2)
		
		for i := 0; i < num / 2; i ++ {
			x += pbf.ReadSVarintPower()
			y += pbf.ReadSVarintPower()
			newlist[i] = []float64{x,y}
		}
		
		/*
		newlist := [][]float64{}
		for pbf.Pos < endpos {
			x += pbf.ReadSVarintPower()
			y += pbf.ReadSVarintPower()
			newlist = append(newlist,[]float64{x,y})
		}
		*/


		return newlist
	} else {
		newlist := make([][]float64,num/2)

		for i := 0; i < num / 2; i ++ {
			x += pbf.ReadSVarintPower()
			y += pbf.ReadSVarintPower()
			newlist[i] = []float64{x,y}

		}
		return newlist
	}
	return [][]float64{}
}

func (pbf *PBF) ReadPolygon(endpos int) [][][]float64 {
	polygon := [][][]float64{}
	for pbf.Pos < endpos {
		num := pbf.ReadVarint()
		polygon = append(polygon,pbf.ReadLine(num,endpos))
	}
	return polygon
}	


func (pbf *PBF) ReadMultiPolygon(endpos int) [][][][]float64 {
	multipolygon := [][][][]float64{}
	for pbf.Pos < endpos {
		num_rings := pbf.ReadVarint()
		polygon := make([][][]float64,num_rings)
		for i := 0; i < num_rings; i++ {
			num := pbf.ReadVarint()
			polygon[i] = pbf.ReadLine(num,endpos)
		}
		multipolygon = append(multipolygon,polygon)
	}
	return multipolygon
}


func (pbf *PBF) ReadBoundingBox() []float64 {
	bb := make([]float64,4)
	pbf.ReadVarint()
	bb[0] = float64(pbf.ReadVarint()) / powerfactor
	bb[1] = float64(pbf.ReadVarint()) / powerfactor
	bb[2] = float64(pbf.ReadVarint()) / powerfactor
	bb[3] = float64(pbf.ReadVarint()) / powerfactor
	return bb

}	



func (pbf *PBF) ReadPackedUInt32() []uint32 {
	//startpos := pbf.Pos

	size := pbf.ReadVarint()
	arr := []uint32{}
	endpos := pbf.Pos + size

	for pbf.Pos < endpos {
		arr = append(arr,pbf.ReadUInt32())
	}

	return arr
}

func (pbf *PBF) ReadPackedUInt32_2() []uint32 {
	size := pbf.ReadVarint()
	endpos := pbf.Pos + size

	arr := make([]uint32,size)
	pos := 0
	for pbf.Pos < endpos {
		arr[pos] = pbf.ReadUInt32()
		pos++
	}
	
	return arr[:pos]
}



func (pbf *PBF) Byte() {
	fmt.Println(pbf.Pbf[pbf.Pos],"current")
	fmt.Println(pbf.Pbf[pbf.Pos:pbf.Pos+5],"next5")
}

func Reverse(val []byte) []byte {
	return []byte{val[7],val[6],val[5],val[4],val[3],val[2],val[1],val[0]}
}



