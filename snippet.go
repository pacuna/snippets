package snippets

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Snippet struct {
	ID        int
	Title     string
	Language  string
	Content   string
	Tags      []string
	CreatedAt time.Time
}

func New(title string, language string, content string, createdAt time.Time, tags []string) *Snippet {
	return &Snippet{Title: title, Language: language, Content: content, CreatedAt: createdAt, Tags: tags}
}

func Encode(s *Snippet, format string) []byte {

	var b bytes.Buffer
	if format == "gob" {
		enc := gob.NewEncoder(&b)
		err := enc.Encode(s)
		if err != nil {
			log.Fatal("encode error:", err)
		}
	}
	return b.Bytes()
}

func Decode(data []byte, format string) *Snippet {

	var s *Snippet
	enc := gob.NewDecoder(bytes.NewReader(data))
	err := enc.Decode(&s)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return s
}
