package projections

type Scanner[T int | float64] interface {
	Points(yield func(y T, x T) bool)
}

type ScannerParams[T int | float64] struct {
	ScanningMode ScanningMode
	Ni           int // total number of points in the x direction
	Nj           int // total number of points in the y direction
	Di           T   // x increment (distance between points)
	Dj           T   // y increment (distance between points)
	I0           T   // starting x point
	J0           T   // starting y point
}

type scanner[T int | float64] struct {
	nDim1     int // total number of points in the first dimension
	nDim2     int // total number of points in the second dimension
	dDim1     T   // dimension 1 increment (distance between points)
	dDim2     T   // dimension 2 increment (distance between points)
	dim1Start T   // starting point of dimension 1
	dim2Start T   // starting point of dimension 2
	iDim      int // indicates which dimension is the X direction
}

func NewScanner[T int | float64](params ScannerParams[T]) Scanner[T] {
	g := &scanner[T]{}
	if params.ScanningMode.RightToLeft && params.Di > 0 {
		params.Di = -params.Di
	}
	if params.ScanningMode.TopToBottom && params.Dj > 0 {
		params.Dj = -params.Dj
	}
	if params.ScanningMode.OverFirst {
		g.nDim1 = params.Nj
		g.nDim2 = params.Ni
		g.dDim1 = params.Dj
		g.dDim2 = params.Di
		g.dim1Start = params.J0
		g.dim2Start = params.I0
		g.iDim = 2
	} else {
		g.nDim1 = params.Ni
		g.nDim2 = params.Nj
		g.dDim1 = params.Di
		g.dDim2 = params.Dj
		g.dim1Start = params.I0
		g.dim2Start = params.J0
		g.iDim = 1
	}

	return g
}

func (s *scanner[T]) Points(yield func(y T, x T) bool) {
	for dim1 := 0; dim1 < s.nDim1; dim1++ {
		dim1Val := s.dim1Start + T(dim1)*s.dDim1
		for dim2 := 0; dim2 < s.nDim2; dim2++ {
			dim2Val := s.dim2Start + T(dim2)*s.dDim2
			if s.iDim == 1 {
				if !yield(dim2Val, dim1Val) {
					return
				}
			} else {
				if !yield(dim1Val, dim2Val) {
					return
				}
			}
		}
	}
}

type ScanningMode struct {
	RightToLeft bool
	TopToBottom bool
	OverFirst   bool
}

func ScanningModeFromByte(b byte) ScanningMode {
	return ScanningMode{
		RightToLeft: (b>>7)&1 == 1,
		TopToBottom: (b>>6)&1 == 0,
		OverFirst:   (b>>5)&1 == 0,
	}
}
