FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go build -o ./main ./src

FROM alpine:latest

WORKDIR /app

COPY --from=0 /app/main /app/main

EXPOSE 8080

CMD [ "/app/main" ]