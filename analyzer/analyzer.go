package analyzer

import (
	"fmt"
	"logAnalyzer/doubleQueue"
	"strings"
	"sync"
	"time"
)

type Data struct {
    Freqs map[string][]int
    Times map[string][]time.Time
	Changes int 
}

type CondVar struct {
	Mu *sync.Mutex
	Cond *sync.Cond
	Count *int 
}

func Analyze(mode string, numThreads int, logs []string) Data {
	if (mode == "s"){
		m := make(map[string][]int)
		times := make(map[string][]time.Time)
		changes := 0

		for i:= 0; i < len(logs); i++ {
			// break up string 
			fields := strings.Fields(logs[i])

			// Get type of HTTP request 
			request := fields[1] 

			if (request != "INFO:") {
				// Get URL request name
				name := fields[2]; 

				// Get HTTP response Code
				code := fields[3]

				// Get + Process Date 
				dateStr := fields[0][1:len(fields[0])-1]
				
				layout := "02/Jan/2006-15:04:05"
				// Parse the date string using the layout
				t, err := time.Parse(layout, dateStr)
				if err != nil {
					fmt.Println("Error parsing date:", err)
				}

				// break string
				if _, exists := times[name]; !exists {
					times[name] = []time.Time{}
				}
				
				times[name] = append(times[name], t)

				// break string
				if _, exists := m[name]; !exists {
					m[name] = []int{0, 0}
				}
				
				if (code[0] == '4'){
					m[name][0] += 1;
				}
				m[name][1] += 1
			} else {
				changes++
			}
		}
		resultsStruct := Data{
			Freqs : m,
			Times : times,
			Changes : changes,
	 	}

		return resultsStruct

	} else {
		var queues []*doubleQueue.DLQueue; 

		for i := 0; i<numThreads; i++ {
			queue := doubleQueue.NewdLQueue()
			queues = append(queues, queue) 
		}

		for i := 0; i < len(logs); i++ {
			queues[i%numThreads].PushTop(logs[i]); 
		}

		
		results := make(chan Data, numThreads)

		

		var mu sync.Mutex
		count := 0 
		cond := CondVar{
			Mu:    &mu,
			Cond:  sync.NewCond(&mu),
			Count: &count,
		}

		for i := 0; i<numThreads; i++ {
			go threadAnalyze(i, cond, queues, results, numThreads)
		}
		
		freqResults := make(map[string][]int)
		timeResults := make(map[string][]time.Time)
		changes := 0 

		for result := range results {
			for key, value := range result.Freqs {
				if _, exists := freqResults[key]; !exists {
					freqResults[key] = []int{0, 0}
				} 
				freqResults[key][0] += value[0]
				freqResults[key][1] += value[1]
			}
			for key, value := range result.Times {
				if _, exists := timeResults[key]; !exists {
					timeResults[key] = []time.Time{}
				} 
				timeResults[key] = append(timeResults[key], value...)
			}
			changes += result.Changes
		}

		resultsStruct := Data{
			Freqs : freqResults,
			Times : timeResults,
			Changes : changes,
	 	}

		return resultsStruct
	}
}

func threadAnalyze(tid int, cond CondVar, queues []*doubleQueue.DLQueue,
				   results chan Data, numThreads int){
	m := make(map[string][]int)
	times := make(map[string][]time.Time)
	changes := 0
	
	for i := 0; i< numThreads; i++ {
		for work := queues[(tid+i)%numThreads].PopBottom(); work != nil; work = queues[tid].PopBottom() {
			work, ok := work.(string)
			if !ok {
				break
			}
			// break up string 
			fields := strings.Fields(work)

			// Get type of HTTP request 
			request := fields[1] 

			if (request != "INFO:") {
				// Get URL request name
				name := fields[2]; 

				// Get HTTP response Code
				code := fields[3]

				// Get + Process Date 
				dateStr := fields[0][1:len(fields[0])-1]
				
				layout := "02/Jan/2006-15:04:05"
				// Parse the date string using the layout
				t, err := time.Parse(layout, dateStr)
				if err != nil {
					fmt.Println("Error parsing date:", err)
				}

				// break string
				if _, exists := times[name]; !exists {
					times[name] = []time.Time{}
				}
				
				times[name] = append(times[name], t)

				// break string
				if _, exists := m[name]; !exists {
					m[name] = []int{0, 0}
				}
				
				if (code[0] == '4'){
					m[name][0] += 1;
				}
				m[name][1] += 1
			} else {
				changes++
			}
		}
	}

	threadResults := Data{
		Freqs : m,
		Times : times,
		Changes : changes,
	}
	
	results <- threadResults

	cond.Mu.Lock()
	*cond.Count += 1
	if *cond.Count == numThreads {
		cond.Cond.Broadcast()
	} else {
		cond.Cond.Wait()
	}
	cond.Mu.Unlock() 

	if tid == 0 {
		close(results)
	}
}