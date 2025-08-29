## Overview

grib2-go is a pure Go library for parsing GRIB2 files. It is currently in development and is missing many features. The goal of this library is to provide a complete solution to read any file compatible with version 33 of the GRIB2 standard.

At the same time, the API is designed to allow users to provide their own implementations of Grid, Product Definition, and Data Representation templates. Thus, if the current version of the library is missing an implementation for a particular template, the user can provide their own while still utilizing the peformant readers provided by this library

Details about the GRIB2 file format can be found at [NCEP](https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/).

## How to use

Install in your project by running `go get github.com/skysparq/grib2-go@latest`

See `main.go` for an example of using this library to parse a GRIB2 file.

The primary entrypoint for parsing GRIB2 files is a `file.GribFile`. You can instantiate a `GribFile` using `file.NewGribFile(r io.Reader, templates interfaces.Template)`. The `interfaces.Template` interface tells the GRIB2 file how long each Grid Definition and Product Definition template is, so it can parse the incoming byte stream from `r`. This API allows you to specify your own custom implementation if the standard templates are not sufficient. Below is an example to instantiate a new `GribFile` from a file on disk using the standard template:

`GribFile` contains a `Records` method that allows the user to iterate through GRIB2 records in the file. Each `Record` contains the 8 sections of the GRIB2 file. For the sections which are driven by templates (Sections 3, 4, and 5), the section contains the raw bytes of the template data. To further parse the templates, there is a helper struct called `Parser` in the `grid`, `product`, and `data_representation` packages. Each `Parser` will parse its section's template using the standard templates included in this package. However, if the user requires a template not included in the package, they can provide their own implementations to the `Templates` member of the `Parser`.

In this way, the grib2-go package immediately provides a high-performance standard framework for parsing GRIB2 files, while allowing future expansion and user-specific implementations.

## Tests

The tests in this package use actual GRIB2 data from various sources (initially, GFS forecast files). Due to the size of GRIB2 files, it is not feasible to keep them in the repo. To run unit tests in this package, download the files from [this link](https://drive.google.com/file/d/1qXFrMPeNCaR7bXzndTsRgKWYMCUL6rgO/view?usp=sharing) and place them in a director called `.test_files` in the root package folder.

## Current Status

The short-term focus of this library is to process GFS GRIB2 files. Currently, the library has implemented the following templates:

### Grid Definition Templates

- Template 3.0 (Longitude / Latitude)

### Product Definition Templates

- Template 4.0 (Analysis or forecast at a horizontal level or in a horizontal layer at a point in time.)
- Template 4.8 (Average, accumulation, extreme values or other statistically processed values at a horizontal level or in a horizontal layer in a continuous or non-continuous time interval.)

### Data Representation Templates

NOTE: Currently, none of the implemented data representation templates can decode the data section of a GRIB2 record. Implementing decoding logic is a high-priority short-term goal.

- Template 5.0 (Grid Point Data - Simple Packing)
- Template 5.3 (Grid Point Data - Complex Packing and Spatial Differencing)

## Roadmap

- [x] Read metadata of GFS GRIB2 records
- [x] Implement iterator to parse multi-record GRIB2 files
- [ ] Read values from GFS GRIB2 records
- [ ] Read metadata for all templates compatible with the version 33 standard
- [ ] Read values for all templates compatible with the version 33 standard
- [ ] Implement current versions of the GRIB2 standard
