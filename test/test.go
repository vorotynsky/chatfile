package test

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type TestFunc func(t *testing.T, input io.Reader, output io.Writer)

// DoTest runs a suite of golden tests by executing the provided test function for each case found under the `/test/cases` folder.
// Each case directory should contain a "chatfile" input file and a validation file specified by the `validation` parameter.
// The name of each test will be derived from the directory path relative to "cases", supporting nested structures.
//
// For each test case:
//   - Opens the "chatfile" as input data.
//   - Executes the `test` function with input and output streams.
//   - Compares the generated output with the validation file.
//   - If the outputs differ, marks the test as failed and saves the output to a "*.failed" file for troubleshooting.
//
// Usage example:
//
//	func myTest(t *testing.T) {
//	    DoTest(t, "validation.txt", func(t *testing.T, input io.Reader, output io.Writer) {
//	        data, err := io.ReadAll(input) // arrange
//	        actual := strings.LastIndexAny(data, "\n\r \t") // act
//	        output.Write(make([]byte, data)) // assert
//	    })
//	}
//
// Notes:
//   - Each case directory must contain a "chatfile" input file and the `validation` file.
//   - It is not necessary to perform all reading and writing separately or in a fixed order.
//   - This test method assumes that the output generated during testing is relatively small and can be safely stored in memory.
func DoTest(t *testing.T, validation string, test TestFunc) {
	_, base, _, ok := runtime.Caller(0)
	if !ok {
		t.Log("Failed to read testing directory")
		t.FailNow()
	}

	base = filepath.Dir(base)
	base = filepath.Join(base, "cases")

	skipped := make(map[string]bool)
	skipped["."] = false

	err := filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			t.Errorf("Failed to access %s: %v", path, err)
			return err
		}

		if d.IsDir() {
			checkIgnore(skipped, path, base)
			return nil
		}

		if d.Name() != "chatfile" {
			return nil
		}

		testAbsDir := filepath.Dir(path)
		testRelDir, _ := filepath.Rel(base, testAbsDir)
		testRelDir = filepath.ToSlash(testRelDir)
		testName := testRelDir

		validationFilePath := filepath.Join(testAbsDir, validation)

		bufferSize, exists := checkValidationFile(t, validationFilePath)
		if !exists {
			return nil
		}

		t.Run(testName, func(t *testing.T) {
			if skip, ok := skipped[testRelDir]; ok {
				if skip {
					t.SkipNow()
				}
				delete(skipped, testRelDir)
			}

			chatFile, err := os.Open(path)
			failNowWithError(t, "Failed to open chatfile %v", err)
			defer closeInTest(t, chatFile)

			validationFile, err := os.Open(validationFilePath)
			failNowWithError(t, "Failed to open validation file %v", err)
			defer closeInTest(t, validationFile)

			outBuffer := &bytes.Buffer{}
			outBuffer.Grow(bufferSize)

			test(t, chatFile, outBuffer)

			bytes := outBuffer.Bytes()
			if !fileCompare(validationFile, outBuffer) {
				t.Fail()
				t.Errorf("Test case %s failed, output doesn't match validation", testRelDir)

				failedOutputPath := validationFilePath + ".failed"
				t.Errorf("Writing output to %s", failedOutputPath)

				err = os.WriteFile(failedOutputPath, bytes, 0666)
				failNowWithError(t, "Failed to write output: %v", err)
			}
		})

		return fs.SkipDir
	})

	if err != nil {
		t.Errorf("Error occured while traversing test cases: %v", err)
		t.Failed()
	}
}

func checkIgnore(checked map[string]bool, path string, base string) {
	rel, _ := filepath.Rel(base, path)
	ignored := false

	dir := rel
	for dir != "." {
		if v, ok := checked[dir]; ok {
			ignored = v
			break
		}

		dir = filepath.Dir(dir)
	}

	name := "TEST_IGNORE"
	if ignored {
		name = "TEST_UNIGNORE"
	}

	ignoreFile := filepath.Join(path, name)
	if _, err := os.Stat(ignoreFile); err == nil {
		ignored = !ignored
	}

	checked[rel] = ignored
}

func checkValidationFile(t *testing.T, validationFilePath string) (int, bool) {
	stat, err := os.Stat(validationFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return -1, false
		} else {
			t.Errorf("Unknown error occured: %v", err)
		}
	}
	return int(stat.Size()), true
}

func failNowWithError(t *testing.T, format string, err error) {
	if err != nil {
		t.Logf(format, err)
		t.FailNow()
	}
}

func closeInTest(t *testing.T, file io.Closer) {
	if err := file.Close(); err != nil {
		t.Log(err)
	}
}

func fileCompare(file1, file2 io.Reader) bool {
	scan1 := bufio.NewScanner(file1)
	scan2 := bufio.NewScanner(file2)

	for scan1.Scan() {
		scan2.Scan()
		if !bytes.Equal(scan1.Bytes(), scan2.Bytes()) {
			return false
		}
	}

	// check that file2 has no additional content
	return !scan2.Scan()
}
