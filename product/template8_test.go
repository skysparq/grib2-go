package product_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/skysparq/grib2-go/product"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestTemplate8(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordProdDef8)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := product.Template8{}.Parse(rec.ProductDefinition)
	if err != nil {
		t.Fatal(err)
	}
	expected := product.Template8{
		ProductDefinitionHeader: record.ProductDefinitionHeader{
			ParameterCategory: 19,
			ParameterNumber:   1,
		},
		GeneratingProcessType:       2,
		BackgroundIdentifier:        0,
		GeneratingProcessIdentifier: 96,
		HoursAfterReference:         0,
		MinutesAfterReference:       0,
		UnitOfTimeRange:             1,
		ForecastTimeInUnits:         42,
		FirstSurfaceType:            1,
		FirstSurfaceScaleFactor:     0,
		FirstSurfaceScaleValue:      0,
		SecondSurfaceType:           255,
		SecondSurfaceScaleFactor:    0,
		SecondSurfaceScaleValue:     0,
		EndYear:                     2025,
		EndMonth:                    3,
		EndDay:                      7,
		EndHour:                     6,
		EndMinute:                   0,
		EndSecond:                   0,
		TotalTimeRanges:             1,
		MissingDataValues:           0,
		TimeRanges: []product.TimeIncrement{
			{
				StatisticalProcess:         0,
				TimeIncrementType:          2,
				StatisticalUnitOfTimeRange: 1,
				StatisticalLengthOfTime:    6,
				SuccessiveUnitOfTimeRange:  255,
				SuccessiveLengthOfTime:     0,
			},
		},
	}
	typed := template.(product.Template8)
	if !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
	if date := time.Date(2025, 3, 7, 6, 0, 0, 0, time.UTC); date != typed.EndTime() {
		t.Fatalf(`expected end time %v, got %v`, date, typed.EndTime())
	}
	encoded, err := json.Marshal(template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(`%s`, encoded)
}
