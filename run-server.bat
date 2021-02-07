@echo off

cd server
swagger generate server -A TodoList -f swagger.yaml
if not %errorlevel% == 0 (
    echo Failed Code Generation
    goto end
)

cd ..

go test .\tests\...
if not %errorlevel% == 0 (
   echo Failed Unit tests
   goto end
)

go install .\server\cmd\todo-list-server\
if not %errorlevel% == 0 (
   echo Failed Install
   goto end
)

todo-list-server --port 8080 --tls-port 8085 --read-timeout 5s --tls-certificate C:\Users\vanme\go\src\flightAPI\openSSL\localhost.crt --tls-key C:\Users\vanme\go\src\flightAPI\openSSL\localhost.key

:End