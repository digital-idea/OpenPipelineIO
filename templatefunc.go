package main

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// projectStatus2color 템플릿 함수는 프로젝트 상태를 받아서 Bootstrap4 컬러로 변환합니다.
func projectStatus2color(status ProjectStatus) string {
	switch status {
	case PreProjectStatus:
		return "text-white bg-info" // 청록
	case PostProjectStatus:
		return "text-white bg-success" // 녹색
	case LayoverProjectStatus:
		return "text-dark bg-warning" // 노랑
	case BackupProjectStatus:
		return "text-white bg-danger" // 빨강
	case ArchiveProjectStatus:
		return "text-dark bg-secondary" // 회색
	case LawsuitProjectStatus:
		return "text-dark bg-waring" // 노랑
	default:
		return "text-white bg-darkmode"
	}
}

// Status2capString 템플릿함수는 status 값을 받아서 대문자를 반환한다.
func Status2capString(num string) string {
	switch num {
	case OMIT:
		return "OMIT"
	case CONFIRM:
		return "CONFIRM"
	case WIP:
		return "WIP"
	case READY:
		return "READY"
	case ASSIGN:
		return "ASSIGN"
	case OUT:
		return "OUT"
	case DONE:
		return "DONE"
	case HOLD:
		return "HOLD"
	case NONE:
		return "NONE"
	default:
		return "NONE"
	}
}

// Status2string 템플릿함수는 status 값을 받아서 소문자를 반환한다.
func Status2string(status string) string {
	switch status {
	case OMIT, "omit":
		return "omit"
	case CONFIRM, "confirm":
		return "confirm"
	case WIP, "wip":
		return "wip"
	case READY, "ready":
		return "ready"
	case ASSIGN, "assign":
		return "assign"
	case OUT, "out":
		return "out"
	case DONE, "done":
		return "done"
	case HOLD, "hold":
		return "hold"
	case NONE, "none":
		return "none"
	default:
		return "none"
	}
}

// StatusString2string 템플릿함수는 status 문자를 받아서 Status 값을 반환한다.
func StatusString2string(status string) string {
	switch status {
	case "omit":
		return OMIT
	case "confirm":
		return CONFIRM
	case "wip":
		return WIP
	case "ready":
		return READY
	case "assign":
		return ASSIGN
	case "out":
		return OUT
	case "done":
		return DONE
	case "hold":
		return HOLD
	case "none":
		return NONE
	default:
		return NONE
	}
}

// name2seq 함수는 "SS_0010"을 받아서 "SS" 이름을 반환한다.
func name2seq(name string) string {
	if strings.Contains(name, "_") {
		return strings.Split(name, "_")[0]
	}
	return name
}

func note2body(note string) string {
	num := strings.LastIndex(note, ";")
	return note[num+1:]
}

func pmnote2body(note string) string {
	var date string
	date = ""
	if strings.Contains(note, ";") {
		date = strings.Split(note, ";")[0]
	}
	num := strings.LastIndex(note, ";")
	return ToShortTime(date) + "-" + note[num+1:]
}

// GetPath 함수는 "컨셉: /show/test.jpg" 문자열을 "/show/test.jpg" 형태로 바꾸어준다. /clib 도 마찬가지로 지원한다.
func GetPath(s string) string {
	for _, i := range strings.Split(s, " ") {
		if strings.HasPrefix(i, "/show") {
			return i
		}
		if strings.HasPrefix(i, "/clib") {
			return i
		}
	}
	return s
}

// RemovePath 함수는 "컨셉:/show/test.jpg" 문자열을 "컨셉:" 형태로 바꾸어 준다. /clib 형태의 클립라이브러리도 바꾸어준다.
func RemovePath(s string) string {
	if strings.Contains(s, "/show") {
		num := strings.Index(s, "/show")
		if num == -1 {
			return s
		}
		return s[0:num]
	}
	if strings.Contains(s, "/clib") {
		num := strings.Index(s, "/clib")
		if num == -1 {
			return s
		}
		return s[0:num]
	}
	return s
}

