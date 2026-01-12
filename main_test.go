package main_test

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "update .golden files")

func TestYamlDiff(t *testing.T) {
	goInstall(t)

	tests := []struct {
		name       string
		file1      string
		file2      string
		goldenFile string
	}{
		{
			name:       "single-doc",
			file1:      "testdata/1.yml",
			file2:      "testdata/2.yml",
			goldenFile: "testdata/diff.golden",
		},
		{
			name:       "multi-doc",
			file1:      "testdata/multi1.yml",
			file2:      "testdata/multi2.yml",
			goldenFile: "testdata/multi-diff.golden",
		},
		{
			name:       "multi-doc-a-shorter",
			file1:      "testdata/multi-short.yml",
			file2:      "testdata/multi-long.yml",
			goldenFile: "testdata/multi-a-shorter.golden",
		},
		{
			name:       "multi-doc-b-shorter",
			file1:      "testdata/multi-long.yml",
			file2:      "testdata/multi-short.yml",
			goldenFile: "testdata/multi-b-shorter.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if *update {
				err := os.WriteFile(tt.goldenFile, runYamldiff(t, tt.file1, tt.file2), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}
			contents, err := os.ReadFile(tt.goldenFile)
			if err != nil {
				t.Fatal(err)
			}
			want := string(contents)
			got := string(runYamldiff(t, tt.file1, tt.file2))
			if got != want {
				t.Errorf("got:\n%v\nwant:\n%v", got, want)
			}
		})
	}
}

func goInstall(t *testing.T) {
	install := exec.Command("go", "build")
	err := install.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func runYamldiff(t *testing.T, file1, file2 string) []byte {
	var out bytes.Buffer
	yamldiff := exec.Command("./yamldiff", file1, file2)
	yamldiff.Stdout = &out

	err := yamldiff.Start()
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan bool)
	go func() {
		err = yamldiff.Wait()
		done <- true
	}()
	timeout := time.Millisecond * 1000
	select {
	case <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(timeout):
		t.Fatalf("timed out after %v", timeout)
	}
	return out.Bytes()
}
