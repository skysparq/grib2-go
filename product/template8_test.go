package product_test

import (
	"github.com/skysparq/grib2-go/product"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"os"
	"testing"
)

func TestTemplate8(t *testing.T) {
	f, err := os.Open("../test_files/single-grib2-record-prod-def-8.grb2")
	if err != nil {
		t.Fatal(err)
	}
	rec, err := record.ParseRecord(f, templates.Revision20120111())
	if err != nil {
		t.Fatal(err)
	}
	template := &product.Template8{}
	err = template.Parse(rec.ProductDefinition)
	if err != nil {
		t.Fatal(err)
	}
}
