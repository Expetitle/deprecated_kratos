# builder image
FROM golang:1.16-alpine as builder
WORKDIR /build
COPY . /build/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o kratos .


# generate clean, final image for end users
FROM alpine:3.13.5
WORKDIR /root/app
COPY --from=builder /build/kratos .

ENTRYPOINT ["/root/app/kratos"]

CMD ["serve", "--config", "/etc/config/kratos/.kratos.yaml", "--dev"]

EXPOSE 4433
EXPOSE 4434
