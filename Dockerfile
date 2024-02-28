# Build Backend
FROM --platform=$BUILDPLATFORM golang:1.22.0-alpine3.19 as builder

ARG TARGETOS
ARG TARGETARCH

ENV GOCACHE=/root/.cache/go-build

# Install build dependencies
RUN apk add --no-cache git make build-base
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomod,target="/go/pkg/mod" go mod download

COPY . ./

# Build Cache
# https://dev.to/jacktt/20x-faster-golang-docker-builds-289n
RUN --mount=type=cache,id=gomod,target="/go/pkg/mod" \
    --mount=type=cache,id=gobuild,target="/root/.cache/go-build" \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o FiberReactTest .

# Build Frontend
FROM --platform=$BUILDPLATFORM node:20-alpine3.19 as frontend-builder

RUN mkdir /frontend
WORKDIR /frontend

COPY ./frontend/package.json ./
COPY ./frontend/package-lock.json ./
RUN --mount=type=cache,id=npmmod,target="/root/.npm" npm ci

COPY ./frontend ./
RUN npm run build

# Frontend Development Image
FROM frontend-builder as frontend-dev
RUN apk update && apk add --no-cache optipng
WORKDIR /
COPY --from=builder /app/FiberReactTest /
COPY --from=builder /app/start_dev.sh /
EXPOSE 3000
EXPOSE 5173
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["/start_dev.sh"]


# Production Image
#FROM alpine:3.19 as production
#RUN apk update && apk add --no-cache optipng
#COPY --from=builder /app/FiberReactTest /
#RUN mkdir -p /frontend/dist
#COPY --from=frontend-builder /frontend/dist /frontend/dist
#
#EXPOSE 8080
#CMD ["./FiberReactTest"]