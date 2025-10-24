package projections_test

import (
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/projections"
)

type point struct {
	x int
	y int
}

type scannerTest struct {
	Name         string
	ScanningMode projections.ScanningMode
	Expected     []point
}

func TestScannerRightThenDown(t *testing.T) {
	for _, test := range tests {
		s := projections.NewScanner(projections.ScannerParams[int]{
			ScanningMode: test.ScanningMode,
			Ni:           2,
			Nj:           2,
			Di:           1,
			Dj:           1,
			I0:           0,
			J0:           0,
		})
		actual := make([]point, 0, 4)
		for y, x := range s.Points {
			actual = append(actual, point{x, y})
		}
		if !reflect.DeepEqual(test.Expected, actual) {
			t.Fatalf("test '%v' failed: expected\n%+v\nbut got\n%+v", test.Name, test.Expected, actual)
		}
	}
}

var tests = []scannerTest{
	{
		Name:         "Right Then Down",
		ScanningMode: projections.ScanningMode{RightToLeft: false, TopToBottom: true, OverFirst: true},
		Expected: []point{
			{0, 0},
			{1, 0},
			{0, -1},
			{1, -1},
		},
	},
	{
		Name:         "Right Then Up",
		ScanningMode: projections.ScanningMode{RightToLeft: false, TopToBottom: false, OverFirst: true},
		Expected: []point{
			{0, 0},
			{1, 0},
			{0, 1},
			{1, 1},
		},
	},
	{
		Name:         "Left Then Down",
		ScanningMode: projections.ScanningMode{RightToLeft: true, TopToBottom: true, OverFirst: true},
		Expected: []point{
			{0, 0},
			{-1, 0},
			{0, -1},
			{-1, -1},
		},
	},
	{
		Name:         "Left Then Up",
		ScanningMode: projections.ScanningMode{RightToLeft: true, TopToBottom: false, OverFirst: true},
		Expected: []point{
			{0, 0},
			{-1, 0},
			{0, 1},
			{-1, 1},
		},
	},
	{
		Name:         "Down Then Right",
		ScanningMode: projections.ScanningMode{RightToLeft: false, TopToBottom: true, OverFirst: false},
		Expected: []point{
			{0, 0},
			{0, -1},
			{1, 0},
			{1, -1},
		},
	},
	{
		Name:         "Up Then Right",
		ScanningMode: projections.ScanningMode{RightToLeft: false, TopToBottom: false, OverFirst: false},
		Expected: []point{
			{0, 0},
			{0, 1},
			{1, 0},
			{1, 1},
		},
	},
	{
		Name:         "Down Then Left",
		ScanningMode: projections.ScanningMode{RightToLeft: true, TopToBottom: true, OverFirst: false},
		Expected: []point{
			{0, 0},
			{0, -1},
			{-1, 0},
			{-1, -1},
		},
	},
	{
		Name:         "Up Then Left",
		ScanningMode: projections.ScanningMode{RightToLeft: true, TopToBottom: false, OverFirst: false},
		Expected: []point{
			{0, 0},
			{0, 1},
			{-1, 0},
			{-1, 1},
		},
	},
}
