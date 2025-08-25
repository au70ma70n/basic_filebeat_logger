FROM golang:1.23 AS builder

WORKDIR /Mythic/

COPY [".", "."]

RUN make build

FROM alpine

RUN apk add --update make curl filebeat

# Create filebeat directory structure
RUN mkdir -p /usr/share/filebeat

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