## Tray
A cross-platform tray for Go.

## Platform implement status
- [x] Windows
- [x] Linux
- [ ] Darwin

## Usage
```shell script
git clone github.com/poohvpn/tray
cd tray
```
Run [example](example/main.go) with
### Linux
```shell script
cd example
go run .
```

### Windows
```shell script
go get github.com/akavel/rsrc
cd example
rsrc -manifest main.exe.manifest -arch amd64 -o rsrc.syso
go run .
```