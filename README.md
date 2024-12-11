
# Stormbreaker

A service fetches electric price in Finland and does handling behind the scene which is written in [Go](https://go.dev/). 

Contribution are welcome. Here is the [development setup](#development-setup) you need to go though before able to use the service

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

3. Run the application
**Note** This service is having authentication handling in middleware layer. 
#### Run the service with database connection 
**Note**: Because this service is required database connection in order to get user's settings for electric price. So you need to have a database connection which is [MongoDB](mongodb.com).  and 1 way for quick verify (no need for db connection):

3.1 Follow and setup your own config file based on config template which is located in following directory: `./internal/config/config.template.yml`

After the first step, there are 2 ways to run this applications with full functionalities (3.2 or 3.3). Please your database configuration needs to be same as the configuration file. 

3.2. [Highly recommended] Run this service and Mongodb as Docker containers by taking [Docker Compose](https://docs.docker.com/compose/) into use. 

3.3 you can just build 2 docker images -> [create a docker network connection](#create-docker-network) -> apply network for those 2 containers

#### Run the service without database connection 
If you just want to run the application locally then you are agreed this service will lack database functionalities.
Ignore setting up `Mongo` in `main.go` and pass `nil` value while initializing `NewHandler`. Then run the following command:
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

#### Notes: create docker network 
If you are running several containers in same machine, in order to make images can accessible between each other. You can consider to use Docker Network 
```bash
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
