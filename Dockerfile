FROM --platform=linux/amd64 golang:alpine

# Set destination inside the container
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY ./cmd/ ./cmd/
COPY ./docs/ ./docs/
COPY ./internal/ ./internal/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /stormbreaker ./cmd/

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 5001

# Run
CMD ["/stormbreaker"]