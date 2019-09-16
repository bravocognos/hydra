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
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"

	log "github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

var userCommand []string
var submodulesNames []string

/*
 * Helpers
 */

// readAndParseGitModulesFile handles opening, reading, and returning the
// `.gitmodules` file.
func readAndParseGitModulesFile() []byte {
	file, err := os.Open(".gitmodules")
	if err != nil {
		log.Fatal("Error: Are you in a folder that has the submodules?")
	}

	defer file.Close()

	// TODO: Create a `.gitmodules` parser to access more information.
	// Example: Extract the path of the submodule.
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error: Failed to read .gitmodules. Is it valid?")
	}

	return b
}

// getSubmodulesNames extracts submodule names and checks for submodules to
// work with. Return the number of submodules and their names.
func getSubmodulesNames(gitModulesFile []byte) ([]string, int) {
	var re = regexp.MustCompile(`(?i)\[submodule "(.*)`)
	var gitModuleFileAsString = string(gitModulesFile)

	allSubmatch := re.FindAllStringSubmatch(gitModuleFileAsString, -1)
	numberOfSubmodules := len(allSubmatch)

	if numberOfSubmodules < 1 {
		log.Fatal("Are there any sudmodules?")
	}

	for _, match := range allSubmatch {
		// Cleanup
		finalMatch := strings.ReplaceAll(match[1], "\"", "")
		finalMatch = strings.ReplaceAll(finalMatch, "]", "")

		submodulesNames = append(submodulesNames, finalMatch)
	}

	return submodulesNames, numberOfSubmodules
}

// createTasks handles task creation. A task is a user command that will be
// executed against all submodules. Each task runs on its own thread,
// concurrently.
func createTasks(submodulesNames []string, tasksQueue chan *exec.Cmd) {
	for _, submoduleName := range submodulesNames {
		// Go to the submodule folder.
		cdIntoFolder := "cd " + submoduleName

		// Add task (user command) to the queue.
		tasksQueue <- exec.Command("sh", "-c", cdIntoFolder+"; "+userCommand[0])
	}
}

/*
 * Starts here
 */

// Runs before main
func init() {
	//
	// Handles CLI arguments.
	// TODO: Should use Cobra package.
	//
	if len(os.Args) > 3 {
		log.Fatal(`Error: Should pass only 1 argument. Example: gsp "npm ci"`)
	}

	userCommand = os.Args[1:] // Get the first argument, the user command.

	//
	// Setup logger.
	// Available levels:
	// - debug
	// - info
	// - warn
	// - warning
	// - error
	// - fatal
	//
	logLevel := os.Getenv("DEBUG")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		fmt.Println("Invalid log level! " +
			`Available: "debug", "info", "warn", "warning", "error", and "fatal"`)
		os.Exit(1)
	}

	log.SetLevel(level)
	log.SetHandler(text.New(os.Stderr))

	// Notifies that application started
	log.Info("Hydra started")
}

func main() {
	gitModulesFile := readAndParseGitModulesFile()
	submodulesNames, numberOfSubmodules := getSubmodulesNames(gitModulesFile)

	//
	// Handles and setup concurrency.
	//
	// The size of the task queue corresponds to the number of submodules.
	tasksQueue := make(chan *exec.Cmd, numberOfSubmodules)

	var wg sync.WaitGroup

	// Spawns goroutines, max out (uses all cores).
	// TODO: Allows user to specify concurrency.
	for i := 0; i < runtime.NumCPU(); i++ {
		// Consumes 1 worker/thread from pool.
		wg.Add(1)

		go func(id int, w *sync.WaitGroup) {
			// Will indicate that task is done.
			// Releases worker/thread.
			defer w.Done()

			for cmd := range tasksQueue {
				log.WithFields(log.Fields{
					"threadId": id,
					"cmd":      cmd.Args[2],
				}).Info("Task started")

				// TODO: Currently, output is quiet. Allows user to see it.
				output, err := cmd.Output()
				if err != nil {
					fmt.Printf(`Error in go routine "%d" running: "%s": %s`+
						"\n",
						id,
						cmd.Args[2],
						err.Error(),
					)
				}

				log.WithFields(log.Fields{
					"threadId": id,
					"cmd":      cmd.Args[2],
					"output":   string(output) + "\n",
				}).Debug("Done")

			}
		}(i, &wg)
	}

	createTasks(submodulesNames, tasksQueue)

	// Close channels.
	close(tasksQueue)

	// Wait workers finish.
	wg.Wait()

	// Notifies that all tasks have been completed.
	log.Info("Done")
}
