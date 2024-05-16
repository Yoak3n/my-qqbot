FROM golang:alpine as builder
WORKDIR /app/my-qqbot
COPY . .
RUN go build -o main


FROM scratch as runtime
LABEL authors="Yoake"
WORKDIR /app
COPY --from=builder /app/my-qqbot/config.example.yaml config.yaml
COPY --from=builder /app/my-qqbot/main main
CMD ["./main"]