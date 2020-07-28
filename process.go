package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

// ProcessMain 함수는 CSI가 실행되면서 처리될 프로세싱을 진행한다.
func ProcessMain() {
	// 버퍼 채널을 만든다.
	jobs := make(chan Review, *flagProcessBufferSize) // 작업을 대기할 버퍼를 만든다.
	// worker 프로세스를 지정한 개수만큼 실행시킨다.
	for w := 1; w <= *flagMaxProcessNum; w++ {
		go worker(jobs)
	}
	// queueingItem을 실행시킨다.
	go queueingItem(jobs)

	// ProcessMain()이 종료되지 않는 역할을 한다.
	select {}
}

// worker 합수는 Review 데이터를 jobs로 보낸다.
func worker(jobs <-chan Review) {
	for j := range jobs {
		processingItem(j)
	}
}

// queueingItem 은 연산할 Review 아이템을 jobs 채널에 전송한다.
func queueingItem(jobs chan<- Review) {
	for {
		if *flagDebug {
			fmt.Println("10 sec")
		}
		time.Sleep(time.Second * 10)
		// ProcessStatus가 wait인 item을 가져온다.
		review, err := GetWaitProcessStatusReview()
		if err != nil {
			// 가지고 올 문서가 없다면 기다렸다가 continue.
			if err == mongo.ErrNoDocuments {
				if *flagDebug {
					log.Println(err)
				}
				continue
			}
			continue
		}
		if *flagDebug {
			fmt.Println(review)
		}
		jobs <- review
	}
}

func processingItem(review Review) {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		return
	}
	defer session.Close()
	// 연산에 필요한 설정을 가지고 온다.
	admin, err := GetAdminSetting(session)
	if err != nil {
		log.Println(err)
		return
	}
	// 연산한다.
	err = genMp4(admin, review)
	if err != nil {
		log.Println(err)
		return
	}
	// 연산 상태를 done 으로 바꾼다.
	reviewID := review.ID.Hex()
	err = setReviewProcessStatus(session, reviewID, "done")
	if err != nil {
		log.Println(err)
		return
	}
	return
}

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
		"-pix_fmt",
		"yuv420p", // 이 옵션이 없다면 Prores로 동영상을 만들때 크롬에서만 재생된다.
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
