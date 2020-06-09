package main

import (
	"os/exec"
)

// genMp4 는 리뷰 아이템 정보를 이용해서 .mp4 동영상을 만든다.
func genMp4(admin Setting, item Review) error {
	args := []string{
		"-y",
		"-i",
		item.Path,
		"-c:v",
		"libx264",
		"-qscale:v",
		"7",
		"-an",
		admin.ReviewDataPath + "/" + item.ID.Hex() + ".mp4",
	}
	err := exec.Command(admin.FFmpeg, args...).Run()
	if err != nil {
		return err
	}
	return nil
}

// genOgg 는 리뷰 아이템 정보를 이용해서 .ogg 동영상을 만든다.
func genOgg(admin Setting, item Review) error {
	args := []string{
		"-y",
		"-i",
		item.Path,
		"-c:v",
		"libtheora",
		"-qscale:v",
		"7",
		"-an",
		admin.ReviewDataPath + "/" + item.ID.Hex() + ".ogg",
	}
	err := exec.Command(admin.FFmpeg, args...).Run()
	if err != nil {
		return err
	}
	return nil
}
