package product_test

import (
	"testing"

	"github.com/skysparq/grib2-go/file"
	"github.com/skysparq/grib2-go/product"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestParseDefinitionOnGfsGrib(t *testing.T) {
	_, r, err := test_files.Load(test_files.FullGfsFile)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	parser := &product.Parser{}
	grib := file.NewGribFile(r, templates.Version33())
	var rec record.Record
	for rec, err = range grib.Records {
		if err != nil {
			t.Fatal(err)
		}
		_, err = parser.ParseDefinition(rec.ProductDefinition)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestParserTemplatesDoNotChange(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordProdDef0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	parser := &product.Parser{}
	_, err = parser.ParseDefinition(rec.ProductDefinition)
	if err != nil {
		t.Fatal(err)
	}
	template, ok := parser.Templates[0].(product.Template0)
	if !ok {
		t.Fatal("template was not a Template0")
	}
	expected := product.Template0{}
	if template != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, template)
	}
}
