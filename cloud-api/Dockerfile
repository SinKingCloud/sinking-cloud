FROM golang as sinking-cloud-api-builder

#作者
MAINTAINER sinkingcloud 1178710004@qq.com

#修改时区
RUN echo "Asia/Shanghai" > /etc/timezone && rm /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

#修改环境变量
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOPROXY=https://goproxy.cn

#复制代码
COPY . /data/cloud-api

#指定目录
WORKDIR /data/cloud-api

#安装依赖
RUN go mod tidy

#编译
RUN go build

FROM alpine as sinking-cloud-api-run

# 设置工作目录
WORKDIR /data/cloud-api

# 复制文件
COPY --from=sinking-cloud-api-builder /data/cloud-api/cloud-api .
COPY --from=sinking-cloud-api-builder /data/cloud-api/config/ ./config/

#映射端口
EXPOSE 1818

#启动
ENTRYPOINT  ["./cloud-api"]