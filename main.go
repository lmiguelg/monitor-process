package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <PID>")
		return
	}

	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid PID:", err)
		return
	}

	// Create a new plot
	p := plot.New()

	// Set plot title and labels
	p.Title.Text = "Memory Usage Over Time"
	p.X.Label.Text = "Time (seconds)"
	p.Y.Label.Text = "Memory Usage (%)"

	// Create plotter for memory usage data points
	points := make(plotter.XYs, 0)

	// Open file for writing memory usage data
	file, err := os.Create("memory_usage.dat")
	if err != nil {
		fmt.Println("Error creating data file:", err)
		return
	}
	defer file.Close()

	// Create a goroutine to monitor memory usage and update the plot
	go func() {
		for {
			// Execute ps command to get memory usage
			out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "%mem=").Output()
			if err != nil {
				fmt.Println("Error executing ps command:", err)
				return
			}

			// Parse memory usage from the output
			memUsageStr := strings.TrimSpace(string(out))
			memUsage, err := strconv.ParseFloat(memUsageStr, 64)
			if err != nil {
				fmt.Println("Error parsing memory usage:", err)
				return
			}

			// Write memory usage data to file
			_, err = file.WriteString(fmt.Sprintf("%d %.2f\n", time.Now().Unix(), memUsage))
			if err != nil {
				fmt.Println("Error writing to data file:", err)
				return
			}

			// Add data point to plotter
			points = append(points, plotter.XY{X: float64(time.Now().Unix()), Y: memUsage})

			// Clear plot and re-plot data points
			err = plotutil.AddLinePoints(p, "Memory Usage", points)
			if err != nil {
				fmt.Println("Error plotting data:", err)
				return
			}

			// Save plot to an image file
			err = p.Save(400, 300, "memory_usage.png")
			if err != nil {
				fmt.Println("Error saving plot:", err)
				return
			}

			time.Sleep(1 * time.Second)
		}
	}()

	// Wait indefinitely
	select {}
}
