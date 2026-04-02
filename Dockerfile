FROM golang:alpine AS builder
WORKDIR /app/my-qqbot
RUN apk --no-cache add ca-certificates  && update-ca-certificates
RUN apk --no-cache add tzdata
COPY . .
RUN go build -o main


FROM scratch AS runtime
LABEL authors="Yoake"
WORKDIR /app
ENV TZ=UTC+8
COPY --from=builder /app/my-qqbot/config.example.yaml config.yaml
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/my-qqbot/main main
CMD ["./main"]