package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pipelinestep struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name        string             `json:"name"`        // PipelineStep 이름, 보통은 팀이름이다.
	Description string             `json:"description"` // 상세설명
}
