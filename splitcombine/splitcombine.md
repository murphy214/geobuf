
# What is a split combine method

Often times I've found when dealing with a large geobuf a little bit of context can go a hell of a long way. I recently constructed a hack that can sniff the first feature in a geobuf file for a certain configuration that significies this features properties contain a metadata gob. This metadata gob can carry useful things like bounding box for the whole dataset but the main reason for its use is indexing which parts of a file pertain to a certain property(s). 

In short this gives me the ability to express something like go to part of the file that contains features that are related to this tile and iterate through all of them. Them all being sequentially together removes the need for slow file.ReadAt and huge feature maps.

Therefore the split combine methods ethos is pretty straight forward as you iterate through a large geobuf file we map a feature to a temporary sub file that is written to until you've gone through all the features. After, write out the metadata generated from the file context we have and create a new geobuf file with that single feature containing the metadata. (effectively a dummy feature)

From there use this bash command 

'''
cat file1,file2 > newfile
'''

This appends all the files to the newgeobuf with just the metadata. Clean up your mess of the split files by removing all of them and finally deleting the orginal file and renaming the new to the original filename. 

TODO:
    - Implement toplevel structure to add files write features to them and dump that context into meta.
