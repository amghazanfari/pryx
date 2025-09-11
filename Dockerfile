FROM golang:1.24.6 AS build

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o pryx cmd/app/main.go
RUN CGO_ENABLED=0 go build -o migrate cmd/migrate/main.go
 
FROM ubuntu:24.04 AS run

RUN apt update && apt install ca-certificates -y

# Copy the application executable from the build image
COPY --from=build /app/pryx /app/pryx
COPY --from=build /app/migrate /app/migrate

WORKDIR /app
EXPOSE 8080
CMD ["/app/pryx"]
