package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type WorkloadType string

const (
	CPUWorkload   WorkloadType = "cpu"
	MixedWorkload WorkloadType = "mixed"
)

type Metrics struct {
	TotalTime       time.Duration
	Throughput      float64
	MeanDelay       time.Duration
	MinDelay        time.Duration
	MaxDelay        time.Duration
	StdDevDelay     time.Duration
	NumBlocking     int64
	Goroutines      int
	GOMAXPROCS      int
	Workload        WorkloadType
}

// CPU-bound task
func cpuHeavyTask(iter int) {
	sum := 0
	for i := 0; i < iter; i++ {
		sum += (i * 2342) % 71
	}
	_ = sum
}

// Mixed workload task, compute + a sleep
func mixedTask(iter int) {
	cpuHeavyTask(iter / 2)
	time.Sleep(20 * time.Millisecond)
	cpuHeavyTask(iter / 2)
}

func runExperiment(numGoroutines int, gomax int, workload WorkloadType, iterations int) Metrics {
	runtime.GOMAXPROCS(gomax)

	var wg sync.WaitGroup
	var blockingCount int64

	// time per each goroutine
	delay := make([]time.Duration, numGoroutines)

	startAll := time.Now()

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			start := time.Now()

			switch workload {
			case CPUWorkload:
				cpuHeavyTask(iterations)
			case MixedWorkload:
				mixedTask(iterations)
				atomic.AddInt64(&blockingCount, 1)
			}

			delay[id] = time.Since(start)
		}(g)
	}

	wg.Wait()
	totalTime := time.Since(startAll)


    // metrics
	var minDelay time.Duration = delay[0]
	var maxDelay time.Duration = delay[0]

	var sum float64 = 0
	for _, d := range delay {
		if d < minDelay {
			minDelay = d
		}
		if d > maxDelay {
			maxDelay = d
		}
		sum += float64(d)
	}

	mean := sum / float64(numGoroutines)

	var variance float64 = 0
	for _, d := range delay {
		diff := float64(d) - mean
		variance += diff * diff
	}
	variance /= float64(numGoroutines)
	stdDev := math.Sqrt(variance)

	throughput := float64(numGoroutines) / totalTime.Seconds()

	return Metrics{
		TotalTime:   totalTime,
		Throughput:  throughput,
		MeanDelay:   time.Duration(mean),
		MinDelay:    minDelay,
		MaxDelay:    maxDelay,
		StdDevDelay: time.Duration(stdDev),
		Goroutines:  numGoroutines,
		GOMAXPROCS:  gomax,
		NumBlocking: blockingCount,
		Workload:    workload,
	}
}
func printMetrics(m Metrics) {
	fmt.Println("===========================================")
	fmt.Println("             EXPERIMENT RESULTS            ")
	fmt.Println("===========================================")
	fmt.Printf("Workload Type     : %v\n", m.Workload)
	fmt.Printf("GOMAXPROCS        : %d\n", m.GOMAXPROCS)
	fmt.Printf("Goroutines        : %d\n", m.Goroutines)
	fmt.Println("-------------------------------------------")
	fmt.Printf("Total Runtime     : %v\n", m.TotalTime)
	fmt.Printf("Throughput        : %.4f ops/sec\n", m.Throughput)
	fmt.Printf("Mean Delay        : %v\n", m.MeanDelay)
	fmt.Printf("Min Delay         : %v\n", m.MinDelay)
	fmt.Printf("Max Delay         : %v\n", m.MaxDelay)
	fmt.Printf("StdDev Delay      : %v\n", m.StdDevDelay)
	fmt.Printf("Blocking Events   : %d\n", m.NumBlocking)
	fmt.Println("===========================================")
	fmt.Println()
}

func main() {
	workloadFlag := flag.String("workload", "cpu", "Workload type: cpu or mixed")
	goroutinesFlag := flag.Int("goroutines", 8, "Number of goroutines to run")
	procsFlag := flag.Int("procs", runtime.NumCPU(), "GOMAXPROCS value")
	iterationsFlag := flag.Int("iterations", 80_000_000, "Complexity of workload computation")

	flag.Parse()

	workloadType := WorkloadType(*workloadFlag)
	if workloadType != CPUWorkload && workloadType != MixedWorkload {
		fmt.Println("Invalid workload type. Use 'cpu' or 'mixed'.")
		return
	}

	metrics := runExperiment(
		*goroutinesFlag,
		*procsFlag,
		workloadType,
		*iterationsFlag,
	)

	printMetrics(metrics)
}
