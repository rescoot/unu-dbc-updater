# DBC Silent Running

A tool for updating the DBC (Dashboard Controller) without putting the scooter in ready-to-drive mode.

## Problem

To update the DBC, we need to turn it on, but that unfortunately also makes the scooter ready to drive, allows opening the seatbox, etc. This is due to two problems:

1. We can only turn on the DBC by sending the scooter an unlock command
2. Whenever the DBC turns on, it signals `dashboard ready true` which lets the scooter drive

## Solution

This tool works around these issues by:

1. Stopping the vehicle service
2. Manually controlling the DBC power via GPIO
3. Monitoring Redis for the dashboard ready signal and immediately setting it to false
4. Setting the BLE pin-code to "UPDATE" during the update process
5. Creating a lock file that external processes can remove when the update is complete
6. Waiting for the lock file to be deleted or timing out
7. Properly restarting the vehicle service when done

## Usage

```
./dbc-updater [options]
```

### Options

- `-update-timeout duration`: Maximum time to wait for update to complete (default 30m)

### Example

```
./dbc-updater -update-timeout 15m
```

## Building

A Makefile is provided with the following targets:

- `make`: Build for both host and target
- `make host`: Build for the host system
- `make target`: Build for the armv7l target with size optimizations
- `make clean`: Clean up build artifacts

## Integration with Update Process

The tool creates a lock file at `/tmp/dbc-update.lock` before starting the update process. The external update process should delete this file when the update is complete. If the lock file is not deleted within the specified timeout, the tool will remove it and proceed with shutdown.
