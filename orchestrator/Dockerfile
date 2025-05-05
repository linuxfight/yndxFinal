ARG GO_VERSION=1.24
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

COPY . .
# RUN go install github.com/swaggo/swag/cmd/swag@latest
# RUN swag init -g cmd/main.go
RUN find ./docs -type f -exec sed -i 's/localhost:8080/calc.lxft.tech/g' {} +

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/main.go

FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        curl \
        && \
        update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=build /src/docs /docs

COPY --from=build /bin/server /bin/

EXPOSE 8080
EXPOSE 9090

ENTRYPOINT [ "/bin/server" ]