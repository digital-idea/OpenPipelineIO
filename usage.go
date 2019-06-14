package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "\nCSI v3.0 Project Management Tool.\n")
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "도움말 자세히보기:\n")
	fmt.Fprintf(os.Stderr, "$ csi3 -help\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "웹서버 실행:\n")
	fmt.Fprintf(os.Stderr, "$ csi3 -http :80\n")
	fmt.Fprintf(os.Stderr, "$ csi3 -http :8080 // 8080 포트로 서버 실행\n")
	fmt.Fprintf(os.Stderr, "$ csi3 -http :8080 -debug // 8080 포트로 디버그 모드를 활성화하여 서버 실행\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "프로젝트 추가:\n")
	fmt.Fprintf(os.Stderr, "$ csi3 -add project -name [projectName]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "프로젝트 삭제:\n")
	fmt.Fprintf(os.Stderr, "# csi3 -rm project -name [projectName]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "샷 추가:\n")
	fmt.Fprintf(os.Stderr, "$ csi3 -add item -project [projectName] -name [SS_0010] -type [org|left]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "에셋 추가:\n")
	fmt.Fprintf(os.Stderr, "$ csi3 -add item -project [projectName] -name [assetName] -type asset\n")
	fmt.Fprintf(os.Stderr, "  -assettype [char|env|global|prop|comp|plant|vehicle|group] -assettags [component|assembly]\n")
	fmt.Fprintf(os.Stderr, "\n")
	flag.PrintDefaults()
	os.Exit(2)
}
