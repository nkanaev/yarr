# Stage 1: Build
FROM golang:1.23-bullseye AS build

# Set the working directory
WORKDIR /src

# Install build dependencies
RUN apt-get update && apt-get install -y build-essential git

# COPY the repository
COPY . .

# Build the application
RUN GOARCH="amd64" GOOS="linux" make build_default

# Stage 2: Create the final image
FROM gcr.io/distroless/base

# Copy the built application from the build stage
COPY --from=build /src/_output/yarr /app/yarr

# Set the volume
VOLUME /app/data

# Set the entry point and default command
ENTRYPOINT ["/app/yarr"]
CMD ["-addr", "0.0.0.0:7070", "-db", "/app/data/yarr.db"]
