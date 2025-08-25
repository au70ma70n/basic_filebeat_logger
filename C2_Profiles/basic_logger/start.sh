#!/bin/sh

# Start Filebeat in the background
echo "Starting Filebeat..."
filebeat -c /usr/share/filebeat/filebeat.yml &

# Start the Mythic logger
echo "Starting Mythic logger..."
make run
