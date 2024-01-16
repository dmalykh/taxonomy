FROM golang:1.21 as builder
WORKDIR /build
COPY ./src /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o taxonomy ./

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/taxonomy .
RUN ln -s /app/taxonomy /usr/bin/taxonomy
