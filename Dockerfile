# Copyright 2023 Undistro Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.22 as builder
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT
ARG DATE

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w -X github.com/undistro/marvin/pkg/version.version=${VERSION:-docker} \
    -X github.com/undistro/marvin/pkg/version.commit=${COMMIT} \
    -X github.com/undistro/marvin/pkg/version.date=${DATE}" -a -o marvin main.go

FROM alpine:3.19.1
RUN apk upgrade && rm /var/cache/apk/*
RUN addgroup -g 8494 -S nonroot && adduser -u 8494 -D -S nonroot -G nonroot
USER 8494:8494

WORKDIR /
COPY --from=builder /workspace/marvin .

ENTRYPOINT ["/marvin"]
