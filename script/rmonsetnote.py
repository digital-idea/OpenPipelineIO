#!/usr/bin/env python
#coding:utf-8

import sys, os, getopt
import json
import urllib2
from datetime import datetime
from CSI2 import *

"""
특정 작성자명의 pmnote를 일괄로 삭제하는 툴
프로젝트 정보와 작성자명을 옵션으로 입력 받는다.
"""
def main():
	# 옵션 설정
	optProject = ""
	optInputUser = ""
	optCheck = ""
	optDel = ""
	try:
		opts, args = getopt.getopt(sys.argv[1:], "p:i:cd")
		for opt,arg in opts:
			if opt in ("-p"): #project
				optProject = arg
			elif opt in ("-i"): #inputUser
				#optInputUser = "%s"%arg.decode('utf8')
				optInputUser = arg
			elif opt in ("-c"): #check
				optCheck = "check"
			elif opt in ("-d"):
				optDel = "del"
	except getopt.GetoptError, err:
		sys.stderr.write(str(err)+"\n")

	if optProject:
		#project = optProject
		project = "TEMP"
	if optInputUser:
		inputUser = optInputUser
	if optCheck:
		print "\n - 체크가 진행됩니다."
	if optDel:
		print "\n - 삭제가 진행됩니다."
	
	print "\n - 처리 결과는 /show/프로젝트명/product/doc에 저장됩니다."
	print " - 저장파일명 : 프로젝트명_removeOnsetnote_날짜시간.txt\n"

	# 삭제된 데이터를 txt로 저장한다. 실수로 삭제할 경우 대비
	path = "/show/%s/product/doc"%project
	time = datetime.today().strftime("%Y%m%d_%H_%M_%S")
	#sys.stdout = open('%s/%s_removeOnsetnote_%s.txt'%(path,project,time),'w')

	# CSI Seq DB
	restURL = "http://10.0.90.251/api/seqs?project=%s"%project
	try:
		seqsCsiDB = json.load(urllib2.urlopen(restURL))
	except:
		print("RestAPI에 연결할 수 없습니다.")

	p = CSI(project) # csi2 project

	s = 0 # Count Shots
	n = 0 # Count Notes
	print " - Remove Note Infomation"
	# 시퀀스의 샷 리스트
	for seq in seqsCsiDB['data']:
		# CSI shots DB
		restURL = "http://10.0.90.251/api/shots?project=%s&seq=%s"%(project,seq)
		try:
			shotsCsiDB = json.load(urllib2.urlopen(restURL))
		except:
			print("RestAPI에 연결할 수 없습니다.")

		# 아이템
		isInputUser = False # check inputUser for count
		for shot in shotsCsiDB['data']:
			# slug 설정
			slug = seq + "_" + shot + "_org" # 입체_left 인 경우 처리 필요
			# CSI item DB
			restURL = "http://10.0.90.251/api/item?project=%s&slug=%s"%(project,slug)
			try:
				itemsCsiDB = json.load(urllib2.urlopen(restURL))
			except:
				print("RestAPI에 연결할 수 없습니다.")

			onsetnote = itemsCsiDB['onsetnote']

			rmOnsetnote = [] # rm note list
			if onsetnote:
				for noteList in onsetnote:
					note = noteList.split(";")
					if inputUser == note[2].encode("utf-8"):
						rmOnsetnote.append(noteList)
						#aa.append(noteList)

						print " X "+slug+" >",noteList.encode("utf-8").replace("\n",",")

						isInputUser = True # InputUser 확인
						n += 1 # count notes

				# remove rmOnsetnote
				if optDel:
					for i in rmOnsetnote:
						err = p.rm_onset(slug,i)
						if err:
							print err

				if isInputUser == True: # count shots
					s += 1

				isInputUser = False # check inputUser for count

	if optCheck:
		print "\n - 결과 : %s샷의 %s개의 노트가 검색되었습니다.\n"%(s,n)
	else:
		print "\n - 결과 : %s샷의 %s개의 노트가 삭제되었습니다.\n"%(s,n)

if __name__ == "__main__":
	main()
