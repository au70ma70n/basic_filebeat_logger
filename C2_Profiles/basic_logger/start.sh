#!/bin/sh

# Start Filebeat in the background
echo "Starting Filebeat..."
/usr/share/filebeat/filebeat -c /Mythic/filebeat_mythic_redelk.yml &

# Start the Mythic logger
echo "Starting Mythic logger..."
make run
