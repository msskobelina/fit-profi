# -------- build stage --------
FROM golang:1.25.0 AS build
ENV CGO_ENABLED=0 GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -buildvcs=false -o /out/app .

# -------- run stage --------
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H -s /sbin/nologin appuser

WORKDIR /srv
COPY --from=build /out/app /srv/app

ENV HTTP_PORT=8080
EXPOSE 8080
USER appuser

CMD ["/srv/app"]
