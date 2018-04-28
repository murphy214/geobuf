# Geobuf - A geojson interchange format based on protocol buffers, for large bulk geojson feature processing

# What is it?

From a top level geobuf provides simple apis to convert from geojson, read, and write geobuf. Geobuf is a custom protobuf implementation of geojson features, its much much faster than json unmarshalling as well as much smaller for geometry heavy features. Performance for reading and writing can be summarized as about 5-10x what your going to see out of plain json. **However with some larger files read concurrently performance for reads could be like 18x (1 gb california roads geojson)** 

# Features 

* One-to-One geojson mapping (no geom collections) in 4326 coordinate system at a 10e-7 precison which is like 30 cm if I remember correctly. 
* Simple Writer Api for both byte buffers and file buffers implemented essentially the same way (other than intialization)
* Simple Reader Api for byte and file buffers implemented in esssentially the same way (other than intialization)
* Command line functions for converting between geobuf and geojson (I still haven't done a javascript impl.)
* Reader / writer implements iterative reading / writing meaning only one feature is being brought into memory at a time.
* Allows for high level abstractions to be applied about a geobuf file features indicies (the positions in a file between where the bytes for a given feature exist) **Allowing for out of memory mapping of a feature properties, types, etc. that can be directly used to create new files or sort geobuf features in place** 
* Tries to implement io.Reader / io.Writer as much as possible, in order to more effecentially write to files a flush is done after several mbs of data is allocated, **unfortunately this means when doing custom writes you will always have to trail a buf.Flush() method to write all existing data to the underlying writer in exchange that data is allocated as a buffer until flushed making writes to file much larger and more effecient.**

# Differences from MB Geobuf 

Mapbox's geobuf was obviously the inspiration for this but there were several issues I found with the implementation (which is probably why its not super active right now) I took several steps to make my life easier when I built my geobuf thing. 

* One resolution, coordinate system and dimmension size is currently supported (none of these are super embedded into the implementation but simplicity > robustness)
* Features are flat, can be read without any underlying structure from within the protobuf (no keys / values array). I found the values array really the most annoying thing with MB's geobuf, for one gzipping nullifies most gains you'd see from the values array, and it introduces the same problem with vector tile implementations: any time you want to abstract anything, you have to drag the values array through your implementation as its a global for your entire file. (Although in vt it makes more sense) 

# Future Features / Half-implemented
* File indicie support (returns positional indexs within file) for future sorting / mapping functionality
	1. You could first apply a mapping function against a feature,ind read looking at properties, geom type etc., then create a map of the indicies you want to group by meaning you could have all the features of certain tile id, in one spot and also return the indicie positions for that entire tile id effectively allowing you to create a entirely new geobuf from the new file. 
	2. Expanding on step 1 this would of course be used with a write for an entire new geobuf (that would replace the underlying geobuf its based on after complete) 
	3. Leaving you with the ability to do things like sort / map fields to an entirely new file, which for me was the main disadvantage of using raw protobufs I had to convert to base64 to use the standard unix sort or simply use json. This allows you to sort the entire file, with the only constraint in memory being your map[field][][2]int which of course is a list of the 2 index positions of each feature. 
	
* A stable concurrent reader / read mapping for an arbitary function
* write in place -> newgeobuf -> several geobufs derived from map
* Could Experiment with implementing different coordinate systems more condusive to creating vector tiles, I've already found that my vector tile writer benifits from the reduced allocations of sending a raw feature byte array into the add feature api, basically reading only the parts I need when I need them encapsolated within the write feature function.  

# Usage (deprecated)

Below I gave a few of the mainline apis for geobuf and how they are used.

# The Reader API

The reader api is pretty simple, nothing to crazy going on here the implementation is pretty simple really.

### Reading a Geobuf File

The reader api is pretty straight forward as you can see below.

```golang
package main

import (
	g "github.com/murphy214/geobuf"
}

func main() {
	// reading geobuf from file
	buf := g.ReaderFile("wv.geobuf")
	for buf.Next() {
		// do something with the feature here
		buf.Feature()
	}
}
```

### Reading from a Geobuf Buffer (byte array)

```golang
// reading geobuf from a []byte array
bytevals,_ := ioutil.ReadFile("wv.geobuf")
buf := g.ReaderBuf(bytevals)
for buf.Next() {
	buf.Feature()
}
```

### Using the Reader To Only Return Byte Arrays of a Feature
```golang
// as you can see both readers are implemented the exact same way but you can also
// retrieve the raw bytevalues of a reader by doing this
buf = g.ReaderBuf(bytevals)
for buf.Next() {
	// do somethign with raw byte values
	buf.Bytes()
}
```

### Resetting a Reader so That it can be read again

As you might have guessed geobuf readers are pretty stateful to some degree so to read on the same reader twice you must reset it. 

```golang
// resetting a reader so that it can be read again
buf = g.ReaderBuf(bytevals)
for buf.Next() {
	// do somethign with raw byte values
	buf.Bytes()
}
// resetting reader here
buf.Reset()

// readiing the same buffer
buf = g.ReaderBuf(bytevals)
for buf.Next() {
	// do somethign with raw byte values
	buf.Bytes()
}
```

# The Writer Api 

There are 4 functions to instantiate the writer api being: ```WriterFile(filename)```,```WriterFileNew(filename)```,```WriterBuf(bytearray)```,```WriterBufNew()``` these methods are pretty self explanatory on what they do and how they differ. 

The writer api has a little more going on in terms of complexity and useful methods so I'll go over some.

### Writing to a Geobuf from Geojson Features

```golang
newbuf := WriterFileNew("newbuf.geobuf")
newbuf.WriteFeature(feature) // where feature is *geojson.Feature
```

### Writing to a Geobuf from a Geobuf Feature Byte Array

```golang
newbuf := WriterFileNew("newbuf.geobuf")
newbuf.Write(bytearray) // where byte array is an array of a geobuf feauture
```

### Adding a buffer Writer to another Writer

Sometimes you may want to write one geobuf to another, while currently I don't implement this for file geobuf (it could be done) it does work if the geobuf your wanting to add is a buffer. 

```golang
newbuf := WriterFileNew("newbuf.geobuf")
newbuf.AddGeobuf(oldbuf) // where old buf is a buf geobuf
```

### Converting a Writer Geobuf to a Reader Geobuf

Sometimes we may want to take something that use to be a writer and immediately use it as a reader this can be done using the reader method.

```golang
newbufwriter := WriterFileNew("newbuf.geobuf")
newbufreader := newbufwriter.Reader()
```

# Differences between previous implementation

* New proto struct
* Writer / Reader methods will be completely independent which will help alot as well use io.Reader / io.Writer
* Reader will use protoscan which hopefully should speed things up quite a bit
* Cleaner abstractions between pure read / writes and file or writer objects
