FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /gateway ./cmd/gateway
FROM gcr.io/distroless/static
COPY --from=build /gateway /gateway
ENTRYPOINT ["/gateway"]
