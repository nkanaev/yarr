FROM golang:1.15 AS build
RUN apt install gcc -y
WORKDIR /src
COPY . .
RUN GOOS=linux go build -tags "sqlite_foreign_keys release linux" -ldflags="-s -w" -o /usr/local/bin/yarr main.go
RUN ls /usr/local/bin

FROM ubuntu:20.04
COPY --from=build /usr/local/bin/yarr /usr/bin/yarr
ENTRYPOINT ["/usr/bin/yarr"]
