FROM golang as sinking-builder

#作者
MAINTAINER sinkingcloud 1178710004@qq.com

#修改环境变量
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOPROXY=https://goproxy.cn

#复制代码
COPY . /data/sinking

#################构建爬虫服务####################
#指定目录
WORKDIR /data/sinking/cloud-service/service-spider
#安装依赖
RUN go mod tidy
#编译
RUN go build

FROM alpine as sinking-spider-run
#修改时区
RUN echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories && apk update && apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
# 设置工作目录
WORKDIR /data/sinking
# 复制文件
COPY --from=sinking-builder /data/sinking/cloud-service/service-spider/service-spider .
COPY --from=sinking-builder /data/sinking/cloud-service/service-spider/config/ ./config/
#映射端口
EXPOSE 1004
#启动
ENTRYPOINT  ["./service-spider"]

