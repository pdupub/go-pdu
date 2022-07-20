## Build for Mac

```
export CGO_ENABLED=1
export GOARCH=amd64
go build -buildmode=c-archive -o pdu.a ./
```

How to fix follow Error:
```
"pdu.a" is missing one or more architectures required by this target: arm64
```
Adding "arm64" to Project -> Build Settings -> Excluded Architecture fixed the issue


## Build for iOS

```
export CGO_ENABLED=1
export GOOS=darwin
export GOARCH=arm64
export SDK=iphoneos
export CC=/usr/local/go/misc/ios/clangwrap.sh
### export CGO_CFLAGS="-fembed-bitcode"
go build -buildmode=c-archive -tags ios -o pdu.a ./
```

Build Settings -> Build Options -> Enable Bitcode : No

