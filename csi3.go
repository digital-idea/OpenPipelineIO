package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/digital-idea/dilog"
	"gopkg.in/mgo.v2"
)

var (
	// DBIP 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	DBIP = "127.0.0.1"
	// DBPORT mongoDB 기본포트.
	DBPORT = ":27017"
	// WFS 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	WFS = "http://127.0.0.1:8081"
	// DILOG 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	DILOG = "http://127.0.0.1:8080"
	// THUMBPATH 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	THUMBPATH = "thumbnail"
	// DNS 값은 서비스 DNS 값입니다.
	DNS = "csi.lazypic.org"
	// MAILDNS 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	MAILDNS = "lazypic.org"
	// COMPANY 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	COMPANY = "lazypic"
	// TEMPLATES 값은 웹서버 실행전 사용할 템플릿이다.
	TEMPLATES = template.New("")
	// SHA1VER  은 Git SHA1 값이다.
	SHA1VER = "26b300a004abae553650c924514dc550e7385c9e" // 첫번째 커밋
	// BUILDTIME 은 빌드타임 시간이다.
	BUILDTIME = "2012-11-08T10:00:00" // 최초로 만든 시간

	// 주요서비스 인수
	flagDBIP           = flag.String("dbip", DBIP+DBPORT, "mongodb ip and port")
	flagMailDNS        = flag.String("maildns", MAILDNS, "mail DNS name")
	flagThumbPath      = flag.String("thumbpath", THUMBPATH, "thumbnail path")
	flagDebug          = flag.Bool("debug", false, "디버그모드 활성화")
	flagDevmode        = flag.Bool("devmode", false, "dev mode")
	flagHTTPPort       = flag.String("http", "", "Web Service Port number.")          // 웹서버 포트
	flagCompany        = flag.String("company", COMPANY, "Web Service Port number.")  // 회사이름
	flagVersion        = flag.Bool("version", false, "Print Version")                 // 버전
	flagCookieAge      = flag.Int64("cookieage", 168, "cookie age (hour)")            // 기본 일주일(168시간)로 설정한다. 참고: MPAA 기준 4시간이다.
	flagThumbnailAge   = flag.Int("thumbnailage", 1, "thumbnail image age (seconds)") // 썸네일 업데이트 시간. 3600초 == 1시간
	flagAuthmode       = flag.Bool("authmode", false, "restAPI authorization active") // restAPI 이용시 authorization 활성화
	flagCertFullchanin = flag.String("certfullchanin", fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", DNS), "certification fullchain path")
	flagCertPrivkey    = flag.String("certprivkey", fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", DNS), "certification privkey path")

	// RV
	flagRV   = flag.String("rvpath", "/opt/rv-Linux-x86-64-7.0.0/bin/rv", "rvplayer path")
	flagPlay = flag.Bool("play", false, "Play RV")
	// Etc Service
	flagDILOG = flag.String("dilog", DILOG, "dilog webserver url and port. ex) "+DILOG)
	flagWFS   = flag.String("wfs", WFS, "wfs webserver url and port. ex) "+WFS)

	// Commandline Args
	flagAdd              = flag.String("add", "", "add project, add item(shot, asset)")
	flagRm               = flag.String("rm", "", "remove project, shot, asset, user")
	flagProject          = flag.String("project", "", "project name")
	flagName             = flag.String("name", "", "name")
	flagType             = flag.String("type", "", "type: org,left,asset,org1,src,src1,lsrc,rsrc")
	flagAssettags        = flag.String("assettags", "", "asset tags, 입력예) prop,char,env,prop,comp,plant,vehicle,component,assembly 형태로 입력")
	flagAssettype        = flag.String("assettype", "", "assettype: char,env,global,prop,comp,plant,vehicle,group") // 추후 삭제예정.
	flagHelp             = flag.Bool("help", false, "자세한 도움말을 봅니다.")
	flagDate             = flag.String("date", "", "Date. ex) 2016-12-06")
	flagThumbnailPath    = flag.String("thumbnailpath", "", "Thumbnail 경로")
	flagThumbnailMovPath = flag.String("thumbnailmovpath", "", "Thumbnail mov 경로")
	flagPlatePath        = flag.String("platepath", "", "Plate 경로")
	// Commandline Args: User
	flagID                = flag.String("id", "", "user id")
	flagInitPass          = flag.String("initpass", "", "initialize user password")
	flagAccessLevel       = flag.Int("accesslevel", -1, "edit user Access Level")
	flagSignUpAccessLevel = flag.Int("signupaccesslevel", 3, "signup access level")
	// scan정보 추가. plate scan tool에서 데이터를 등록할 때 활용되는 옵션
	flagPlatesize       = flag.String("platesize", "", "스캔 플레이트 사이즈")
	flagTask            = flag.String("task", "", "태스크 이름. 예) model,mm.layout,ani,fx,mg,fur,sim,crowd,light,comp,matte,env")
	flagScantimecodein  = flag.String("scantimecodein", "00:00:00:00", "스캔 Timecode In")
	flagScantimecodeout = flag.String("scantimecodeout", "00:00:00:00", "스캔 Timecode Out")
	flagJusttimecodein  = flag.String("justtimecodein", "00:00:00:00", "Just구간 Timecode In")
	flagJusttimecodeout = flag.String("justtimecodeout", "00:00:00:00", "Just구간 Timecode Out")
	flagScanframe       = flag.Int("scanframe", 0, "스캔 총 프레임수")
	flagScanname        = flag.String("scanname", "", "스캔 폴더명")
	flagScanin          = flag.Int("scanin", -1, "스캔 In Frame")
	flagScanout         = flag.Int("scanout", -1, "스캔 Out Frame")
	flagJustin          = flag.Int("justin", -1, "Just In Frame")
	flagJustout         = flag.Int("justout", -1, "Just Out Frame")
	flagPlatein         = flag.Int("platein", -1, "플레이트 In Frame")
	flagPlateout        = flag.Int("plateout", -1, "플레이트 Out Frame")
	flagUpdateParent    = flag.Bool("updateparent", false, "org1,org2 형태의 재스캔 항목이라면 원본 org 정보를 업데이트 한다.")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("csi: ")
	flag.Usage = usage
	flag.Parse()
	if *flagVersion {
		fmt.Println("buildTime:", BUILDTIME)
		fmt.Println("git SHA1:", SHA1VER)
		os.Exit(0)
	}

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	ip, err := serviceIP()
	if err != nil {
		log.Fatal(err)
	}
	if *flagAdd == "project" && *flagName != "" { //프로젝트 추가
		addProjectCmd(*flagName)
		return
	} else if *flagRm == "project" && *flagName != "" { //프로젝트 삭제
		rmProjectCmd(*flagName)
		return
	} else if *flagAccessLevel != -1 && *flagID != "" {
		if user.Username != "root" {
			log.Fatal(errors.New("사용자의 레벨을 수정하기 위해서는 root 권한이 필요합니다"))
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		u, err := getUser(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		err = rmToken(session, u.ID)
		if err != nil {
			log.Fatal(err)
		}
		u.AccessLevel = AccessLevel(*flagAccessLevel)
		err = setUser(session, u)
		if err != nil {
			log.Fatal(err)
		}
		err = addToken(session, u)
		if err != nil {
			log.Fatal(err)
		}
		return

	} else if *flagInitPass != "" && *flagID != "" {
		if user.Username != "root" {
			log.Fatal(errors.New("사용자의 비밀변호를 변경하기 위해서는 root 권한이 필요합니다"))
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		err = initPassUser(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		u, err := getUser(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		err = rmToken(session, u.ID)
		if err != nil {
			log.Println(err)
		}
		err = addToken(session, u)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagRm == "division" && *flagID != "" { // division 삭제
		if user.Username != "root" {
			log.Fatal(errors.New("사용자를 삭제하기 위해서는 root 권한이 필요합니다"))
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		err = rmDivision(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagRm == "department" && *flagID != "" { // department 삭제
		if user.Username != "root" {
			log.Fatal(errors.New("사용자를 삭제하기 위해서는 root 권한이 필요합니다"))
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		err = rmDepartment(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagRm == "team" && *flagID != "" { // team 삭제
		if user.Username != "root" {
			log.Fatal(errors.New("사용자를 삭제하기 위해서는 root 권한이 필요합니다"))
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		err = rmTeam(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagRm == "role" && *flagID != "" { // role 삭제
		if user.Username != "root" {
			log.Fatal(errors.New("사용자를 삭제하기 위해서는 root 권한이 필요합니다"))
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		err = rmRole(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagRm == "position" && *flagID != "" { // position 삭제
		if user.Username != "root" {
			log.Fatal(errors.New("사용자를 삭제하기 위해서는 root 권한이 필요합니다"))
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		err = rmPosition(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagRm == "user" && *flagID != "" { // 사용자 삭제
		if user.Username != "root" {
			log.Fatal(errors.New("사용자를 삭제하기 위해서는 root 권한이 필요합니다"))
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		u, err := getUser(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		err = rmToken(session, u.ID)
		if err != nil {
			log.Fatal(err)
		}
		err = rmUser(session, u)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagAdd == "item" && *flagName != "" && *flagProject != "" && *flagType != "" { //아이템 추가
		switch *flagType {
		case "org", "left": // 일반영상은 org가 샷 타입이다. 입체프로젝트는 left가 샷타입이다.
			addShotItemCmd(*flagProject, *flagName, *flagType, *flagPlatesize, *flagScanname, *flagScantimecodein, *flagScantimecodeout, *flagJusttimecodein, *flagJusttimecodeout, *flagScanframe, *flagScanin, *flagScanout, *flagPlatein, *flagPlateout, *flagJustin, *flagJustout)
			dilog.Add(*flagDBIP, ip, "샷 생성되었습니다.", *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			dilog.Add(*flagDBIP, ip, "스캔이름 : "+*flagScanname, *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			dilog.Add(*flagDBIP, ip, fmt.Sprintf("스캔타임코드 : %s(%d) / %s(%d) (총%df)", *flagScantimecodein, *flagScanin, *flagScantimecodeout, *flagScanout, *flagScanframe), *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			dilog.Add(*flagDBIP, ip, fmt.Sprintf("플레이트 구간 : %d - %d", *flagPlatein, *flagPlateout), *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			dilog.Add(*flagDBIP, ip, "플레이트 사이즈 : "+*flagPlatesize, *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			return
		case "asset": //에셋 추가
			addAssetItemCmd(*flagProject, *flagName, *flagType, *flagAssettype, *flagAssettags)
			dilog.Add(*flagDBIP, ip, "에셋이 생성되었습니다.", *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			dilog.Add(*flagDBIP, ip, fmt.Sprintf("에셋타입 : %s, 에셋태그 : %s", *flagAssettype, *flagAssettags), *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			return
		default: //소스, 재스캔 추가
			addOtherItemCmd(*flagProject, *flagName, *flagType, *flagPlatesize, *flagScanname, *flagScantimecodein, *flagScantimecodeout, *flagJusttimecodein, *flagJusttimecodeout, *flagScanframe, *flagScanin, *flagScanout, *flagPlatein, *flagPlateout, *flagJustin, *flagJustout)
			dilog.Add(*flagDBIP, ip, "아이템이 생성되었습니다.", *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			dilog.Add(*flagDBIP, ip, "스캔이름 : "+*flagScanname, *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			dilog.Add(*flagDBIP, ip, fmt.Sprintf("스캔타임코드 : %s(%d) / %s(%d) (총%df)", *flagScantimecodein, *flagScanin, *flagScantimecodeout, *flagScanout, *flagScanframe), *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			dilog.Add(*flagDBIP, ip, fmt.Sprintf("플레이트 구간 : %d - %d", *flagPlatein, *flagPlateout), *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			dilog.Add(*flagDBIP, ip, "플레이트 사이즈 : "+*flagPlatesize, *flagProject, *flagName+"_"+*flagType, "csi3", user.Username, 180)
			if *flagUpdateParent {
				// updateParent 옵션이 활성화되어있고, org, left가 재스캔이라면..
				// 원본플레이트의 정보를 업데이트한다.
				if (*flagType != "org" && strings.Contains(*flagType, "org")) || (*flagType != "left" && strings.Contains(*flagType, "left")) {
					session, err := mgo.Dial(*flagDBIP)
					if err != nil {
						log.Fatal(err)
					}
					defer session.Close()
					err = SetImageSize(session, *flagProject, *flagName, "platesize", *flagPlatesize)
					if err != nil {
						log.Fatal(err)
					}
					err = SetTimecode(session, *flagProject, *flagName, "scantimecodein", *flagScantimecodein)
					if err != nil {
						log.Fatal(err)
					}
					err = SetTimecode(session, *flagProject, *flagName, "scantimecodeout", *flagScantimecodeout)
					if err != nil {
						log.Fatal(err)
					}
					if *flagJusttimecodein != "" {
						err = SetTimecode(session, *flagProject, *flagName, "justtimecodein", *flagJusttimecodein)
						if err != nil {
							log.Fatal(err)
						}
					}
					if *flagJusttimecodeout != "" {
						err = SetTimecode(session, *flagProject, *flagName, "justtimecodeout", *flagJusttimecodeout)
						if err != nil {
							log.Fatal(err)
						}
					}
					err = SetFrame(session, *flagProject, *flagName, "scanin", *flagScanin)
					if err != nil {
						log.Fatal(err)
					}
					err = SetFrame(session, *flagProject, *flagName, "scanout", *flagScanout)
					if err != nil {
						log.Fatal(err)
					}
					err = SetFrame(session, *flagProject, *flagName, "scanframe", *flagScanframe)
					if err != nil {
						log.Fatal(err)
					}
					err = SetFrame(session, *flagProject, *flagName, "platein", *flagPlatein)
					if err != nil {
						log.Fatal(err)
					}
					err = SetFrame(session, *flagProject, *flagName, "plateout", *flagPlateout)
					if err != nil {
						log.Fatal(err)
					}
					// Just In/Out 등록
					if *flagJustin > 0 {
						err = SetFrame(session, *flagProject, *flagName, "justin", *flagJustin)
						if err != nil {
							log.Fatal(err)
						}
					}
					if *flagJustout > 0 {
						err = SetFrame(session, *flagProject, *flagName, "justout", *flagJustout)
						if err != nil {
							log.Fatal(err)
						}
					}
					err = SetUseType(session, *flagProject, *flagName, *flagType)
					if err != nil {
						log.Fatal(err)
					}
					if strings.Contains(*flagType, "org") { // 일반샷이 재스캔 되었을 때 로그처리
						dilog.Add(*flagDBIP, ip, fmt.Sprintf("스캔타임코드 업데이트 : %s(%d) / %s(%d) (총%df)", *flagScantimecodein, *flagScanin, *flagScantimecodeout, *flagScanout, *flagScanframe), *flagProject, *flagName+"_"+"org", "csi3", user.Username, 180)
						dilog.Add(*flagDBIP, ip, fmt.Sprintf("플레이트 구간 업데이트 : %d - %d", *flagPlatein, *flagPlateout), *flagProject, *flagName+"_"+"org", "csi3", user.Username, 180)
						dilog.Add(*flagDBIP, ip, "플레이트 사이즈 업데이트 : "+*flagPlatesize, *flagProject, *flagName+"_"+"org", "csi3", user.Username, 180)
					}
					if strings.Contains(*flagType, "left") { // 입체샷이 재스캔 되었을 때 로그처리
						dilog.Add(*flagDBIP, ip, fmt.Sprintf("스캔타임코드 업데이트 : %s(%d) / %s(%d) (총%df)", *flagScantimecodein, *flagScanin, *flagScantimecodeout, *flagScanout, *flagScanframe), *flagProject, *flagName+"_"+"left", "csi3", user.Username, 180)
						dilog.Add(*flagDBIP, ip, fmt.Sprintf("플레이트 구간 업데이트 : %d - %d", *flagPlatein, *flagPlateout), *flagProject, *flagName+"_"+"left", "csi3", user.Username, 180)
						dilog.Add(*flagDBIP, ip, "플레이트 사이즈 업데이트 : "+*flagPlatesize, *flagProject, *flagName+"_"+"left", "csi3", user.Username, 180)
					}
				}
			}
			return
		}
	} else if *flagRm == "item" && *flagName != "" && *flagProject != "" && *flagType != "" { //아이템 삭제
		rmItemCmd(*flagProject, *flagName, *flagType)
		return
	} else if *flagHTTPPort != "" {
		if _, err := os.Stat(*flagThumbPath); err != nil {
			os.Stderr.WriteString("CSI에 사용되는 썸네일 경로가 존재하지 않습니다.\n")
			os.Stderr.WriteString("csi에서 생성되는 이미지, 사용자 프로필사진을 저장할 thumbnail 경로가 필요합니다.\n")
			os.Stderr.WriteString("명령어를 실행하는 곳에 thumbnail 폴더를 생성하거나,\n")
			os.Stderr.WriteString("-thumbpath 옵션을 이용하여 thumbnail로 사용될 경로를 지정하여 csi를 실행해주세요.\n")
			os.Exit(1)
		}
		// 만약 프로젝트가 하나도 없다면 "TEMP" 프로젝트를 생성한다. 프로젝트가 있어야 템플릿이 작동하기 때문이다.
		session, err := mgo.DialWithTimeout(*flagDBIP, 2*time.Second)
		if err != nil {
			log.Fatal("DB가 실행되고 있지 않습니다.")
		}
		plist, err := Projectlist(session)
		if err != nil {
			log.Fatal(err)
		}
		if len(plist) == 0 {
			p := *NewProject("TEMP")
			err = addProject(session, p)
			if err != nil {
				log.Fatal(err)
			}
		}
		session.Close()
		if *flagHTTPPort == ":80" {
			fmt.Printf("Service start: http://%s\n", ip)
		} else if *flagHTTPPort == ":443" {
			fmt.Printf("Service start: https://%s\n", ip)
		} else {
			fmt.Printf("Service start: http://%s%s\n", ip, *flagHTTPPort)
		}
		vfsTempates, err := LoadTemplates()
		if err != nil {
			log.Fatal(err)
		}
		TEMPLATES = vfsTempates
		webserver(*flagHTTPPort)
	} else if MatchNormalTime.MatchString(*flagDate) {
		// date 값이 데일리 형식이면 해당 날짜에 업로드된 mov를 RV를 통해 플레이한다.
		// 예: $ csi3 -date 2016-12-05 -play
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		dbProjectlist, err := Projectlist(session)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		// 해당프로젝트만 데일리를 위한 옵션
		// 만약 담당 프로젝트 감독님이 오면 해당 프로젝트 영상만 띄운다.
		reviewProjectlist := dbProjectlist
		if *flagProject != "" {
			err := HasProject(session, *flagProject)
			if err != nil {
				log.Fatalln(err)
			}
		}

		var playlist []string
		for _, project := range reviewProjectlist {
			op := SearchOption{
				Project:    project,
				Searchword: *flagDate,
				Assign:     true,
				Ready:      true,
				Wip:        true,
				Confirm:    true,
				Done:       true,
				Out:        true,
				Sortkey:    "slug",
			}
			items, err := Searchv2(session, op)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			// 만약 태스크명을 입력받았다면, 태스크명이 유효한지 체크하는 부분.
			if *flagTask != "" {
				hastask := false
				tasks, err := TasksettingNames(session)
				if err != nil {
					log.Fatal(err)
				}
				for _, t := range tasks {
					if *flagTask == strings.ToLower(t) {
						hastask = true
					}
				}
				if !hastask {
					log.Fatalf("%s Task 이름은 사용할 수 없습니다.\n", *flagTask)
				}
			}
			// 검색 옵션을 이용해서 daily 리스트를 만든다.
			for _, item := range items {
				tasks, err := TasksettingNames(session)
				if err != nil {
					log.Println(err)
				}
				for _, t := range tasks {
					if _, found := item.Tasks[t]; !found {
						continue
					}
					if *flagDate == ToNormalTime(item.Tasks[t].Mdate) && isMov(item.Tasks[t].Mov) {
						playlist = append(playlist, item.Tasks[t].Mov)
					}
				}
			}
		}
		// -play 인수가 붙어있다면, RV를 이용해서 플레이한다.
		if *flagPlay {
			if _, err := os.Stat(*flagRV); os.IsNotExist(err) {
				fmt.Println("RV가 rvpath에 존재하지 않습니다.")
				os.Exit(1)
			}
			out, err := exec.Command(*flagRV, playlist...).CombinedOutput()
			if err != nil {
				fmt.Printf("%v: %s\n", err, string(out))
				os.Exit(1)
			}
		}
		// -play 인수가 없다면, mov경로만 출력한다.
		for _, mov := range playlist {
			fmt.Println(mov)
		}
		return
	}
	if *flagHelp {
		flag.Usage()
		return
	}
	flag.Usage()
}
