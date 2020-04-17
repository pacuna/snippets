package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"snippets"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

func main() {

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createClipboard := createCmd.Bool("c", true, "Copy snippet from clipboard")
	createLang := createCmd.String("l", "", "Language of the snippet")
	createTitle := createCmd.String("t", "", "Title for the snippet")
	createFilePath := createCmd.String("f", "", "Path to file containing the snippet. Overrides clipboard")
	createTags := createCmd.String("tags", "", "Tags for the snippet")
	createdbPath := createCmd.String("dbpath", "", "Optional dbpath")

	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	viewLang := viewCmd.String("l", "", "Filter snippets by language")
	viewId := viewCmd.Int("id", 0, "ID of snippet to display")
	viewTag := viewCmd.String("tag", "", "Filter snippets by tag")
	viewdbPath := viewCmd.String("dbpath", "", "Optional dbpath")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteId := deleteCmd.Int("id", 0, "ID of snippet to delete")
	deletedbPath := deleteCmd.String("dbpath", "", "Optional dbpath")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'create' or 'view' subcommands")
		os.Exit(1)
	}

	user, err := user.Current()
	if err != nil {
		fmt.Println("Error initializing db")
		os.Exit(1)
	}

	// intialize db in user's home dir unless path is given
	store := &snippets.Store{}
	dbpath := user.HomeDir + "/snippets.db"

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *createdbPath != "" {
			dbpath = *createdbPath
		}
		store = snippets.NewStore(dbpath, "snippets")
	case "view":
		viewCmd.Parse(os.Args[2:])
		if *viewdbPath != "" {
			dbpath = *viewdbPath
		}
		store = snippets.NewStore(dbpath, "snippets")
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deletedbPath != "" {
			dbpath = *deletedbPath
		}
		store = snippets.NewStore(dbpath, "snippets")
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if createCmd.Parsed() {
		// no clipboard and no file passed is error
		if *createClipboard == false && *createFilePath == "" {
			createCmd.PrintDefaults()
			os.Exit(1)
		}

		if *createLang == "" {
			createCmd.PrintDefaults()
			os.Exit(1)
		}
		if *createTitle == "" {
			createCmd.PrintDefaults()
			os.Exit(1)
		}

		var content string

		// read from clipboard
		if *createClipboard {
			clipboard, err := clipboard.ReadAll()
			if err != nil {
				log.Fatal("Error while reading from clipboard: ", err)
				os.Exit(1)
			}
			content = clipboard
		}

		// read from file overrides clipboard which is default to true
		if *createFilePath != "" {
			byteContent, err := ioutil.ReadFile(*createFilePath)
			if err != nil {
				log.Fatal("File reading error: ", err)
				os.Exit(1)
			}
			content = string(byteContent)
		}

		// parse tags and trim whitespaces
		var tags []string
		if len(*createTags) > 0 {
			tags = strings.Split(*createTags, ",")
		}
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
		}

		s := snippets.New(*createTitle, *createLang, content, time.Now(), tags)
		store.CreateSnippet(s)
	}

	if viewCmd.Parsed() {
		if *viewId != 0 {
			s, _ := store.GetSnippetByID(*viewId)
			fmt.Println(s.Content)
			return
		}

		if *viewTag != "" {
			sl, _ := store.ListSnippetsByTag(*viewTag)
			for _, s := range sl {
				fmt.Printf("[%d] - %s\n", s.ID, s.Title)
			}
			return
		}

		if *viewLang != "" {
			sl, _ := store.ListSnippetsByLang(*viewLang)
			for _, s := range sl {
				fmt.Printf("[%d] - %s\n", s.ID, s.Title)
			}
			return
		}

		if *viewId == 0 && *viewLang == "" && *viewTag == "" {
			viewCmd.PrintDefaults()
			os.Exit(2)
		}
	}

	if deleteCmd.Parsed() {
		if *deleteId == 0 {
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}
		err := store.DeleteSnippet(*deleteId)
		if err != nil {
			fmt.Println("Error deleting snippet: ", err)
			os.Exit(1)
		}
		return
	}

}