// ReverseStringSlice 함수는 받아들인 string 슬라이스의 아이템 순서를 역순으로 변경한다.
// 받아들인 값이 nil이라면 빈 string 슬라이스를 반환하는데,
// 그 이유는 template 파싱중 패닉이 일어나지 않게 하기 위해서다.
func ReverseStringSlice(lists []string) []string {
	if lists == nil {
		return []string{}
	}
	var result []string
	for i := len(lists); i > 0; i-- {
		result = append(result, lists[i-1])
	}
	return result
}

// ReverseIntSlice 함수는 받아들인 int 슬라이스의 아이템 순서를 역순으로 변경한다.
// 받아들인 값이 nil이라면 빈 int 슬라이스를 반환한다.
func ReverseIntSlice(lists []int) []int {
	if lists == nil {
		return []int{}
	}
	var result []int
	for i := len(lists); i > 0; i-- {
		result = append(result, lists[i-1])
	}
	return result
}

// GenPageNums 함수는 숫자를 받아서 그 숫자의 크기만큼의 int 슬라이스(Pages)를 만든다.
func GenPageNums(size int) []int {
	result := make([]int, size)
	for i := 0; i < size; i++ {
		result[i] = i + 1
	}
	return result
}

// Publishes 자료구조를 정리하기 위해 사용하는 자료구조
type Publishes []Publish

