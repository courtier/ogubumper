GOOS=windows GOARCH=amd64 go build -o builds/ogubumperWindows.exe
GOOS=linux GOARCH=amd64 go build -o builds/ogubumperLinux
GOOS=darwin GOARCH=amd64 go build -o builds/ogubumperMac