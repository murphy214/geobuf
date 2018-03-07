package geobuf_stdout

import (
	"github.com/paulmach/go.geojson"
	"sync"
	"io"

	"os"
	g "github.com/murphy214/geobuf"
	"fmt"
	"strings"
)

type Column struct {
	Column_Map map[string]int
	Column []string
	M *sync.Mutex
}

func (column *Column) Add_Column(i string) {
	//*newcolumn = *column
	column.Column = append(column.Column,i)
	column.Column_Map[i] = len(column.Column) - 1		
	
}

func Get_Row(feature  *geojson.Feature,column *Column) {
	feature.Properties[`Type`] = string(feature.Geometry.Type)

	switch feature.Geometry.Type {
	case "Point":
		feature.Properties[`Geometry`] = strings.Replace(fmt.Sprintf(`"%v"`,feature.Geometry.Point)," ",",",-1)
	
	case "LineString":
		feature.Properties[`Geometry`] = strings.Replace(fmt.Sprintf(`"%v"`,feature.Geometry.LineString)," ",",",-1)
	case "Polygon":
		feature.Properties[`Geometry`] = strings.Replace(fmt.Sprintf(`"%v"`,feature.Geometry.Polygon)," ",",",-1)
	case "MultiPoint":
		feature.Properties[`Geometry`] = strings.Replace(fmt.Sprintf(`"%v"`,feature.Geometry.MultiPoint)," ",",",-1)
	case "MultiLineString":
		feature.Properties[`Geometry`] = strings.Replace(fmt.Sprintf(`"%v"`,feature.Geometry.MultiLineString)," ",",",-1)
	case "MultiPolygon":
		feature.Properties[`Geometry`] = strings.Replace(fmt.Sprintf(`"%v"`,feature.Geometry.MultiPolygon)," ",",",-1)
	}

	column.M.Lock()
	
	for k := range feature.Properties {

		_,boolval := column.Column_Map[k]
		if boolval ==  false {
			column.Add_Column(k)
		}

	}
	column.M.Unlock()

	//newlist := []interface{}{}
	newlist := []string{}
	for _,column := range column.Column {
		val,boolval := feature.Properties[column]
		if boolval == false {
			val = ""
		} 
		newlist = append(newlist,fmt.Sprint(val))
	}

	io.WriteString(os.Stdout,strings.Join(newlist,"\t") + "\n")

}

func ReadGeobuf(geobuf *g.Reader) {
	var mm sync.Mutex
	column := &Column{Column_Map:map[string]int{},M:&mm}
	for geobuf.Next() {
		Get_Row(geobuf.Feature(),column)
	}
	io.WriteString(os.Stdout,strings.Join(column.Column,"\t"))
}



