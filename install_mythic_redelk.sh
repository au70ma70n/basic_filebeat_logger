#!/bin/bash
#
# Mythic RedELK Integration Setup Script
# This script helps set up the basic_filebeat_logger with RedELK
#
# Usage: ./install_mythic_redelk.sh <mythic-instance-name> <attack-scenario> <redelk-host:port>
#
# Example: ./install_mythic_redelk.sh mythic-server-1 "red-team-2024" "192.168.1.100:5045"

if [ $# -ne 3 ]; then
    echo "Usage: $0 <mythic-instance-name> <attack-scenario> <redelk-host:port>"
    echo "Example: $0 mythic-server-1 \"red-team-2024\" \"192.168.1.100:5045\""
    exit 1
fi

MYTHIC_HOSTNAME=$1
ATTACK_SCENARIO=$2
REDELK_HOSTPORT=$3

echo "Setting up Mythic RedELK integration..."
echo "Mythic Hostname: $MYTHIC_HOSTNAME"
echo "Attack Scenario: $ATTACK_SCENARIO"
echo "RedELK Host:Port: $REDELK_HOSTPORT"

# Create the filebeat configuration
echo "Creating Filebeat configuration..."
sed -e "s/@@HOSTNAME@@/$MYTHIC_HOSTNAME/g" \
    -e "s/@@ATTACKSCENARIO@@/$ATTACK_SCENARIO/g" \
    -e "s/@@HOSTANDPORT@@/$REDELK_HOSTPORT/g" \
    filebeat_mythic_redelk.yml > filebeat_mythic_redelk_configured.yml

echo "Configuration created: filebeat_mythic_redelk_configured.yml"
echo ""
echo "Next steps:"
echo "1. Install this C2 profile in Mythic:"
echo "   sudo ./mythic-cli install folder $(pwd)"
echo ""
echo "2. Copy the Filebeat configuration to your Mythic server:"
echo "   cp filebeat_mythic_redelk_configured.yml /etc/filebeat/filebeat.yml"
echo ""
echo "3. Install and start Filebeat on your Mythic server"
echo "4. Start the C2 profile in Mythic:"
echo "   sudo ./mythic-cli c2 start basic_logger"
echo ""
echo "The logs will be sent to RedELK on port 5045 (non-TLS)"
