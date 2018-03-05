# geobuf_new

# What is it?

Protobuf representation of geojson features is something I've been trying for a while now. Previously things like not being able to implement bufio and having a bad protobuf struct to start with sort of limited the project. This is another attempt at creating such a structure and hopefully will be worth it. The protobuf structure looks like this: 

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

**As you can see it takes a few cues from vector-tile protobuf spec, this protobuf structure is essentially an intermediate feature object for a regular geojson object.** 

#### How do Geojson Features and Geobuf Features Differ

Geobuf is mean to be as close to a 1 to 1 geojson mapping as possible. Geobuf's geometry implementation uses delta encoding at a precision of 10e-7 which is like a few cm I think. Your file size mileage may very depending on your number of fields if you have a large number of fields like 50 in each feature your going to probably have a slighly bigger file than geojson but if you have like 5 and largeish geometries (not points) you should see a pretty signicant file size decrease. 

#### If it is based on a protobuf where is it? 

I've taken steps to remove the protobuf implementation. I still utilize my own reader and writer which is a little faster 30-50% but its mainly so that I can wrap the methods for creating the geometries in such within the read and the write. In the previous implementation you had to go from geojson to geobuf feature, although the end user couldn't see it this was a pretty needless allocation. Also Implementing my own reader / writer will allow me to do pretty cool thing with creating vector tiles which I'll detail somewhere else.

# Performance

Obviously at a single feature rate I/O is much much faster the problem previously was reading from a file iteratively was extremely hacky and I ended up splitting up code where I should have used an io.Reader. My new repo [protoscan](github.com/murphy214/protoscan) fixes this. 

Still as you can see as for a single feature read currently its > 10x faster. Of course this could vary drastically based on number of features vs. size of geometry etc but its a same assumption thats its much much faster for single feature reads.

However the FeatureCollection aren't exactly a one to one comparison as one is reading iteratively (geobuf) and another is reading the entire collection an once. A more close comparison would be line delimited geojson to geobuf which I may do later. 

```
goos: darwin
goarch: amd64
pkg: github.com/murphy214/geobuf
Benchmark_Read_FeatureCollection_Old-8    	      10	 137046994 ns/op	21370684 B/op	  638656 allocs/op
Benchmark_Read_FeatureCollection_New-8    	   50000	     24610 ns/op	   78248 B/op	       9 allocs/op
Benchmark_Read_Feature_Benchmark_Old-8    	    5000	    360952 ns/op	   70968 B/op	    2240 allocs/op
Benchmark_Read_Feature_Benchmark_New-8    	   50000	     27324 ns/op	   11856 B/op	     282 allocs/op
Benchmark_Write_Feature_Benchmark_Old-8   	    5000	    246050 ns/op	   26776 B/op	      23 allocs/op
Benchmark_Write_Feature_Benchmark_New-8   	   30000	     52173 ns/op	   19952 B/op	     564 allocs/op
PASS
ok  	github.com/murphy214/geobuf	9.893s
```


# Differences between previous implementation

* New proto struct
* Writer / Reader methods will be completely independent which will help alot as well use io.Reader / io.Writer
* Reader will use protoscan which hopefully should speed things up quite a bit
* Cleaner abstractions between pure read / writes and file or writer objects
