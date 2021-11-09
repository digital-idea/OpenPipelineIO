package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dchest/captcha"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"gopkg.in/mgo.v2"
)

// MaxFileSize 사이즈는 웹에서 전송할 수 있는 최대 사이즈를 2기가로 제한한다.(인트라넷)
const MaxFileSize = 2000 * 1000 * 1024

// LoadTemplates 함수는 템플릿을 로딩합니다.
func LoadTemplates() (*template.Template, error) {
	t := template.New("").Funcs(funcMap)
	t, err := vfstemplate.ParseGlob(assets, t, "/template/*.html")
	return t, err
}

//템플릿 함수를 로딩합니다.
var funcMap = template.FuncMap{
	"AddProductionStartFrame":      AddProductionStartFrame,
	"title":                        strings.Title,
	"Split":                        strings.Split,
	"Join":                         strings.Join,
	"Parentpath":                   filepath.Dir,
	"projectStatus2color":          projectStatus2color,
	"Status2capString":             Status2capString, // regacy
	"Status2string":                Status2string,
	"name2seq":                     name2seq,
	"note2body":                    note2body,
	"pmnote2body":                  pmnote2body,
	"GetPath":                      GetPath,
	"ReverseStringSlice":           ReverseStringSlice,
	"ReverseCommentSlice":          ReverseCommentSlice,
	"SortByCreatetimeForPublishes": SortByCreatetimeForPublishes,
	"CutStringSlice":               CutStringSlice,
	"CutCommentSlice":              CutCommentSlice,
	"ToShortTime":                  ToShortTime,
	"ToNormalTime":                 ToNormalTime,
	"List2str":                     List2str,
	"CheckDate":                    CheckDate,
	"CheckUpdate":                  CheckUpdate,
	"CheckDdline":                  CheckDdline,
	"CheckDdlinev2":                CheckDdlinev2,
	"ToHumantime":                  ToHumantime,
	"Framecal":                     Framecal,
	"Add":                          Add,
	"Minus":                        Minus,
	"Scanname2RollMedia":           Scanname2RollMedia,
	"AddTagColon":                  AddTagColon, //Hashtag2tag,
	"Username2Elements":            Username2Elements,
	"RemovePath":                   RemovePath,
	"ShortPhoneNum":                ShortPhoneNum,
	"TaskStatus":                   TaskStatus,
	"TaskUser":                     TaskUser,
	"TaskDate":                     TaskDate,
	"TaskPredate":                  TaskPredate,
	"GetTaskLevel":                 GetTaskLevel,
	"ProductionVersionFormat":      ProductionVersionFormat,
	"Protocol":                     Protocol,
	"RmProtocol":                   RmProtocol,
	"ProtocolTarget":               ProtocolTarget,
	"userInfo":                     userInfo,
	"onlyID":                       onlyID,
	"mapToSlice":                   mapToSlice,
	"hasStatus":                    hasStatus,
	"GenPageNums":                  GenPageNums,
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "NotFound 404")
	}
}

// 도움말 페이지 입니다.
func handleHelp(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type recipe struct {
		Wfs     string
		User    User
		Devmode bool
		SearchOption
		Sha1ver   string
		BuildTime string
		Status    []Status
		Stages    []Stage
		DBIP      string
		DBVer     string
		ServerIP  string
	}

	rcp := recipe{}
	rcp.Sha1ver = SHA1VER
	rcp.BuildTime = BUILDTIME
	rcp.DBIP = *flagDBIP
	info, err := session.BuildInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.DBVer = info.Version
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Devmode = *flagDevmode
	rcp.User = u
	rcp.Wfs = *flagWFS
	rcp.Stages, err = AllStages(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Status, err = AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ip, err := serviceIP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.ServerIP = ip
	err = t.ExecuteTemplate(w, "help", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	type recipe struct {
		IP         string
		Web        bool
		DB         bool
		MountPoint bool
		All        bool
	}
	rcp := recipe{}
	// IP구하기
	ip, err := serviceIP()
	if err != nil {
		rcp.IP = ""
	}
	rcp.IP = ip
	// DB 채크

	// 웹서버 체크
	rcp.Web = true

	// Mount Point 경로 존재하는지 체크
	_, err = os.Stat(CachedAdminSetting.RootPath)
	if os.IsNotExist(err) {
		rcp.MountPoint = false
	}
	rcp.MountPoint = true

	// json 으로 결과 전송
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 전송되는 컨텐츠의 캐쉬 수명을 설정하는 핸들러입니다.
func maxAgeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", *flagThumbnailAge))
		h.ServeHTTP(w, r)
	})
}

