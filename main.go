package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	fmt.Println("DBC Updater Tool")

	// Step 1: Stop vehicle service
	fmt.Println("Stopping vehicle service...")
	err := stopVehicleService()
	if err != nil {
		fmt.Printf("Error stopping vehicle service: %v\n", err)
		return
	}
	fmt.Println("Vehicle service stopped.")

	// Step 2: Prepare GPIO for DBC power
	fmt.Println("Preparing GPIO for DBC power...")
	err = prepareGPIOPower()
	if err != nil {
		fmt.Printf("Error preparing GPIO: %v\n", err)
		return
	}
	fmt.Println("GPIO prepared.")

	// Step 3: Turn on DBC
	fmt.Println("Turning on DBC...")
	err = turnOnDBC()
	if err != nil {
		fmt.Printf("Error turning on DBC: %v\n", err)
		return
	}
	fmt.Println("DBC turned on.")

	// Step 4: Monitor and reset dashboard ready state
	fmt.Println("Monitoring and resetting dashboard ready state...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Monitor for a limited time
	defer cancel()
	err = monitorAndResetDashboardReady(ctx)
	if err != nil {
		fmt.Printf("Error monitoring/resetting dashboard ready state: %v\n", err)
		// Continue with the update process even if monitoring fails? Or exit?
		// For now, let's just log the error and continue.
	}
	fmt.Println("Dashboard ready state monitoring finished.")

	// Step 5: Set BLE pin-code to UPDATE
	fmt.Println("Setting BLE pin-code to UPDATE...")
	err = setBLEPinCodeToUpdate(ctx)
	if err != nil {
		fmt.Printf("Error setting BLE pin-code: %v\n", err)
		// Continue with the update process even if this fails
	}
	fmt.Println("BLE pin-code set to UPDATE.")

	// Step 6: Run the update process
	fmt.Println("Starting update process (Placeholder)...")
	// TODO: Integrate SMUT or other update mechanism

	// Simulate update process
	time.Sleep(5 * time.Second)
	fmt.Println("Update process finished (Simulated).")

	// Step 7: Clear BLE pin-code
	fmt.Println("Clearing BLE pin-code...")
	err = clearBLEPinCode(ctx)
	if err != nil {
		fmt.Printf("Error clearing BLE pin-code: %v\n", err)
		// Continue with the shutdown process even if this fails
	}
	fmt.Println("BLE pin-code cleared.")

	// Step 8: Turn off DBC
	fmt.Println("Turning off DBC...")
	err = turnOffDBC()
	if err != nil {
		fmt.Printf("Error turning off DBC: %v\n", err)
		return
	}
	fmt.Println("DBC turned off.")

	// Step 9: Restart vehicle service
	fmt.Println("Restarting vehicle service...")
	err = restartVehicleService()
	if err != nil {
		fmt.Printf("Error restarting vehicle service: %v\n", err)
		return
	}
	fmt.Println("Vehicle service restarted.")

	fmt.Println("DBC Updater Tool finished.")
}

func stopVehicleService() error {
	cmd1 := exec.Command("systemctl", "stop", "librescoot-vehicle")
	err1 := cmd1.Run()
	cmd2 := exec.Command("systemctl", "stop", "unu-vehicle")
	err2 := cmd2.Run()

	// If both commands failed, return an error. Otherwise, consider it successful.
	if err1 != nil && err2 != nil {
		return fmt.Errorf("failed to stop both librescoot-vehicle (%v) and unu-vehicle (%v)", err1, err2)
	}
	return nil
}

func prepareGPIOPower() error {
	// Assuming GPIO 50 is available and not used by the vehicle service
	// This might require running as root or with appropriate permissions

	// Using os.WriteFile is more reliable than exec.Command with redirection
	err := os.WriteFile("/sys/class/gpio/export", []byte("50"), 0644)
	if err != nil {
		// Ignore "device or resource busy" error if already exported
		if !os.IsExist(err) {
			return fmt.Errorf("failed to export GPIO 50: %w", err)
		}
	}

	// Give the system a moment to set up the GPIO
	time.Sleep(100 * time.Millisecond)

	err = os.WriteFile("/sys/class/gpio/gpio50/direction", []byte("out"), 0644)
	if err != nil {
		return fmt.Errorf("failed to set GPIO 50 direction: %w", err)
	}
	return nil
}

func turnOnDBC() error {
	err := os.WriteFile("/sys/class/gpio/gpio50/value", []byte("1"), 0644)
	if err != nil {
		return fmt.Errorf("failed to turn on DBC: %w", err)
	}
	return nil
}

func turnOffDBC() error {
	err := os.WriteFile("/sys/class/gpio/gpio50/value", []byte("0"), 0644)
	if err != nil {
		return fmt.Errorf("failed to turn off DBC: %w", err)
	}
	return nil
}

func restartVehicleService() error {
	cmd1 := exec.Command("systemctl", "start", "librescoot-vehicle")
	err1 := cmd1.Run()
	cmd2 := exec.Command("systemctl", "start", "unu-vehicle")
	err2 := cmd2.Run()

	// If both commands failed, return an error. Otherwise, consider it successful.
	if err1 != nil && err2 != nil {
		return fmt.Errorf("failed to start both librescoot-vehicle (%v) and unu-vehicle (%v)", err1, err2)
	}
	return nil
}

// Create a Redis client with the standard configuration
func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "192.168.7.1:6379", // Redis address
		Password: "",                 // no password set
		DB:       0,                  // use default DB
	})
}

func setBLEPinCodeToUpdate(ctx context.Context) error {
	rdb := createRedisClient()
	defer rdb.Close()

	// Set the BLE pin-code to UPDATE
	err := rdb.HSet(ctx, "ble", "pin-code", "UPDATE").Err()
	if err != nil {
		return fmt.Errorf("failed to set BLE pin-code to UPDATE: %w", err)
	}
	return nil
}

func clearBLEPinCode(ctx context.Context) error {
	rdb := createRedisClient()
	defer rdb.Close()

	// Delete the pin-code field from the ble hash
	err := rdb.HDel(ctx, "ble", "pin-code").Err()
	if err != nil {
		return fmt.Errorf("failed to clear BLE pin-code: %w", err)
	}
	return nil
}

func monitorAndResetDashboardReady(ctx context.Context) error {
	rdb := createRedisClient()
	defer rdb.Close()

	// Use a Pub/Sub subscriber to listen for changes to the dashboard key
	pubsub := rdb.Subscribe(ctx, "dashboard") // Subscribe to the dashboard channel
	defer pubsub.Close()

	// Wait for confirmation that the subscription is active
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to dashboard channel: %w", err)
	}

	fmt.Println("Subscribed to dashboard channel.")

	// Loop to receive messages or until context is cancelled
	for {
		select {
		case msg := <-pubsub.Channel():
			fmt.Printf("Received message on dashboard channel: %s\n", msg.Payload)
			// Check if the message indicates dashboard ready
			if msg.Payload == "ready" { // Check for the "ready" payload
				fmt.Println("Dashboard ready detected. Resetting state...")
				// Set dashboard ready to false using HSet
				err := rdb.HSet(ctx, "dashboard", "ready", "false").Err()
				if err != nil {
					fmt.Printf("Error setting dashboard ready to false: %v\n", err)
					// Decide if this is a critical error or just log and continue
				} else {
					fmt.Println("Dashboard ready state reset to false.")
				}
				// We might want to stop monitoring after the first detection and reset
				// return nil // Exit after first reset
			}
		case <-ctx.Done():
			fmt.Println("Monitoring context cancelled.")
			return ctx.Err()
		}
	}
}
