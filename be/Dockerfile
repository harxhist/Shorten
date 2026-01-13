# Build the application from source
FROM golang:1.23.4 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . ./

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /sh10

# Use a minimal base image, e.g., distroless or alpine
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /sh10 /sh10
COPY --from=build-stage /app/resource /resource

EXPOSE 8080

ENTRYPOINT ["/sh10"]