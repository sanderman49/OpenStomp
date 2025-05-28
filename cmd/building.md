For Linux: `go build .`
For Windows:
- Building from Linux: `env GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ go build -ldflags -H=windowsgui .`
- Building from Windows: `go build -ldflags -H=windowsgui .`

