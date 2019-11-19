FROM golang:1.13 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -a -o k8s-dashboard-auth-proxy .


FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/k8s-dashboard-auth-proxy /bin/

EXPOSE 9443 9999

ENTRYPOINT ["/bin/k8s-dashboard-auth-proxy"]
