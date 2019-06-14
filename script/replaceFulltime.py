#coding:utf8
import CSI2

p = CSI2.CSI("thesea")
for i in list(p.search("")):
	if i["ddline2d"] != "":
		if len(i["ddline2d"]) == 4:
			slug = i["slug"]
			oldtime = i["ddline2d"]
			newtime = CSI2.ToFullTime(i["ddline2d"])
			print slug,oldtime,newtime
			p.set_ddline2d(slug, newtime)
