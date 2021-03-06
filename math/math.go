package cvmath

import (
	"image"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/digisan/go-generics/f64"
	"github.com/digisan/go-generics/f64i64"
	"github.com/digisan/go-generics/f64v4i64"
	"github.com/digisan/go-generics/pti64"
)

func Max(data ...float64) float64 {
	return f64.Max(data...)
}

func Min(data ...float64) float64 {
	return f64.Min(data...)
}

func MaxIdx(data ...float64) (float64, int) {
	return f64.MaxIdx(data...)
}

func MinIdx(data ...float64) (float64, int) {
	return f64.MinIdx(data...)
}

func MaxMinAbs(data ...float64) (float64, float64) {
	tempAbs := f64.FM(data, nil, func(i int, e float64) float64 {
		if e < 0 {
			return -e
		}
		return e
	})
	return Max(tempAbs...), Min(tempAbs...)
}

func StepSelect(step, offset int, data ...float64) []float64 {
	return f64.FM(data, func(i int, e float64) bool { return i%step == offset }, nil)
}

func Sum(data ...float64) (sum float64) {
	for _, v := range data {
		sum += v
	}
	return
}

func SumStep4(data ...float64) (sum [4]float64) {
	for i := 0; i < len(data); i += 4 {
		sum[0] += data[i]
		sum[1] += data[i+1]
		sum[2] += data[i+2]
		sum[3] += data[i+3]
	}
	return
}

func Mean(data ...float64) float64 {
	return Sum(data...) / float64(len(data))
}

func MeanStep4(data ...float64) (mean [4]float64) {
	sums := SumStep4(data...)
	nVec := len(data) / 4
	for i := 0; i < 4; i++ {
		mean[i] = sums[i] / float64(nVec)
	}
	return
}

func Median(data ...float64) float64 {
	sort.Float64s(data)
	pos := len(data) / 2
	if len(data)%2 == 0 {
		return data[pos-1]
	}
	return data[pos]
}

func MedianStep4(data ...float64) [4]float64 {
	d0 := StepSelect(4, 0, data...)
	d1 := StepSelect(4, 1, data...)
	d2 := StepSelect(4, 2, data...)
	d3 := StepSelect(4, 3, data...)
	return [4]float64{Median(d0...), Median(d1...), Median(d2...), Median(d3...)}
}

func Mode(data ...float64) float64 {
	m := make(map[float64]int)
	for _, d := range data {
		m[d]++
	}
	ks, _ := f64i64.Map2KVs(m, nil, func(i int, j int) bool {
		return i > j
	})
	return ks[0]
}

func ModeVec4(data ...float64) [4]float64 {
	m := make(map[[4]float64]int)
	for i := 0; i < len(data); i += 4 {
		k := [4]float64{data[i], data[i+1], data[i+2], data[i+3]}
		m[k]++
	}
	ks, _ := f64v4i64.Map2KVs(m, nil, func(i, j int) bool {
		return i > j
	})
	return ks[0]
}

func ModeStep4(data ...float64) [4]float64 {
	d0 := StepSelect(4, 0, data...)
	d1 := StepSelect(4, 1, data...)
	d2 := StepSelect(4, 2, data...)
	d3 := StepSelect(4, 3, data...)
	return [4]float64{Mode(d0...), Mode(d1...), Mode(d2...), Mode(d3...)}
}

func StdDev(data ...float64) float64 {
	sum2 := 0.0
	ave := Mean(data...)
	for _, v := range data {
		d2 := (v - ave) * (v - ave)
		sum2 += d2
	}
	return math.Sqrt(sum2 / float64(len(data)))
}

func DotProduct(v1, v2 []float64) (dp float64) {
	if len(v1) != len(v2) {
		log.Fatalf("DotProduct vector dimensions error")
	}
	for i := 0; i < len(v1); i++ {
		dp += v1[i] * v2[i]
	}
	return
}

func Derivative(data ...float64) (ret []float64) {
	dp := []float64{86, -142, -193, -126, 126, 193, 142, -86}
	ret = make([]float64, len(data))
	for i := 4; i < len(data)-4; i++ {
		// ret[i] = (86*temp[i-4] - 142*temp[i-3] - 193*temp[i-2] - 126*temp[i-1] + 126*temp[i+1] + 193*temp[i+2] + 142*temp[i+3] - 86*temp[i+4]) / 1188
		temp := data[i-4 : i+5]
		ret[i] = DotProduct(dp, temp) / 1188
	}
	return
}

