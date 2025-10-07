package main

import (
	"encoding/json"
	"fmt"

	"github.com/skysparq/grib2-go/file"
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

		emitProduct(rec.Product)
		emitGrid(rec.Grid)
		emitDataRepresentation(rec.DataRepresentation, rec)
	}
	println("\nNo errors")
}

func emitProduct(def record.Section4) {
	prodDef, recErr := def.Definition()
	if recErr != nil {
		panic(recErr.Error())
	}
	emitJson(prodDef)
}

func emitGrid(def record.Section3) {
	gridDef, recErr := def.Definition()
	if recErr != nil {
		panic(recErr.Error())
	}
	emitJson(gridDef)
}

func emitDataRepresentation(def record.Section5, rec record.Record) {
	drDef, err := def.Definition()
	if err != nil {
		panic(err.Error())
	}
	emitJson(drDef)
	data, err := drDef.GetValues(rec)
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
