# 构建阶段
FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/gin_kubelet main.go
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gin_kubelet main.go


# 运行阶段
FROM alpine:3.20

ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /bin/gin_kubelet .

# 设置时区
RUN apk add --no-cache tzdata curl wget iputils \
   && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
   && echo $TZ > /etc/timezone \
   && echo "Asia/Shanghai" > /etc/timezone

# 给应用添加所有权限
RUN chmod 775 /app/gin_kubelet

