
### Design Considerations

I am trying to add support for multiple coordinate dimensions but I am struggling with how to do it intuitively. The problem is I can hard code the coordinate sizes and blow out a function for each one or I can do a for loop within each point iteration and carry that through. This would allow me to use the same functions but the for loop probably has a performance cost. 

While I'm here I figure I should atleast add support for m-values as well. 

UPDATE: I ended up doing it generically with little performance impact from what I can tell so far. 

# Geometry Type Code

To deal with backwards compability I used a geometry code that encodes the dim_size into its value. The first four bits are the geometry type integer and last four are the dimension size. 

If this code embedded is less than or equal to 6 we just use that value as the geometry type integer and assume the dim_size is 2. 

