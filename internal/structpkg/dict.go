package structpkg

import (
	"fmt"
	"regexp"
	"synodict-go/internal/stgpkg"
)

var WordRegex = regexp.MustCompile(`^[\p{L}\s-]+$`)

type Dict struct {
	graph *Graph
}

func NewDict() *Dict {
	return &Dict{graph: NewGraph()}
}

func getFormatSerializator(d *Dict, format string) func() []byte {
	formatHandlers := map[string]func() []byte{
		"gob":  d.graph.SerializeGob,
		"csv":  d.graph.SerializeCsv,
		"csvc": d.graph.SerializeCsvCondensed,
	}

	handler := formatHandlers[format]

	return handler
}

func getFormatDeserializator(format string) func(data []byte) (*Graph, error) {
	formatHandlers := map[string]func(data []byte) (*Graph, error){
		"gob":  DeserializeGob,
		"csv":  DeserializeCsv,
		"csvc": DeserializeCsvCondensed,
	}

	handler := formatHandlers[format]

	return handler
}

func logWordAlreadyExists(d *Dict, word string, log *[]error) bool {
	if d.WordExists(word) {
		*log = append(
			*log,
			fmt.Errorf("dictionary: word \"%s\" already exists", word),
		)

		return false
	}

	return true
}

func logWordNotFound(d *Dict, word string, log *[]error) bool {
	if !d.WordExists(word) {
		*log = append(
			*log,
			fmt.Errorf("dictionary: word \"%s\" does not exist", word),
		)

		return false
	}

	return true
}

func logWordNotMatch(word string, log *[]error) bool {
	if !WordRegex.MatchString(word) {
		*log = append(
			*log,
			fmt.Errorf("dictionary: word \"%s\" does not match the conditions", word),
		)

		return false
	}

	return true
}

func logAlreadyDirectSynonyms(d *Dict, a, b string, log *[]error) bool {
	if are, _ := d.AreDirectSynonyms(a, b); are {
		*log = append(
			*log,
			fmt.Errorf("dictionary: words \"%s\" and \"%s\" already are direct-linked synonyms", a, b),
		)

		return false
	}

	return true
}

func logWordsNotLinked(d *Dict, a, b string, log *[]error) bool {
	if are, _ := d.AreDirectSynonyms(a, b); !are {
		*log = append(
			*log,
			fmt.Errorf("dicationary: words \"%s\" and \"%s\" are not direct-linked synonyms", a, b),
		)

		return false
	}

	return true
}

func (d *Dict) AddSynonyms(words ...string) []error {
	filtered := []string{}
	var errs []error

	for i, word := range words {
		ok := logWordNotMatch(word, &errs)

		if i < len(words)-1 {
			logAlreadyDirectSynonyms(d, word, words[i+1], &errs)
		}

		if ok {
			filtered = append(filtered, words[i])
		}
	}

	switch len(filtered) {
	case 0:
		return errs

	case 1:
		d.AddWords(filtered[0])

	default:
		for i := 0; i < len(filtered)-1; i++ {
			d.graph.AddEdge(words[i], words[i+1])
		}
	}

	return errs
}

func (d *Dict) AddWords(words ...string) []error {
	var errs []error

	for _, word := range words {
		ok := logWordNotMatch(word, &errs) && logWordAlreadyExists(d, word, &errs)

		if ok {
			d.graph.AddVertex(word)
		}
	}

	return errs
}

func (d *Dict) RemoveWords(words ...string) []error {
	var errs []error

	for _, word := range words {
		ok := logWordNotFound(d, word, &errs)

		if ok {
			d.graph.RemoveVertex(word)
		}
	}

	return errs
}

func (d *Dict) UnlinkSynonyms(a, b string) []error {
	var errs []error
	ok := logWordNotFound(d, a, &errs) && logWordNotFound(d, b, &errs) && logWordsNotLinked(d, a, b, &errs)

	if ok {
		d.graph.RemoveEdge(a, b)
	}

	return errs
}

func (d *Dict) UnlinkSynonymsAndCleanup(a, b string) []error {
	var errs []error
	ok := logWordNotFound(d, a, &errs) && logWordNotFound(d, b, &errs) && logWordsNotLinked(d, a, b, &errs)

	if ok {
		d.graph.RemoveEdgeAndCleanup(a, b)
	}

	return errs
}

func (d *Dict) GetDirectSynonyms(word string) ([]string, error) {
	var errs []error
	ok := logWordNotFound(d, word, &errs)
	var result []string

	if ok {
		result = d.graph.GetNeighbors(word)
	}

	var err error

	if errs != nil {
		err = errs[0]
	}

	return result, err
}

func (d *Dict) GetSynomyms(word string) ([]string, error) {
	var errs []error
	ok := logWordNotFound(d, word, &errs)
	var result []string

	if ok {
		result = d.graph.GetConnectedVertices(word)
	}

	var err error

	if errs != nil {
		err = errs[0]
	}

	return result, err
}

func (d *Dict) SynonymCount(word string) (int, error) {
	var errs []error
	ok := logWordNotFound(d, word, &errs)
	var result int

	if ok {
		result = d.graph.ConnectedVertexCount(word)
	}

	var err error

	if errs != nil {
		err = errs[0]
	}

	return result, err
}

func (d *Dict) AreSynonyms(a, b string) (bool, []error) {
	var errs []error
	ok := logWordNotFound(d, a, &errs) && logWordNotFound(d, b, &errs)
	var result bool

	if ok {
		result = d.graph.AreConnected(a, b)
	}

	return result, errs
}

func (d *Dict) AreDirectSynonyms(a, b string) (bool, []error) {
	var errs []error
	ok := logWordNotFound(d, a, &errs) && logWordNotFound(d, b, &errs)
	var result bool

	if ok {
		result = d.graph.HasEdge(a, b)
	}

	return result, errs
}

func (d *Dict) WordExists(word string) bool {
	return d.graph.HasVertex(word)
}

func (d *Dict) GetWords() []string {
	return d.graph.GetVertices()
}

func (d *Dict) WordCount() int {
	return d.graph.Order()
}

func (d *Dict) GetSynonymGroups() [][]string {
	return d.graph.GetConnectivityGroups()
}

func (d *Dict) SynonymGroupCount() int {
	return len(d.GetSynonymGroups())
}

func (d *Dict) Clear() {
	d.graph = NewGraph()
}

func (d *Dict) Cleanup() {
	d.graph.Cleanup()
}

func (d *Dict) IsEmpty() bool {
	return d.WordCount() == 0
}

func (d *Dict) Export(path, format string) error {
	serializator := getFormatSerializator(d, format)

	if serializator == nil {
		return fmt.Errorf("export failed: format %s is not supported", format)
	}

	data := serializator()
	err := stgpkg.Write(data, path)

	return err
}

func (d *Dict) Import(path, format string) error {
	deserializator := getFormatDeserializator(format)

	if deserializator == nil {
		return fmt.Errorf("import failed: format %s is not supported", format)
	}

	data, err := stgpkg.Read(format)

	if err != nil {
		return err
	}

	graph, err := deserializator(data)

	if err != nil {
		return err
	}

	if d.graph.IsEmpty() {
		d.graph.FromGraphUnsafe(graph)
	} else {
		d.graph.MergeUnsafe(graph)
	}
	return nil
}
