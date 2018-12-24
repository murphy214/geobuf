# Geobuf - A GeoJSON Interchange Format 
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/murphy214/geobuf)

#### What is it?

From a top level geobuf provides simple apis to convert from geojson, read, and write geobuf. Geobuf is a custom protobuf implementation of geojson features, its much much faster than json unmarshalling as well as much smaller for geometry heavy features. Performance for reading and writing can be summarized as about 5-10x what your going to see out of plain json. **However with some larger files reading concurrently is like 18x faster (1 gb california roads geojson)** 

#### Why Should I consider Using This?

Beyond being much faster for serialization reads can be done **iteratively** and more importantly piece wise, so one could do a partial read of just the values needed for a filter than read the entire feature if those conditions are satisified. 

I think thats the main deviation from mapbox's geobuf is the use of flat geojson features at the top level. The problem with mapbox's geobuf format is it can't be stream because in order to assemble an entire geojson feature from MB's geobuf you need the end of the file being the key, value lists, and I also recall there being another structure you needed to interate in parallel in order to assemble the geojson feature. That effectively relegates it to a faster & smaller geojson feature collection implementation, but it really doesn't solve the problem I was having, **needing to deserialize every feature in the feature collection in memory before any of the features can be operated on.** 

As for the pretty beefy properties sizes that my features have, instead of using the value / key lists basically without any super formal testing the differences between the two are pretty minimal when gzipped which is about what you would expect anyway.

#### Features 

* Straightforward Reader / Writer methods for everything thing that is done to geobufs 
* **Expect at least 5-10x performance gains in both read / write against line-delimited geojson**
* CLI executables to convert to and from geobuf from geojson 
* Inplace geobuf sorts to do things like feature mapping about a file out of memory for tiling (AKA mapreduce)
* Backend for filters for quick visualizations

#### Internals

Below is the geobuf proto file that this implementation is based. 

```
syntax = "proto3";

// Variant type encoding
// The use of values is described in section 4.1 of the specification
message Value {
	// Exactly one of these values must be present in a valid message
	string string_value = 1;
	float float_value = 2;
	double double_value = 3;
	int64 int_value = 4;
	uint64 uint_value = 5;
	sint64 sint_value = 6;
	bool bool_value = 7;
}

// GeomType is described in section 4.3.4 of the specification
enum GeomType {
	UNKNOWN = 0;
	POINT = 1;
	LINESTRING = 2;
	POLYGON = 3;
	MULTIPOINT = 4;
	MULTILINESTRING = 5;
	MULTIPOLYGON = 6;
}

message Feature {
	uint64 Id = 1;
	map<string, Value>  Properties = 2;
 	GeomType type = 3;
	repeated uint64 Geometry = 4 [ packed = true ];
	repeated int64 BoundingBox = 5 [ packed = true ]; // N,S,E,W
}

message FeatureCollection {
	repeated Feature Features = 1;
}
```


#### Usage 

Reading and writing geojson features is pretty easy there are two types of data stores, raw byte buffers or file buffers. They basically work the same way but one reads / writes from an underlying file while the other reads / writes from a buffer.

**Currently I don't suggest using concurrent reads because the speed gains while signicant don't make up for the brittleness of the API currently more work has to be done to make this more useful.** 

That being said plain reads are far faster than line-delimitted JSON and still offer a unified API. 