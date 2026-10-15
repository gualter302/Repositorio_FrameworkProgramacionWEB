FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /servidor .

FROM alpine:3.23
COPY --from=builder /servidor /servidor
EXPOSE 3000
CMD ["/servidor"]
