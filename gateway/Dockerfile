# Start from golang:1.18-alpine base image
FROM golang:1.24-alpine AS build-env

# Add Maintainer Info
LABEL maintainer="Aburavi Arrazi <aburavi.arrazi@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/gateway


FROM alpine:latest

# The latest alpine images don't have some tools like (`git` and `bash`).
#Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
     apk add --no-cache bash git openssh busybox-extras

WORKDIR /app

COPY --from=build-env /build/main main
COPY --from=build-env /build/config config.env

# Expose port 50051 to the outside world
EXPOSE 8080

# Run the executable

ENTRYPOINT [ "./main" ]

