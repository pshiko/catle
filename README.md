# catle

head command for csv

![sample](https://github.com/pshiko/catle/blob/for_readme/tes2.gif)

Install
```bash
go get github.com/pshiko/catle
go install github.com/pshiko/catle/cmd/csvhead
```

[Usage]

	csvhead <csvpath>
  
[Options]

	-nh:  without header
  
	-t: set '\t' as delimiter
  
	-s: set white space as delimiter
  
	-n <number>: skip <number> row
  
  


# Command

hjkl -> Move next

\<C-hjkl\> -> Move half window size

s -> Sort ascending

S -> Sort descending

\<C-I\> -> Change current column to Int Column

\<C-S\> -> Change current column to String Column

\<space\> -> Shrink current column.
q -> quite window
