csi3 -add project -name TEMP
csi3 -add project -name circle

# set asset mov
csi3 -add item -type asset -name mamma -project TEMP -assettype char -assettags char,component
curl -d "project=TEMP&name=mamma&task=concept&mov=/show/TEMP/concept.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=model&mov=/show/TEMP/model.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=mm&mov=/show/TEMP/mm.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=layout&mov=/show/TEMP/layout.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=ani&mov=/show/TEMP/ani.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=fx&mov=/show/TEMP/fx.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=mg&mov=/show/TEMP/mg.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=previz&mov=/show/TEMP/previz.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=fur&mov=/show/TEMP/fur.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=rig&mov=/show/TEMP/rig.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=crowd&mov=/show/TEMP/crowd.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=look&mov=/show/TEMP/look.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=comp&mov=/show/TEMP/comp.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=matte&mov=/show/TEMP/matte.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=mamma&task=env&mov=/show/TEMP/env.mov" http://127.0.0.1/api/setmov

# set shot mov
csi3 -add item -type org -name SS_0010 -project TEMP
curl -d "project=TEMP&name=SS_0010&task=concept&mov=/show/TEMP/concept.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=model&mov=/show/TEMP/model.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=mm&mov=/show/TEMP/mm.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=layout&mov=/show/TEMP/layout.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=ani&mov=/show/TEMP/ani.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=fx&mov=/show/TEMP/fx.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=mg&mov=/show/TEMP/mg.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=previz&mov=/show/TEMP/previz.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=fursim&mov=/show/TEMP/fursim.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=sim&mov=/show/TEMP/sim.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=crowd&mov=/show/TEMP/crowd.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=light&mov=/show/TEMP/light.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=comp&mov=/show/TEMP/comp.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=matte&mov=/show/TEMP/matte.mov" http://127.0.0.1/api/setmov
curl -d "project=TEMP&name=SS_0010&task=env&mov=/show/TEMP/env.mov" http://127.0.0.1/api/setmov
