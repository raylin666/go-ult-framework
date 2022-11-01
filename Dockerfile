FROM raylin666/golang:1.17-wire AS builder

COPY . /go
WORKDIR /go

EXPOSE 10001

CMD ["/bin/bash", "-c", "make init && make generate && make run"]
