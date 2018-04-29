package geobuf

import (
	//g "github.com/murphy214/geobuf"
	"fmt"
	//"encoding/csv"
	"github.com/paulmach/go.geojson"
	"io"
	"os"
	"strings"
	//"time"
	//"io/ioutil"
)

func GetKeys(buf *Reader) ([]string, int) {
	keymap := map[string]string{}
	totalkeys := []string{}
	i := 0
	for buf.Next() {
		keys := ReadKeys(buf.Bytes())
		for _, key := range keys {
			_, boolval := keymap[key]
			if !boolval {
				keymap[key] = ""
				totalkeys = append(totalkeys, key)
			}
		}
		i++
	}
	totalkeys = append(totalkeys, []string{"Type", "Geometry"}...)
	buf.Reset()
	return totalkeys, i
}

func WriteRow(feature *geojson.Feature, keys []string) {
	feature.Properties["Type"] = string(feature.Geometry.Type)
	s, _ := feature.Geometry.MarshalJSON()

	feature.Properties["Geometry"] = string(s)
	newrow := make([]string, len(keys))
	for pos, key := range keys {
		val, boolval := feature.Properties[key]
		//fmt.Println(fmt.Sprint(val), val, feature.Properties)

		if !boolval {
			val = ""
		}
		newrow[pos] = fmt.Sprint(val)
	}
	io.WriteString(os.Stdout, strings.Join(newrow, ",")+"\n")
}

func ReadGeobufCSV(filename string) {
	buf := ReaderFile(filename)
	keys, _ := GetKeys(buf)
	io.WriteString(os.Stdout, strings.Join(keys, ",")+"\n")
	myfunc := func(feature *geojson.Feature) interface{} {
		WriteRow(feature, keys)
		return ""
	}

	buf2 := NewGeobufReaderFunction(buf, myfunc)
	for buf2.Next() {
		buf2.Value()
	}
}
