@echo off

cd server
swagger generate server -A TodoList -f swagger.yaml
cd ..
go test .\tests\...
if %errorlevel% == 0 (
    go install .\server\cmd\todo-list-server\
    todo-list-server --port 8080 --tls-port 8085 --read-timeout 5s --tls-certificate C:\Users\vanme\go\src\flightAPI\openSSL\localhost.crt --tls-key C:\Users\vanme\go\src\flightAPI\openSSL\localhost.key
) else (
    echo Unit tests Failed
)