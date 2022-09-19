# WebAuthn

```
主程序: webauthn
安装位置: /srun3/bin
主程序需要执行权限
chmod +x webauthn
配置文件: webauthn.yaml
安装位置: /srun3/etc
```

```
配置文件修改说明
app:
  name: "webauthn"
  mode: "prod"
  protocol: "https" // 不用修改 不被浏览器信任的证书不知道可不可以,需要测试
  host: "idp.srun.com" // 修改为当前所在服务器域名, IP不知道可不可以,需要测试
  port: 18080
  cert_file: "/srun3-bak0819/httpd/etc/ssl.key/server8080.crt" // 需要修改为真实证书路径
  key_file: "/srun3-bak0819/httpd/etc/ssl.key/server8080.key" // 需要修改为真实证书路径
log:
  level: "warn"
  filename: "webauthn.log"
  max_size: 100
  max_age: 30
  max_backups: 7
mysql:
  ip: "127.0.0.1" // 需要根据真实情况修改
  port: 3506 // 需要根据真实情况修改
  pwd: "Srun4000@srun.com" // 需要根据真实情况修改
  user: "icc" // 需要根据真实情况修改
  dbname: "srun4k" // 需要根据真实情况修改
  max_life_time: 5
  max_open: 100
  max_idle: 50
redis:
  ip: "127.0.0.1" // 需要根据真实情况修改
  port: 16384 // 需要根据真实情况修改
  pwd: "srun_3000@redis" // 需要根据真实情况修改
  index: 0
  pool_size: 100
sso:
  url: "http://127.0.0.1" // 需要根据真实情况修改
  secret: "123456" // 需要根据真实情况修改
```

```
登录MySQL进入srun4k库,执行sql文件: webauthn.sql
```

```
启动方式
- 调试启动: 
cd /srun3/bin
./webauthn

- 后台启动: 
cd /srun3/bin
nohup ./webauthn &
```

```
sso配置说明
8082上的微信临时放行key(必须核实)
也可在服务器文件srun4kauth.xml中ApiAuthSecret字段获得,需修改EnableAPIAuth=1
然后重启srun3kauth
重启命令: service srun3kauth restart
```

