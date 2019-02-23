package main

import (
	"fmt"
	"log"
	"os"

	"github.com/processone/dpk/pkg/httpmock"
)

// Fully record HTTP request interactions, including redirects
// The result can be used as an input for testing.
// It is used to help simulating test involving HTTP requests.
// The result is saved as a .gob file, a native format to serialize Go struct.
// It also save a .url file containing the original URLs. It can be useful to regenerate
// the data if we need to change the structure of the Sequence.

// params:
// - dir: Fixture dir (default: fixtures)
// - name: Test fixture name
// - url: URL retrieve to record

// Example:
// httprec https://pic.twitter.com/OPYJZhQ9ih test
// httprec https://pbs.twimg.com/media/DeywOwwWsAIij8t.jpg:large test
// httprec https://pic.twitter.com/ncJzTbz3dT scenario1
// httprec https://pbs.twimg.com/media/DuIZsfQX4AAZFbs.png:large scenario1

// TODO(mr): Use Cobra for parameters formatting
// TODO(mr): Support method to be able to record http "POST" requests

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Missing url or name")
		os.Exit(1)
	}

	uri := args[0]
	scnName := args[1]

	recorder := httpmock.Recorder{Logger: log.New(os.Stderr, "", 0)}
	if err := recorder.Record(uri, scnName); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
