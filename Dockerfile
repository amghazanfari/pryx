FROM golang:1.25 AS build

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

COPY . .

RUN go build -o pryx main.go
 
FROM ubuntu:latest AS run

# Copy the application executable from the build image
COPY --from=build /app/pryx /app/pryx

WORKDIR /app
EXPOSE 8080
CMD ["/app/pryx"]