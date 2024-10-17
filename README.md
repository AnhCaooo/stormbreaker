
# Stormbreaker

A service fetches electric price in Finland and does handling behind the scene which is written in [Go](https://go.dev/)

## Development Setup
### Prerequisite
- Make sure you have Go installed on your machine. If not, you can install from [Golang official page](https://go.dev/doc/install) 
- Make sure you set up the environment variables by specifying it in the `config/config.yml`

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

3. Run the application locally
```bash
go run cmd/main.go
```

### Build Docker image

Build image locally
```bash
# Option 1 
docker build --tag <image-name>:<number-version> .

# Option 2
docker build -t <image-name>:<number-version> .

# Option 3
# this command by default will build image with tag version 'latest'. 
# this is an enhancement when before the image is built, all unit tests will be executed
make docker 
```

### Run Docker image locally
```bash
docker run --name <image-name> -d <container-name>:<tagged-image-version>
```
Terms explanation:
`--name`: specify the name of image while running container

`-d`: detached mode

**Notes**: if you are running several containers in same machine, in order to make images can accessible between each other. You can consider to use Docker Network 
```bash
# create docker network 
docker network create <network_name>

# include docker network while run image
docker run --name <image-name> -d --network <network_name> <container-name>:<tagged-image-version>
```

### Tag and push image to Docker hub
**Note**: firstly need to log in Docker (only in case you have not logged in)

```bash
docker login -u "<docker_username>" -p "<docker_password>" docker.io
```

#### Step 1
Tag image
```bash
docker tag <image-name>:<number-version> <docker_username>/<image-name>:<tagged-version-number> 
```

#### Step 2
Push image to Docker hub
```bash
docker push <docker_username>/<image-name>:<tagged-version-number> 
```
