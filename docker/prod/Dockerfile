FROM golang:latest as builder
LABEL maintainer="Ilya Romanov <romanov.i.u@yandex.ru>"
WORKDIR /app
COPY go.mod /app
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app/main.go

FROM alpine:latest
LABEL maintainer="Ilya Romanov <romanov.i.u@yandex.ru>"
RUN apk --no-cache add ca-certificates && mkdir /app
WORKDIR /app
COPY --from=builder /app/main .
COPY ./output ./output

CMD ["./main"]