func (a Publishes) Len() int           { return len(a) }
func (a Publishes) Less(i, j int) bool { return a[i].Createtime > a[j].Createtime }
func (a Publishes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// SortByCreatetimeForPublishes 함수는 Publish 슬라이스의 아이템 순서를 역순으로 변경한다.
func SortByCreatetimeForPublishes(lists []Publish) []Publish {
	sort.Sort(Publishes(lists))
	return lists
}

// ReverseCommentSlice 함수는 받아들인 string 슬라이스의 아이템 순서를 역순으로 변경한다.
// 받아들인 값이 nil이라면 빈 string 슬라이스를 반환하는데,
// 그 이유는 template 파싱중 패닉이 일어나지 않게 하기 위해서다.
func ReverseCommentSlice(lists []Comment) []Comment {
	if lists == nil {
		return []Comment{}
	}
	var result []Comment
	for i := len(lists); i > 0; i-- {
		result = append(result, lists[i-1])
	}
	return result
}

// CutCommentSlice 템플릿 함수는 3개의 리스트만 반환한다.
func CutCommentSlice(lists []Comment) []Comment {
	if len(lists) < 4 {
		return lists
	}
	return lists[:3]
}

// CutStringSlice 템플릿 함수는 3개의 리스트만 반환한다.
func CutStringSlice(lists []string) []string {
	if len(lists) < 4 {
		return lists
	}
	return lists[:3]
}

// MatchShortTime 은 "1019" 형식의 레귤러 익스프레션
var MatchShortTime = regexp.MustCompile(`^\d{4}$`)

// MatchNormalTime 은 "2016-10-19" 형식의 레귤러 익스프레션
var MatchNormalTime = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// MatchFullTime 은 "2016-10-19T16:41:24+09:00" 형식의 레귤러 익스프레션
var MatchFullTime = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[-+]\d{2}:\d{2}$`)

// MatchUTCFullTime 은 "2016-10-19T16:41:24Z" 형식의 레귤러 익스프레션
var MatchUTCFullTime = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

// ToShortTime 템플릿함수는 시간문자열을 일반사용자의 가독성을 위해서 4자리 문자열로 바꾸어주는 함수.
//1019, 2016-10-19, 2016-10-19T16:41:24+09:00 을 전부 1019로 바꾼다.
func ToShortTime(t string) string {
	if MatchNormalTime.MatchString(t) {
		return strings.Replace(t[5:10], "-", "", -1)
	}
	if MatchFullTime.MatchString(t) || MatchUTCFullTime.MatchString(t) {
		zonetime, err := time.Parse(time.RFC3339, t) // 현지시간으로 바꾼다.
		if err != nil {
			return t
		}
		// 현지시간을 줄여서 출력해준다.
		return strings.Replace(zonetime.Format(time.RFC3339)[5:10], "-", "", -1)
	}
	return t
}

// ToNormalTime 시간문자열을 Normal 타입으로 바꾸어주는 함수.
// 1019, 2016-10-19, 2016-10-19T16:41:24+09:00 을 전부 2016-10-19로 바꾼다.
func ToNormalTime(t string) string {
	if MatchNormalTime.MatchString(t) {
		return t
	}
	if MatchFullTime.MatchString(t) || MatchUTCFullTime.MatchString(t) {
		return t[0:10]
	}
	if MatchShortTime.MatchString(t) {
		return fmt.Sprintf("%d-%s-%s", time.Now().Year(), t[0:2], t[2:])
	}
	return t
}

// ToFullTime 템플릿함수는 시간문자열을 RFC3339 시간포멧으로 바꾸는 함수이다.
// 시간은 퇴근시간으로 맞추어져 있다.
func ToFullTime(t string) string {
	if MatchShortTime.MatchString(t) {
		return fmt.Sprintf("%d-%s-%sT19:00:00%s", time.Now().Year(), t[0:2], t[2:4], time.Now().Format(time.RFC3339)[19:])
	}
	if MatchNormalTime.MatchString(t) {
		return fmt.Sprintf("%sT19:00:00%s", t, time.Now().Format(time.RFC3339)[19:])
	}
	return t
}

// List2str 템플릿함수는 리스트를 ,로 붙힌 문자열로 만든다.
func List2str(items []string) string {
	var lists []string
	for _, i := range items {
		if i != "" {
			lists = append(lists, i)
		}
	}
	return strings.Join(lists, ",")
}

// Str2List 템플릿함수는 태그 문자열을 리스트로 만든다.
func Str2List(str string) []string {
	if str == "" {
		return []string{}
	}
	var lists []string
	cleanup := strings.Replace(str, " ", "", -1)
	for _, tag := range strings.Split(cleanup, ",") {
		if tag != "" {
			lists = append(lists, tag)
		}
	}
	return lists
}

// CheckDate 함수는 1차마감, 2차마감, mov업데이이트날짜와 검색어를 비교하여 statusbox의 색깔을 반영
// 검색어에 날짜형식의 검색어만 적용된다.(2017-12-05, 1205)
func CheckDate(predate, date, mdate, searches string) string {
	for _, search := range strings.Split(searches, " ") {
		if MatchNormalTime.MatchString(search) && ToNormalTime(mdate) == search {
			return "_daily" // 파랑색
		}
		if MatchFullTime.MatchString(predate) && MatchShortTime.MatchString(search) {
			head := strings.Split(predate, "T")[0]
			if strings.Replace(head[5:], "-", "", -1) == search {
				return "_deadline" // 빨강색
			}
		}
		if MatchFullTime.MatchString(date) && MatchShortTime.MatchString(search) {
			head := strings.Split(date, "T")[0]
			if strings.Replace(head[5:], "-", "", -1) == search {
				return "_deadline"
			}
		}
	}
	return ""
}

// CheckUpdate 함수는 시간을 입력받아서 그 시간이 72시간이 지났는지 판단하는 템플릿함수이다.
func CheckUpdate(t string) bool {
	if !MatchFullTime.MatchString(t) && !MatchUTCFullTime.MatchString(t) {
		return false
	}
	utime := str2time(t)
	if 72 < time.Now().Sub(utime).Hours() {
		return false
	}
	return true
}

// ToHumantime 템플릿함수는 RFC3339 포맷의 시간문자열을 받아서 현재시간과의 차이를 읽기 쉽게 출력해주는 함수이다.
func ToHumantime(timestring string) string {
	// 2017-02-20T18:15:43+09:00 형태의 시간을 time.Duration으로 바꾸고,
	// (-)48h26m17.016071813s 형태의 time.Duration 형식을 loop 하면서 처리한다.
	// 읽기 쉽게 4개의 패턴으로만 나온다.
	// 예) "2d", "1h30m", "30m분25s", "15s"
	// 72시간이 지나면 아무 문자열도 나오지 않는다.
	// time.Duration이 음수라면 아무 문자열도 나오지 않는다.
	// Template function이라서 실패해도 아무 문자열도 나오지 않는다.
	timevalue, err := time.Parse(time.RFC3339, timestring)
	if err != nil {
		return ""
	}
	due := time.Now().Sub(timevalue)
	if due.Hours() > 72.0 {
		return ""
	}
	if due.Hours() < 0 {
		return ""
	}
	if due.Hours() > 24.0 {
		return fmt.Sprintf("%dd", int(due.Hours()/24.0))
	}
	var timeStr string
	section := 0 // 마디가 몇개인지 체크한다.
	for _, r := range due.String() {
		if section == 2 {
			break
		}
		switch r {
		case 'h':
			timeStr += "h"
			section++
		case 'm':
			timeStr += "m"
			section++
		case '.':
			timeStr += "s"
			section = 2 // 만약 초만 있다면 마디가 1개여도 빠져나간다.
			break
		default:
			timeStr += string(r)
		}
	}
	return timeStr
}

// CheckDdline 템플릿함수는 해당 시간이 이번주에 해당하는지 다음주에 해당하는지 판단한다.
func CheckDdline(t string) string {
	// 검색 태그가 있을 수 있다. 제거한다.
	if strings.HasPrefix(t, "ddline2d:") {
		t = strings.TrimLeft(t, "ddline2d:")
	}
	if strings.HasPrefix(t, "ddline3d:") {
		t = strings.TrimLeft(t, "ddline3d:")
	}
	if !MatchNormalTime.MatchString(t) {
		return ""
	}
	// 1124을 시간형태로 바꾸고 7일안에 포함하는 형태인지 체크한다.
	// 현재시간에서 체크할 시간을 뺀 수
	offset := str2time(ToFullTime(t)).Sub(time.Now()).Hours()
	thisline := time.Hour * 24 * 7
	nextline := time.Hour * 24 * 14
	switch {
	// 시간이 지나면 빈 문자열을 출력한다.
	case offset < 0:
		return ""
	// 7일 이전의 날짜라면 _this 문자를 반환한다.
	case offset <= thisline.Hours():
		return "_this"
	// 7일보다 크고, 14일 이전의 날짜라면 _next 문자를 반환한다.
	case thisline.Hours() < offset && offset < nextline.Hours():
		return "_next"
	// 14일보다 크면 빈 문자열을 출력한다.
	default:
		return ""
	}
}

// CheckDdlinev2 템플릿함수는 해당 시간이 이번주에 해당하는지 다음주에 해당하는지 판단한다.
func CheckDdlinev2(t string) string {
	// 검색 태그가 있을 수 있다. 제거한다.
	if strings.HasPrefix(t, "ddline2d:") {
		t = strings.TrimLeft(t, "ddline2d:")
	}
	if strings.HasPrefix(t, "ddline3d:") {
		t = strings.TrimLeft(t, "ddline3d:")
	}
	if !(MatchNormalTime.MatchString(t) || MatchShortTime.MatchString(t) || MatchFullTime.MatchString(t)) {
		return "darkmode"
	}
	// 1124을 시간형태로 바꾸고 7일안에 포함하는 형태인지 체크한다.
	// 현재시간에서 체크할 시간을 뺀 수
	offset := str2time(ToFullTime(t)).Sub(time.Now()).Hours()
	thisline := time.Hour * 24 * 7
	nextline := time.Hour * 24 * 14
	switch {
	// 지난시간
	case offset < 0:
		return "fade"
	// 이번주: 7일 이전의 날짜
	case offset <= thisline.Hours():
		return "danger"
	// 다음주: 7일보다 크고, 14일 이전의 날짜
	case thisline.Hours() < offset && offset < nextline.Hours():
		return "warning"
	// 일반모드: 14일보다 클때
	default:
		return "darkmode"
	}
}

// Framecal 템플릿함수는 in, out 프레임을 받아서 총 프레임수를 문자로 반환한다.
func Framecal(in int, out int) string {
	// DB의 초기값은 0이다.
	// 기본적으로 회사에서 사용하는 시작 프레임은 1001이다.
	if in <= 0 || out <= 0 {
		return ""
	}
	if in > out {
		return ""
	}
	return fmt.Sprintf("%d", out-in+1)
}

// Add 함수는 template안에서 두 수를 더한다.
func Add(a, b int) int {
	return a + b
}

// Minus 함수는 template안에서 두 수를 뺀다.
func Minus(a, b int) int {
	return a - b
}

// Protocol 템플릿 함수는 프로토콜 문자열을 반환한다.
func Protocol(s string) string {
	if strings.HasPrefix(s, "http://") {
		return "http"
	}
	if strings.HasPrefix(s, "https://") {
		return "https"
	}
	if strings.HasPrefix(s, "ftp://") {
		return "ftp"
	}
	if strings.HasPrefix(s, "sftp://") {
		return "sftp"
	}
	if strings.HasPrefix(s, "mailto://") {
		return "mailto"
	}
	if strings.HasPrefix(s, "git://") {
		return "git"
	}
	if strings.HasPrefix(s, "slack://") {
		return "slack"
	}
	return "dilink"
}

// RmProtocol 템플릿 함수는 프로토콜 문자열을 반환한다.
func RmProtocol(s string) string {
	if strings.HasPrefix(s, "http://") {
		return strings.TrimLeft(s, "http://")
	}
	if strings.HasPrefix(s, "https://") {
		return strings.TrimLeft(s, "https://")
	}
	if strings.HasPrefix(s, "ftp://") {
		return strings.TrimLeft(s, "ftp://")
	}
	if strings.HasPrefix(s, "sftp://") {
		return strings.TrimLeft(s, "sftp://")
	}
	if strings.HasPrefix(s, "mailto://") {
		return strings.TrimLeft(s, "mailto://")
	}
	if strings.HasPrefix(s, "dilink://") {
		return strings.TrimLeft(s, "dilink://")
	}
	if strings.HasPrefix(s, "git://") {
		return strings.TrimLeft(s, "git://")
	}
	if strings.HasPrefix(s, "slack://") {
		return strings.TrimLeft(s, "slack://")
	}
	return s
}

// ProtocolTarget 템플릿 함수는 프로토콜 문자열을 반환한다.
func ProtocolTarget(s string) string {
	if strings.HasPrefix(s, "http://") {
		return "_blank"
	}
	if strings.HasPrefix(s, "https://") {
		return "_blank"
	}
	return "_self"
}

// Scanname2RollMedia 템플릿함수는 스캔이름으로 RollMedia를 반환한다.
// 예) 22_A039C002_150916_R529 문자를 받아 A039C002를 반환한다.
// 인수가 Scanname 형태가 아니면 "" 문자를 반환한다.
func Scanname2RollMedia(scanname string) string {
	if !regexpRollMedia.MatchString(scanname) {
		return ""
	}
	return strings.Split(scanname, "_")[1]
}

// AddTagColon 템플릿함수는 태그문자에 "tag:태그"를 붙혀 반환한다.
func AddTagColon(tag string) string {
	return "tag:" + tag
}

// Username2Elements 함수는 element 이름들을 반환한다.
// 김한웅(fire),박지섭(smoke) 형태의 문자라면 fire,smoke를 반환한다.
// 이 함수는 유연한 task 구조가 되기 전까지만 활용한다.
func Username2Elements(str string) string {
	var elements []string
	re := regexp.MustCompile(`\(([0-9a-z]+)\)`)
	results := re.FindAllStringSubmatch(str, -1)
	for _, e := range results {
		elements = append(elements, e[1])
	}
	return strings.Join(elements, ",")
}

// ShortPhoneNum 함수는 내선번호 길다면 마지막 4자리만 출력한다.
func ShortPhoneNum(str string) string {
	if len(str) < 5 {
		return str
	}
	return str[len(str)-4:]
}

// TaskStatus 템플릿 함수는 아이템과 Task 문자를 받아서 상태를 반환한다.
func TaskStatus(i Item, task string) string {
	return reflect.ValueOf(i).FieldByName(strings.Title(task)).FieldByName("Status").String()
}

// TaskUser 템플릿 함수는 아이템과 Task 문자를 받아서 User를 반환한다.
func TaskUser(i Item, task string) string {
	return reflect.ValueOf(i).FieldByName(strings.Title(task)).FieldByName("User").String()
}

// TaskDate 템플릿 함수는 아이템과 Task 문자를 받아서 Date를 반환한다.
func TaskDate(i Item, task string) string {
	return reflect.ValueOf(i).FieldByName(strings.Title(task)).FieldByName("Date").String()
}

// TaskPredate 템플릿 함수는 아이템과 Task 문자를 받아서 Predate를 반환한다.
func TaskPredate(i Item, task string) string {
	return reflect.ValueOf(i).FieldByName(strings.Title(task)).FieldByName("Predate").String()
}

// GetTaskLevel 템플릿 함수는 아이템과 Task 문자를 받아서 Tasklevel을 반환한다.
func GetTaskLevel(i Item, task string) TaskLevel {
	return TaskLevel(reflect.ValueOf(i).FieldByName(strings.Title(task)).FieldByName("TaskLevel").Int())
}

// userInfo 템플릿 함수는 id(name,team) 문자열을 받아서 name,team만 반환한다.
func userInfo(userInfo string) string {
	if regexpUserInfo.MatchString(userInfo) {
		re := regexp.MustCompile(`\(.+\)`)
		f := re.FindString(userInfo)
		return f[1 : len(f)-1]
	}
	return userInfo
}

func onlyID(userInfo string) string {
	if regexpUserInfo.MatchString(userInfo) {
		return strings.Split(userInfo, "(")[0]
	}
	return userInfo
}

func mapToSlice(tasks map[string]Task, ordermap map[string]float64) []Task {
	var results []Task
	for _, value := range tasks {
		results = append(results, value)
	}
	// results를 BubbleSort 한다.
	for i := len(results); i > 0; i-- {
		for j := 1; j < i; j++ {
			// ordermap 에 키가 없다면 정렬하지 않는다.
			if _, found := ordermap[results[j-1].Title]; !found {
				continue
			}
			if _, found := ordermap[results[j].Title]; !found {
				continue
			}
			// float64 값을 비교하여 정렬
			if ordermap[results[j-1].Title] > ordermap[results[j].Title] {
				intermediate := results[j]
				results[j] = results[j-1]
				results[j-1] = intermediate
			}
		}
	}
	return results
}

// hasStatus는 statuslist에 status가 존재하는지 체크한다. 존재하면 true를 반환한다.
func hasStatus(statuslist []string, status string) bool {
	for _, s := range statuslist {
		if s == status {
			return true
		}
	}
	return false
}

// AddProductionStartFrame 템플릿함수는 프레임에 프로덕션 시작 프레임을 더한다.
func AddProductionStartFrame(frame int) int {
	return CachedAdminSetting.ProductionStartFrame + frame - 1
}

// ProductionVersionFormat 템플릿 함수는 숫자를 받아서 프로덕션 버전 자리수 만큼 "0"문자를 붙혀서 문자열로 반환한다.
func ProductionVersionFormat(version int) string {
	n := strconv.Itoa(version)
	for len(n) < CachedAdminSetting.ProductionPaddingVersionNumber {
		n = "0" + n
	}
	return n
}
