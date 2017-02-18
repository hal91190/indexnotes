package notes

import (
	"os"
	"sync"
	"log"
	"bufio"
	"strings"
	"io/ioutil"
	"sort"
)

// Préfixes utilisés lors de l'analyse des fichiers pour
// identifier les champs à traiter.
const (
	titlePrefix = "= "
	contextPrefix = ":context: "
	projectPrefix = ":project: "
	contextSeparator = ", "
)

// Réalise l'analyse du répertoire des notes et
// retourne un slice de Note.
func ReadNotesFromDir(noteDirectory string) []Note {
	filenames, err := ioutil.ReadDir(noteDirectory)
	if err != nil {
		log.Fatal(err)
	}

	numberOfNotes := len(filenames)
	notes := make([]Note, numberOfNotes)

	var wg sync.WaitGroup
	wg.Add(numberOfNotes)

	for idx, noteFilename := range filenames {
		go readNoteFromFile(noteDirectory, noteFilename, &notes[idx], &wg)
	}
	wg.Wait()

	return notes
}

// Extrait d'un fichier de note les champs de metadonnées.
func readNoteFromFile(noteDirectory string, noteFilename os.FileInfo, note *Note, wg *sync.WaitGroup) {
	defer wg.Done()

	note.filename = noteFilename.Name()
	note.modTime = noteFilename.ModTime()

	noteFile, err := os.Open(noteDirectory + string(os.PathSeparator) + noteFilename.Name())
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(noteFile)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		if strings.HasPrefix(line, titlePrefix) {
			note.title = getTitle(line)
		} else if strings.HasPrefix(line, contextPrefix) {
			note.contexts = getContexts(line)
		} else if strings.HasPrefix(line, projectPrefix) {
			note.project = getProject(line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getTitle(line string) string {
	return strings.TrimPrefix(line, titlePrefix)
}

func getContexts(line string) []string {
	contexts := strings.ToLower(line)
	contexts = strings.TrimPrefix(contexts, contextPrefix)
	contextsAsSlice := strings.Split(contexts, contextSeparator)
	sort.Strings(contextsAsSlice)
	return contextsAsSlice
}

func getProject(line string) string {
	project := strings.ToLower(line)
	return strings.TrimPrefix(project, projectPrefix)
}
