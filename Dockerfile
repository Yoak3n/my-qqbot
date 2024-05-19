FROM golang:alpine as builder
WORKDIR /app/my-qqbot
RUN apk --no-cache add ca-certificates  && update-ca-certificates
COPY . .
RUN go build -o main


FROM scratch as runtime
LABEL authors="Yoake"
WORKDIR /app
COPY --from=builder /app/my-qqbot/config.example.yaml config.yaml
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=builder /app/my-qqbot/main main
CMD ["./main"]