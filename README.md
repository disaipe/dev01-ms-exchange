# Dev01 MS Exchange statistics module helper

```shell
dev01-ms-exchange.exe -addr=[:8090] -app.url=http://dev01.com/api/fw -app.secret=[secret word]
```

```shell
# install as Windows service
dev01-ms-exchange.exe -srv.install

# start as Windows service
dev01-ms-exchange.exe -srv ...

# uninstall Windows service
dev01-ms-exchange.exe -srv.uninstall
```

# Credentials

Execute command to generate credentials file:

```shell
powershell.exe -F generateCredentials.ps1
```