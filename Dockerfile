# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24.2 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /app .

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot
ENV PORT=8080
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app","service"]
