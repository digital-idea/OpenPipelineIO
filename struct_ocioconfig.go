package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Colorspace 자료구조는 OCIO Color 자료구조이다.
type Colorspace struct {
	Name     string `yaml:"name"`
	Family   string `yaml:"family"`
	Bitdepth string `yaml:"bitdepth"`
}

// Displays 자료구조는 OCIO Display 자료구조이다.
type Displays struct {
	ACES []View `yaml:"ACES"`
}

// View 자료구조는 OCIO Display > ACES: 자료구조이다.
type View struct {
	Name       string `yaml:"name"`
	Colorspace string `yaml:"colorspace"`
}

// OCIOConfig 자료구조는 config.ocio 파일 자료구조이다.
type OCIOConfig struct {
	OCIOProfileVersion string       `yaml:"ocio_profile_version"`
	Colorspaces        []Colorspace `yaml:"colorspaces"`
	Displays           `yaml:"displays"`
	Roles              map[string]string `yaml:"roles"`
}

func loadOCIOConfig() ([]string, error) {
	var results []string
	if CachedAdminSetting.OCIOConfig == "" {
		return results, nil
	}
	// 파일이 존재하는지 체크한다.
	_, err := os.Stat(CachedAdminSetting.OCIOConfig)
	if err != nil {
		return nil, err
	}
	// OCIO.config 파일을 불러온다.
	dat, err := ioutil.ReadFile(CachedAdminSetting.OCIOConfig)
	if err != nil {
		return nil, err
	}
	// 존재하면, 해당파일을 파싱한다.
	var oc OCIOConfig
	err = yaml.Unmarshal(dat, &oc)
	if err != nil {
		return nil, err
	}
	// color 리스트만 results에 넣는다.
	for _, c := range oc.Colorspaces {
		results = append(results, c.Name)
	}
	return results, nil
}
