FROM --platform=amd64 centos:centos7

RUN mkdir -p /root/app

WORKDIR /root/app
COPY bin/xuanwu-agent ./
COPY ./config ./config