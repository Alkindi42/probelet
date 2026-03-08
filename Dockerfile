ARG ALPINE_VERSION=3.23
ARG GOLANG_VERSION=1.25.6

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
  go build -trimpath -ldflags="-s -w \
  -X github.com/Alkindi42/probelet/internal/version.Version=$VERSION \
  -X github.com/Alkindi42/probelet/internal/version.Commit=$COMMIT \
  -X github.com/Alkindi42/probelet/internal/version.BuildDate=$BUILD_DATE" \
  -o /out/probelet ./main.go

FROM alpine:${ALPINE_VERSION} as runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates wget \
  && addgroup -S app \
  && adduser -S app -G app

USER app

COPY --from=build /out/probelet .

ENTRYPOINT ["/app/probelet"]
CMD ["serve"]

HEALTHCHECK --interval=30s --timeout=2s --retries=3 \
  CMD wget -qO- http://localhost:8000/healthz || exit 1
