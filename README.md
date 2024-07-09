# automation_appWeb_test_golang
<br> <br> 

This project aims to contain an API-REST, using the Go programming language (Golang) and automating test scenarios using GO + Testing* . It uses the Visual Studio Code (VSCode) development environment and Go in version go1.22.3 for the Darwin/amd64 platform.

<br> <br> 

<hr>

<br> <br> 

Example - Unit test reports:
<br>
go test -coverprofile=coverage.out
<br>
Generates a coverage.out file with the coverage data.
<br><br>
go tool cover -func=coverage.out
<br>
 Displays a coverage summary in the terminal.
 <br><br>
go tool cover -html=coverage.out
<br>
Opens a coverage report in HTML format.
<br><br>
Tag for running unit tests and integration tests
<br>
go test -v -tags=integration ./…
<br> <br> 

<hr>

<br> <br> 
The database must be created locally
<br> 
docker-compose up -d
<br> <br> 
The service must be carried out
<br> 
go run ./cmd/api
<hr>

<br> <br> 
CURL - Successful token generation
<br> 
curl http://localhost:8090/auth -X POST -H "Content-Type: application/json" -d '{"email":"admin@example.com","password":"secret"}'
<br> <br>