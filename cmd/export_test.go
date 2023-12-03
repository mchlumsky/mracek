package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testItem struct {
	name        string
	env         map[string]string
	args        []string
	expected    string
	expectedErr error
}

func testExecuteFunc(item testItem) func(t *testing.T) {
	return func(t *testing.T) {
		t.Cleanup(cleanEnv)

		for k, v := range item.env {
			t.Setenv(k, v)
		}

		rootCommand := buildRootCommand("")

		buffer := bytes.NewBufferString("")
		rootCommand.SetOut(buffer)
		rootCommand.SetArgs(item.args)

		if err := rootCommand.Execute(); err != nil {
			if diff := cmp.Diff(item.expectedErr.Error(), err.Error()); diff != "" {
				t.Error(diff)
			}
		}

		out, err := io.ReadAll(buffer)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(item.expected, string(out)); diff != "" {
			t.Error(diff)
		}
	}
}
