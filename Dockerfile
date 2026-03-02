# build frontend
FROM node AS web_image

# 华为源
# RUN npm config set registry https://repo.huaweicloud.com/repository/npm/

RUN npm install pnpm -g

WORKDIR /build

COPY ./package.json /build

COPY ./pnpm-lock.yaml /build

RUN pnpm install

COPY . /build

RUN pnpm run build

# build backend
# 最新alpine3.19导致sqlite3编译失败(https://github.com/mattn/go-sqlite3/issues/1164，
# 临时解决方案:https://github.com/mattn/go-sqlite3/pull/1177)
# sun-panel暂时解决方案使用golang:1.21-alpine3.18（因旧版本使用没问题，短期内较稳定） 
FROM golang:1.21-alpine3.18 as server_image

WORKDIR /build

# 仅复制 Go 后端相关目录（Go 项目已上移到仓库根目录）
COPY ./go.mod ./go.sum ./main.go ./
COPY ./api ./api
COPY ./assets ./assets
COPY ./conf ./conf
COPY ./global ./global
COPY ./initialize ./initialize
COPY ./lang ./lang
COPY ./lib ./lib
COPY ./models ./models
COPY ./router ./router
COPY ./runtime ./runtime
COPY ./structs ./structs
COPY ./uploads ./uploads

# 中国国内源
# RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories \
#     && go env -w GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache bash curl gcc git musl-dev

RUN go env -w GO111MODULE=on \
    && go build -o ange-panel --ldflags="-X sun-panel/global.RUNCODE=release -X sun-panel/global.ISDOCKER=docker" ./main.go


# run_image
FROM alpine

WORKDIR /app

COPY --from=web_image /build/dist /app/web

COPY --from=server_image /build/ange-panel /app/ange-panel

# 中国国内源
# RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories

EXPOSE 3002

RUN apk add --no-cache bash ca-certificates su-exec tzdata \
    && chmod +x ./ange-panel \
    && ./ange-panel -config

CMD ./ange-panel
