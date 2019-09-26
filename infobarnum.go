package main

import (
	"math"
)

// Infobarnum 검색결과에 대한 상태별 갯수를 담기위한 자료구조이다.
type Infobarnum struct {
	Assign  int
	Ready   int
	Wip     int
	Confirm int
	Done    int
	Omit    int
	Hold    int
	Out     int
	None    int
	Total   int
	Search  int
	Shot    int
	Shot2d  int
	Shot3d  int
	Assets  int
	Percent float64
}

// Percent 메소드는 Infobarnum 자료구조를 분석해서 진행률을 계산한다.
func (i *Infobarnum) calculatePercent() {
	if i.Total == 0 {
		i.Percent = 0.0
		return
	}
	i.Percent = math.Round(float64(i.Done+i.Hold) / float64(i.Total-i.None-i.Omit) * 100)
}
