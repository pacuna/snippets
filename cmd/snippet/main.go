package main

import (
	"flag"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"log"
	"os"
	"snippets"
	"time"
)

type Client struct {
	db *bolt.DB
}

// CreateSnippet stores a snippet using $lang:$id:$title as the key and the whole obj as the val
func (c *Client) CreateSnippet(s *snippets.Snippet) error {

	// create bucket for this specific language if it doesn't exist
	c.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(s.Language)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.Language))
		id, _ := b.NextSequence()
		s.ID = int(id)
		key := string(s.ID)
		return b.Put([]byte(key), snippets.Encode(s, "gob"))
	})
}

func (c *Client) ListSnippetsByLang(lang string) []*snippets.Snippet {
	var sl []*snippets.Snippet
	c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(lang))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			s := snippets.Decode(v, "gob")
			sl = append(sl, s)
		}
		return nil
	})
	return sl
}

func (c *Client) GetSnippetByID(lang string, id int) *snippets.Snippet{
	var s *snippets.Snippet
	c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(lang))
		v := b.Get([]byte(string(id)))
		s = snippets.Decode(v, "gob")
		return nil
	})
	return s
}

func main() {

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createLang := createCmd.String("l", "", "Language")
	createTitle := createCmd.String("t", "", "Title")
	createFilePath := createCmd.String("f", "", "File path")

	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	viewLang := viewCmd.String("l", "", "Language")
	viewId := viewCmd.Int("id", 0, "ID")

	if len(os.Args) < 2{
		fmt.Println("Expected 'create' or 'view' subcommands")
		os.Exit(1)
	}

	db, err := bolt.Open("snippets.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	client := &Client{db: db}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *createFilePath == "" {
			fmt.Println("File path is mandatory")
			os.Exit(1)
		}
		if *createLang == "" {
			fmt.Println("Language is mandatory")
			os.Exit(1)
		}
		if *createTitle == "" {
			fmt.Println("Title is mandatory")
			os.Exit(1)
		}
		content, err := ioutil.ReadFile(*createFilePath)
		if err != nil {
			log.Fatal("File reading error: ", err)
			os.Exit(1)
		}
		s := snippets.New(*createTitle, *createLang, string(content), time.Now())
		client.CreateSnippet(s)
	case "view":
		viewCmd.Parse(os.Args[2:])
		if *viewLang == "" {
			fmt.Println("Language is mandatory")
			os.Exit(1)
		}

		// No ID provided. List all for this language
		if *viewId == 0 {
			sl := client.ListSnippetsByLang(*viewLang)
			for _, s := range sl {
				fmt.Printf("[%d] - %s\n", s.ID, s.Title)
			}
		}else{
			s := client.GetSnippetByID(*viewLang, *viewId)
			fmt.Println(s.Content)
		}
	default:
		fmt.Println("expected 'create' or 'view' subcommands")
		os.Exit(1)
	}

}
