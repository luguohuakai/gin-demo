#### 编译为Linux可执行文件

```
SET CGO_ENABLED=0
SET GOARCH=amd64
SET GOOS=linux
go build -o webauthn

```

#### 编译为windows可执行文件

```
SET CGO_ENABLED=1
SET GOARCH=
SET GOOS=windows
go build -o webauthn.exe

```