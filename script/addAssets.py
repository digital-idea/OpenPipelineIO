#!/usr/bin/env python                                                       
#coding:utf-8

# 엑셀의 어셋들을 자동으로 OpenPipelineIO에 등록하는 툴
# 입력 정보: 타입,에셋명,설명

import sys, os
import json
import urllib2
sys.path.append("/lustre/INHouse/Tool/pmutil")
import xlrd

def checkProjectNameFromFileName(path):
	"""
	엑셀 파일의 영문프로젝트명이 존재하는지 확인한다.
	"""
	err = ""
	for f in os.listdir(path):
		fileProject = f.split(".")[0].split("_")[0] # 파일명에서 프로젝트명을 구함

	# OpenPipelineIO 프로젝트 리스트
	restURL = "http://10.0.90.251/api/projects" # 기본적으로 현재 작업중인 프로젝트를 가지고옵니다.(pre + post + backup상태)
	try:
		projects = json.load(urllib2.urlopen(restURL))
	except:
		print("RestAPI에 연결할 수 없습니다.")
	# 에러처리
	if projects["error"]:
		print(projects["error"])

	if fileProject not in projects["data"]: # 프로젝트 리스트에서 확인
		err = "\n - %s 는 등록되지 않은 프로젝트명입니다.\n - 파일명의 프로젝트명을 확인해주세요.\n"%fileProject

	return fileProject,err

def getAssetListFromExcel(excel):
	"""
	엑셀에서 Asset들의 정보를 가지고 온다.
	"""
	# 엑셀파일의 정보를 리스트로 만든다.
	assetList = {}
	try:
		wb = xlrd.open_workbook(excel)
	except ValueError:
		sys.stderr.write(" - %s : 파일에 오류가 있거나 엑셀형식이 아닙니다.\n"%f)
		return assetList,1

	ws = wb.sheet_by_index(0) # 첫번째 시트 참조

	# assetName, assetType, checkComponent
	nIndex = ""
	tIndex = ""
	cIndex = ""
	for c in range(0, ws.ncols):
		for r in range(0, ws.nrows):
			if ws.col_values(c)[r] == "assetName":
				nIndex = c
				rowIndex = r
				continue # 찾으면 loop종료
		for r in range(0, ws.nrows):
			if ws.col_values(c)[r] == "assetType":
				tIndex = c
				continue # 찾으면 loop종료
		for r in range(0, ws.nrows):
			if ws.col_values(c)[r] == "checkComponent":
				cIndex = c
				continue # 찾으면 loop종료

	if nIndex == "" or tIndex == "" or cIndex == "":
		if nIndex == "": # 에셋네임 열 정보
			sys.stderr.write(" - 에셋네임 열 제목을 찾을 수 없습니다.\n")
		if tIndex == "": # 에셋타입 열 정보
			sys.stderr.write(" - 에셋타입 열 제목을 찾을 수 없습니다.\n")
		if cIndex == "": # CheckComponent
			sys.stderr.write(" - checkComponent 열 제목을 찾을 수 없습니다.\n")

		return assetList,1

	# 엑셀 정보를 리스트로 만듭니다.
	r = rowIndex + 1 # 정보가 시작하는 index
	while r < ws.nrows:
		assetName = str(ws.col_values(nIndex)[r])
		assetType = str(ws.col_values(tIndex)[r])
		component = str(ws.col_values(cIndex)[r])

		# 정보가 없다면 오류 출력 후 종료
		if assetName == "" or assetType == "" or component == "":
			if assetName == "":
				sys.stderr.write(" - 어셋네임 정보가 없습니다.\n")
			if assetType == "":
				if assetName:
					sys.stderr.write(" -%s의 어셋타입 정보가 없습니다.\n"%assetName)
				else:
					sys.stderr.write(" - 어셋타입 정보가 없습니다.\n")
			if component == "":
				if assetName:
					sys.stderr.write(" - %s어셋에 component,assembly 정보가 없습니다.\n"%assetName)
				else:
					sys.stderr.write(" - 어셋에 component,assembly 정보가 없습니다.\n")

			return assetList,1

		# {assetName:assetType,component,...}
		assetList[assetName] = [assetType,component]

		r += 1

	del wb

	return assetList,""

def addAssets(proj,name,typ,component):
	"""
	asset을 등록한다. component,assembly정보에 맞춰 각각 등록한다.
	"""
	if component == "assembly":
		os.system("/lustre/INHouse/CentOS/bin/openpipelineio -add item -project %s -name %s -type asset -assettype %s -assettags %s,assembly" % (proj, name, typ, typ))
	elif component == "component":
		os.system("/lustre/INHouse/CentOS/bin/openpipelineio -add item -project %s -name %s -type asset -assettype %s -assettags %s,component" % (proj, name, typ, typ))

	result = "AssetName : %s\nAssetType : %s\nComponent : %s\n"%(name,typ,component)

	return result

def assetTypeCheck(typ):
	"""
	등록된 Asset Type인지 체크한다. 
	"""
	if typ in ["char", "comp", "env", "fx", "global", "matte", "prop", "plant", "vehicle", "concept"]:
		err = ""
	else:
		err = "등록된 에셋 타입이 아닙니다. 필요하다면 등록을 요청해주세요."
	return err

def main():
	cwdpath = os.getcwd() # 현재경로
	# 실행 경로에 파일이 하나만 존재해야한다.(오류 방지)
	files = len(os.listdir(cwdpath))
	if 1 < files or 0 == files:
		sys.stderr.write("파일이 하나만 존재해야 합니다. 파일이 없거나 하나 이상입니다.\n")
		sys.exit(1)

	# 엑셀파일 여부 체크
	for f in os.listdir(cwdpath):
		fileName, ext = os.path.splitext(f)
		if not any([ext == e for e in [".xls",".xlsx"]]): # 확장자로 엑셀파일 확인
			sys.stderr.write("- %s : .xls,.xlsx 확장자를 가진 엑셀 파일만 처리 가능합니다.\n"%f)
			continue
		excel = f

	# 프로젝트 리스트에 있는 영문 프로젝트명을 파일명에서 확인한다.
	project,err = checkProjectNameFromFileName(cwdpath)
	if err:
		print err
		sys.exit(1)

	# 엑셀 파일에서 assetList 를 만든다.
	assetList,err = getAssetListFromExcel(excel)
	if err:
		sys.stderr.write(" - 엑셀파일 에러 : 엑셀파일의 정보를 확인할 수 없습니다.\n - 열 제목을 확인해주세요.\n")
		sys.exit(1)

	n = 0
	for assetName in assetList:
		assetType = assetList[assetName][0] # 타입정보
		checkComponent = assetList[assetName][1] # Component, Assembly 체크

		err = assetTypeCheck(assetType)
		if err:
			continue
		
		# 체크 옵션이 있다면 등록하지 않는다.
		r = addAssets(project,assetName,assetType,checkComponent)

		print "\n - 어셋이 등록되었습니다.\n%s"%(r)
		n += 1

	print "\n총 %s개의 어셋이 등록되었습니다."%(n)
	print " * 총 개수가 맞지 않는 경우 엑셀파일에 중복된 어셋네임이 있는지 확인해주세요."

if __name__ == "__main__":
	main()
