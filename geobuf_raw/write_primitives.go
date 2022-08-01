package geobuf_raw

import (
	"github.com/murphy214/pbf"
	"math"
	"reflect"
)

var powerfactor = math.Pow(10.0, 7.0)

// encodes a var int for 32 bit number
func EncodeVarint32(x uint32) []byte {
	var buf [4]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

func AppendAll(b ...[]byte) []byte {
	total := 0
	for _, i := range b {
		total += len(i)
	}
	pos := 0
	totalbytes := make([]byte, total)
	for _, i := range b {
		for _, byteval := range i {
			totalbytes[pos] = byteval
			pos += 1
		}
	}
	return totalbytes
}

// writes a packed uint32 number
// this function was benchmarked against several implementations
// to reduce allocations, i found this one to be the best
func WritePackedUint64_2(geom []uint64) []byte {
	buf := make([]byte, len(geom)*8+8)
	pos := 8
	for _, x := range geom {
		for x > 127 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
		}
		buf[pos] = uint8(x)
		pos++
	}
	beg := pbf.EncodeVarint(uint64(pos - 8))
	startpos := 8 - len(beg)
	currentpos := startpos
	i := 0
	for currentpos < 8 {
		buf[currentpos] = beg[i]
		i++
		currentpos++
	}

	return buf[startpos:pos]
}

// writes a packed uint32 number
// this function was benchmarked against several implementations
// to reduce allocations, i found this one to be the best
func WritePackedUint64(geom []uint64) []byte {
	buf := make([]byte, len(geom)*8+8)
	pos := 8
	for _, x := range geom {
		if x < 128 {
			buf[pos] = uint8(x)
			x >>= 7
			pos++
		} else if x < 16384 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		} else if x < 2097152 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		} else if x < 268435456 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		} else if x < 34359738368 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		} else if x < 4398046511104 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		} else if x < 562949953421312 {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		} else {
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = 0x80 | uint8(x&0x7F)
			x >>= 7
			pos++
			buf[pos] = uint8(x)
			pos++
		}

	}
	beg := pbf.EncodeVarint(uint64(pos - 8))
	startpos := 8 - len(beg)
	currentpos := startpos
	i := 0
	for currentpos < 8 {
		buf[currentpos] = beg[i]
		i++
		currentpos++
	}

	return buf[startpos:pos]
}

const maxVarintBytes = 10

// encodes are var int value
func EncodeVarint_Value(x uint64, typeint int) []byte {
	var buf [maxVarintBytes]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	total := []byte{18, byte(n + 1), byte(typeint)}
	return append(total, buf[0:n]...)
}

// writes a float32 into a float value
func FloatVal32(f float32) []byte {
	buf := make([]byte, 4)
	n := math.Float32bits(f)
	buf[3] = byte(n >> 24)
	buf[2] = byte(n >> 16)
	buf[1] = byte(n >> 8)
	buf[0] = byte(n)
	return append([]byte{18, 5, 21}, buf...)
}

// writes a float64 into a double value
func FloatVal64(f float64) []byte {
	buf := make([]byte, 8)
	n := math.Float64bits(f)
	buf[7] = byte(n >> 56)
	buf[6] = byte(n >> 48)
	buf[5] = byte(n >> 40)
	buf[4] = byte(n >> 32)
	buf[3] = byte(n >> 24)
	buf[2] = byte(n >> 16)
	buf[1] = byte(n >> 8)
	buf[0] = byte(n)
	return append([]byte{18, 9, 25}, buf...)
}

// writes a value and returns the bytes of such value
// does not implement write sint currently
func WriteValue(value interface{}) []byte {
	vv := reflect.ValueOf(value)
	kd := vv.Kind()

	// switching for each type
	switch kd {
	case reflect.String:
		if len(vv.String()) >= 0 {
			size := uint64(len(vv.String()))
			size_bytes := pbf.EncodeVarint(size)
			bytevals := []byte{10}
			bytevals = append(bytevals, size_bytes...)
			bytevals = append(bytevals, []byte(vv.String())...)
			bytevals = append(pbf.EncodeVarint(uint64(len(bytevals))), bytevals...)
			return append([]byte{18}, bytevals...)
		}
	case reflect.Float32:
		return FloatVal32(float32(vv.Float()))
	case reflect.Float64:
		return FloatVal64(float64(vv.Float()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return EncodeVarint_Value(uint64(vv.Int()), 32)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return EncodeVarint_Value(uint64(vv.Uint()), 40)
	case reflect.Bool:
		if vv.Bool() == true {
			return []byte{18, 2, 56, 1}
		} else if vv.Bool() == false {
			return []byte{18, 2, 56, 0}
		}
	}

	return []byte{}
}