// webserver함수는 웹서버의 URL을 선언하는 함수입니다.
func webserver(port string) {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assets)))
	http.Handle("/thumbnail/", maxAgeHandler(http.StripPrefix("/thumbnail/", http.FileServer(http.Dir(CachedAdminSetting.ThumbnailRootPath)))))
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight)) // Captcha
	// Item
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/searchsubmit", handleSearchSubmit)
	http.HandleFunc("/searchsubmitv2", handleSearchSubmitV2)
	http.HandleFunc("/help", handleHelp)
	http.HandleFunc("/setellite", handleSetellite)
	http.HandleFunc("/uploadsetellite", handleUploadSetellite)
	http.HandleFunc("/addshot", handleAddShot)
	http.HandleFunc("/addshot_submit", handleAddShotSubmit)
	http.HandleFunc("/addasset", handleAddAsset)
	http.HandleFunc("/addasset_submit", handleAddAssetSubmit)
	http.HandleFunc("/detail", handleItemDetail)

	// Review
	http.HandleFunc("/daily", handleDaily)
	http.HandleFunc("/review", handleReview)
	http.HandleFunc("/reviewdata", handleReviewData)
	http.HandleFunc("/reviewdrawingdata", handleReviewDrawingData)
	http.HandleFunc("/review-submit", handleReviewSubmit)
	http.HandleFunc("/upload-reviewfile", handleUploadReviewFile)

	// Project
	http.HandleFunc("/projectinfo", handleProjectinfo)
	http.HandleFunc("/addproject", handleAddProject)
	http.HandleFunc("/addproject_submit", handleAddProjectSubmit)
	http.HandleFunc("/editproject", handleEditProject)
	http.HandleFunc("/editproject_submit", handleEditProjectSubmit)
	http.HandleFunc("/rmproject", handleRmProject)
	http.HandleFunc("/rmproject_submit", handleRmProjectSubmit)
	http.HandleFunc("/noonproject", handleNoOnProject)

	// User
	http.HandleFunc("/signup", handleSignup)
	http.HandleFunc("/signup_submit", handleSignupSubmit)
	http.HandleFunc("/signin", handleSignin)
	http.HandleFunc("/signin_submit", handleSigninSubmit)
	http.HandleFunc("/signin_success", handleSigninSuccess)
	http.HandleFunc("/signout", handleSignout)
	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/updatepassword", handleUpdatePassword)
	http.HandleFunc("/updatepassword_submit", handleUpdatePasswordSubmit)
	http.HandleFunc("/edituser", handleEditUser)
	http.HandleFunc("/edituser-submit", handleEditUserSubmit)
	http.HandleFunc("/replacetag", handleReplaceTag)
	http.HandleFunc("/replacetag_submit", handleReplaceTagSubmit)
	http.HandleFunc("/invalidaccess", handleInvalidAccess)
	http.HandleFunc("/invalidpass", handleInvalidPass)
	http.HandleFunc("/nouser", handleNoUser)

	// Admin Setting
	http.HandleFunc("/adminsetting", handleAdminSetting)
	http.HandleFunc("/adminsetting_submit", handleAdminSettingSubmit)
	http.HandleFunc("/setadminsetting", handleSetAdminSetting)

	// Organization
	http.HandleFunc("/divisions", handleDivisions)
	http.HandleFunc("/departments", handleDepartments)
	http.HandleFunc("/teams", handleTeams)
	http.HandleFunc("/roles", handleRoles)
	http.HandleFunc("/positions", handlePositions)
	http.HandleFunc("/adddivision", handleAddOrganization)
	http.HandleFunc("/editdivision", handleEditDivision)
	http.HandleFunc("/editdivisionsubmit", handleEditDivisionSubmit)
	http.HandleFunc("/adddepartment", handleAddOrganization)
	http.HandleFunc("/editdepartment", handleEditDepartment)
	http.HandleFunc("/editdepartmentsubmit", handleEditDepartmentSubmit)
	http.HandleFunc("/addteam", handleAddOrganization)
	http.HandleFunc("/editteam", handleEditTeam)
	http.HandleFunc("/editteamsubmit", handleEditTeamSubmit)
	http.HandleFunc("/addrole", handleAddOrganization)
	http.HandleFunc("/editrole", handleEditRole)
	http.HandleFunc("/editrolesubmit", handleEditRoleSubmit)
	http.HandleFunc("/addposition", handleAddOrganization)
	http.HandleFunc("/editposition", handleEditPosition)
	http.HandleFunc("/editpositionsubmit", handleEditPositionSubmit)
	http.HandleFunc("/adddivisionsubmit", handleAddDivisionSubmit)
	http.HandleFunc("/adddepartmentsubmit", handleAddDepartmentSubmit)
	http.HandleFunc("/addteamsubmit", handleAddTeamSubmit)
	http.HandleFunc("/addrolesubmit", handleAddRoleSubmit)
	http.HandleFunc("/addpositionsubmit", handleAddPositionSubmit)
	http.HandleFunc("/rmorganization", handleRmOrganization)
	http.HandleFunc("/rmorganization-submit", handleRmOrganizationSubmit)

	// Export: Excel, Json, Csv, Dump
	http.HandleFunc("/importexcel", handleImportExcel)
	http.HandleFunc("/importjson", handleImportJSON)
	http.HandleFunc("/exportexcel", handleExportExcel)
	http.HandleFunc("/exportjson", handleExportJSON)
	http.HandleFunc("/reportexcel", handleReportExcel)
	http.HandleFunc("/reportjson", handleReportJSON)
	http.HandleFunc("/excel-submit", handleExcelSubmit)
	http.HandleFunc("/json-submit", handleJSONSubmit)
	http.HandleFunc("/exportexcel-submit", handleExportExcelSubmit)
	http.HandleFunc("/exportjson-submit", handleExportJSONSubmit)
	http.HandleFunc("/upload-excel", handleUploadExcel)
	http.HandleFunc("/upload-json", handleUploadJSON)
	http.HandleFunc("/download-excel-template", handleDownloadExcelTemplate)
	http.HandleFunc("/download-excel-file", handleDownloadExcelFile)
	http.HandleFunc("/download-json-file", handleDownloadJSONFile)
	http.HandleFunc("/download-csv-file", handleDownloadCsvFile)
	http.HandleFunc("/export-dump-project", handleExportDumpProject)

	// Task
	http.HandleFunc("/tasksettings", handleTasksettings)
	http.HandleFunc("/addtasksetting", handleAddTasksetting)
	http.HandleFunc("/rmtasksetting", handleRmTasksetting)
	http.HandleFunc("/edittasksetting", handleEditTasksetting)
	http.HandleFunc("/addtasksetting-submit", handleAddTasksettingSubmit)
	http.HandleFunc("/rmtasksetting-submit", handleRmTasksettingSubmit)
	http.HandleFunc("/edittasksetting-submit", handleEditTasksettingSubmit)

	// Status
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/addstatus", handleAddStatus)
	http.HandleFunc("/addstatus-submit", handleAddStatusSubmit)
	http.HandleFunc("/editstatus", handleEditStatus)
	http.HandleFunc("/editstatus-submit", handleEditStatusSubmit)
	http.HandleFunc("/rmstatus", handleRmStatus)
	http.HandleFunc("/rmstatus-submit", handleRmStatusSubmit)

	// Stage
	http.HandleFunc("/stage", handleStage)
	http.HandleFunc("/addstage", handleAddStage)
	http.HandleFunc("/addstage-submit", handleAddStageSubmit)
	http.HandleFunc("/editstage", handleEditStage)
	http.HandleFunc("/editstage-submit", handleEditStageSubmit)
	http.HandleFunc("/rmstage", handleRmStage)
	http.HandleFunc("/rmstage-submit", handleRmStageSubmit)

	// Publish Key
	http.HandleFunc("/publishkey", handlePublishKey)
	http.HandleFunc("/addpublishkey", handleAddPublishKey)
	http.HandleFunc("/addpublishkey-submit", handleAddPublishKeySubmit)
	http.HandleFunc("/editpublishkey", handleEditPublishKey)
	http.HandleFunc("/editpublishkey-submit", handleEditPublishKeySubmit)
	http.HandleFunc("/rmpublishkey", handleRmPublishKey)
	http.HandleFunc("/rmpublishkey-submit", handleRmPublishKeySubmit)

	// Input
	http.HandleFunc("/inputmode", handleInputMode)

	// Error
	http.HandleFunc("/error-captcha", handleErrorCaptcha)

	//Health
	http.HandleFunc("/health", handleHealth)

	// restAPI Project
	http.HandleFunc("/api/project", handleAPIProject)
	http.HandleFunc("/api/projects", handleAPIProjects)
	http.HandleFunc("/api/addproject", handleAPIAddproject)
	http.HandleFunc("/api/projecttags", handleAPIProjectTags)
	http.HandleFunc("/api/projectassettags", handleAPIProjectAssetTags)

	// restAPI Onset(Setellite)
	http.HandleFunc("/api/setellite", handleAPISetelliteItems)
	http.HandleFunc("/api/setellitesearch", handleAPISetelliteSearch)

	// restAPI Item
	http.HandleFunc("/api/timeinfo", handleAPITimeinfo)
	http.HandleFunc("/api/item", handleAPIItem) // legacy
	http.HandleFunc("/api2/item", handleAPI2Item)
	http.HandleFunc("/api/rmitem", handleAPIRmItem) // legacy
	http.HandleFunc("/api/rmitemid", handleAPIRmItemID)
	http.HandleFunc("/api/items", handleAPI2Items)  // legacy
	http.HandleFunc("/api2/items", handleAPI2Items) // legacy
	http.HandleFunc("/api3/items", handleAPI3Items)
	http.HandleFunc("/api/searchname", handleAPISearchname)
	http.HandleFunc("/api/seqs", handleAPISeqs)
	http.HandleFunc("/api/allshots", handleAPIAllShots)
	http.HandleFunc("/api/shots", handleAPIShots) // legacy
	http.HandleFunc("/api2/shots", handleAPI2Shots)
	http.HandleFunc("/api/shot", handleAPIShot)
	http.HandleFunc("/api/asset", handleAPIAsset)
	http.HandleFunc("/api/assets", handleAPIAssets)
	http.HandleFunc("/api/setplatesize", handleAPISetPlateSize)
	http.HandleFunc("/api/setundistortionsize", handleAPISetUnDistortionSize)
	http.HandleFunc("/api/setrendersize", handleAPISetRenderSize) // legacy
	http.HandleFunc("/api2/setrendersize", handleAPI2SetRenderSize)
	http.HandleFunc("/api/setoverscanratio", handleAPISetOverscanRatio)
	http.HandleFunc("/api/setcamerapubpath", handleAPISetCameraPubPath)
	http.HandleFunc("/api/setcamerapubtask", handleAPISetCameraPubTask)
	http.HandleFunc("/api/setcameralensmm", handleAPISetCameraLensmm)
	http.HandleFunc("/api/setcameraprojection", handleAPISetCameraProjection)
	http.HandleFunc("/api/setseq", handleAPISetSeq)
	http.HandleFunc("/api/setseason", handleAPISetSeason)
	http.HandleFunc("/api/setepisode", handleAPISetEpisode)
	http.HandleFunc("/api/setplatepath", handleAPISetPlatePath)
	http.HandleFunc("/api/setthummov", handleAPISetThummov) // legacy
	http.HandleFunc("/api2/setthummov", handleAPI2SetThummov)
	http.HandleFunc("/api/setbeforemov", handleAPISetBeforemov)
	http.HandleFunc("/api/setaftermov", handleAPISetAftermov)
	http.HandleFunc("/api/seteditmov", handleAPISetEditmov)
	http.HandleFunc("/api/settaskstatus", handleAPISetTaskStatus) // legacy
	http.HandleFunc("/api2/settaskstatus", handleAPI2SetTaskStatus)
	http.HandleFunc("/api/taskstatusnum", handleAPITaskStatusNum)
	http.HandleFunc("/api/taskanduserstatusnum", handleAPITaskAndUserStatusNum)
	http.HandleFunc("/api/userstatusnum", handleAPIUserStatusNum)
	http.HandleFunc("/api/statusnum", handleAPIStatusNum)
	http.HandleFunc("/api/addtask", handleAPIAddTask)
	http.HandleFunc("/api/rmtask", handleAPIRmTask)
	http.HandleFunc("/api/settaskuser", handleAPISetTaskUser)
	http.HandleFunc("/api/settaskusercomment", handleAPISetTaskUserComment)
	http.HandleFunc("/api/setplatein", handleAPISetPlateIn)
	http.HandleFunc("/api/setplateout", handleAPISetPlateOut)
	http.HandleFunc("/api/setjustin", handleAPISetJustIn)
	http.HandleFunc("/api/setjustout", handleAPISetJustOut)
	http.HandleFunc("/api/setscanin", handleAPISetScanIn)
	http.HandleFunc("/api/setscanout", handleAPISetScanOut)
	http.HandleFunc("/api/setscanframe", handleAPISetScanFrame)
	http.HandleFunc("/api/sethandlein", handleAPISetHandleIn)
	http.HandleFunc("/api/sethandleout", handleAPISetHandleOut)
	http.HandleFunc("/api/setshottype", handleAPISetShotType)
	http.HandleFunc("/api/setusetype", handleAPISetUseType)
	http.HandleFunc("/api/setassettype", handleAPISetAssetType)
	http.HandleFunc("/api/setoutputname", handleAPISetOutputName)
	http.HandleFunc("/api/setrnum", handleAPISetRnum)
	http.HandleFunc("/api/setdeadline2d", handleAPISetDeadline2D)
	http.HandleFunc("/api/setdeadline3d", handleAPISetDeadline3D)
	http.HandleFunc("/api/setscantimecodein", handleAPISetScanTimecodeIn)
	http.HandleFunc("/api/setscantimecodeout", handleAPISetScanTimecodeOut)
	http.HandleFunc("/api/setjusttimecodein", handleAPISetJustTimecodeIn)
	http.HandleFunc("/api/setjusttimecodeout", handleAPISetJustTimecodeOut)
	http.HandleFunc("/api/setfinver", handleAPISetFinver)
	http.HandleFunc("/api/setfindate", handleAPISetFindate)
	http.HandleFunc("/api/addtag", handleAPIAddTag)
	http.HandleFunc("/api/renametag", handleAPIRenameTag)
	http.HandleFunc("/api/rmtag", handleAPIRmTag)
	http.HandleFunc("/api/settags", handleAPISetTags)
	http.HandleFunc("/api/setnote", handleAPISetNote)
	http.HandleFunc("/api/addcomment", handleAPIAddComment)
	http.HandleFunc("/api/editcomment", handleAPIEditComment)
	http.HandleFunc("/api/rmcomment", handleAPIRmComment)
	http.HandleFunc("/api/addsource", handleAPIAddSource)
	http.HandleFunc("/api/rmsource", handleAPIRmSource)
	http.HandleFunc("/api/addreference", handleAPIAddReference)
	http.HandleFunc("/api/rmreference", handleAPIRmReference)
	http.HandleFunc("/api/search", handleAPISearch)
	http.HandleFunc("/api/deadline2d", handleAPIDeadline2D)
	http.HandleFunc("/api/deadline3d", handleAPIDeadline3D)
	http.HandleFunc("/api/setstatus", handleAPISetTaskStatus)
	http.HandleFunc("/api/settaskmov", handleAPISetTaskMov) // legacy
	http.HandleFunc("/api2/settaskmov", handleAPI2SetTaskMov)
	http.HandleFunc("/api/settaskusernote", handleAPISetTaskUserNote)
	http.HandleFunc("/api/setretimeplate", handleAPISetRetimePlate)
	http.HandleFunc("/api/settasklevel", handleAPISetTaskLevel)
	http.HandleFunc("/api/setobjectid", handleAPISetObjectID)
	http.HandleFunc("/api/setociocc", handleAPISetOCIOcc)
	http.HandleFunc("/api/setrollmedia", handleAPISetRollmedia)
	http.HandleFunc("/api/setscanname", handleAPISetScanname)
	http.HandleFunc("/api/settaskdate", handleAPISetTaskDate)
	http.HandleFunc("/api/settaskexpectday", handleAPISetTaskExpectDay)
	http.HandleFunc("/api/settaskresultday", handleAPISetTaskResultDay)
	http.HandleFunc("/api/settaskpredate", handleAPISetTaskPredate)
	http.HandleFunc("/api/settaskstartdate", handleAPISetTaskStartdate)
	http.HandleFunc("/api/task", handleAPITask)
	http.HandleFunc("/api/shottype", handleAPIShottype)
	http.HandleFunc("/api/setcrowdasset", handleAPISetCrowdAsset)
	http.HandleFunc("/api/mailinfo", handleAPIMailInfo)
	http.HandleFunc("/api/usetypes", handleAPIUseTypes)
	http.HandleFunc("/api/publish", handleAPIAddTaskPublish) // legacy
	http.HandleFunc("/api/addpublish", handleAPIAddTaskPublish)
	http.HandleFunc("/api/setpublishstatus", handleAPISetTaskPublishStatus)
	http.HandleFunc("/api/rmpublish", handleAPIRmTaskPublish)
	http.HandleFunc("/api/rmpublishkey", handleAPIRmTaskPublishKey)
	http.HandleFunc("/api/uploadthumbnail", handleAPIUploadThumbnail)
	http.HandleFunc("/api/setnetflixid", handleAPISetNetflixID)

	// restAPI USER
	http.HandleFunc("/api/user", handleAPIUser) // legacy
	http.HandleFunc("/api2/user", handleAPI2User)
	http.HandleFunc("/api/users", handleAPISearchUser)
	http.HandleFunc("/api/validuser", handleAPIValidUser) // 보안취약점 이슈가 있다. 다른 툴과 쉽게 연동할 때 편리하다. 보안레벨을 높게 올릴때는 허용하지 않도록 한다.
	http.HandleFunc("/api/setleaveuser", handleAPISetLeaveUser)
	http.HandleFunc("/api/autocompliteusers", handleAPIAutoCompliteUsers)
	http.HandleFunc("/api/initpassword", handleAPIInitPassword)
	http.HandleFunc("/api/ansiblehosts", handleAPIAnsibleHosts)

	// restAPI Organization
	http.HandleFunc("/api/teams", handleAPIAllTeams)

	// restAPI Tasksetting
	http.HandleFunc("/api/tasksetting", handleAPITasksetting)
	http.HandleFunc("/api/shottasksetting", handleAPIShotTasksetting)
	http.HandleFunc("/api/assettasksetting", handleAPIAssetTasksetting)
	http.HandleFunc("/api/categorytasksettings", handleAPICategoryTasksettings)

	// restAPI Status
	http.HandleFunc("/api/status", handleAPIStatus)
	http.HandleFunc("/api/addstatus", handleAPIAddStatus)

	// restAPI PublishKey
	http.HandleFunc("/api/publishkeys", handleAPIPublishKeys)
	http.HandleFunc("/api/getpublish", handleAPIGetPublish)

	// restAPI Review
	http.HandleFunc("/api/addreview", handleAPIAddReview)
	http.HandleFunc("/api/review", handleAPIReview)
	http.HandleFunc("/api/searchreview", handleAPISearchReview)
	http.HandleFunc("/api/setreviewstatus", handleAPISetReviewStatus)
	http.HandleFunc("/api/setreviewstage", handleAPISetReviewStage)
	http.HandleFunc("/api/setreviewnextstatus", handleAPISetReviewNextStatus)
	http.HandleFunc("/api/setreviewnextstage", handleAPISetReviewNextStage)
	http.HandleFunc("/api/addreviewcomment", handleAPIAddReviewComment)
	http.HandleFunc("/api/editreviewcomment", handleAPIEditReviewComment)
	http.HandleFunc("/api/rmreviewcomment", handleAPIRmReviewComment)
	http.HandleFunc("/api/rmreview", handleAPIRmReview)
	http.HandleFunc("/api/setreviewproject", handleAPISetReviewProject)
	http.HandleFunc("/api/setreviewtask", handleAPISetReviewTask)
	http.HandleFunc("/api/setreviewname", handleAPISetReviewName)
	http.HandleFunc("/api/setreviewpath", handleAPISetReviewPath)
	http.HandleFunc("/api/setreviewcreatetime", handleAPISetReviewCreatetime)
	http.HandleFunc("/api/setreviewmainversion", handleAPISetReviewMainVersion)
	http.HandleFunc("/api/setreviewsubversion", handleAPISetReviewSubVersion)
	http.HandleFunc("/api/setreviewfps", handleAPISetReviewFps)
	http.HandleFunc("/api/setreviewdescription", handleAPISetReviewDescription)
	http.HandleFunc("/api/setreviewcamerainfo", handleAPISetReviewCameraInfo)
	http.HandleFunc("/api/setreviewprocessstatus", handleAPISetReviewProcessStatus)
	http.HandleFunc("/api/uploadreviewdrawing", handleAPIUploadReviewDrawing)
	http.HandleFunc("/api/rmreviewdrawing", handleAPIRmReviewDrawing)
	http.HandleFunc("/api/reviewdrawingframe", handleAPIReviewDrawingFrame)
	http.HandleFunc("/api/setreviewagainforwaitstatustoday", handleAPISetReviewAgainForWaitStatusToday)
	http.HandleFunc("/api/reviewoutputdatapath", handleAPIReviewOutputDataPath)

	// Deprecated: 사용하지 않는 url, 과거호환성을 위해서 남겨둠
	http.HandleFunc("/edititem", handleEditItem)                    // legacy
	http.HandleFunc("/editeditem", handleEditedItem)                // legacy
	http.HandleFunc("/api/setmov", handleAPISetTaskMov)             // legacy
	http.HandleFunc("/api/setstartdate", handleAPISetTaskStartdate) // legacy
	http.HandleFunc("/edititem-submit", handleEditItemSubmitv2)     // legacy

	if port == ":443" || port == ":8443" { // https ports
		err := http.ListenAndServeTLS(port, *flagCertFullchanin, *flagCertPrivkey, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
