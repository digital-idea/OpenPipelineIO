package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pipelinestep struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name        string             `json:"name"`        // PipelineStep 이름, 보통은 팀이름이다.
	Shortcode   string             `json:"shortcode"`   // 짧은 이름
	Description string             `json:"description"` // 상세설명
	TextColor   string             `json:"textcolor"`   // TEXT 색상
	BGColor     string             `json:"bgcolor"`     // BG 상태 색상
	BorderColor string             `json:"bordercolor"` // Border 색상
}
