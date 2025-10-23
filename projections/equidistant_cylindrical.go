package projections

type EquidistantCylindricalParams struct {
	RightToLeft bool
	TopToBottom bool
	OverFirst   bool
	Ni          int // total number of points in the x direction
	Nj          int // total number of points in the y direction
	Di          int // x increment (distance between points)
	Dj          int // y increment (distance between points)
	I0          int // starting x point
	J0          int // starting y point
}

type equidistantCylindrical struct {
	nDim1     int // total number of points in the first dimension
	nDim2     int // total number of points in the second dimension
	dDim1     int // dimension 1 increment (distance between points)
	dDim2     int // dimension 2 increment (distance between points)
	dim1Start int // starting point of dimension 1
	dim2Start int // starting point of dimension 2
	iDim      int // indicates which dimension is the X direction
}

func ExtractEquidistantCylindricalGrid(params EquidistantCylindricalParams) (lats []float32, lngs []float32) {
	g := equidistantCylindrical{}
	if params.RightToLeft && params.Di > 0 {
		params.Di = -params.Di
	}
	if params.TopToBottom && params.Dj > 0 {
		params.Dj = -params.Dj
	}
	if params.OverFirst {
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
	lats = make([]float32, 0, g.nDim1*g.nDim2)
	lngs = make([]float32, 0, g.nDim1*g.nDim2)

	for dim1 := 0; dim1 < g.nDim1; dim1++ {
		dim1Val := g.dim1Start + dim1*g.dDim1
		for dim2 := 0; dim2 < g.nDim2; dim2++ {
			dim2Val := g.dim2Start + dim2*g.dDim2
			if g.iDim == 1 {
				lngs = append(lngs, float32(normalizeLng(dim1Val))*1e-6)
				lats = append(lats, float32(dim2Val)*1e-6)
			} else {
				lngs = append(lngs, float32(normalizeLng(dim2Val))*1e-6)
				lats = append(lats, float32(dim1Val)*1e-6)
			}
		}
	}
	return lats, lngs
}
