FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/nerdgate-hub ./cmd/nerdgate-hub

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=build /out/nerdgate-hub /usr/local/bin/nerdgate-hub

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/nerdgate-hub"]
