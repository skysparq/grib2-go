package record

import (
	"fmt"
	"io"
)

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
	//if len(values.Values) != len(values.Lngs) || len(values.Values) != len(values.Lats) {
	//	return values, fmt.Errorf("error getting gridded values: the length of Values, Lngs, and Lats do not match")
	//}
	return values, nil
}
