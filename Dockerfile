# build stage
FROM golang:1.15.3-alpine3.12 as builder

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o postcodes .

# final stage
FROM alpine:3.2
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/postcodes /app/
COPY --from=builder /app/area/file /area/file
ENTRYPOINT ["/app/postcodes"]