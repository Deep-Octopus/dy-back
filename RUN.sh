#!/bin/bash

# 删除容器
docker rm -f dy

# 删除镜像
docker rmi -f dy

# 构建镜像
docker build -t dy .

# 运行容器
docker run -d -p 8082:8082 --name dy dy
