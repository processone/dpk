package main

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
// httprec add fixtures/test -u https://pic.twitter.com/OPYJZhQ9ih
// httprec add fixtures/test -u https://pbs.twimg.com/media/DeywOwwWsAIij8t.jpg:large
// httprec add fixtures/scenario1 -u https://pic.twitter.com/ncJzTbz3dT
// httprec add fixtures/scenario1 -u https://pbs.twimg.com/media/DuIZsfQX4AAZFbs.png:large

// TODO(mr): Support method to be able to record http "POST" requests

import "github.com/processone/dpk/pkg/httpmock/cmd/httprec/cmd"

func main() {
	cmd.Execute()
}
