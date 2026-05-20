FROM ubuntu:24.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    curl \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# mcp[cli] brings in fastmcp; --break-system-packages is required on Ubuntu 24.04
RUN pip3 install "mcp[cli]" --break-system-packages

COPY test_add.sh /usr/local/bin/test_add.sh
RUN chmod +x /usr/local/bin/test_add.sh

CMD ["/usr/local/bin/test_add.sh"]
