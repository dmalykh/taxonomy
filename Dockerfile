FROM golang:1.20 as builder
WORKDIR /build
COPY ./src /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o tagservice ./

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/tagservice .
RUN ln -s /app/tagservice /usr/bin/tagservice
