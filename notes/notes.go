/*
Ce package permet le chargement des notes.
 */
package notes

import (
	"time"
	"fmt"
	"sort"
	"strings"
	"bytes"
)

// Représente les métadonnées sur une note.
type Note struct {
	filename string
	modTime  time.Time
	title    string
	project  string
	contexts []string
}

func (n Note) String() string {
	return fmt.Sprintf("%s (%s), %s, %s, %s", n.title, n.filename, n.modTime.Format(time.UnixDate), n.project, n.contexts)
}

func (n Note) AdocString() string {
	var result bytes.Buffer
	result.WriteString(fmt.Sprintf("link:++%s++[%s]", n.filename, n.title))
	result.WriteString(", " + n.modTime.Format(time.RFC822))
	if len(n.project) > 0 {
		result.WriteString(", *" + strings.ToTitle(n.project) + "*")
	}
	for _, context := range n.contexts {
		result.WriteString(", __" + context + "__")
	}
	return result.String()
}

type SortCriteria struct {
	criteria string
	less func(n1, n2 *Note) bool
}

func (sc SortCriteria) Criteria() string {
	return strings.ToTitle(sc.criteria)
}

func (sc SortCriteria) Sort(notes []Note) {
	ps := &noteSorter{
		notes: notes,
		SortCriteria: SortCriteria{
			criteria: sc.criteria,
			less: sc.less,
		},
	}
	sort.Sort(ps)
}

// Predefined sort functions
var (
	byTitle = SortCriteria{
		criteria: "par titre",
		less: func(n1, n2 *Note) bool {return n1.title < n2.title},
	}
	byContext = SortCriteria{
		criteria: "par contexte",
		less: func(n1, n2 *Note) bool {
			return strings.Join(n1.contexts, ", ") < strings.Join(n2.contexts, ", ")
		},
	}
	byProject = SortCriteria{
		criteria: "par projet",
		less: func(n1, n2 *Note) bool {return n1.project < n2.project},
	}
	byDate = SortCriteria{
		criteria: "par date",
		less: func(n1, n2 *Note) bool {return n1.modTime.Before(n2.modTime)},
	}
)

type noteSorter struct {
	notes []Note
	SortCriteria
}

func (s *noteSorter) Len() int {
	return len(s.notes)
}

func (s *noteSorter) Swap(i, j int) {
	s.notes[i], s.notes[j] = s.notes[j], s.notes[i]
}

func (s *noteSorter) Less(i, j int) bool {
	return s.less(&s.notes[i], &s.notes[j])
}
