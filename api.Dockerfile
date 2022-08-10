# BUILDER Stage
FROM golang:1.18.5-alpine3.16 AS builder

ARG CMD_BUILD_DIR

# Setting necessary environment variables needed to build the app in the image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copying and downloading dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copying the code into the container
COPY . .

# Checking Unit Tests before build
RUN go test --cover -v ./...

# Building app
RUN go build -o main ./cmd/$CMD_BUILD_DIR

###############################################

# RUNNER Stage
FROM alpine:3.14 AS runner

COPY --from=builder /app/main /

EXPOSE 8001
CMD ["/main"]