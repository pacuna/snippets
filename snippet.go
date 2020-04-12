package snippets

import "time"

type Snippet struct {
	Id        string
	Title     string
	Language  string
	Content   string
	CreatedAt time.Time
	Tags      []string
}
