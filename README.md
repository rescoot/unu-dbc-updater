# Silent Running

A Go utility for managing DBC (Dashboard Controller) updates on unu electric scooters.

## Overview

This tool controls the dashboard update process by:

1. Stopping vehicle services
2. Controlling dashboard power via GPIO
3. Managing BLE pin-codes
4. Coordinating with the update mechanism via lock files
5. Monitoring dashboard state via Redis
6. Restarting vehicle services after update

## Features

- GPIO-based power control for the dashboard
- Redis integration for dashboard state monitoring
- BLE pin-code management
- Lock file mechanism for update coordination
- Configurable update timeout

## Usage

```bash
# Run with default timeout (30 minutes)
./dbc-updater

# Run with custom timeout
./dbc-updater --update-timeout=1h
```

## Build

```bash
# Build the binary
make

# or build only for arm
make target
```

## License

[Creative Commons Attribution-NonCommercial 4.0 International (CC-BY-NC-4.0)](LICENSE)

## Notes

- Requires root permissions to control GPIO
- Expects Redis server at 192.168.7.1:6379