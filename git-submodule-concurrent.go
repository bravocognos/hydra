// The MIT License
//
// Copyright (c) 2019 Bravo Cognos, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

func main() {
	// Set max number of concurrency
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Handle arguments
	if len(os.Args) > 3 {
		log.Fatalln(`Error: should pass only 1 argument. Example: gsp "npm ci"`)
	}

	// Get the first argument
	userCommand := os.Args[1:]

	// Handle file opening and reading
	file, err := os.Open(".gitmodules")
	if err != nil {
		log.Fatalln(`Error: Failed to load .gitmodules`)
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(`Error: Failed to read .gitmodules`)
	}

	// Matches services names
	var re = regexp.MustCompile(`(?i)\[submodule "(.*)`)
	var gitModuleFileAsString = string(b)

	// Creates a channel with up 64 "tasks"/goroutines
	tasks := make(chan *exec.Cmd, 64)

	//Spawning N CPU Cores goroutines
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		go func(num int, w *sync.WaitGroup) {
			defer w.Done()

			var err error

			for cmd := range tasks { // this will exit the loop when the channel closes
				fmt.Printf(`Go routine id "%d" is running: "%s"`+"\n", num, cmd.Args[2])

				_, err = cmd.Output()
				if err != nil {
					fmt.Println("can't get stdout:", err)
				}
			}
		}(i, &wg)
	}

	//Generate Tasks
	// Iterate through all matches
	for _, match := range re.FindAllStringSubmatch(gitModuleFileAsString, -1) {
		// Clean up
		finalMatch := strings.ReplaceAll(match[1], "\"", "")
		finalMatch = strings.ReplaceAll(finalMatch, "]", "")

		// Get into the folder to run the command
		cdIntoFolder := "cd " + finalMatch

		// Runs the command
		tasks <- exec.Command("sh", "-c", cdIntoFolder+"; "+userCommand[0])
	}

	close(tasks)

	// Wait workers finish
	wg.Wait()

	fmt.Println("Done")
}
