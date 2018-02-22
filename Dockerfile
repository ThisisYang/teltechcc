From golang:1.9.3 as builder
WORKDIR /go/src/teltechcc/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o teltechcc .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
EXPOSE 8000
COPY --from=builder /go/src/teltechcc/teltechcc .

ENTRYPOINT ["./teltechcc"]
