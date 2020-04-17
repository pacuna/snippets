package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/atotto/clipboard"
)

var (
	binName  = "snippet"
	fileName = "snippets_test.db"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestSnippetsCLI(t *testing.T) {

	content := `
		import "fmt"

		func main(){
		  fmt.Println("Hello world")
		}
		`

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, binName)

	clipboard.WriteAll(content)

	t.Run("CreateSnippetFromClipboard", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "create", "-c", "-t", "Hello world", "-l", "go", "-tags", "basic,begginer", "-dbpath", fileName)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CreateSnippetNoInput", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "create", "-c", "false", "-t", "Hello world", "-l", "go", "-tags", "basic,begginer", "-dbpath", fileName)
		if err := cmd.Run(); err == nil {
			t.Fatal("Did not get an error")
		}
	})

	t.Run("CreateSnippetWithBothInputs", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "create", "-c", "false", "-t", "Hello world", "-l", "go", "-tags", "basic,begginer", "-dbpath", fileName)
		if err := cmd.Run(); err == nil {
			t.Fatal("Did not get an error")
		}
	})

	t.Run("ViewByLang", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "view", "-l", "go", "-dbpath", fileName)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := "[1] - Hello world\n"

		if expected != string(out) {
			t.Errorf("Expect %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("ViewByTag", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "view", "-tag", "basic", "-dbpath", fileName)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := "[1] - Hello world\n"

		if expected != string(out) {
			t.Errorf("Expect %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("ViewById", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "view", "-id", "1", "-dbpath", fileName)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		if content+"\n" != string(out) {
			t.Errorf("Expect %q, got %q instead\n", content, string(out))
		}
	})

	t.Run("DeleteSnippet", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "delete", "-id", "1", "-dbpath", fileName)
		if err := cmd.Run(); err != nil {
			t.Fatal("err")
		}
	})
}
