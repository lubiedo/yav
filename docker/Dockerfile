FROM golang:1.16-alpine

ARG PORT=8080
ARG ADDR=0.0.0.0
ENV PORT=${PORT}
ENV ADDR=${ADDR}

ADD . /opt/yav
WORKDIR /opt/yav/src

RUN go build -o ../yav

WORKDIR /opt/yav
ENTRYPOINT /opt/yav/yav -use-https=false -verbose -port $PORT -addr $ADDR
