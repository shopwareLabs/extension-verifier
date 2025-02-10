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

FROM base AS phpstan

COPY tools/phpstan /phpstan
WORKDIR /phpstan
RUN composer install

FROM base AS php-cs-fixer
COPY tools/php-cs-fixer /php-cs-fixer
WORKDIR /php-cs-fixer
RUN composer install

FROM base AS eslint
COPY tools/eslint /eslint
WORKDIR /eslint
RUN npm install

FROM base AS stylelint
COPY tools/stylelint /stylelint
WORKDIR /stylelint
RUN npm install

FROM golang:alpine AS executor

COPY go.* /app/
COPY *.go /app/
COPY configs /app/configs
COPY internal /app/internal

WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o /app/executor -ldflags "-s -w" .

FROM base AS final
WORKDIR /opt/

COPY --from=phpstan /phpstan /opt/tools/phpstan
COPY --from=eslint /eslint /opt/tools/eslint
COPY --from=stylelint /stylelint /opt/tools/stylelint
COPY --from=php-cs-fixer /php-cs-fixer /opt/tools/php-cs-fixer
COPY --from=executor /app/executor /opt/executor

ENTRYPOINT ["/opt/executor"]
