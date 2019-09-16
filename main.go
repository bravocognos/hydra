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
	//
	// Handles CLI arguments.
	// TODO: Use Cobra package.
	//
    if len(os.Args) > 3 {
		log.Fatalln(`Error: Should pass only 1 argument. Example: gsp "npm ci"`)
	}

    userCommand := os.Args[1:] // Get the first argument, the user command.

	//
	// Handles opening and reading files.
	//
	file, err := os.Open(".gitmodules")
	if err != nil {
		log.Fatalln("Error: Are you in a folder that has the submodules?")
	}

	defer file.Close()

    // TODO: Create a .gitmodules parser to access more information.
    // Example: Extract submodule path.
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln("Error: Failed to read .gitmodules. Is it valid?")
    }
    
    //
    // Extract sudmodule names, and check if there are submodules to work on.
    //
	var re = regexp.MustCompile(`(?i)\[submodule "(.*)`)
	var gitModuleFileAsString = string(b)

    allSubmatch := re.FindAllStringSubmatch(gitModuleFileAsString, -1)
    numberOfSubmodules := len(allSubmatch);
    
    if numberOfSubmodules < 1 {
        log.Fatalln("Are there any sudmodules?");
    }

    //
    // Handles and setup concurrency.
    // The size of the task queue corresponds to the number of submodules.
    //
	tasks := make(chan *exec.Cmd, numberOfSubmodules)

    var wg sync.WaitGroup
    
    // Spawns goroutines, max out (uses all cores).
    // TODO: Allows user to specify concurrency.
	for i := 0; i < runtime.NumCPU(); i++ {
        // Consumes 1 worker/thread from pool.
        wg.Add(1)

		go func(num int, w *sync.WaitGroup) {
            // Will indicate that task is done.
            // Releases worker/thread.
			defer w.Done()

			for cmd := range tasks {
                // Notify user about the task that's starting.
                fmt.Printf(
                    `Go routine "%d" is running: "%s"`+"\n",
                    num,
                    cmd.Args[2],
                )

                // TODO: Currently, output is quiet. Allows user to see it.
				_, err := cmd.Output()
				if err != nil {
					fmt.Println("can't get stdout:", err)
				}
			}
		}(i, &wg)
    }

    //
    // Handles submodule names, and task creation.
    //
	for _, match := range allSubmatch {
		finalMatch := strings.ReplaceAll(match[1], "\"", "")
		finalMatch = strings.ReplaceAll(finalMatch, "]", "")

		// Go to the submodule folder.
		cdIntoFolder := "cd " + finalMatch

		// Add task (user command) to the queue.
		tasks <- exec.Command("sh", "-c", cdIntoFolder+"; "+userCommand[0])
	}

    // Close channel.
	close(tasks)

	// Wait workers finish.
	wg.Wait()

    // Notifies that all tasks have been completed.
	fmt.Println("Done")
}
