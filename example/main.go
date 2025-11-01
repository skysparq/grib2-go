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
	var indexed file.IndexedRecord
	i := 0
	for indexed, err = range grib.Records {
		if err != nil {
			panic(err.Error())
		}
		rec := indexed.Record

		emitProduct(rec.Product)
		processAndEmitGrid(rec.Grid)
		if i == 0 {
			processGridPoints(rec.Grid)
		}
		processAndEmitDataRepresentation(rec.DataRepresentation, rec)
		i++
	}
	println(fmt.Sprintf("\nNo errors parsing %v records\n", i))
}

func emitProduct(def record.Section4) {
	prodDef, recErr := def.Definition()
	if recErr != nil {
		panic(recErr.Error())
	}
	emitJson(prodDef)
}

func processAndEmitGrid(def record.Section3) {
	gridDef, recErr := def.Definition()
	if recErr != nil {
		panic(recErr.Error())
	}
	emitJson(gridDef)
}

func processGridPoints(def record.Section3) {
	gridDef, recErr := def.Definition()
	if recErr != nil {
		panic(recErr.Error())
	}
	points, err := gridDef.Points()
	if err != nil {
		panic(err.Error())
	}
	println(fmt.Sprintf(`points length=%d`, len(points.Lats)))
}

func processAndEmitDataRepresentation(def record.Section5, rec record.Record) {
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
