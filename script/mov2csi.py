#!/usr/bin/env python
#coding:utf-8
# 이 코드는 넘벳을 통하지 않고 mov파일을 가지고 CSI에 샷을 등록하는 코드입니다.
# 드라마의 경우 외주가 많이 나가기 때문에 샷 관리를 위해서 제작되었습니다.

import sys
import os
import re
sys.path.append("/lustre/INHouse/Tool/dipath")
import dipath
sys.path.append("/lustre/INHouse/Tool/csi3/pythonapi")
import csi3
FFMPEG="/lustre/INHouse/Tool/ffmpeg_lib/ffmpeg/bin/ffmpeg"
from CSI2 import *
sys.path.append("/lustre/INHouse/Tool/pymediainfo")
from pymediainfo import MediaInfo
import subprocess

def typ2linktyp(typ):
    if typ.startswith("org"):
        return "org"
    if typ.startswith("left"):
        return "left"
    return typ

def duration2frame(duration, fps):
    """
    duration, fps를 받아서 프레임수를 반환한다.
    """
    # 만약 fps가 24라면, 1/24 second, (1000 ms / 24 = 41.67 ms)
    return int(round(float(duration)*float(fps)*0.001))

def main():
  cwdpath = os.getcwd()
  # 프로젝트값 체크
  project,err = dipath.Project(cwdpath)
  if err:
    print(err)
    sys.exit(1)

  for i in sorted(os.listdir(cwdpath)):
    f = cwdpath + "/" + i
    filename, ext = os.path.splitext(i)

    # 에러처리
    if ext.lower() != ".mov":
      print("파일을 처리할 수 없습니다. : " + f + "\n")
      continue

    # 네이밍 에러 처리
    regex = re.compile("(\D+\d{2}S\d{3}_\d{4})(.*$)")
    r = regex.match(filename)

    shotname =  r.group(1)
    others = r.group(2)
    typ = ""

    if "_" in others:
      typ =  re.findall("_(.+)$", others)[0]
    if not regex.findall(shotname):
      print("mov파일명이 EP01S036_0010 형태가 아닙니다. : " + i)
      continue

    if typ == "":
      typ = "org"

    # MediaInfo를 이용해 mov의 정보를 구한다
    v = MediaInfo.parse(i).tracks[1] # MediaInfo Video 항목 지정
    scanframe = str(duration2frame(v.duration, v.frame_rate)) # frame형태의 duration을 구한다
    platesize = str(v.width)+'x'+str(v.height) # platesize를 구한다

    # CSI등록
    cmd = ["csi3","-add","item","-project",project,"-name",shotname,"-type",typ,"-scanname",filename,"-scanframe",scanframe,"-platesize",platesize] # 참고:filename을 scanname으로 등록한다
    addCsi = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, err = addCsi.communicate()
    if err != "":
        sys.stderr.write(err)
        continue

    os.system("curl -X POST -d 'project=%s&name=%s&type=%s&path=%s' http://10.0.90.251/api/setthummov" % (project, shotname, typ, f)) # mov링크 등록
    if (typ.startswith("org") and len(typ) > 3) or (typ.startswith("left") and len(typ) > 4):
      os.system("curl -X POST -d 'project=%s&name=%s&type=%s&path=%s' http://10.0.90.251/api/setthummov" % (project, shotname, typ2linktyp(typ), f)) # org1,left1형태인경우 mov링크 업데이트
    thumbJpgPath = "/lustre/INHouse/Tool/csidata/thumbnail/%s/%s_%s.jpg" % (project, shotname, typ) # 썸네일 경로
    os.system("%s -y -v 0 -i %s -vframes 1 -vf scale=410:231 %s" % (FFMPEG, f, thumbJpgPath)) # 썸네일 생성

    # org1,src1,ref형태인 경우 org에 left1,lsrc,lref 형태인 경우 left에 mov 링크 추가
    if typ not in ["org", "left", "asset"]:
      p = CSI(project)
      if typ.startswith("src") or typ.startswith("ref"):
        typ = "org"
      elif typ.startswith("lsrc") or typ.startswith("lref"):
        typ = "left"
      else:
        typ = typ2linktyp(typ)
      p.set_link(shotname+"_%s"%typ, "", dipath.Rmlustre(f))

    # mov2csi으로 shot 추가한 log 기록
    slug = filename
    if type == "":
      slug = filename + "_org"
    log = "shot added."

    dilog = "/lustre/INHouse/CentOS/bin/dilog -tool mov2csi -project %s -slug %s -log '%s'" % (project, slug, log)

    os.system(dilog)

if __name__ == "__main__":
       main()
