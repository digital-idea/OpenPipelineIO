#!/bin/bash
Token="JDJhJDEwJEdDWm1lVXEzZm9vNWVIQXhBV21iZnU2OFl6OHlBVTRRY1lEL0JwekdHWFlCNFBTNDQvTkVD"
# 토큰키가 첫번째 인수로 오면 덮어쓰기한다.
if [ $1 != '' ]
then
	Token=$1
fi
echo "Token:" $Token
curl -X POST -H "Authorization: Basic $Token" -d "id=none&description=none&textcolor=#FFFFFF&bgcolor=#3D3B3B&bordercolor=#787474&order=0&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=hold&description=hold&textcolor=#FFFFFF&bgcolor=#606161&bordercolor=#606161&order=1&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=done&description=done&textcolor=#000000&bgcolor=#E4D2B7&bordercolor=#E4D2B7&order=2&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=out&description=out&textcolor=#000000&bgcolor=#EEA4F1&bordercolor=#EEA4F1&order=3&defaulton=false" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=assign&description=assign&textcolor=#000000&bgcolor=#FFF76B&bordercolor=#FFF76B&order=4&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=ready&description=ready&textcolor=#000000&bgcolor=#BEEF37&bordercolor=#BEEF37&order=5&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=wip&description=wip&textcolor=#000000&bgcolor=#77BB40&bordercolor=#77BB40&order=6&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=confirm&description=confirm&textcolor=#000000&bgcolor=#59C6FD&bordercolor=#59C6FD&order=7&defaulton=true" "http://localhost/api/addstatus"
curl -X POST -H "Authorization: Basic $Token" -d "id=omit&description=omit&textcolor=#000000&bgcolor=#FC9F55&bordercolor=#FC9F55&order=8&defaulton=false" "http://localhost/api/addstatus"
