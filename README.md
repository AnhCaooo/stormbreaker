
# Stormbreaker

A service fetches electric price in Finland and does handling behind the scene which is written in [Go](https://go.dev/)

## Development Setup
### Prerequisite
- Make sure you have Go installed on your machine. If not, you can install from [Golang official page](https://go.dev/doc/install) 

### Getting started
1. Fetch all dependencies listed in the `go.mod` file and remove unused dependencies from repository
```bash
  go mod tidy
```

2. Take total control over the code
```bash
  go mod vendor
```    
This command will give you total control of all libraries which are used in this repo. Even if the owner of some libraries change, archive or remove the code, you still have the code running there. Learn more: [understand "go mod vendor"](https://stackoverflow.com/questions/76705408/understanding-go-mod-vendor) 

3. Run the application 
```bash
go run cmd/main.go
```

### Build Docker image

**Note**: in case you are planning to push your docker image, you first need to log in Docker (only in case you have not logged in)

```bash
docker login -u "<docker_username>" -p "<docker_password>" docker.io
```

#### Step 1
Build image locally
```bash
# Option 1 (not recommended)
docker build --tag stormbreaker:<number-version> .

# Option 2 (not recommended)
docker build -t stormbreaker<number-version> .

# Option 3 (recommended)
# this command by default will build image with tag version 'latest'. 
# this is an enhancement when before the image is built, all unit tests will be executed
make docker 
```

#### Step 2
Tag image
```bash
docker tag stormbreaker:<number-version> anhcaoo/stormbreaker:<tagged-version-number> 
```

#### Step 3
Push image to Docker hub
```bash
docker push anhcaoo/stormbreaker:<tagged-version-number> 
```