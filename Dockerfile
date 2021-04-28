# builder image
FROM golang:1.16-alpine as builder
WORKDIR /root/app
COPY . /root/app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o kratos .


# generate clean, final image for end users
FROM alpine:3.13.5
WORKDIR /root/app
ENV GOPATH /root/app
COPY --from=builder /root/app .
# TODO: find a way to set the app to read those files from /root/app and not from / folder
RUN ln -s /root/app/persistence /persistence
RUN ln -s /root/app/selfservice /selfservice
RUN ln -s /root/app/schema /schema
RUN ln -s /root/app/session /session
RUN ln -s /root/app/identity /identity

ENTRYPOINT ["./kratos"]

CMD ["serve", "--config", "/etc/config/kratos/.kratos.yaml", "--dev"]

EXPOSE 4433
EXPOSE 4434
