package main

import (
	"math"
	"testing"
	"time"
)

func TestCalculateSpeed(t *testing.T) {
	p1 := Path{1, 37.966660, 23.728308, 1405594957}
	p2 := Path{1, 37.966627, 23.728263, 1405594966}
	s := CalculateFareInfo(p1, p2).speed
	if s != 2.1550435800609105 {
		t.Error("Expected 2.1550435800609105, got ", s)
	}
	p1 = Path{3, 37.994643, 23.872791, 1405592339}
	p2 = Path{3, 37.994889, 23.872559, 1405592339}
	s = CalculateFareInfo(p1, p2).speed
	if !math.IsInf(s, 1) {
		t.Error("Expected +Inf, got ", s)
	}
}

func TestFilterPaths(t *testing.T) {
	paths := map[int]Path{
		1: Path{1, 37.955217, 23.714548, 1405595237},
		2: Path{1, 37.954302, 23.713370, 1405595284},
		3: Path{1, 37.938042, 23.692308, 1405595362}, //invalid
		4: Path{1, 37.938985, 23.690435, 1405595371}, //invalid
		5: Path{1, 37.940058, 23.688853, 1405595379}, //invalid
		6: Path{1, 37.940872, 23.687423, 1405595387},
		7: Path{1, 37.941705, 23.685902, 1405595396},
	}
	filtered := FilterPaths(&paths)
	if len(filtered) != 3 {
		t.Error("Expected 3 to be filtered, got ", len(filtered))
	}
}

func TestEstimateFare(t *testing.T) {
	infos := []FareInfo{
		{0.08769823737089351, 26.30947121126805, time.Unix(63541188618, 0), 12},
		{0.08387324892913653, 37.74296201811144, time.Unix(63541188626, 0), 8},
		{0.06619058325347046, 19.857174976041136, time.Unix(63541188638, 0), 12},
		{0.0, 9.0, time.Unix(63541188646, 0), 8},
	}
	if p := EstimateFare(infos, 1); p != 347 {
		t.Error("Expected 5.00, got ", p)
	}
}

func TestApplyPrice(t *testing.T) {
	if p := ApplyPrice(5.0, 0.0, 0.0, 1); p != 500 {
		t.Error("Expected 5.00, got ", p)
	}
	if p := ApplyPrice(0.0, 9.0, 0.0, 2); p != 1300 {
		t.Error("Expected 13.00, got ", p)
	}
	if p := ApplyPrice(0.0, 0.0, 3*3600, 3); p != 3700 {
		t.Error("Expected 37.00, got ", p)
	}
	if p := ApplyPrice(1.0, 1.0, 0, 4); p != 347 {
		t.Error("Expected 15.24, got ", p)
	}
	if p := ApplyPrice(1.0, 1.0, 3600, 5); p != 1524 {
		t.Error("Expected 15.24, got ", p)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
