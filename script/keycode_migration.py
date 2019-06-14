#coding:utf8
# 이 코드는 마이그레이션을 돕는 스크립트이며 추후 폐기된다.
import pymongo
import CSI2
import time

#DBIP = "10.0.90.251"
DBIP = "10.0.92.135"
Runmode = False #DB를 돌리지 않고 Print할 때 사용하는 옵션 Print시에는 False로 수동으로 바꾸어서 확인할 때 사용하는 용도.

def main():
	db = pymongo.Connection(DBIP, 27017).project
	for project in CSI2.PROJECTLIST():
		if project in ["jangmi","shs","gunhamdo","gunghop"]:
			docs = db[project].find()
			for doc in docs:
				if doc["jin"] != "":
					print project, doc["slug"], doc["jin"]
					if Runmode:
						db[project].update({"slug": doc["slug"]}, {"$set": {"justkeycodein": doc["jin"]}})
						time.sleep(0.0001)
				if doc["jout"] != "":
					print project, doc["slug"], doc["jout"]
					if Runmode:
						db[project].update({"slug": doc["slug"]}, {"$set": {"justkeycodeout": doc["jout"]}})
						time.sleep(0.0001)

if __name__ == "__main__":
	main()
