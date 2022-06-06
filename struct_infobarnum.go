package main

import (
	"math"
)

// Infobarnum 검색결과에 대한 상태별 갯수를 담기위한 자료구조이다.
type Infobarnum struct {
	Assign    int            `json:"assign"`  // legacy
	Ready     int            `json:"ready"`   // legacy
	Wip       int            `json:"wip"`     // legacy
	Confirm   int            `json:"confirm"` // legacy
	Done      int            `json:"done"`    // legacy
	Omit      int            `json:"omit"`    // legacy
	Hold      int            `json:"hold"`    // legacy
	Out       int            `json:"out"`     // legacy
	None      int            `json:"none"`    // legacy
	StatusNum map[string]int `json:"statusnum"`
	Total     int            `json:"total"`
	Search    int            `json:"search"`
	Shot      int            `json:"shot"`
	Shot2d    int            `json:"shot2d"`
	Shot3d    int            `json:"shot3d"`
	Assets    int            `json:"assets"`
	Percent   float64        `json:"percent"`
}

// Percent 메소드는 Infobarnum 자료구조를 분석해서 진행률을 계산한다.
func (i *Infobarnum) calculatePercent() {
	if i.Total == 0 {
		i.Percent = 0.0
		return
	}
	i.Percent = math.Round(float64(i.Done+i.Hold) / float64(i.Total-i.None-i.Omit) * 100)
}
