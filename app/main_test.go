package main

import (
	"chatfile/test"
	"io"
	"testing"
)

func TestDummy(t *testing.T) {
	test.DoTest(t, "chatfile", func(t *testing.T, input io.Reader, output io.Writer) {
		data, err := io.ReadAll(input)
		if err != nil {
			t.Errorf("Reading test data failed %v", err)
		}

		_, err = output.Write(data)
		if err != nil {
			t.Errorf("Writing test data failed %v", err)
		}
	})
}
