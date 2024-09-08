# Start from golang base image
FROM golang:1.22.6-alpine3.20 as builder

# Install necessary packages
RUN apk update && apk add --no-cache \
    git python3 py3-pip bash \
    build-base linux-headers \
    libcap-dev asciidoc gcc make pkgconfig \
    eudev-dev \
    sudo \
    libcap

# Set the current working directory inside the container
WORKDIR /app

# Install CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon@latest

# Clone isolate
RUN git clone https://github.com/ioi/isolate.git /isolate

# Modify isolate to work without systemd
RUN sed -i 's/#include <systemd\/sd-daemon.h>/\/\/#include <systemd\/sd-daemon.h>/' /isolate/isolate-cg-keeper.c && \
    sed -i 's/sd_notify/\/\/sd_notify/' /isolate/isolate-cg-keeper.c

# Build and install isolate
WORKDIR /isolate
RUN make && make install

# Set up isolate
RUN addgroup -S isolate && adduser -S -G isolate isolate \
    && mkdir -p /var/local/lib/isolate \
    && chown -R isolate:isolate /var/local/lib/isolate \
    && chmod 777 /var/local/lib/isolate

# Set capabilities for isolate
RUN setcap cap_sys_admin,cap_sys_chroot,cap_mknod,cap_net_admin+ep /usr/local/bin/isolate

# Return to the app directory
WORKDIR /app

# Copy your Go application files
COPY . .

# Command to run the executable
CMD CompileDaemon --build="go build main.go" --command="./main" --color