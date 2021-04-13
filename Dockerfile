FROM golang:1.16-alpine

WORKDIR /root/app

COPY . /root/app/

RUN cd /root/app && go build

ENTRYPOINT ["/root/app/kratos"]

CMD ["serve", "--config", "/etc/config/kratos/.kratos.yaml", "--dev"]

EXPOSE 4433
EXPOSE 4434

