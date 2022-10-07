go run assets/asset_generate.go
go install -ldflags "-X main.SHA1VER=`git rev-parse HEAD` -X main.BUILDTIME=`date -u +%Y-%m-%dT%H:%M:%S`"
$HOME/go/bin/openpipelineio -http :80 -devmode -reviewrender -scanplaterender -debug
