FROM golang:1.21 AS backend

WORKDIR /builder
COPY ./ /builder

RUN go build -o ./dist/bin/clickhouse-data-clone ./cmd/data-clone/main.go

FROM alpine:3.14.0

WORKDIR /app

COPY --from=backend /builder/dist /app/

CMD /app/bin/clickhouse-data-clone
