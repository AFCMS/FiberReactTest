# syntax=docker/dockerfile:1

# Build Backend
FROM --platform=$BUILDPLATFORM golang:1.22.0-alpine3.19 as builder

ARG TARGETOS
ARG TARGETARCH

ENV GOCACHE=/root/.cache/go-build

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomod,target="/go/pkg/mod" go mod download
RUN --mount=type=cache,id=gomod,target="/go/pkg/mod" go mod verify

COPY . ./

# Build Cache
# https://dev.to/jacktt/20x-faster-golang-docker-builds-289n
RUN --mount=type=cache,id=gomod,target="/go/pkg/mod" \
    --mount=type=cache,id=gobuild,target="/root/.cache/go-build" \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o FiberReactTest .

# Build Frontend
FROM --platform=$BUILDPLATFORM node:20-alpine3.19 as frontend-builder

RUN mkdir /frontend
WORKDIR /frontend

COPY ./frontend/package.json ./
COPY ./frontend/package-lock.json ./
RUN --mount=type=cache,id=npmmod,target="/root/.npm" npm ci

COPY ./frontend ./
RUN npm run build

# Production Image
FROM alpine:3.19 as production
RUN apk update && apk add --no-cache optipng
COPY --from=builder /app/FiberReactTest /
COPY --from=builder /app/index.html /
RUN mkdir -p /frontend/dist
COPY --from=frontend-builder /frontend/dist /frontend/dist

EXPOSE 8080
CMD ["./FiberReactTest"]