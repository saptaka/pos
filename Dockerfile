
FROM golang:1.16-alpine

WORKDIR /app

COPY . ./

RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

RUN apk --no-cache add ca-certificates

EXPOSE 3030

CMD ["./app"]  

