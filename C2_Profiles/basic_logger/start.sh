#!/bin/sh

# Start Filebeat in the background
echo "Starting Filebeat..."
/usr/share/filebeat/filebeat -c /tmp/filebeat.yml &

# Start the Mythic logger
echo "Starting Mythic logger..."
make run
