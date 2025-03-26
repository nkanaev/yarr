FROM golang:alpine3.21 AS build
RUN apk add build-base git
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg \
  make host

FROM alpine:latest
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=build /src/out/yarr /usr/local/bin/yarr
EXPOSE 7070
ENTRYPOINT ["/usr/local/bin/yarr"]
CMD ["-addr", "0.0.0.0:7070", "-db", "/data/yarr.db"]
