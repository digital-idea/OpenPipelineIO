go run assets/asset_generate.go
go install
# 테스트시에는 썸네일이 실시간으로 바뀌어야 한다.
csi3 -http :80 -thumbnailage 1
