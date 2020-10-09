FROM golang:alpine AS build
RUN apk add build-base git
WORKDIR /src
COPY . .
RUN make build_linux

FROM alpine:latest
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates
COPY --from=build /src/_output/linux/yarr /usr/local/bin/yarr
EXPOSE 7070
CMD ["/usr/local/bin/yarr", "-addr", "0.0.0.0:7070", "-db", "/data/yarr.db"]
