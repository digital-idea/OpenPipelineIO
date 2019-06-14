package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func itemStatus2color(num string) string {
	switch num {
	case OMIT:
		return "#FC8F55"
	case CONFIRM:
		return "#54D6FD"
	case WIP:
		return "#77BB40"
	case READY:
		return "#BEEF37"
	case ASSIGN:
		return "#FFF76B"
	case OUT:
		return "#EEA4F1"
	case DONE:
		return "#F0F1F0"
	case HOLD:
		return "#989898"
	case NONE:
		return "#CCCDCC"
	default:
		return "#CCCDCC"
	}
}

func projectStatus2color(status ProjectStatus) string {
	switch status {
	case PreProjectStatus:
		return "#E3EB8D"
	case PostProjectStatus:
		return "#9ABF62"
	case LayoverProjectStatus:
		return "#B2B084"
	case BackupProjectStatus:
		return "#7EAECB"
	case ArchiveProjectStatus:
		return "#A3A3A3"
	case LawsuitProjectStatus:
		return "#CC9372"
	default:
		return "#EBEBEB"
	}
}

func statusnum2string(num string) string {
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

// GetPath 함수는 "컨셉 : /show/test.jpg" 문자열을 "/show/test.jpg" 형태로 바꾸어주는 함수.
func GetPath(s string) string {
	for _, i := range strings.Split(s, " ") {
		if strings.HasPrefix(i, "/show") {
			return i
		}
	}
	return s
}

// RemovePath 함수는 "컨셉:/show/test.jpg"문자열을"컨셉:"형태로 바꾸어 준다.
func RemovePath(s string) string {
	num := strings.Index(s, "/show")
	if num == -1 {
		return s
	}
	return s[0:num]
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
	if MatchFullTime.MatchString(t) {
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
	if MatchFullTime.MatchString(t) {
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

// Tags2str 템플릿함수는 태그리스트를 ,로 붙힌 문자열로 만든다.
func Tags2str(tags []string) string {
	var newtags []string
	for _, tag := range tags {
		if tag != "" {
			newtags = append(newtags, tag)
		}
	}
	return strings.Join(newtags, ",")
}

// Str2Tags 템플릿함수는 태그 문자열을 리스트로 만든다.
func Str2Tags(tags string) []string {
	var newtags []string
	cleanup := strings.Replace(tags, " ", "", -1)
	for _, tag := range strings.Split(cleanup, ",") {
		if tag != "" {
			newtags = append(newtags, tag)
		}
	}
	return newtags
}

// CheckDate 함수는 1차마감, 2차마감, mov업데이이트날짜와 검색어를 비교하여 statusbox의 색깔을 반영
// 검색어에 날짜형식의 검색어만 적용된다.(2017-12-05, 1205)
func CheckDate(predate, date, mdate, searchs string) string {
	for _, search := range strings.Split(searchs, " ") {
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
	// 예) "2일전", "1시간30분전", "30분25초전", "15초전"
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
		return fmt.Sprintf("%d일전", int(due.Hours()/24.0))
	}
	var kortime string
	section := 0 // 마디가 몇개인지 체크한다.
	for _, r := range due.String() {
		if section == 2 {
			break
		}
		switch r {
		case 'h':
			kortime += "시간"
			section++
		case 'm':
			kortime += "분"
			section++
		case '.':
			kortime += "초"
			section = 2 // 만약 초만 있다면 마디가 1개여도 빠져나간다.
			break
		default:
			kortime += string(r)
		}
	}
	kortime += "전"
	return kortime
}

// CheckDdline 템플릿함수는 해당 시간이 이번주에 해당하는지 다음주에 해당하는지 판단한다.
func CheckDdline(t string) string {
	// 1124을 시간형태로 바꾸고 7일안에 포함하는 형태인지 체크한다.
	// 현재시간에서 체크할 시간을 뺀 수
	if !MatchShortTime.MatchString(t) {
		return ""
	}
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

// Review 템플릿 함수는 리뷰툴 URL을 리턴한다. ily는 현재 사용하지 않는다.
func Review(project, name, typ, task string) string {
	reviewURL := fmt.Sprintf("http://10.0.90.252:8000/filter?project=%s&name=%s", project, name)
	switch typ {
	case "shot", "org", "left":
		reviewURL += "&type=shot"
	case "asset":
		reviewURL += "&type=asset"
	}
	if task != "" {
		reviewURL += "&task=" + task
	}
	return reviewURL
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

// Hashtag2tag 템플릿함수는 태그명을 `#태그명`으로 변경한다.
func Hashtag2tag(tag string) string {
	var hastag string
	if tag != "" {
		hastag = `#` + tag
	}
	return hastag
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
