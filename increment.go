package geobuf

import (
	"bytes"
	"fmt"
)

// increments returning the given bytes of a feature collection
func Increment(buf *Reader, increment int) ([]byte, bool) {

	start := []byte(`{"type": "FeatureCollection", "features": [`)
	end := []byte(`]}`)
	bytebuf := bytes.NewBuffer(start)
	i := 0
	for buf.Next() && i < increment {
		feature := buf.Feature()
		feature.Properties["COLORKEY"] = "purple"
		bytevals, err := feature.MarshalJSON()
		if err != nil {
			fmt.Println(err)
		}
		bytebuf.Write(bytevals)
		bytebuf.WriteByte(',')

		i++
	}
	buf.FeatureCount--
	bytevals := bytebuf.Bytes()
	bytevals = bytevals[:len(bytevals)-1]

	return append(bytevals, end...), increment == i
}
