FROM golang:1.16.10-alpine3.14 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build


FROM alpine:3.14 as output
COPY --from=builder /app/operator /usr/local/bin/operator

CMD [ "operator" ]