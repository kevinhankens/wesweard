FROM golang:1.5-alpine

RUN mkdir /wesweard

COPY . /wesweard

RUN apk add --no-cache --virtual .build-deps \
        git \
        bash && \
    cd /wesweard && \
    go get github.com/constabulary/gb/... && \
    gb build
