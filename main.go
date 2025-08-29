package main

import (
	"encoding/json"
	"fmt"

	"github.com/skysparq/grib2-go/data_representation"
	"github.com/skysparq/grib2-go/file"
	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/product"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func main() {
	_, r, err := test_files.Load(test_files.FullGfsFile)
	if err != nil {
		panic(err.Error())
	}
	defer func() { _ = r.Close() }()

	grib := file.NewGribFile(r, templates.Version33())
	var rec record.Record
	for rec, err = range grib.Records {
		if err != nil {
			panic(err.Error())
		}

		emitProduct(rec.ProductDefinition)
		emitGrid(rec.GridDefinition)
		emitDataRepresentation(rec.DataRepresentation, rec)
	}
	println("\nNo errors")
}

func emitProduct(def record.Section4) {
	parser := &product.Parser{}
	prod, recErr := parser.ParseDefinition(def)
	if recErr != nil {
		panic(recErr.Error())
	}
	emitJson(prod)
}

func emitGrid(def record.Section3) {
	parser := &grid.Parser{}
	prod, recErr := parser.ParseDefinition(def)
	if recErr != nil {
		panic(recErr.Error())
	}
	emitJson(prod)
}

func emitDataRepresentation(def record.Section5, rec record.Record) {
	parser := &data_representation.Parser{}
	prod, err := parser.ParseDefinition(def)
	if err != nil {
		panic(err.Error())
	}
	emitJson(prod)
	data, err := prod.GetValues(rec)
	if err != nil {
		panic(err.Error())
	}
	println(fmt.Sprintf(`data length=%d`, len(data)))
}

func emitJson(obj any) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}
	println(string(jsonBytes))
}
