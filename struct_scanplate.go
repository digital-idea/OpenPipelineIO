package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type ScanPlate struct {
	// 사용자 입력값
	ID                    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CreateTime            string             `json:"createtime" bson:"createtime"`                       // CreateTime
	ProcessStatus         string             `json:"processstatus" bson:"processstatus"`                 // 연산상태
	InColorspace          string             `json:"incolorspace" bson:"incolorspace"`                   // 시퀀스의 IN 컬러스페이스
	OutColorspace         string             `json:"outcolorspace" bson:"outcolorspace"`                 // 시퀀스의 OUT 컬러스페이스
	RenderIn              int                `json:"renderin" bson:"renderin"`                           // 렌더링 할 시작 IN 프레임
	RenderOut             int                `json:"renderout" bson:"renderout"`                         // 렌더링 할 끝 OUT 프레임
	Searchpath            string             `json:"searchpath" bson:"searchpath"`                       // 시퀀스 검색을 시작한 Endpoint
	ConvertExt            string             `json:"convertext" bson:"convertext"`                       // 만약 소스를 저장할 때 변환하여 저장한다면 사용할 확장자
	Type                  string             `json:"type" bson:"type"`                                   // 만약 소스를 저장할 때 변환하여 저장한다면 사용할 확장자
	Name                  string             `json:"name" bson:"name"`                                   // 등록이름
	Project               string             `json:"project" bson:"project"`                             // 프로젝트
	RenderWidth           int                `json:"renderwidth" bson:"renderwidth"`                     // 렌더링 가로길이
	RenderHeight          int                `json:"renderheight" bson:"renderheight"`                   // 렌더링 세로길이
	UndistortionWidth     int                `json:"undistortionwidth" bson:"undistortionwidth"`         // Undistortion 가로길이
	UndistortionHeight    int                `json:"undistortionheight" bson:"undistortionheight"`       // Undistortion 세로길이
	Dir                   string             `json:"dir" bson:"dir"`                                     // 시퀀스 디렉토리
	Base                  string             `json:"base" bson:"base"`                                   // 파일명(시퀀스 숫자 제외)
	Ext                   string             `json:"ext" bson:"ext"`                                     // 확장자
	Digitnum              int                `json:"digitnum" bson:"digitnum"`                           // 시퀀스 자릿수
	FrameIn               int                `json:"framein" bson:"framein"`                             // 시작프레임
	FrameOut              int                `json:"frameout" bson:"frameout"`                           // 끝프레임
	Width                 int                `json:"width" bson:"width"`                                 // 가로길이
	Height                int                `json:"height" bson:"height"`                               // 세로길이
	TimecodeIn            string             `json:"timecodein" bson:"timecodein"`                       // 시작 타임코드
	TimecodeOut           string             `json:"timecodeout" bson:"timecodeout"`                     // 마지막 타임코드
	Length                int                `json:"length" bson:"length"`                               // 소스 길이
	InputCodec            string             `json:"inputcodec" bson:"inputcodec"`                       // 소스 코덱
	OutputCodec           string             `json:"outputcodec" bson:"outputcodec"`                     // 설정하는 아웃풋 코덱. 웹이라서 일반적으로 H.264를 사용한다.
	Fps                   string             `json:"fps" bson:"fps"`                                     // FPS
	Rollmedia             string             `json:"rollmedia" bson:"rollmedia"`                         // 촬영한 데이터라면, 카메라에서 생성된 데이터 이름
	Error                 string             `json:"error" bson:"error"`                                 // 에러기록
	GenPlatePath          bool               `json:"genplatepath" bson:"genplatepath"`                   // 플레이트 경로 생성
	GenMov                bool               `json:"genmov" bson:"genmov"`                               // mov 생성
	GenMovSlate           bool               `json:"genmovslate" bson:"genmovslate"`                     // mov 슬레이트 생성
	CopyPlate             bool               `json:"copyplate" bson:"copyplate"`                         // 플레이트 복사여부
	ProxyJpg              bool               `json:"proxyjpg" bson:"proxyjpg"`                           // Proxy .jpg 생성여부
	ProxyHalfJpg          bool               `json:"proxyhalfjpg" bson:"proxyhalfjpg"`                   // Proxy half .jpg 생성여부
	ProxyHalfExr          bool               `json:"proxyhalfexr" bson:"proxyhalfexr"`                   // Proxy half .exr 생성여부
	SetFrame              bool               `json:"setframe" bson:"setframe"`                           // 프레임 설정여부
	SetTimecode           bool               `json:"settimecode" bson:"settimecode"`                     // Timecode 설정여부
	LutPath               string             `json:"lutpath" bson:"lutpath"`                             // Lut Path
	ProxyRatio            int                `json:"proxyratio" bson:"proxyratio"`                       // Proxy Ratio
	DNS                   string             `json:"dns" bson:"dns"`                                     // OpenPipelineIO DNS
	Token                 string             `json:"token" bson:"token"`                                 // OpenPipelineIO Token
	UseOriginalNameForMov bool               `json:"useoriginalnameformov" bson:"useoriginalnameformov"` // mov생성시 오리지널 이름을 사용할 것 인지 체크
}
