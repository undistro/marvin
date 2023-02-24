FROM golang:1.19 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o marvin main.go

FROM alpine:3.17.2

RUN addgroup -g 8494 -S nonroot && adduser -u 8494 -D -S nonroot -G nonroot
USER 8494:8494

WORKDIR /
COPY --from=builder /workspace/marvin .

ENTRYPOINT ["/marvin"]
