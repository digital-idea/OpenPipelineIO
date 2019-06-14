#coding:utf8
# 이 코드는 마이그레이션을 돕는 스크립트이며 추후 폐기 된다.

import pymongo
import time

DBIP = "10.0.90.251"
#DBIP = "10.0.92.135"

def scanlength2scanframe():
	"""
	군함도에서 앞으로 도입될 스켄길이 검색을 미리 사용하고 싶다고 요청이 옴.
	"""
	db = pymongo.Connection(DBIP, 27017).project
	docs = db["gunhamdo"].find()
	for doc in docs:
		print doc["slug"], doc["scanlength"]
		if doc["scanlength"] != "":
			db["gunhamdo"].update({"slug": doc["slug"]}, {"$set" : {"scanframe" : int(doc["scanlength"])}})
		else:
			db["gunhamdo"].update({"slug": doc["slug"]}, {"$set" : {"scanframe" : 0}})
		time.sleep(0.0001)


if __name__ == "__main__":
	scanlength2scanframe()

