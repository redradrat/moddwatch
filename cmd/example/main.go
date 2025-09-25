package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cortesi/moddwatch"
)

func main() {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	fmt.Printf("Watching for changes in: %s\n", wd)
	fmt.Println("Press Ctrl+C to stop...")

	// Create a channel to receive file change notifications
	modChan := make(chan *moddwatch.Mod)

	// Watch current directory with common include patterns
	// Watch all files except common ignore patterns
	includes := []string{"**"}
	excludes := []string{
		".git/**",
		"*.tmp",
		"*.swp",
		"*~",
		".DS_Store",
	}

	// Start watching with 100ms lull time
	watcher, err := moddwatch.Watch(wd, includes, excludes, 100*time.Millisecond, modChan)
	if err != nil {
		log.Fatal("Failed to start watcher:", err)
	}
	defer watcher.Stop()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Process file change events
	for {
		select {
		case mod := <-modChan:
			if mod != nil && !mod.Empty() {
				fmt.Printf("\n=== Files changed at %s ===\n", time.Now().Format("15:04:05"))

				if len(mod.Added) > 0 {
					fmt.Printf("Added: %v\n", mod.Added)
				}
				if len(mod.Changed) > 0 {
					fmt.Printf("Changed: %v\n", mod.Changed)
				}
				if len(mod.Deleted) > 0 {
					fmt.Printf("Deleted: %v\n", mod.Deleted)
				}
			}
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

