FROM ghcr.io/shyim/wolfi-php/base:latest AS base

RUN <<EOF
set -eo pipefail
apk add --no-cache \
    php-8.3 \
    php-8.3-bz2 \
    php-8.3-curl \
    php-8.3-gd \
    php-8.3-gmp \
    php-8.3-ldap \
    php-8.3-mysqlnd \
    php-8.3-openssl \
    php-8.3-pdo_mysql \
    php-8.3-soap \
    php-8.3-sodium \
    php-8.3-exif \
    php-8.3-gettext \
    php-8.3-intl \
    php-8.3-mbstring \
    php-8.3-opcache \
    php-8.3-pcntl \
    php-8.3-pdo \
    php-8.3-phar \
    php-8.3-sockets \
    php-8.3-bcmath \
    php-8.3-ctype \
    php-8.3-iconv \
    php-8.3-dom \
    php-8.3-posix \
    php-8.3-simplexml \
    php-8.3-xml \
    php-8.3-xmlreader \
    php-8.3-xmlwriter \
    php-8.3-fileinfo \
    php-8.3-zip \
    php-8.3-ffi \
    php-8.3-ftp \
    composer \
    nodejs-23 \
    npm

    npm install --location=global @biomejs/biome
EOF

FROM --platform=$BUILDPLATFORM base AS php
COPY tools/php /php
WORKDIR /php
RUN composer install

FROM base AS js
COPY tools/js /js
WORKDIR /js
RUN npm install

FROM --platform=$BUILDPLATFORM golang:alpine AS executor

COPY go.* /app/
COPY *.go /app/
COPY internal /app/internal

WORKDIR /app
RUN go mod download
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o /app/executor -ldflags "-s -w" .

FROM base AS final
WORKDIR /opt/

COPY --from=php /php /opt/tools/php
COPY --from=js /js /opt/tools/js
COPY --from=executor /app/executor /opt/executor

ENTRYPOINT ["/opt/executor"]
