FROM golang:1.15 AS build
RUN apt install gcc -y
WORKDIR /src
COPY . .
RUN make build_linux

FROM ubuntu:20.04
COPY --from=build /src/_output/linux/yarr /usr/bin/yarr
ENTRYPOINT ["/usr/bin/yarr", "-addr", "0.0.0.0:7070"]