func Smooth9(pts []float64) (ret []float64) {
	sp9 := []float64{-21, 14, 39, 54, 59, 54, 39, 14, -21}
	ret = make([]float64, len(pts))
	copy(ret, pts[:4])
	copy(ret[len(ret)-4:], pts[len(pts)-4:])
	for i := 4; i < len(pts)-4; i++ {
		// ret[i] = ((-21)*pts[i-4] + 14*pts[i-3] + 39*pts[i-2] + 54*pts[i-1] + 59*pts[i] + 54*pts[i+1] + 39*pts[i+2] + 14*pts[i+3] - 21*pts[i+4]) / 231
		pts9 := pts[i-4 : i+5]
		ret[i] = DotProduct(sp9, pts9) / 231
	}
	for i := 0; i < len(ret); i++ {
		if ret[i] < 0 {
			ret[i] = 0
		}
	}
	return
}

func Dis(a, b float64) float64 {
	return math.Abs(a - b)
}

func DisInt(a, b int) int {
	fa, fb := float64(a), float64(b)
	return int(Dis(fa, fb))
}

func DisByte(a, b byte) byte {
	fa, fb := float64(a), float64(b)
	return byte(Dis(fa, fb))
}

func DisPt(pt1, pt2 image.Point) float64 {
	dx := float64(DisInt(pt1.X, pt2.X))
	dy := float64(DisInt(pt1.Y, pt2.Y))
	return math.Sqrt(dx*dx + dy*dy)
}

func DisPtX(pt1, pt2 image.Point) float64 {
	return float64(DisInt(pt1.X, pt2.X))
}

func DisPtY(pt1, pt2 image.Point) float64 {
	return float64(DisInt(pt1.Y, pt2.Y))
}

func NearFarPoint(toPt image.Point, pts ...image.Point) (nearest, farest image.Point) {
	disGrp := []float64{}
	for _, pt := range pts {
		disGrp = append(disGrp, DisPt(pt, toPt))
	}
	dMin := Min(disGrp...)
	dMax := Max(disGrp...)
	for _, pt := range pts {
		dis := DisPt(pt, toPt)
		if dis == dMin {
			nearest = pt
		}
		if dis == dMax {
			farest = pt
		}
	}
	return
}

func PointsRect(pts ...image.Point) (left, top, right, bottom float64) {
	xs := I64ToF64(pti64.FM(pts, nil, func(i int, e image.Point) int { return e.X })...)
	ys := I64ToF64(pti64.FM(pts, nil, func(i int, e image.Point) int { return e.Y })...)
	xMin, xMax := Min(xs...), Max(xs...)
	yMin, yMax := Min(ys...), Max(ys...)
	return xMin, yMin, xMax, yMax
}

func RandomPoint(pts ...image.Point) image.Point {

	xMin, yMin, xMax, yMax := PointsRect(pts...)

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed).Float64()

	x := r*(xMax-xMin+0.5) + xMin
	y := r*(yMax-yMin+0.5) + yMin
	rPt := image.Pt(int(x), int(y))

	near, _ := NearFarPoint(rPt, pts...)
	return near
}

func CentrePoint(pts ...image.Point) image.Point {

	xMin, yMin, xMax, yMax := PointsRect(pts...)

	x := (xMax + xMin) / 2
	y := (yMax + yMin) / 2
	cPt := image.Pt(int(x), int(y))

	near, _ := NearFarPoint(cPt, pts...)
	return near
}

func SplitPtsXY(pts []image.Point, xy string, n int) [][]image.Point {

	xMin, yMin, xMax, yMax := PointsRect(pts...)
	areas := make([][]image.Point, n)

	switch xy {
	case "X", "x":
		span := (xMax - xMin) / float64(n)
		for _, pt := range pts {
			i := int(math.Round((float64(pt.X) - xMin) / span))
			areas[i] = append(areas[i], pt)
		}
	case "Y", "y":
		span := (yMax - yMin) / float64(n)
		for _, pt := range pts {
			i := int(math.Round((float64(pt.Y) - yMin) / span))
			areas[i] = append(areas[i], pt)
		}
	default:
		log.Fatalln("[xy] can only set as 'X' or 'Y'")
	}

	return areas
}
