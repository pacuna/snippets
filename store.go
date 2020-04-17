package snippets

import (
	"bytes"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type Store struct {
	client *bolt.DB
	bucket []byte
	path   string
}

func NewStore(path, bucket string) *Store {
	var (
		db *bolt.DB
	)

	s := &Store{
		client: db,
		bucket: []byte(bucket),
		path:   path,
	}
	return s
}

func (s *Store) getHandle() (*bolt.DB, error) {
	var (
		db  *bolt.DB
		err error
	)
	if db, err = bolt.Open(s.path, 0600, nil); err != nil {
		return nil, err
	}
	s.client = db

	return s.client, nil
}

func (s *Store) releaseHandle() {
	s.client.Close()
}

func (s *Store) CreateSnippet(snippet *Snippet) error {

	var (
		db  *bolt.DB
		err error
	)

	if db, err = s.getHandle(); err != nil {
		return err
	}
	defer s.releaseHandle()

	// create bucket for this specific language if it doesn't exist
	db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("snippets")); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("snippets"))
		id, _ := b.NextSequence()
		snippet.ID = int(id)

		// keys for different queries (by id, by lang, by tag and by lang+tag)
		var keys []string
		keys = append(keys, string(snippet.ID))
		keys = append(keys, snippet.Language+":"+string(snippet.ID))
		for _, tag := range snippet.Tags {
			keys = append(keys, tag+":"+string(snippet.ID))
			//keys = append(keys, snippet.Language + ":" + tag + ":" + string(snippet.ID))
		}

		for _, key := range keys {
			if err := b.Put([]byte(key), Encode(snippet, "gob")); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *Store) ListSnippetsByLang(lang string) ([]*Snippet, error) {
	var (
		db  *bolt.DB
		err error
		sl  []*Snippet
	)

	if db, err = s.getHandle(); err != nil {
		return nil, err
	}
	defer s.releaseHandle()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("snippets")).Cursor()
		prefix := []byte(lang)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			s := Decode(v, "gob")
			sl = append(sl, s)
		}
		return nil
	})
	return sl, nil
}

func (s *Store) ListSnippetsByTag(tag string) ([]*Snippet, error) {
	var (
		db  *bolt.DB
		err error
		sl  []*Snippet
	)

	if db, err = s.getHandle(); err != nil {
		return nil, err
	}
	defer s.releaseHandle()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("snippets")).Cursor()
		prefix := []byte(tag)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			s := Decode(v, "gob")
			sl = append(sl, s)
		}
		return nil
	})
	return sl, nil
}

func (s *Store) GetSnippetByID(id int) (*Snippet, error) {
	var (
		db      *bolt.DB
		err     error
		snippet *Snippet
	)

	if db, err = s.getHandle(); err != nil {
		return nil, err
	}
	defer s.releaseHandle()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("snippets"))
		v := b.Get([]byte(string(id)))
		if v != nil {
			snippet = Decode(v, "gob")
		}
		return nil
	})
	return snippet, nil
}

func (s *Store) DeleteSnippet(id int) error {

	snippet, _ := s.GetSnippetByID(id)
	var (
		db  *bolt.DB
		err error
	)

	if db, err = s.getHandle(); err != nil {
		return err
	}
	defer s.releaseHandle()

	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("snippets"))

		// keys to delete
		var keys []string
		keys = append(keys, string(snippet.ID))
		keys = append(keys, snippet.Language+":"+string(snippet.ID))
		for _, tag := range snippet.Tags {
			keys = append(keys, tag+":"+string(snippet.ID))
		}

		for _, key := range keys {
			if err := b.Delete([]byte(key)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
