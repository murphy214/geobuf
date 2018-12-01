## Geobuf Demo 

Geobuf is a file format built for processing large geojson files. Below is a demo of one of the potential uses of geobuf.

### Installing proper packages

Copy/Paste the following get get command into your terminal.

```
go get -u github.com/murphy214/lazyosm
go get -u github.com/murphy214/lazyosm/lazyosm
go get -u github.com/murphy214/geobuf/cmd/geobuf2geojson
go get -u github.com/murphy214/geobuf/cmd/geojson2geobuf
```

A few more repositories may have to be retrieved. Just use ```go get ...```


### Downloading the OSM Data

For this example were going to be using a relatively large dataset being osm feature in Texas.

```
curl -o texas.pbf http://download.geofabrik.de/north-america/us/texas-latest.osm.pbf
```

### Generating End User Features Into The Geobuf Format

This CLI command generates all the OSM geojsoon features with Texas to the best of the pbf parsers ability. It basically returns as generic OSM geojson features as possible. The ```-l``` 
argument is for how many node map blocks can be opened at once which is basically the constraint on memory. I would start with 3000 depending on your machine's memory.

```
lazyosm make -f texas.pbf -o texas.geobuf -l 5000
```

### Mapping All The Features Into A Sub-File Structure

This is where the cool stuff starts to happen (I think at least). The following script maps the entire texas geobuf file of 5-6 million features and effectively adds an index of size 12 TileID to the geobuf file. 

```go
package main

import (
	"github.com/murphy214/tileclip"
	"github.com/murphy214/geobuf/splitcombine"
	"github.com/paulmach/go.geojson"
	m "github.com/murphy214/mercantile"
)

func main() {
	// creating the tile configuration
	tileconfig := &splitcombine.TileConfig{
		Bounds:m.Extrema{N:90.0,S:0.0,E:0.0,W:-180.0},
		OutputFileName:"new_texas.geobuf",
		Zoom:12,
		Concurrent:10000,
	}

	// creating the tile mapping function
	mymap := func(feature *geojson.Feature) map[m.TileID]*geojson.Feature {
		return tileclip.ClipFeature(feature,tileconfig.Zoom,false)
	}

	splitcombine.MapGeobuf("texas.geobuf", mymap, tileconfig)
}
```

### Using the Sub File Structure 

The example code below shows how you would use the sub file structure within the geobuf to get all the features within a particular tile.

```go
package main

import (
	"github.com/murphy214/tileclip"
	g "github.com/murphy214/geobuf"
	"github.com/paulmach/go.geojson"
	m "github.com/murphy214/mercantile"
)

func main() {
	buf := g.ReaderFile("new_texas.geobuf")
	k := m.TileID{943,1652,12}
	
	feats := []*geojson.Feature{}
	buf.SubFileSeek(m.TilestrFile(k))
	for buf.SubFileNext() {
		feats = append(feats,buf.Feature())
	}
	
	tileclip.MakeFeatures(feats,"a.geojson")
}
```