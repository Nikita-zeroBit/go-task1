package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	serverURL         = "http://srv.msk01.gigacorp.local/_stats"
	maxFetchErrors    = 3
	pollInterval      = 10 * time.Second
	loadAverageLimit  = 30
	memoryUsageLimit  = 0.8
	diskUsageLimit    = 0.9
	networkUsageLimit = 0.9
)

func main() {
	errorCount := 0

	for {
		resp, err := http.Get(serverURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			errorCount++
			if errorCount >= maxFetchErrors {
				fmt.Println("Unable to fetch server statistic.")
				return
			}
			time.Sleep(pollInterval)
			continue
		}

		errorCount = 0

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		statistics := strings.Split(strings.TrimSpace(string(body)), ",")

		loadAverage, _ := strconv.ParseFloat(statistics[0], 64)
		totalMemory, _ := strconv.ParseFloat(statistics[1], 64)
		usedMemory, _ := strconv.ParseFloat(statistics[2], 64)
		totalDisk, _ := strconv.ParseFloat(statistics[3], 64)
		usedDisk, _ := strconv.ParseFloat(statistics[4], 64)
		totalNetwork, _ := strconv.ParseFloat(statistics[5], 64)
		usedNetwork, _ := strconv.ParseFloat(statistics[6], 64)

		if loadAverage > loadAverageLimit {
			fmt.Printf("Load Average is too high: %.0f\n", loadAverage)
		}

		if totalMemory > 0 {
			memoryUsage := math.Round((usedMemory / totalMemory) * 100)
			if memoryUsage > memoryUsageLimit {
				fmt.Printf("Memory usage too high: %.0f%%\n", memoryUsage*100)
			}
		}

		if totalDisk > 0 {
			freeDisk := math.Floor((totalDisk - usedDisk) / (1024 * 1024))
			if usedDisk/totalDisk > diskUsageLimit {
				fmt.Printf("Free disk space is too low: %.0f Mb left\n", freeDisk)
			}
		}

		if totalNetwork > 0 {
			freeNetwork := math.Floor((totalNetwork - usedNetwork) * 8 / (1024 * 1024))
			if usedNetwork/totalNetwork > networkUsageLimit {
				fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", freeNetwork)
			}
		}

		time.Sleep(pollInterval)
	}
}
