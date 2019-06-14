#coding:utf-8
# 이 코드는 마이그레이션을 돕는 일회성 스크립트이다.
import sys
sys.path.append("/lustre/INHouse/Tool/csi3/pythonapi")
sys.path.append("/lustre/INHouse/CentOS/python2.7.5/lib/python2.7/site-packages")
import csi3
import pymongo

DBIP = "10.0.90.251"
#DBIP = "10.0.92.135"

def main():
	projects, err = csi3.Projects()
	if err:
		sys.stderr.write(err+"\n")
		sys.exit(1)
	for project in projects:
		db = pymongo.Connection(DBIP, 27017).project
		docs = db[project].find()
		for doc in docs:
			if doc["type"] == "asset":
				continue
			# platein, plateout 값이 없다면 셋팅한다.
			print project, doc["slug"], doc.get("scanframe",0)
			length = doc.get("scanframe",0)
			if length == 0: # Full3D샷
				db[project].update({"slug": doc["slug"]}, {"$set": {"platein": 1001}})
				db[project].update({"slug": doc["slug"]}, {"$set": {"plateout": 1001}})
			else: # 일반샷
				db[project].update({"slug": doc["slug"]}, {"$set": {"platein": 1001}})
				db[project].update({"slug": doc["slug"]}, {"$set": {"plateout": 1001 + length -1}})

if __name__ == "__main__":
	main()
