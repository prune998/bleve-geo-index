package main

import (
	"flag"
	"log"
	"os"

	"github.com/blevesearch/bleve/geo"

	"github.com/blevesearch/bleve"
)

var (
	blevePath = flag.String("blevepath", "index.bleve", "bleve index path")
	distance  = flag.String("distance", "1km", "search distance")
	search    = flag.String("search", "cafe", "search text")
	debug     = flag.Bool("debug", false, "show debug")
)

type Place struct {
	ID       string
	Name     string `json:"name" bleve:"name"`
	Amenity  string
	City     string
	Location []float64 `json:"location"`
}

func main() {
	flag.Parse()

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// cretae a new index
	mapping := bleve.NewIndexMapping()
	placeMapping := bleve.NewDocumentMapping()

	// add a mapping for places
	mapping.AddDocumentMapping("place", placeMapping)
	nameMapping := bleve.NewTextFieldMapping()
	nameMapping.Store = true
	placeMapping.AddFieldMappingsAt("name", nameMapping)

	// add a geo mapping which will hold the geographic data
	geoMapping := bleve.NewGeoPointFieldMapping()
	geoMapping.Store = true
	placeMapping.AddFieldMappingsAt("location", geoMapping)

	// define our mapping as the default mapping
	mapping.DefaultMapping = placeMapping

	// As this is a demo, remove the index at each run
	os.RemoveAll(*blevePath)

	// open the index
	index, err := bleve.New(*blevePath, mapping)
	if err != nil {
		log.Fatal(err)
	}

	// create a place and add it to the index
	place := &Place{
		ID:       "1",
		Name:     "School Secondary Les Etchemins",
		City:     "levis",
		Amenity:  "school",
		Location: []float64{-71.26917, 46.72009},
	}

	if *debug {
		log.Println(geo.ExtractGeoPoint(place.Location))
	}
	err = index.Index(place.ID, place)
	if err != nil {
		log.Fatal(err)
	}

	// create another place
	place = &Place{
		ID:       "2",
		Name:     "Quebec High School",
		City:     "quebec",
		Amenity:  "school",
		Location: []float64{-71.23966, 46.79775},
	}
	err = index.Index(place.ID, place)
	if err != nil {
		log.Fatal(err)
	}

	// for the demo purpose, close the DB
	index.Close()
	log.Println("write completed")

	// re open the index
	index, _ = bleve.Open(*blevePath)

	// search by term
	query1 := bleve.NewQueryStringQuery(*search)
	searchRequest1 := bleve.NewSearchRequest(query1)
	// grab sall fields from the entries
	searchRequest1.Fields = []string{"*"}
	searchResult1, _ := index.Search(searchRequest1)

	// for each result, display the score and the fields
	log.Printf("--------------------- searching for term ---------------------")
	for _, res := range searchResult1.Hits {
		log.Printf("ID: %v. Score %v.\n", res.ID, res.Score)

		for k, v := range res.Fields {
			log.Printf("Field %v. Value %v.\n", k, v)
		}
	}

	// search by distance and term
	lon, lat := -71.26050, 46.79049

	// search for term
	textQuery := bleve.NewMatchQuery(*search)
	// search by distance
	distanceQuery := bleve.NewGeoDistanceQuery(lon, lat, *distance)
	// set the distance search to use the right index field
	distanceQuery.SetField("location")

	// create an AND query with the two search criteria
	conRequest := bleve.NewConjunctionQuery()
	conRequest.AddQuery(textQuery)
	conRequest.AddQuery(distanceQuery)

	// do the search
	searchRequest := bleve.NewSearchRequest(conRequest)
	searchRequest.Fields = []string{"*"}
	searchResults, _ := index.Search(searchRequest)

	// display each hits with score and fields
	log.Printf("--------------------- searching for term and distance ---------------------")
	for _, res := range searchResults.Hits {
		log.Printf("distance %v - query : ID: %v. Score %v.\n", *distance, res.ID, res.Score)

		for k, v := range res.Fields {
			log.Printf("distance %v - Field %v. Value %v.\n", *distance, k, v)
		}
	}
}
