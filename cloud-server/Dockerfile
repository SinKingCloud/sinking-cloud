FROM golang:1.19 as sinking-builder

#作者
MAINTAINER sinkingcloud 1178710004@qq.com

#修改环境变量
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
ENV GOPRIVATE=github.com/SinKingCloud
RUN git config --global http.proxy socks5://admin:123456@154.204.60.224:8787
RUN git config --global https.proxy socks5://admin:123456@154.204.60.224:8787
RUN git config --global url."https://SinKingCloud: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@github.com".insteadOf "https://github.com"
RUN apt install apt-transport-https

#复制代码
COPY . /data/sinking

#################构建网关服务####################
WORKDIR /data/sinking/cloud-gateway
RUN go mod tidy
RUN go build
#################构建推送服务####################
WORKDIR /data/sinking/cloud-websocket
RUN go mod tidy
RUN go build
#################构建授权服务####################
WORKDIR /data/sinking/cloud-service/service-auth
RUN go mod tidy
RUN go build
#################构建后台服务####################
WORKDIR /data/sinking/cloud-service/service-admin
RUN go mod tidy
RUN go build
#################构建用户服务####################
WORKDIR /data/sinking/cloud-service/service-user
RUN go mod tidy
RUN go build
#################构建开放服务####################
WORKDIR /data/sinking/cloud-service/service-open
RUN go mod tidy
RUN go build
#################构建爬虫服务####################
WORKDIR /data/sinking/cloud-service/service-spider
RUN go mod tidy
RUN go build

#################网关服务镜像####################
FROM alpine as cloud-gateway
RUN echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories && apk update && apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /data/sinking
COPY --from=sinking-builder /data/sinking/cloud-gateway/cloud-gateway .
COPY --from=sinking-builder /data/sinking/cloud-gateway/config/ ./config/
EXPOSE 1000
ENTRYPOINT  ["./cloud-gateway"]

#################推送服务镜像####################
FROM alpine as cloud-websocket
RUN echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories && apk update && apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /data/sinking
COPY --from=sinking-builder /data/sinking/cloud-websocket/cloud-websocket .
COPY --from=sinking-builder /data/sinking/cloud-websocket/config/ ./config/
EXPOSE 1005
ENTRYPOINT  ["./cloud-websocket"]

#################用户服务镜像####################
FROM alpine as service-user
RUN echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories && apk update && apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /data/sinking
COPY --from=sinking-builder /data/sinking/cloud-service/service-user/service-user .
COPY --from=sinking-builder /data/sinking/cloud-service/service-user/config/ ./config/
EXPOSE 1002
ENTRYPOINT  ["./service-user"]

#################后台服务镜像####################
FROM alpine as service-admin
RUN echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories && apk update && apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /data/sinking
COPY --from=sinking-builder /data/sinking/cloud-service/service-admin/service-admin .
COPY --from=sinking-builder /data/sinking/cloud-service/service-admin/config/ ./config/
EXPOSE 1001
ENTRYPOINT  ["./service-admin"]

#################授权服务镜像####################
FROM alpine as service-auth
RUN echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories && apk update && apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /data/sinking
COPY --from=sinking-builder /data/sinking/cloud-service/service-auth/service-auth .
COPY --from=sinking-builder /data/sinking/cloud-service/service-auth/config/ ./config/
EXPOSE 1003
ENTRYPOINT  ["./service-auth"]

#################开放服务镜像####################
FROM alpine as service-open
RUN echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories && apk update && apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /data/sinking
COPY --from=sinking-builder /data/sinking/cloud-service/service-open/service-open .
COPY --from=sinking-builder /data/sinking/cloud-service/service-open/config/ ./config/
EXPOSE 1004
ENTRYPOINT  ["./service-open"]

#################爬虫服务镜像####################
#FROM centos:centos7.9.2009 as service-spider
#ENV TIME_ZONE Asia/Shanghai
#WORKDIR /data/sinking
#COPY --from=sinking-builder /data/sinking/cloud-service/service-spider/service-spider .
#COPY --from=sinking-builder /data/sinking/cloud-service/service-spider/config/ ./config/
#RUN yum install wget -y && wget https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm && yum localinstall -y google-chrome-stable_current_x86_64.rpm
#EXPOSE 1006
#ENTRYPOINT  ["./service-spider"]

