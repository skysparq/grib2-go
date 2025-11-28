package record

import (
	"fmt"
	"io"
)

// Record represents all 7 sections of a GRIB record, excluding section 8 which only contains '7777' as a marker for the end of the record.
// Record is the main data structure for processing GRIB2 records. From here, the metadata, geospatial points, and data values can be retrieved.
type Record struct {
	Indicator          Section0
	Identification     Section1
	LocalUse           Section2
	Grid               Section3
	Product            Section4
	DataRepresentation Section5
	BitMap             Section6
	Data               Section7
}

// ParseRecord parses a GRIB record from an io.Reader and the provided templates.
// It parses all 8 sections of the record and returns a Record.
func ParseRecord(r io.Reader, templates Templates) (record Record, err error) {
	record.Indicator, err = ParseSection0(r)
	if err != nil {
		return record, fmt.Errorf("error parsing record: %w", err)
	}

	totalLength := record.Indicator.GribLength
	readLength := 16
	var data SectionData
	for {
		data, err = readVariableLengthSection(r)
		if err != nil {
			return record, fmt.Errorf("error parsing record: %w", err)
		}
		readLength += data.Length
		switch data.SectionNumber {
		case 1:
			record.Identification, err = ParseSection1(data)
		case 2:
			record.LocalUse, err = ParseSection2(data)
		case 3:
			record.Grid, err = ParseSection3(data, templates)
		case 4:
			record.Product, err = ParseSection4(data, templates)
		case 5:
			record.DataRepresentation, err = ParseSection5(data, templates)
		case 6:
			record.BitMap, err = ParseSection6(data)
		case 7:
			record.Data, err = ParseSection7(data)
		default:
			err = nil
		}
		if readLength > totalLength-4 {
			return record, fmt.Errorf("error parsing record: the GRIB record appears to be malformed")
		}
		if readLength == totalLength-4 {
			var section8 []byte
			section8, err = readFixedLengthSection(r, 4)
			if err != nil {
				return record, fmt.Errorf("error parsing record: %w", err)
			}
			if string(section8) != "7777" {
				return record, fmt.Errorf("error parsing record: the GRIB record does not end with 7777")
			}
			break
		}
		if err != nil {
			return record, fmt.Errorf("error parsing record: %w", err)
		}
	}

	return record, nil
}

// GetGriddedValues returns the longitude, latitude, and decoded data value of every data point in the record.
// It is a convenience method to retrieve the entirety of the geospatial data in a single call.
//
// Note: When processing multiple records that use the same grid, this method adds CPU and memory overhead that can be
// avoided by retrieving the data values directly from the DataRepresentationDefinition. The longitude and latitude
// should be retrieved once from the GridDefinition and cached for the remaining records.
func (r Record) GetGriddedValues() (GriddedValues, error) {
	var values GriddedValues
	grid, err := r.Grid.Definition()
	if err != nil {
		return values, fmt.Errorf("error getting gridded values: %w", err)
	}
	values.GridPoints, err = grid.Points()
	if err != nil {
		return values, fmt.Errorf("error getting gridded values: %w", err)
	}
	dataRep, err := r.DataRepresentation.Definition()
	if err != nil {
		return values, fmt.Errorf("error getting gridded values: %w", err)
	}
	values.Values, err = dataRep.GetValues(r)
	if err != nil {
		return values, fmt.Errorf("error getting gridded values: %w", err)
	}
	values.XVals = grid.XVals()
	values.YVals = grid.YVals()
	if len(values.Values) != len(values.Lngs) || len(values.Values) != len(values.Lats) {
		return values, fmt.Errorf("error getting gridded values: the length of ValuesIterator, Lngs, and Lats do not match")
	}
	return values, nil
}
