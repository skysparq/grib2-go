// grib2-go is a pure Go library for reading GRIB2 files. It has been designed to avoid external dependencies and
// while offering high performance with low memory usage.
//
// GRIB2 files are composed of one or more independent GRIB2 records. Each record is composed of eight sections defining
// the metadata, grid, and data values. The specification for GRIB2 records is complex by nature of being very flexible.
// Three metadata sections (grid, product, and data representation) of the grib record can each be one of many possible
// structures. Thus, a proper implementation of a GRIB2 parser must handle large variety in the structure of the record.
//
// This package is intended to provide a reader that can handle all the standardized structures of a GRIB2 record.
// Currently, it is able to handle common grib2 records from NOAA's GFS, GDAS, HRRR, and MRMS products.
//
// The starting point for grib2-go is file.NewGribFile. This function requires an io.Reader and a record.Templates.
// The io.Reader should be a file-based or in-memory reader of grib2 data. The record.Templates is a definition of
// templates to use for section and template parsing. This package includes a standard set of templates that can be
// accessed by calling templates.Version33.

package main
