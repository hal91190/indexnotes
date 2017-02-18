package notes

import (
	"bytes"
	"os"
	"log"
	"bufio"
)

// Nom du fichier d'index des notes
const (
	indexFilename = "00-index.adoc"
	indexTitle = "= Index des notes\n\n"
)

// Écrit un index des notes du répertoire.
func WriteNotesToDir(notes []Note, noteDirectory string) {
	noteFile, err := os.Create(noteDirectory + string(os.PathSeparator) + indexFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer noteFile.Close()

	noteWriter := bufio.NewWriter(noteFile)
	noteWriter.WriteString(indexTitle)

	c := make(chan string)

	noteSorters := [...]SortCriteria{ byTitle, byProject, byContext, byDate }
	for _, sorter := range noteSorters {
		go toStringBy(notes, sorter, c)
	}
	for range noteSorters {
		noteWriter.WriteString(<-c)
	}
	noteWriter.Flush()
}

func toStringBy(notes []Note, by SortCriteria, c chan string) {
	sortedNotes := make([]Note, len(notes))
	copy(sortedNotes, notes)
	by.Sort(sortedNotes)

	var result bytes.Buffer
	result.WriteString("== " + by.Criteria() + "\n")
	for _, note := range sortedNotes {
		result.WriteString("* ")
		result.WriteString(note.AdocString())
		result.WriteString("\n")
	}
	result.WriteString("\n")
	c <- result.String()
}
