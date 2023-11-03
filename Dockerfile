FROM golang:alpine as builder
RUN apk add --no-cache unzip
WORKDIR /go/src/

COPY . .
RUN go get .
RUN CGO_ENABLED=0 go build -o /go/bin/etlwakatime .

RUN ARCH=$([ $(uname -m) != "aarch64" ] && echo "amd64" || echo "aarch64") && \
  wget "https://github.com/duckdb/duckdb/releases/download/v0.9.1/duckdb_cli-linux-$ARCH.zip" -O "/opt/duckdb_cli-linux-$ARCH.zip" && \
  unzip -d /opt "/opt/duckdb_cli-linux-$ARCH.zip"

FROM postgres:16-alpine
WORKDIR /app

COPY --from=builder /go/bin/etlwakatime ./
COPY --from=builder /opt/duckdb /usr/local/sbin/duckdb

ENV PGHOST=
ENV PGPORT=
ENV PGDATABASE=
ENV PGUSER=postgres
ENV PGPASSWORD=
ENV PGSSLMODE=disable
ENV PGAPPNAME=etl-wakatime


ENTRYPOINT [ "/app/etlwakatime", "--help"]
