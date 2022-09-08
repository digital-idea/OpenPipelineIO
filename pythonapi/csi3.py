#coding:utf-8
import json
import urllib2
import urllib

# python3.x부터는 import.requests 로 POST를 구현하면 코드가 더욱 짧아진다.

def Project(name):
	"""
	프로젝트 정보를 가지고 오는 함수.
	(딕셔너리, err) 값을 반환한다.
	"""
	restURL = "https://openpipeline.io/api/project?id=%s" %  name
	try:
		data = json.load(urllib2.urlopen(restURL))
	except:
		return {}, "RestAPI에 연결할 수 없습니다."
	if "error" in data:
		return {}, data["error"]
	return data, None

def Shot(project, name):
	"""
	샷 정보를 가지고 오는 함수.
	(딕셔너리, err)값을 반환한다.
	"""
	restURL = "https://openpipeline.io/api/shot?project=%s&name=%s" % (project, name)
	try:
		data = json.load(urllib2.urlopen(restURL))
	except:
		return {}, "RestAPI에 연결할 수 없습니다."
	if "error" in data:
		return {}, data["error"]
	return data["data"], None

def Setmov(project, name, task, mov):
	"""
	mov를 등록한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&task=%s&mov=%s" % (project, name, task, mov)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setmov", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetPlateSize(project, name, size):
	"""
	샷 플레이트 사이즈를 셋팅한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&size=%s" % (project, name, size)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setplatesize", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetDistortionSize(project, name, size):
	"""
	샷 디스토션사이즈를 셋팅한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&size=%s" % (project, name, size)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setdistortionsize", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetRenderSize(project, name, size):
	"""
	샷 렌더사이즈를 셋팅한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&size=%s" % (project, name, size)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setrendersize", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetJustIn(project, name, frame):
	"""
	JustIn 값을 셋팅한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&frame=%s" % (project, name, frame)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setjustin", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetJustOut(project, name, frame):
	"""
	JustOut 값을 셋팅한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&frame=%s" % (project, name, frame)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setjustout", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetCameraPubPath(project, name, path):
	"""
	샷 카메라 퍼블리쉬경로를 설정한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&path=%s" % (project, name, path)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setcamerapubpath", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetCameraPubTask(project, name, task):
	"""
	샷 카메라 퍼블리쉬 Task를 설정한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&task=%s" % (project, name, task)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setcamerapubtask", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetCameraProjection(project, name, status):
	"""
	샷 카메라 프로젝션 유무를 설정한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	if status:
		status = "true"
	data = "project=%s&name=%s&projection=%s" % (project, name, status)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setcameraprojection", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetThummov(project, name, movpath):
	"""
	아이템의 썸네일 mov를 설정한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&path=%s" % (project, name, movpath)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setthummov", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def Projects(status=None):
	"""
	프로젝트 리스트 정보를 가지고 오는 함수.
	(리스트, err) 값을 반환한다.
	"""
	restURL = "https://openpipeline.io/api/projects"
	if status:
		restURL += "?status=%s" % status
	try:
		projects = json.load(urllib2.urlopen(restURL))
	except:
		return [], "RestAPI에 연결할 수 없습니다."
	if projects["error"]:
		return [], projects["error"]
	return projects["data"], None

def Seqs(project):
	"""
	프로젝트의 시퀀스 리스트 정보를 가지고 오는 함수.
	(리스트, err) 값을 반환한다.
	"""
	if not project:
		return [], "프로젝트 정보가 없습니다."
	restURL = "https://openpipeline.io/api/seqs?project={0}".format(project)
	try:
		seqs = json.load(urllib2.urlopen(restURL))
	except:
		return [], "RestAPI에 연결할 수 없습니다."
	if "error" in seqs:
		return [], seqs["error"]
	return seqs["data"], None

def Shots(project, seq):
	"""
	프로젝트와 시퀀스에 대한 샷 리스트 정보를 가지고 오는 함수.
	(리스트, err) 값을 반환한다.
	"""
	if not project:
		return [], "프로젝트 정보가 없습니다."
	if not seq:
		return [], "시퀀스 정보가 없습니다."
	restURL = "https://openpipeline.io/api/shots?project={0}&seq={1}".format(project, seq)
	try:
		shots = json.load(urllib2.urlopen(restURL))
	except:
		return [], "RestAPI에 연결할 수 없습니다."
	if "error" in shots:
		return [], shots["error"]
	return shots["data"], None

def SetStatus(project, name, task, status):
	"""
	아이템의 task에 대한 상태를 설정한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&task=%s&status=%s" % (project, name, task, status)
	try:
		request = urllib2.Request("https://openpipeline.io/api/setstatus", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetStartdate(project, name, task, startdate):
	"""
	아이템의 task에 대한 시작일을 설정한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&task=%s&startdate=%s" % (project, name, task, urllib2.quote(startdate))
	try:
		request = urllib2.Request("https://openpipeline.io/api/setstartdate", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def SetPredate(project, name, task, predate):
	"""
	아이템의 task에 대한 1차마감일을 설정한다.
	에러문자열을 반환한다. 에러가 없다면 ""문자를 반환한다.
	"""
	data = "project=%s&name=%s&task=%s&predate=%s" % (project, name, task, urllib2.quote(predate))
	try:
		request = urllib2.Request("https://openpipeline.io/api/setpredate", data)
		err = urllib2.urlopen(request).read()
	except:
		err = "restAPI에 접근할 수 없습니다."
	return err

def searchwordItems(project, searchword, status=[]):
	"""
	아이템들을 검색어(searchword)를 통해 찾는 함수이다.
	status는 리스트 형태로 입력받고, 값이 없다면 기본설정을 해준다.
	(리스트, err) 값을 반환한다.
	"""
	values = {}
	values["project"] = project
	values["searchword"] = searchword

	if not status:
		values["assign"] = "true"
		values["ready"] = "true"
		values["wip"] = "true"
		values["confirm"] = "true"
		values["done"] = "true"
	else:
		for i in status:
			if i not in ["assign", "ready", "wip", "confirm", "done", "omit", "hold", "out", "none"]:
				continue
			values[i] = "true"
	url = "https://openpipeline.io/api2/items"
	query = urllib.urlencode(values)
	restURL = url + "?" + query
	try:
		data = json.load(urllib2.urlopen(restURL))
	except:
		return [], "RestAPI에 연결할 수 없습니다."
	if "error" in data:
		return [], data["error"]
	return data["data"], None
