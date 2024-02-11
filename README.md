
# Stormbreaker

A service fetches electric price in Finland and does handling behind the scene which is written in [Go](https://go.dev/)

## Development Setup
### Prerequisite
- Make sure you have Go installed on your machine. If not, you can install from [Golang official page](https://go.dev/doc/install) 

### Getting started
Fetch all dependencies listed in the `go.mod` file and remove unused dependencies from repository
```bash
  go mod tidy
```

Take total control over the code
```bash
  go mod vendor
```    
This command will give you total control of all libraries which are used in this repo. Even if the owner of some libraries change, archive or remove the code, you still have the code running there. Learn more: [understand "go mod vendor"](https://stackoverflow.com/questions/76705408/understanding-go-mod-vendor) 

- Run the application 
```bash
go run cmd/main.go
```