package snippets_test

import (
	"io/ioutil"
	"os"
	"snippets"
	"testing"
	"time"
)

func tempFile() string {
	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}

func TestStore_CreateSnippet(t *testing.T) {

	var tags = []string{
		"go",
	}
	s := &snippets.Snippet{
		ID:       1,
		Title:    "Hello world",
		Language: "Go",
		Content: `
import "fmt"

func main(){
  fmt.Println("Hello world")
}
`,
		CreatedAt: time.Now(),
		Tags:      tags,
	}

	store := snippets.NewStore(tempFile(), "snippets")
	err := store.CreateSnippet(s)

	if err != nil {
		t.Errorf("There was an error creating snippet: %v", err)
	}
}

func TestStore_GetSnippetByID(t *testing.T) {
	var tags = []string{
		"go",
	}
	s := &snippets.Snippet{
		ID:       1,
		Title:    "Hello world",
		Language: "Go",
		Content: `
import "fmt"

func main(){
  fmt.Println("Hello world")
}
`,
		CreatedAt: time.Now(),
		Tags:      tags,
	}

	store := snippets.NewStore(tempFile(), "snippets")
	err := store.CreateSnippet(s)

	if err != nil {
		t.Errorf("There was an error creating snippet: %v", err)
	}

	got, _ := store.GetSnippetByID(1)
	if got.ID != s.ID {
		t.Errorf("Got ID %d, expected %d", got.ID, s.ID)
	}
	if got.Title != s.Title {
		t.Errorf("Got Title %s, expected %s", got.Title, s.Title)
	}
	if got.Language != s.Language {
		t.Errorf("Got Language %s, expected %s", got.Language, s.Language)
	}
	if got.Content != s.Content {
		t.Errorf("Got Content %s, expected %s", got.Content, s.Content)
	}

}

func TestStore_ListSnippetsByLang(t *testing.T) {
	var tags = []string{
		"go",
	}
	s := &snippets.Snippet{
		ID:       1,
		Title:    "Hello world",
		Language: "Go",
		Content: `
import "fmt"

func main(){
  fmt.Println("Hello world")
}
`,
		CreatedAt: time.Now(),
		Tags:      tags,
	}

	store := snippets.NewStore(tempFile(), "snippets")
	err := store.CreateSnippet(s)

	if err != nil {
		t.Errorf("There was an error creating snippet: %v", err)
	}

	sl, _ := store.ListSnippetsByLang("Go")
	if len(sl) == 0 {
		t.Errorf("No results were found")
	} else {
		if sl[0].ID != s.ID {
			t.Errorf("Got ID %d, expected %d", sl[0].ID, s.ID)
		}
		if sl[0].Title != s.Title {
			t.Errorf("Got Title %s, expected %s", sl[0].Title, s.Title)
		}
		if sl[0].Language != s.Language {
			t.Errorf("Got Language %s, expected %s", sl[0].Language, s.Language)
		}
		if sl[0].Content != s.Content {
			t.Errorf("Got Content %s, expected %s", sl[0].Content, s.Content)
		}
	}

}

func TestStore_ListSnippetsByTag(t *testing.T) {
	var tags = []string{
		"basic",
	}
	s := &snippets.Snippet{
		ID:       1,
		Title:    "Hello world",
		Language: "Go",
		Content: `
import "fmt"

func main(){
  fmt.Println("Hello world")
}
`,
		CreatedAt: time.Now(),
		Tags:      tags,
	}

	store := snippets.NewStore(tempFile(), "snippets")
	err := store.CreateSnippet(s)

	if err != nil {
		t.Errorf("There was an error creating snippet: %v", err)
	}

	sl, _ := store.ListSnippetsByTag("basic")

	if len(sl) == 0 {
		t.Fatalf("No results were found")
	} else {
		if sl[0].ID != s.ID {
			t.Errorf("Got ID %d, expected %d", sl[0].ID, s.ID)
		}
		if sl[0].Title != s.Title {
			t.Errorf("Got Title %s, expected %s", sl[0].Title, s.Title)
		}
		if sl[0].Language != s.Language {
			t.Errorf("Got Language %s, expected %s", sl[0].Language, s.Language)
		}
		if sl[0].Content != s.Content {
			t.Errorf("Got Content %s, expected %s", sl[0].Content, s.Content)
		}
	}

}

func TestStore_DeleteSnippet(t *testing.T) {

	s := &snippets.Snippet{
		ID:       1,
		Title:    "Hello world",
		Language: "Go",
		Content: `
import "fmt"

func main(){
  fmt.Println("Hello world")
}
`,
		CreatedAt: time.Now(),
	}

	store := snippets.NewStore(tempFile(), "snippets")
	err := store.CreateSnippet(s)

	if err != nil {
		t.Errorf("There was an error creating snippet: %v", err)
	}

	snippet, err := store.GetSnippetByID(1)
	if err != nil {
		t.Errorf("There was an error retrieving snippet: %v", err)
	}

	if snippet.Title == "" {
		t.Error("Snippet was not create")
	}

	// Now delete it
	err = store.DeleteSnippet(1)
	if err != nil {
		t.Error("Error deleting snippet: ", err)
	}

	snippet, err = store.GetSnippetByID(1)
	if err != nil {
		t.Errorf("There was an error retrieving snippet: %v", err)
	}

	if snippet != nil {
		t.Error("Snippet was not deleted")
	}

}
