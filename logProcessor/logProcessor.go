package main

import (
	"bufio"
	"fmt"
	"logAnalyzer/analyzer"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)
func main() {
	startTime := time.Now()

	// get command line arguments 
	args := os.Args
	
	scanner := bufio.NewScanner(os.Stdin)
	var inputLines []string

	// Read in the input from stdin, line by line
	for scanner.Scan() {
		line := scanner.Text()
		inputLines = append(inputLines, line) 
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		return 
	}

	var numThreads int 
	mode := "s"

	if len(args) > 1{
		numThreads, _ = strconv.Atoi(args[1])
		mode = "p"
	}

	results := analyzer.Analyze(mode, numThreads, inputLines); 

	plotAPICallTimes(&results)
	
	for api, stats := range results.Freqs {
		totalCalls := stats[1]
		numErrors := stats[0]
		var errorRate float64
		if totalCalls > 0 {
			errorRate = float64(numErrors) / float64(totalCalls) * 100
		} else {
			errorRate = 0 // Avoid division by zero
		}

		fmt.Printf("%s: Total Calls = %d, Errors = %d, Error Rate = %.2f%%\n", api, totalCalls, numErrors, errorRate)
	}

	fmt.Println("Number of Code Changes:", results.Changes);  

	duration := time.Since(startTime)
	fmt.Println(duration)
}


func plotAPICallTimes(data *analyzer.Data) error {
	p := plot.New()
	p.Title.Text = "API Call Times"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Count"

	var minTime, maxTime time.Time
	first := true

	for _, times := range data.Times {
		for _, t := range times {
			if first {
				minTime, maxTime = t, t
				first = false
			} else {
				if t.Before(minTime) {
					minTime = t
				}
				if t.After(maxTime) {
					maxTime = t
				}
			}
		}
	}

	if first {
		return fmt.Errorf("no time data available")
	}

	
	p.X.Min = float64(minTime.Unix())
	p.X.Max = float64(maxTime.Unix())

	
	for api, times := range data.Times {
		pts := make(plotter.XYs, len(times))
		for i, t := range times {
			pts[i].X = float64(t.Unix())
			pts[i].Y = 1 
		}

		scatter, err := plotter.NewScatter(pts)
		if err != nil {
			return err
		}
		scatter.GlyphStyle.Radius = vg.Points(float64(len(times))) 
		
		p.Legend.Add(api, scatter)

		p.Add(scatter)
	}

	
	p.X.Tick.Marker = plot.TimeTicks{Format: "Jan 02 15:04"}

	
	if err := p.Save(10*vg.Inch, 4*vg.Inch, "api_call_times.png"); err != nil {
		return err
	}

	return nil
}