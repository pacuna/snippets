package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"log"
	"snippets"
	"strings"
	"time"
)

type Client struct {
	db *bolt.DB
}

func (c *Client) CreateSnippet(s *snippets.Snippet) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("snippets"))

		id, _ := b.NextSequence()
		s.ID = int(id)

		return b.Put(itob(s.ID), snippets.Encode(s, "gob"))
	})
}

func (c *Client) GetSnippetByID(id int) *snippets.Snippet {
	var s = &snippets.Snippet{}
	c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("snippets"))
		v := b.Get(itob(id))
		s = snippets.Decode(v, "gob")
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})
	return s
}

func (c *Client) GetAllSnippets() []*snippets.Snippet {
	var sl []*snippets.Snippet
	c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("snippets"))

		b.ForEach(func(k, v []byte) error {
			sl = append(sl, snippets.Decode(v, "gob"))
			return nil
		})
		return nil
	})
	return sl
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
func main() {

	var filePath = flag.String("f", "", "filepath for the snippet")
	var lang = flag.String("l", "", "language of the snippet")
	var title = flag.String("t", "", "title for the snippet")
	var tags = flag.String("tags", "", "tags for the snippet separated by comma (no spaces)")
	var op = flag.String("op", "", "create|view")
	flag.Parse()

	// open connection and make sure snippets buckets exists or else gets created
	db, err := bolt.Open("snippets.db", 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("snippets"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	client := &Client{db: db}

	if *op == "create" {
		content, err := ioutil.ReadFile(*filePath)
		if err != nil {
			log.Fatal("File reading error: ", err)
			return
		}

		var t []string

		if len(*tags) != 0 {
			t = strings.Split(*tags, ",")
		}

		s := snippets.New(*title, *lang, string(content), time.Now(), t)
		client.CreateSnippet(s)
	} else if *op == "view" {
		sl := client.GetAllSnippets()

		for _, s := range sl {
			fmt.Println(s)
		}
	}

}
