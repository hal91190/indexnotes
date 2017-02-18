package main

import (
	"os"
	"github.com/hal91190/indexnotes/notes"
)

func main() {
	noteDirectory := os.Args[1]
	noteSet := notes.ReadNotesFromDir(noteDirectory)
	notes.WriteNotesToDir(noteSet, noteDirectory)
}
