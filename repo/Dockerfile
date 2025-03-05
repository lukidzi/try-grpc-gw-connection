FROM ubuntu:24.04

ARG TARGETARCH

RUN apt-get update && apt-get install -y \
    curl \
    vim \
    iptables \
    iproute2 \
    net-tools \
    iputils-ping && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY main-${TARGETARCH} /app/main
RUN chmod +x /app/main

CMD ["/bin/bash"]