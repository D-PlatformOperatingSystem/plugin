#!/usr/bin/env bash
/root/dplatformos -f /root/dplatformos.toml &
# to wait nginx start
sleep 15
/root/dplatformos -f "$PARAFILE"
