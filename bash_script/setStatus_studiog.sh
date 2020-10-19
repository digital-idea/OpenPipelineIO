#!/bin/bash
# 2020.10.19 Status List
#13-omit
#12-hold
#11-keep
#10-send
#9-done
#8-ok
#7-qc
#6-ing
#5-retake 
#4-confirm
#3-active
#2-ready
#1-wait
#0-none

Token="JDJhJDEwJEdDWm1lVXEzZm9vNWVIQXhBV21iZnU2OFl6OHlBVTRRY1lEL0JwekdHWFlCNFBTNDQvTkVD"
# 토큰키가 첫번째 인수로 오면 덮어쓰기한다.
if [ $1 != '' ]
then
	Token=$1
fi
echo "Token:" $Token
curl -X POST -H "Authorization: Basic $Token" -d "id=omit&description=omit&textcolor=#000000&bgcolor=#FC9F55&bordercolor=#FC9F55&order=13&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=hold&description=hold&textcolor=#FFFFFF&bgcolor=#606161&bordercolor=#606161&order=12&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=keep&description=keep&textcolor=#000000&bgcolor=#49A0BB&bordercolor=#49A0BB&order=11&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=send&description=send&textcolor=#000000&bgcolor=#E374B3&bordercolor=#E374B3&order=10&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=done&description=done&textcolor=#000000&bgcolor=#AB9878&bordercolor=#AB9878&order=9&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=ok&description=ok&textcolor=#000000&bgcolor=#E4D2B7&bordercolor=#E4D2B7&order=8&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=qc&description=qc&textcolor=#000000&bgcolor=#6DC7A3&bordercolor=#6DC7A3&order=7&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=ing&description=ing&textcolor=#000000&bgcolor=#BEEF37&bordercolor=#BEEF37&order=6&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=retake&description=retake&textcolor=#FFFFFF&bgcolor=#FF430A&bordercolor=#FF430A&order=5&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=confirm&description=confirm&textcolor=#000000&bgcolor=#59C6FD&bordercolor=#59C6FD&order=4&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=active&description=active&textcolor=#000000&bgcolor=#FD7060&bordercolor=#FD7060&order=3&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=ready&description=ready&textcolor=#000000&bgcolor=#FFCC40&bordercolor=#FFCC40&order=2&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=wait&description=wait&textcolor=#000000&bgcolor=#FFF76B&bordercolor=#FFF76B&order=1&defaulton=true&initstatus=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=none&description=none&textcolor=#FFFFFF&bgcolor=#3D3B3B&bordercolor=#787474&order=0&defaulton=false" "http://localhost/api/addstatus"