FROM --platform=amd64 centos:centos7

RUN mkdir -p /app
COPY ./bin/xuanwu-agent /app/xuanwu-agent
COPY ./config /app/config