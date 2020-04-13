# snippets

Small utility to manage code snippets.

Commands:

To create a snippet from the system's clipboard content:

`snippet create -c -t "Cool snippet" -l "python" -tags strings,algorithms`

Or from a file:

`snippet create -f /path/to/file -t "Cool snippet" -l "python" -tags strings,algorithms`

View all snippets for a language:

`snippet view -l python`

View all snippets for a tag:

`snippet view -tag strings`

View the content of a snippet:

`snippet view -id snippet_id`
