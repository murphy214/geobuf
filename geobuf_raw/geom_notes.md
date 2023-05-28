
### Design Considerations

I am trying to add support for multiple coordinate dimensions but I am struggling with how to do it intuitively. The problem is I can hard code the coordinate sizes and blow out a function for each one or I can do a for loop within each point iteration and carry that through. This would allow me to use the same functions but the for loop probably has a performance cost. 

While I'm here I figure I should atleast add support for m-values as well. 