version 1.0

Build:
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go

Start:
./main
or

pm2 start ecosystem_gpark.config.js