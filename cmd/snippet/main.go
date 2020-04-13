package main

import (
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"log"
	"os"
	"snippets"
	"strings"
	"time"
)

func main() {

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createLang := createCmd.String("l", "", "Language")
	createTitle := createCmd.String("t", "", "Title")
	createFilePath := createCmd.String("f", "", "File path")
	createTags := createCmd.String("tags", "", "Tags")
	createClipboard := createCmd.Bool("c", true, "Copy from clipboard")

	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	viewLang := viewCmd.String("l", "", "Language")
	viewId := viewCmd.Int("id", 0, "ID")
	viewTag := viewCmd.String("tag", "", "Tag")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'create' or 'view' subcommands")
		os.Exit(1)
	}

	store := snippets.NewStore("snippets.db", "snippets")

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if !*createClipboard && *createFilePath == "" {
			fmt.Println("Clipboard copy (-c) or file path (-f) is mandatory")
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

		var content string
		if *createFilePath != "" {
			byteContent, err := ioutil.ReadFile(*createFilePath)
			if err != nil {
				log.Fatal("File reading error: ", err)
				os.Exit(1)
			}
			content = string(byteContent)
		}

		if *createClipboard {
			clipboard, err := clipboard.ReadAll()
			if err != nil {
				log.Fatal("Error while reading from clipboard: ", err)
				os.Exit(1)
			}
			content = clipboard
		}

		// parse tags
		var tags []string
		if len(*createTags) > 0 {
			tags = strings.Split(*createTags, ",")
		}
		// trim whitespaces
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
		}
		s := snippets.New(*createTitle, *createLang, content, time.Now(), tags)

		store.CreateSnippet(s)
	case "view":
		viewCmd.Parse(os.Args[2:])

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
			fmt.Println("You need to provide id (-id), lang (l) or tag (-t)")
			os.Exit(2)
		}
	default:
		fmt.Println("expected 'create' or 'view' subcommands")
		os.Exit(1)
	}

}
