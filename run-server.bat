@echo off

swagger generate server -A TodoList -f server\swagger.yaml
if not %errorlevel% == 0 (
    echo Failed Code Generation
    goto end
)
go test .\tests\...
if %errorlevel% == 0 (
    go install .\server\cmd\todo-list-server\
    todo-list-server --port 8080 --tls-port 8085 --read-timeout 5s --tls-certificate C:\Users\vanme\go\src\flightAPI\openSSL\localhost.crt --tls-key C:\Users\vanme\go\src\flightAPI\openSSL\localhost.key
) else (
    echo Failed Unit tests
    goto end
)

:End