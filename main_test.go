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
	goldenfile := "testdata/diff.golden"
	if *update {
		err := os.WriteFile(goldenfile, runYamldiff(t), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	contents, err := os.ReadFile(goldenfile)
	want := string(contents)
	if err != nil {
		t.Fatal(err)
	}
	got := string(runYamldiff(t))
	if got != want {
		t.Errorf("got: %v want: %v", got, want)
	}
}

func goInstall(t *testing.T) {
	install := exec.Command("go", "install")
	err := install.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func runYamldiff(t *testing.T) []byte {
	var out bytes.Buffer
	yamldiff := exec.Command("yamldiff", "--file1=testdata/1.yml", "--file2=testdata/2.yml")
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
	timeout := time.Millisecond * 100
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
