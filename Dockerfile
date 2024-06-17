FROM golang:1.22 as build

WORKDIR /go/src/github.com/thoughtgears/run-service-discovery
COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/app .

FROM golang:1.22-alpine as artifact
WORKDIR /app
COPY --from=build /go/src/github.com/thoughtgears/run-service-discovery/builds/app .

EXPOSE 8080
CMD ["./app"]
