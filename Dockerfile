
FROM golang:1.16
WORKDIR /go/src/
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/app ./
CMD ["./app"]  
