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
	"github.com/gorilla/mux"
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
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
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
	// DB 체크
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	err = session.Ping()
	if err != nil {
		rcp.DB = false
	}
	rcp.DB = true

	// 웹서버 체크
	rcp.Web = true

	// Mount Point 경로 존재하는지 체크
	_, err = os.Stat(CachedAdminSetting.RootPath)
	if os.IsNotExist(err) {
		rcp.MountPoint = false
	}
	rcp.MountPoint = true

	// 모든 요소가 true인지 체크
	isIPexist := false
	if rcp.IP != "" {
		isIPexist = true
	}
	rcp.All = isIPexist && rcp.Web && rcp.DB && rcp.MountPoint

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

func helpMethodOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
}

// webserver함수는 웹서버의 URL을 선언하는 함수입니다.
func webserver(port string) {
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(assets)))
	r.PathPrefix("/thumbnail/").Handler(maxAgeHandler(http.StripPrefix("/thumbnail/", http.FileServer(http.Dir(CachedAdminSetting.ThumbnailRootPath)))))
	r.PathPrefix("/captcha/").Handler(captcha.Server(captcha.StdWidth, captcha.StdHeight)) // Captcha
	// Item
	r.HandleFunc("/", handleIndex)
	r.HandleFunc("/inputmode", handleInputMode)
	r.HandleFunc("/searchsubmit", handleSearchSubmit)
	r.HandleFunc("/searchsubmitv2", handleSearchSubmitV2)
	r.HandleFunc("/help", handleHelp)
	r.HandleFunc("/setellite", handleSetellite)
	r.HandleFunc("/uploadsetellite", handleUploadSetellite)
	r.HandleFunc("/addshot", handleAddShot)
	r.HandleFunc("/addshot_submit", handleAddShotSubmit)
	r.HandleFunc("/addasset", handleAddAsset)
	r.HandleFunc("/addasset_submit", handleAddAssetSubmit)
	r.HandleFunc("/detail", handleItemDetail)

	// Review
	r.HandleFunc("/daily-review-stage", handleDailyReviewStage)
	r.HandleFunc("/daily-review-status", handleDailyReviewStatus)
	r.HandleFunc("/reviewstage", handleReviewStage)
	r.HandleFunc("/reviewstatus", handleReviewStatus)
	r.HandleFunc("/reviewdata", handleReviewData)
	r.HandleFunc("/reviewdrawingdata", handleReviewDrawingData)
	r.HandleFunc("/review-stage-submit", handleReviewStageSubmit)
	r.HandleFunc("/review-status-submit", handleReviewStatusSubmit)
	r.HandleFunc("/upload-reviewfile", handleUploadReviewFile)

	// Project
	r.HandleFunc("/projectinfo", handleProjectinfo)
	r.HandleFunc("/addproject", handleAddProject)
	r.HandleFunc("/addproject_submit", handleAddProjectSubmit)
	r.HandleFunc("/editproject", handleEditProject)
	r.HandleFunc("/editproject_submit", handleEditProjectSubmit)
	r.HandleFunc("/rmproject", handleRmProject)
	r.HandleFunc("/rmproject_submit", handleRmProjectSubmit)
	r.HandleFunc("/noonproject", handleNoOnProject)

	// User
	r.HandleFunc("/signup", handleSignup)
	r.HandleFunc("/signup_submit", handleSignupSubmit)
	r.HandleFunc("/signin", handleSignin)
	r.HandleFunc("/signin_submit", handleSigninSubmit)
	r.HandleFunc("/signin_success", handleSigninSuccess)
	r.HandleFunc("/signout", handleSignout)
	r.HandleFunc("/user", handleUser)
	r.HandleFunc("/users", handleUsers)
	r.HandleFunc("/updatepassword", handleUpdatePassword)
	r.HandleFunc("/updatepassword_submit", handleUpdatePasswordSubmit)
	r.HandleFunc("/edituser", handleEditUser)
	r.HandleFunc("/edituser-submit", handleEditUserSubmit)
	r.HandleFunc("/replacetag", handleReplaceTag)
	r.HandleFunc("/replacetag_submit", handleReplaceTagSubmit)
	r.HandleFunc("/invalidaccess", handleInvalidAccess)
	r.HandleFunc("/invalidpass", handleInvalidPass)
	r.HandleFunc("/nouser", handleNoUser)

	// Partner
	r.HandleFunc("/partners", handlePartners)
	r.HandleFunc("/addpartner", handleAddPartner)
	r.HandleFunc("/addpartner_submit", handleAddPartnerSubmit)

	// Admin Setting
	r.HandleFunc("/adminsetting", handleAdminSetting)
	r.HandleFunc("/adminsetting_submit", handleAdminSettingSubmit)
	r.HandleFunc("/setadminsetting", handleSetAdminSetting)

	// Organization
	r.HandleFunc("/divisions", handleDivisions)
	r.HandleFunc("/departments", handleDepartments)
	r.HandleFunc("/teams", handleTeams)
	r.HandleFunc("/roles", handleRoles)
	r.HandleFunc("/positions", handlePositions)
	r.HandleFunc("/adddivision", handleAddOrganization)
	r.HandleFunc("/editdivision", handleEditDivision)
	r.HandleFunc("/editdivisionsubmit", handleEditDivisionSubmit)
	r.HandleFunc("/adddepartment", handleAddOrganization)
	r.HandleFunc("/editdepartment", handleEditDepartment)
	r.HandleFunc("/editdepartmentsubmit", handleEditDepartmentSubmit)
	r.HandleFunc("/addteam", handleAddOrganization)
	r.HandleFunc("/editteam", handleEditTeam)
	r.HandleFunc("/editteamsubmit", handleEditTeamSubmit)
	r.HandleFunc("/addrole", handleAddOrganization)
	r.HandleFunc("/editrole", handleEditRole)
	r.HandleFunc("/editrolesubmit", handleEditRoleSubmit)
	r.HandleFunc("/addposition", handleAddOrganization)
	r.HandleFunc("/editposition", handleEditPosition)
	r.HandleFunc("/editpositionsubmit", handleEditPositionSubmit)
	r.HandleFunc("/adddivisionsubmit", handleAddDivisionSubmit)
	r.HandleFunc("/adddepartmentsubmit", handleAddDepartmentSubmit)
	r.HandleFunc("/addteamsubmit", handleAddTeamSubmit)
	r.HandleFunc("/addrolesubmit", handleAddRoleSubmit)
	r.HandleFunc("/addpositionsubmit", handleAddPositionSubmit)
	r.HandleFunc("/rmorganization", handleRmOrganization)
	r.HandleFunc("/rmorganization-submit", handleRmOrganizationSubmit)

	// Export: Excel, Json, Csv, Dump
	r.HandleFunc("/importexcel", handleImportExcel)
	r.HandleFunc("/importjson", handleImportJSON)
	r.HandleFunc("/exportexcel", handleExportExcel)
	r.HandleFunc("/exportjson", handleExportJSON)
	r.HandleFunc("/reportexcel", handleReportExcel)
	r.HandleFunc("/reportjson", handleReportJSON)
	r.HandleFunc("/excel-submit", handleExcelSubmit)
	r.HandleFunc("/json-submit", handleJSONSubmit)
	r.HandleFunc("/exportexcel-submit", handleExportExcelSubmit)
	r.HandleFunc("/exportjson-submit", handleExportJSONSubmit)
	r.HandleFunc("/upload-excel", handleUploadExcel)
	r.HandleFunc("/upload-json", handleUploadJSON)
	r.HandleFunc("/download-excel-template", handleDownloadExcelTemplate)
	r.HandleFunc("/download-excel-file", handleDownloadExcelFile)
	r.HandleFunc("/download-json-file", handleDownloadJSONFile)
	r.HandleFunc("/download-csv-file", handleDownloadCsvFile)
	r.HandleFunc("/export-dump-project", handleExportDumpProject)

	// Task
	r.HandleFunc("/tasksettings", handleTasksettings)
	r.HandleFunc("/addtasksetting", handleAddTasksetting)
	r.HandleFunc("/rmtasksetting", handleRmTasksetting)
	r.HandleFunc("/edittasksetting", handleEditTasksetting)
	r.HandleFunc("/addtasksetting-submit", handleAddTasksettingSubmit)
	r.HandleFunc("/rmtasksetting-submit", handleRmTasksettingSubmit)
	r.HandleFunc("/edittasksetting-submit", handleEditTasksettingSubmit)

	// Status
	r.HandleFunc("/status", handleStatus)
	r.HandleFunc("/addstatus", handleAddStatus)
	r.HandleFunc("/addstatus-submit", handleAddStatusSubmit)
	r.HandleFunc("/editstatus", handleEditStatus)
	r.HandleFunc("/editstatus-submit", handleEditStatusSubmit)
	r.HandleFunc("/rmstatus", handleRmStatus)
	r.HandleFunc("/rmstatus-submit", handleRmStatusSubmit)

	// Stage
	r.HandleFunc("/stage", handleStage)
	r.HandleFunc("/addstage", handleAddStage)
	r.HandleFunc("/addstage-submit", handleAddStageSubmit)
	r.HandleFunc("/editstage", handleEditStage)
	r.HandleFunc("/editstage-submit", handleEditStageSubmit)
	r.HandleFunc("/rmstage", handleRmStage)
	r.HandleFunc("/rmstage-submit", handleRmStageSubmit)

	// Publish Key
	r.HandleFunc("/publishkey", handlePublishKey)
	r.HandleFunc("/addpublishkey", handleAddPublishKey)
	r.HandleFunc("/addpublishkey-submit", handleAddPublishKeySubmit)
	r.HandleFunc("/editpublishkey", handleEditPublishKey)
	r.HandleFunc("/editpublishkey-submit", handleEditPublishKeySubmit)
	r.HandleFunc("/rmpublishkey", handleRmPublishKey)
	r.HandleFunc("/rmpublishkey-submit", handleRmPublishKeySubmit)

	// Error
	r.HandleFunc("/error-captcha", handleErrorCaptcha)

	//Health & Help
	r.HandleFunc("/health", handleHealth)
	r.HandleFunc("/api/statusinfo", handleAPIStatusInfo).Methods(http.MethodGet, http.MethodOptions)

	// Statistics
	r.HandleFunc("/api/statistics/projectnum", handleAPIStatisticsProjectnum).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/shot", handleAPI1StatisticsShot).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api2/statistics/shot", handleAPI2StatisticsShot).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/asset", handleAPI1StatisticsAsset).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api2/statistics/asset", handleAPI2StatisticsAsset).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/task", handleAPI1StatisticsTask).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/tag", handleAPI1StatisticsTag).Methods(http.MethodGet, http.MethodOptions)

	// restAPI Project
	r.HandleFunc("/api/project", handleAPIProject)
	r.HandleFunc("/api/projects", handleAPIProjects)
	r.HandleFunc("/api/addproject", handleAPIAddproject)
	r.HandleFunc("/api/projecttags", handleAPIProjectTags)
	r.HandleFunc("/api/projectassettags", handleAPIProjectAssetTags)

	// restAPI Onset(Setellite)
	r.HandleFunc("/api/setellite", handleAPISetelliteItems)
	r.HandleFunc("/api/setellitesearch", handleAPISetelliteSearch)

	// restAPI Item
	r.HandleFunc("/api/timeinfo", handleAPITimeinfo)
	r.HandleFunc("/api/item", handleAPIItem) // legacy
	r.HandleFunc("/api2/item", handleAPI2Item)
	r.HandleFunc("/api/rmitem", handleAPIRmItem) // legacy
	r.HandleFunc("/api/rmitemid", handleAPIRmItemID)
	r.HandleFunc("/api/items", handleAPI2Items)  // legacy
	r.HandleFunc("/api2/items", handleAPI2Items) // legacy
	r.HandleFunc("/api3/items", handleAPI3Items)
	r.HandleFunc("/api/searchname", handleAPISearchname)
	r.HandleFunc("/api/seqs", handleAPISeqs)
	r.HandleFunc("/api/allshots", handleAPIAllShots)
	r.HandleFunc("/api/shots", handleAPIShots) // legacy
	r.HandleFunc("/api2/shots", handleAPI2Shots)
	r.HandleFunc("/api/shot", handleAPIShot)
	r.HandleFunc("/api/asset", handleAPIAsset)
	r.HandleFunc("/api/assets", handleAPIAssets)
	r.HandleFunc("/api/setplatesize", handleAPISetPlateSize)
	r.HandleFunc("/api/setundistortionsize", handleAPISetUnDistortionSize)
	r.HandleFunc("/api/setrendersize", handleAPISetRenderSize) // legacy
	r.HandleFunc("/api2/setrendersize", handleAPI2SetRenderSize)
	r.HandleFunc("/api/setoverscanratio", handleAPISetOverscanRatio)
	r.HandleFunc("/api/setcamerapubpath", handleAPISetCameraPubPath)
	r.HandleFunc("/api/setcamerapubtask", handleAPISetCameraPubTask)
	r.HandleFunc("/api/setcameralensmm", handleAPISetCameraLensmm)
	r.HandleFunc("/api/setcameraprojection", handleAPISetCameraProjection)
	r.HandleFunc("/api/setseq", handleAPISetSeq)
	r.HandleFunc("/api/setseason", handleAPISetSeason)
	r.HandleFunc("/api/setepisode", handleAPISetEpisode)
	r.HandleFunc("/api/setplatepath", handleAPISetPlatePath)
	r.HandleFunc("/api/setthummov", handleAPISetThummov) // legacy
	r.HandleFunc("/api2/setthummov", handleAPI2SetThummov)
	r.HandleFunc("/api/setbeforemov", handleAPISetBeforemov)
	r.HandleFunc("/api/setaftermov", handleAPISetAftermov)
	r.HandleFunc("/api/seteditmov", handleAPISetEditmov)
	r.HandleFunc("/api/settaskstatus", handleAPISetTaskStatus) // legacy
	r.HandleFunc("/api2/settaskstatus", handleAPI2SetTaskStatus)
	r.HandleFunc("/api/taskstatusnum", handleAPITaskStatusNum)
	r.HandleFunc("/api/taskanduserstatusnum", handleAPITaskAndUserStatusNum)
	r.HandleFunc("/api/userstatusnum", handleAPIUserStatusNum)
	r.HandleFunc("/api/statusnum", handleAPIStatusNum)
	r.HandleFunc("/api/addtask", handleAPIAddTask)
	r.HandleFunc("/api/rmtask", handleAPIRmTask)
	r.HandleFunc("/api/settaskuser", handleAPISetTaskUser) // legacy
	r.HandleFunc("/api2/settaskuser", handleAPI2SetTaskUser)
	r.HandleFunc("/api/settaskusercomment", handleAPISetTaskUserComment)
	r.HandleFunc("/api/setplatein", handleAPISetPlateIn)
	r.HandleFunc("/api/setplateout", handleAPISetPlateOut)
	r.HandleFunc("/api/setjustin", handleAPISetJustIn)
	r.HandleFunc("/api/setjustout", handleAPISetJustOut)
	r.HandleFunc("/api/setscanin", handleAPISetScanIn)
	r.HandleFunc("/api/setscanout", handleAPISetScanOut)
	r.HandleFunc("/api/setscanframe", handleAPISetScanFrame)
	r.HandleFunc("/api/sethandlein", handleAPISetHandleIn)
	r.HandleFunc("/api/sethandleout", handleAPISetHandleOut)
	r.HandleFunc("/api/setshottype", handleAPISetShotType)
	r.HandleFunc("/api/setusetype", handleAPISetUseType)
	r.HandleFunc("/api/setassettype", handleAPISetAssetType)
	r.HandleFunc("/api/setoutputname", handleAPISetOutputName)
	r.HandleFunc("/api/setrnum", handleAPISetRnum) // legacy
	r.HandleFunc("/api2/setrnum", handleAPI2SetRnum)
	r.HandleFunc("/api/setdeadline2d", handleAPISetDeadline2D)
	r.HandleFunc("/api/setdeadline3d", handleAPISetDeadline3D)
	r.HandleFunc("/api/setscantimecodein", handleAPISetScanTimecodeIn)
	r.HandleFunc("/api/setscantimecodeout", handleAPISetScanTimecodeOut)
	r.HandleFunc("/api/setjusttimecodein", handleAPISetJustTimecodeIn)
	r.HandleFunc("/api/setjusttimecodeout", handleAPISetJustTimecodeOut)
	r.HandleFunc("/api/setfinver", handleAPISetFinver)
	r.HandleFunc("/api/setfindate", handleAPISetFindate)
	r.HandleFunc("/api/addtag", handleAPIAddTag)
	r.HandleFunc("/api/addassettag", handleAPIAddAssetTag)
	r.HandleFunc("/api/renametag", handleAPIRenameTag)
	r.HandleFunc("/api/rmtag", handleAPIRmTag)
	r.HandleFunc("/api/rmassettag", handleAPIRmAssetTag)
	r.HandleFunc("/api/settags", handleAPISetTags)
	r.HandleFunc("/api/setnote", handleAPISetNote)
	r.HandleFunc("/api/addcomment", handleAPIAddComment)
	r.HandleFunc("/api/editcomment", handleAPIEditComment)
	r.HandleFunc("/api/rmcomment", handleAPIRmComment)
	r.HandleFunc("/api/addsource", handleAPIAddSource)
	r.HandleFunc("/api/rmsource", handleAPIRmSource)
	r.HandleFunc("/api/addreference", handleAPIAddReference)
	r.HandleFunc("/api/rmreference", handleAPIRmReference)
	r.HandleFunc("/api/search", handleAPISearch)
	r.HandleFunc("/api/deadline2d", handleAPIDeadline2D)
	r.HandleFunc("/api/deadline3d", handleAPIDeadline3D)
	r.HandleFunc("/api/setstatus", handleAPISetTaskStatus)
	r.HandleFunc("/api/settaskmov", handleAPISetTaskMov) // legacy
	r.HandleFunc("/api2/settaskmov", handleAPI2SetTaskMov)
	r.HandleFunc("/api/settaskusernote", handleAPISetTaskUserNote)
	r.HandleFunc("/api/setretimeplate", handleAPISetRetimePlate)
	r.HandleFunc("/api/settasklevel", handleAPISetTaskLevel)
	r.HandleFunc("/api/setobjectid", handleAPISetObjectID)
	r.HandleFunc("/api/setociocc", handleAPISetOCIOcc)
	r.HandleFunc("/api/setrollmedia", handleAPISetRollmedia)
	r.HandleFunc("/api/setscanname", handleAPISetScanname)
	r.HandleFunc("/api/settaskdate", handleAPISetTaskDate)
	r.HandleFunc("/api/settaskexpectday", handleAPISetTaskExpectDay)
	r.HandleFunc("/api/settaskresultday", handleAPISetTaskResultDay)
	r.HandleFunc("/api/settaskpredate", handleAPISetTaskPredate)
	r.HandleFunc("/api/settaskstartdate", handleAPISetTaskStartdate)
	r.HandleFunc("/api/task", handleAPITask)
	r.HandleFunc("/api/shottype", handleAPIShottype)
	r.HandleFunc("/api/setcrowdasset", handleAPISetCrowdAsset)
	r.HandleFunc("/api/mailinfo", handleAPIMailInfo)
	r.HandleFunc("/api/usetypes", handleAPIUseTypes)
	r.HandleFunc("/api/publish", handleAPIAddTaskPublish) // legacy
	r.HandleFunc("/api/addpublish", handleAPIAddTaskPublish)
	r.HandleFunc("/api/setpublishstatus", handleAPISetTaskPublishStatus)
	r.HandleFunc("/api/rmpublish", handleAPIRmTaskPublish)
	r.HandleFunc("/api/rmpublishkey", handleAPIRmTaskPublishKey)
	r.HandleFunc("/api/uploadthumbnail", handleAPIUploadThumbnail)
	r.HandleFunc("/api/setnetflixid", handleAPISetNetflixID)

	// restAPI USER
	r.HandleFunc("/api/user", handleAPIUser) // legacy
	r.HandleFunc("/api2/user", handleAPI2User)
	r.HandleFunc("/api/users", handleAPISearchUser)
	r.HandleFunc("/api/validuser", handleAPIValidUser) // 보안취약점 이슈가 있다. 다른 툴과 쉽게 연동할 때 편리하다. 보안레벨을 높게 올릴때는 허용하지 않도록 한다.
	r.HandleFunc("/api/setleaveuser", handleAPISetLeaveUser)
	r.HandleFunc("/api/autocompliteusers", handleAPIAutoCompliteUsers)
	r.HandleFunc("/api/initpassword", handleAPIInitPassword)
	r.HandleFunc("/api/ansiblehosts", handleAPIAnsibleHosts)

	// restAPI Organization
	r.HandleFunc("/api/teams", handleAPIAllTeams)

	// restAPI Tasksetting
	r.HandleFunc("/api/tasksetting", handleAPITasksetting)
	r.HandleFunc("/api/shottasksetting", handleAPIShotTasksetting)
	r.HandleFunc("/api/assettasksetting", handleAPIAssetTasksetting)
	r.HandleFunc("/api/categorytasksettings", handleAPICategoryTasksettings)

	// restAPI Status
	r.HandleFunc("/api/status", handleAPIStatus)
	r.HandleFunc("/api/addstatus", handleAPIAddStatus)

	// restAPI PublishKey
	r.HandleFunc("/api/publishkeys", handleAPIPublishKeys)
	r.HandleFunc("/api/getpublish", handleAPIGetPublish)

	// restAPI Review
	r.HandleFunc("/api/addreview", handleAPIAddReviewStageMode) // legacy
	r.HandleFunc("/api/addreviewstagemode", handleAPIAddReviewStageMode)
	r.HandleFunc("/api/addreviewstatusmode", handleAPIAddReviewStatusMode)
	r.HandleFunc("/api/review", handleAPIReview)
	r.HandleFunc("/api/searchreview", handleAPISearchReview)
	r.HandleFunc("/api/setreviewstatus", handleAPISetReviewStatus)
	r.HandleFunc("/api/setreviewitemstatus", handleAPISetReviewItemStatus)
	r.HandleFunc("/api/setreviewstage", handleAPISetReviewStage)
	r.HandleFunc("/api/setreviewnextstatus", handleAPISetReviewNextStatus)
	r.HandleFunc("/api/setreviewnextstage", handleAPISetReviewNextStage)
	r.HandleFunc("/api/addreviewcomment", handleAPIAddReviewComment)
	r.HandleFunc("/api/addreviewstatusmodecomment", handleAPIAddReviewStatusModeComment)
	r.HandleFunc("/api/editreviewcomment", handleAPIEditReviewComment)
	r.HandleFunc("/api/rmreviewcomment", handleAPIRmReviewComment)
	r.HandleFunc("/api/rmreview", handleAPIRmReview)
	r.HandleFunc("/api/setreviewproject", handleAPISetReviewProject)
	r.HandleFunc("/api/setreviewtask", handleAPISetReviewTask)
	r.HandleFunc("/api/setreviewname", handleAPISetReviewName)
	r.HandleFunc("/api/setreviewpath", handleAPISetReviewPath)
	r.HandleFunc("/api/setreviewcreatetime", handleAPISetReviewCreatetime)
	r.HandleFunc("/api/setreviewmainversion", handleAPISetReviewMainVersion)
	r.HandleFunc("/api/setreviewsubversion", handleAPISetReviewSubVersion)
	r.HandleFunc("/api/setreviewfps", handleAPISetReviewFps)
	r.HandleFunc("/api/setreviewdescription", handleAPISetReviewDescription)
	r.HandleFunc("/api/setreviewcamerainfo", handleAPISetReviewCameraInfo)
	r.HandleFunc("/api/setreviewprocessstatus", handleAPISetReviewProcessStatus)
	r.HandleFunc("/api/uploadreviewdrawing", handleAPIUploadReviewDrawing)
	r.HandleFunc("/api/rmreviewdrawing", handleAPIRmReviewDrawing)
	r.HandleFunc("/api/reviewdrawingframe", handleAPIReviewDrawingFrame)
	r.HandleFunc("/api/setreviewagainforwaitstatustoday", handleAPISetReviewAgainForWaitStatusToday)
	r.HandleFunc("/api/reviewoutputdatapath", handleAPIReviewOutputDataPath)

	// REST API Partner
	r.HandleFunc("/api/partner", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/partner", postPartnerHandler).Methods("POST")
	r.HandleFunc("/api/partner/{id}", getPartnerHandler).Methods("GET")
	r.HandleFunc("/api/partner/{id}", putPartnerHandler).Methods("PUT")
	r.HandleFunc("/api/partner/{id}", deletePartnerHandler).Methods("DELETE")

	// Deprecated: 사용하지 않는 url, 과거호환성을 위해서 남겨둠
	r.HandleFunc("/edititem", handleEditItem)                    // legacy
	r.HandleFunc("/editeditem", handleEditedItem)                // legacy
	r.HandleFunc("/api/setmov", handleAPISetTaskMov)             // legacy
	r.HandleFunc("/api/setstartdate", handleAPISetTaskStartdate) // legacy
	r.HandleFunc("/edititem-submit", handleEditItemSubmitv2)     // legacy
	r.Use(mux.CORSMethodMiddleware(r))
	http.Handle("/", r)

	if port == ":443" || port == ":8443" { // https ports
		err := http.ListenAndServeTLS(port, *flagCertFullchanin, *flagCertPrivkey, r)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := http.ListenAndServe(port, r)
		if err != nil {
			log.Fatal(err)
		}
	}
}
