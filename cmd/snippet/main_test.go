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
	fileName = "snippets.db"
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
	t.Run("CreateSnippet", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "create", "-c", "-t", "Hello world", "-l", "go", "-tags", "basic,begginer")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ViewByLang", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "view", "-l", "go")
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
		cmd := exec.Command(cmdPath, "view", "-tag", "basic")
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
		cmd := exec.Command(cmdPath, "view", "-id", "1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		if content+"\n" != string(out) {
			t.Errorf("Expect %q, got %q instead\n", content, string(out))
		}
	})
}
