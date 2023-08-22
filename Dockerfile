FROM golang:1.20.7-bookworm

RUN ARCH="$(arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/)" \
    && wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-3/wkhtmltox_0.12.6.1-3.bookworm_${ARCH}.deb \
    && apt update \
    && apt install -y ./wkhtmltox_0.12.6.1-3.bookworm_${ARCH}.deb

