# 第一阶段：构建 (Node.js)
FROM node:18-alpine as builder

WORKDIR /app

# 单独复制 package.json 安装依赖，利用 Docker 缓存
COPY package.json package-lock.json ./
# 配置国内淘宝镜像加速
RUN npm config set registry https://registry.npmmirror.com
RUN npm install

# 复制所有源代码
COPY . .

# 执行构建，生成 dist 目录
RUN npm run build

# 第二阶段：运行 (Nginx)
FROM nginx:alpine

# 复制我们在 frontend 目录下准备好的 customized nginx.conf
# 注意：这一步需要我们在本地先创建好 frontend/nginx.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 从构建阶段复制 dist 产物到 Nginx 默认目录
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
