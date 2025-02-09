FROM golang:1.23-alpine

WORKDIR /app
COPY .. .

RUN CGO_ENABLED=0 GOOS=linux
RUN go build -mod=vendor -o main .

EXPOSE 8080
CMD ["./main"]