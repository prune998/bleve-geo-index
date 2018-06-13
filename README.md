# bleve-geo-index

## description
This `program` use the [Bleve DB](http://www.blevesearch.com/) to index geographic data

This is a pure demo to test geo capabilities. It is not a production code and may not use all the features Bleve can offer.

## build
just `cd` in this directory then install Bleve and build
```
go get -u github.com/blevesearch/bleve
go build
```

## run
use `./bleve-geo-index -h` to get info of the command line

- -blevepath: path of the index file to create. If existing, this folder will be removed ! 
- -distance: the distance for the search. Ex : 1km, 1000m
- -search: the keyword to search in the index. As we just index 2 schools, try searching for `school`
- -debug: add some debug (really few for the moment)