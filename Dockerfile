FROM golang:1.21 AS backend

WORKDIR /builder
COPY ./ /builder

RUN CGO_ENABLED=0 go build -o ./dist/bin/data-clone ./cmd/data-clone/main.go

FROM alpine:3.14.0

WORKDIR /app

COPY --from=backend /builder/dist /app/

ENTRYPOINT /app/bin/data-clone
