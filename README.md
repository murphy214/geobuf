# Geobuf - A geojson interchange format based on protocol buffers, for large bulk geojson feature processing
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/murphy214/vector-tile-go)

#### What is it?

From a top level geobuf provides simple apis to convert from geojson, read, and write geobuf. Geobuf is a custom protobuf implementation of geojson features, its much much faster than json unmarshalling as well as much smaller for geometry heavy features. Performance for reading and writing can be summarized as about 5-10x what your going to see out of plain json. **However with some larger files reading concurrently is like 18x faster (1 gb california roads geojson)** 

#### Why Should I consider Using This?

Beyond being much faster for serialization reads can be done **interatively** and more importantly piece wise, so one could do a partial read of just the values needed for a filter than read the entire feature if those conditions are satisified. 

I think thats the main deviation from mapbox's geobuf is the use of flat geojson features at the top level. The problem with mapbox's geobuf format is it can't be stream because in order to assemble an entire geojson feature from MB's geobuf you need the end of the file being the key, value lists, and I also recall there being another structure you needed to interate in parallel in order to assemble the geojson feature. That effectively relegates it to a faster & smaller geojson feature collection implementation, but it really doesn't solve the problem I was having, **needing to deserialize every feature in the feature collection in memory before any of the features can be operated on.** 

As for the pretty beefy properties sizes that my features have, instead of using the value / key lists basically without any super formal testing the differences between the two are pretty minimal when gzipped which is about what you would expect anyway.

