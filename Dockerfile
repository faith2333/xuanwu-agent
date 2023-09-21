FROM --platform=amd64 centos:centos7

COPY ./bin/xuanwu-agent /app/xuanwu-agent
COPY ./config /app/config