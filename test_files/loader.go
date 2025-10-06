package test_files

import (
	"fmt"
	"os"
	"path/filepath"
)

type TestFile string

const (
	SingleRecordProdDef8  TestFile = `.test_files/single-grib2-record-prod-def-8.grb2`
	SingleRecordProdDef0  TestFile = `.test_files/single-grib2-record-prod-def-0.grb2`
	FullGfsFile           TestFile = `.test_files/full-gfs-file.grb2`
	SingleRecordDataDef40 TestFile = `.test_files/single-grib2-record-data-def-40.grb2`
	SingleRecordDataDef41 TestFile = `.test_files/single-grib2-record-data-def-41.grb2`
	SingleRecordGridDef40 TestFile = `.test_files/single-grib2-record-grid-def-40.grb2`
)

func Load(file TestFile) (int, *os.File, error) {
	wd, err := os.Getwd()
	if err != nil {
		return 0, nil, fmt.Errorf("error loading test file: %v", err)
	}
	for {
		if wd == `/` {
			return 0, nil, fmt.Errorf("did not find project root directory, project directory must be called 'grib2-go'")
		}
		if filepath.Base(wd) == `grib2-go` {
			break
		}
		wd = filepath.Dir(wd)
	}

	path := filepath.Join(wd, string(file))

	stat, err := os.Stat(path)
	if err != nil {
		return 0, nil, fmt.Errorf("error loading test file: %v", err)
	}
	size := int(stat.Size())

	f, err := os.Open(path)
	if err != nil {
		return size, f, fmt.Errorf("error loading test file: %v", err)
	}
	return size, f, nil
}
