# 基础镜像
FROM ubuntu:20.04

# 安装必要依赖（可选，看你程序是否需要）
# RUN apt-get update && apt-get install -y ca-certificates

# 拷贝可执行文件到镜像中
COPY webook /app/webook

# 设置工作目录
WORKDIR /app

# 设置容器启动命令
ENTRYPOINT ["./webook"]
