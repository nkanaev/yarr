FROM golang:1.15 AS build
RUN apt install gcc -y
WORKDIR /src/
COPY . /src/
RUN go build -tags "sqlite_foreign_keys release linux" -ldflags="-s -w" -o /bin/yarrd yarrd.go

FROM ubuntu:20.04
COPY --from=build /bin/yarrd /bin/yarrd
ENTRYPOINT ["/bin/yarrd"]
