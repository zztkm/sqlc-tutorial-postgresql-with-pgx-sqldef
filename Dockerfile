FROM golang:1.21-alpine as build
ARG TARGETARCH
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go version && \
    go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags "-w -s" -trimpath -o server cmd/server/main.go

FROM scratch
COPY --from=build /app/server /go/bin/server
ENTRYPOINT ["/go/bin/server"]
