package main

import (
	"encoding/json"
	"os"

	"github.com/skysparq/grib2/file"
	"github.com/skysparq/grib2/grid"
	"github.com/skysparq/grib2/product"
	"github.com/skysparq/grib2/record"
	"github.com/skysparq/grib2/templates"
)

func main() {
	r, err := os.Open(`./.large_test_files/full_gfs_file.grb2`)
	if err != nil {
		panic(err.Error())
	}
	defer func() { _ = r.Close() }()

	grib := file.NewGribFile(r, templates.Revision20120111())
	var rec record.Record
	for rec, err = range grib.Records {
		if err != nil {
			panic(err.Error())
		}

		emitProduct(rec.ProductDefinition)
		emitGrid(rec.GridDefinition)
	}
	println("\nNo errors")
}

func emitProduct(def record.Section4) {
	prod, recErr := product.ParseDefinition(def)
	if recErr != nil {
		panic(recErr.Error())
	}
	emitJson(prod)
}

func emitGrid(def record.Section3) {
	prod, recErr := grid.ParseDefinition(def)
	if recErr != nil {
		panic(recErr.Error())
	}
	emitJson(prod)
}

func emitJson(obj any) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}
	println(string(jsonBytes))
}
