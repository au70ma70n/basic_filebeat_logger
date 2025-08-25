FROM golang:1.23 AS builder

WORKDIR /Mythic/

COPY [".", "."]

RUN make build

FROM alpine

RUN apk add --update make curl

# Install Filebeat - download and install the correct binary
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-7.17.9-linux-x86_64.tar.gz && \
    tar xzvf filebeat-7.17.9-linux-x86_64.tar.gz && \
    mv filebeat-7.17.9-linux-x86_64 /usr/share/filebeat && \
    rm filebeat-7.17.9-linux-x86_64.tar.gz && \
    # Install required dependencies for the binary
    apk add --no-cache libc6-compat

COPY --from=builder /main /main

WORKDIR /Mythic/

COPY [".", "."]

# Create log directory for filebeat integration
RUN mkdir -p /var/log/mythic && chmod 755 /var/log/mythic

# Copy Filebeat configuration
COPY filebeat_mythic_redelk.yml /usr/share/filebeat/filebeat.yml

# Make the startup script executable
RUN chmod +x /Mythic/start.sh
RUN chmod +x start.sh


ENTRYPOINT ["/Mythic/start.sh"]