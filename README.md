## Overview

grib2-go is a pure Go library for parsing GRIB2 files. It is currently in development and is missing many features. The goal of this library is to provide a complete solution to read any file compatible with version 33 of the GRIB2 standard.

At the same time, the API is designed to allow users to provide their own implementations of Grid, Product Definition, and Data Representation templates. Thus, if the current version of the library is missing an implementation for a particular template, the user can provide their own while still utilizing the peformant readers provided by this library

Details about the GRIB2 file format can be found at [NCEP](https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/).

NOTE: This package is not yet version 1. The API is subject to change until the v1 release.

## How to use

Install in your project by running `go get github.com/skysparq/grib2-go@latest`

See `main.go` for an example of using this library to parse a GRIB2 file.

The primary entrypoint for parsing GRIB2 files is a `file.GribFile`. You can instantiate a `GribFile` using `file.NewGribFile(r io.Reader, templates record.Templates)`. The `record.Templates` interface tells the GRIB2 file about each section of a record, so it can parse the incoming byte stream from `r`. This API allows you to specify your own custom implementation if the standard templates are insufficient.

`GribFile` contains a `Records` method that allows the user to iterate through GRIB2 records in the file. Each `Record` contains the 8 sections of the GRIB2 file. For the sections which are driven by templates (Sections 3, 4, and 5), the section contains the raw bytes of the template data. To further parse the templates, there is a helper method on the section called `Definition`, which will parse its section's template using the templates passed to the `GribFile`.

In this way, the grib2-go package immediately provides a high-performance standard framework for parsing GRIB2 files, while allowing future expansion and user-specific implementations.

## Tests

The tests in this package use actual GRIB2 data from various sources (initially, GFS forecast files). Due to the size of GRIB2 files, it is not feasible to keep them in the repo. To run unit tests in this package, download the files from [this link](https://drive.google.com/file/d/1qXFrMPeNCaR7bXzndTsRgKWYMCUL6rgO/view?usp=sharing) and place them in a directory called `.test_files` in the root package folder.

## Current Status

The short-term focus of this library is to process GFS GRIB2 files. Currently, the library has implemented the following templates:

### Grid Definition Templates

- Template 3.0 (Longitude / Latitude)
- Template 3.30 (Lambert Conformal)
- Template 3.40 (Gaussian Latitude/Longitude)

### Product Definition Templates

- Template 4.0 (Analysis or forecast at a horizontal level or in a horizontal layer at a point in time.)
- Template 4.1 (Individual ensemble forecast, control and perturbed, at a horizontal level or in a horizontal layer at a point in time.)
- Template 4.2 (Derived forecasts based on all ensemble members at a horizontal level or in a horizontal layer at a point in time.)
- Template 4.3 (Derived forecasts based on a cluster of ensemble members over a rectangular area at a horizontal level or in a horizontal layer at a point in time.)
- Template 4.4 (Derived forecasts based on a cluster of ensemble members over a circular area at a horizontal level or in a horizontal layer at a point in time.)
- Template 4.5 (Probability forecasts at a horizontal level or in a horizontal layer at a point in time.)
- Template 4.6 (Percentile forecasts at a horizontal level or in a horizontal layer at a point in time.)
- Template 4.7 (Analysis or forecast error at a horizontal level or in a horizontal layer at a point in time.)
- Template 4.8 (Average, accumulation, extreme values or other statistically processed values at a horizontal level or in a horizontal layer in a continuous or non-continuous time interval.)

### Data Representation Templates

NOTE: Currently, all data representations listed below, except for JPEG 2000, are implemented. The implementations are still under testing. If you find they do not accurately unpack your GRIB2 data, open a PR and share the file so it can be fixed.

- Template 5.0 (Grid Point Data - Simple Packing)
- Template 5.2 (Grid Point Data - Complex Packing)
- Template 5.3 (Grid Point Data - Complex Packing and Spatial Differencing)
- Template 5.40 (Grid point data - JPEG 2000 code stream format)
- Template 5.41 (Grid point data - PNG)

## Roadmap

- [x] Read metadata of GFS GRIB2 records
- [x] Implement iterator to parse multi-record GRIB2 files
- [x] Read values and lat/lon points from GFS GRIB2 records
- [ ] Read metadata for all templates compatible with the version 33 standard
- [ ] Read values for all templates compatible with the version 33 standard
- [ ] Implement current versions of the GRIB2 standard
