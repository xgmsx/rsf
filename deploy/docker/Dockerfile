ARG GOLANG_VERSION=1.24
ARG ALPINE_VERSION=3.21

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder
ENV APP=/app
WORKDIR $APP

# Copy project files
COPY . .
ARG SERVICE_PATH
WORKDIR $APP/${SERVICE_PATH}

# Make binary files
RUN mkdir -p $APP/bin && \
    go mod download && \
    go build -v -o $APP/bin/app ./cmd/main.go

FROM alpine:${ALPINE_VERSION} AS rutime
ENV APP=/app
WORKDIR $APP

# Copy project binaries from builder
ENV PATH="$APP:$PATH"
COPY --from=builder $APP/bin $APP
COPY --from=builder $APP/shared/api $APP/api

# Create app user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && chown appuser:appgroup $APP
USER appuser

ENTRYPOINT ["app"]
