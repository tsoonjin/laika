FROM golang:1.14 AS builder
ADD . /app
WORKDIR /app
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -a -o /main ./cmd/sebastian
# final stage
FROM alpine:3.10
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]
EXPOSE 3000
