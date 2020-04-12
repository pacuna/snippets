package snippets_test

import (
	"snippets"
	"testing"
	"time"
)

func TestEncodeDecodeGob(t *testing.T) {
	var tags = []string{
		"go",
	}
	s := &snippets.Snippet{
		Id:       "0",
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

	enc := snippets.Encode(s, "gob")
	dec := snippets.Decode(enc, "gob")
	if dec.Content != s.Content {
		t.Errorf("Decoded content: %s, different than original content: %s", dec.Content, s.Content)
	}
}
