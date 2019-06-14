#coding:utf8
# 이 코드는 마이그레이션을 돕는 스크립트이며 추후 폐기 된다.

import pymongo
import time

DBIP = "10.0.90.251"
#DBIP = "10.0.92.135"

def replacelist(taglist, before, after):
	newlist = []
	for i in taglist:
		if i == before:
			newlist.append(after)
		else:
			newlist.append(i)
	return newlist

def renametag(before, after):
	"""
	군함도에서 사용되는 태그문자를 수정.
	"""
	db = pymongo.Connection(DBIP, 27017).project
	docs = db["gunhamdo"].find()
	for doc in docs:
		if before in doc["tag"]:
			#print doc["tag"], before, after
			#print replacelist(doc["tag"], before, after)
			db["gunhamdo"].update({"slug": doc["slug"]}, {"$set" : {"tag" : replacelist(doc["tag"], before, after) }})
			time.sleep(0.0001)


if __name__ == "__main__":
	renametag("CA","LvA")
	renametag("CB","LvB")
	renametag("CC","LvC")
	renametag("S2","Dep2")
	renametag("1S","Dep1")
	renametag("SA","ZA")
	renametag("SB","ZB")
	renametag("SC","ZC")
	renametag("SD","ZD")
	renametag("SE","ZE")
	renametag("SF","ZF")
	renametag("SG","ZG")
	renametag("SH","ZH")
	renametag("SI","ZI")
	renametag("SJ","ZJ")
	renametag("SK","ZK")
	renametag("SL","ZL")
	renametag("SM","ZM")
	renametag("SN","ZN")
	renametag("SO","ZO")
	renametag("SP","ZP")
	renametag("SQ","ZQ")

