FROM golang:1.23.5-alpine3.21 AS builder

COPY . /rest-wallets
WORKDIR /rest-wallets

RUN go mod download
RUN go build -o ./bin/app cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /rest-wallets/bin/app .
COPY --from=0 /rest-wallets/configs configs/
COPY --from=0 /rest-wallets/config.env config.env
COPY --from=builder /rest-wallets/wait-for-postgres.sh wait-for-postgres.sh

RUN apk update && apk add --no-cache postgresql-client

RUN chmod +x wait-for-postgres.sh

CMD ["./app"]