package abook

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	programKey = "program"
	versionKey = "version"
	nameKey    = "name"
	emailKey   = "email"
)

//A block is inteded as all the info
//about someone (name not included)
type AbookBlock map[string]string

//An entry is formed by the name of
//a contact (which is the index) and
//an AbookBlock
type AbookDB map[string]AbookBlock

//The Whole Abook object
type Abook struct {
	DB               AbookDB
	r                io.Reader
	currentName      string
	program, version string
}

//Create an abook object
func NewAbook(r io.Reader) *Abook {
	db := make(AbookDB)
	abook := &Abook{db, r, "", "", ""}
	return abook
}

//Split the line with '='
func splitKeyValue(r rune) bool {
	if r == '=' {
		return true
	} else {
		return false
	}
}

//Prints the structure as the standard abook format
func (abook *Abook) String() string {
	out := fmt.Sprintf("[format]\nprogram=%s\nversion=%s\n\n",
		abook.program, abook.version)
	counter := 0
	for k, v := range abook.DB {
		out += fmt.Sprintf("[%d]\n", counter)
		out += fmt.Sprintf("name=%s\n", k)
		counter++
		for kblock, vblock := range v {
			out += fmt.Sprintf("%s=%s\n", kblock, vblock)
		}
		out += "\n"
	}
	return out
}

//Parses a single line of the abook DB
func (abook *Abook) parseValue(line *string) {
	fields := strings.FieldsFunc(*line, splitKeyValue)
	if len(fields) != 2 {
		return
	}
	key, val := fields[0], fields[1]
	switch key {
	case programKey:
		abook.program = val
	case versionKey:
		abook.version = val
	case nameKey:
		if _, ok := abook.DB[val]; !ok {
			abook.DB[val] = make(AbookBlock)
			abook.currentName = val
		}
	//add the value to the abookDB
	default:
		if oldval, ok := abook.DB[abook.currentName][key]; ok {
			abook.DB[abook.currentName][key] = oldval + ", " + val
		} else {
			abook.DB[abook.currentName][key] = val
		}
	}
}

//Populates the Abook structure
func (abook *Abook) Parse() error {
	scanner := bufio.NewScanner(abook.r)
	for scanner.Scan() {
		line := strings.TrimLeft(scanner.Text(), " ")
		if len(line) > 0 && line[0] != '#' && line[0] != '[' {
			abook.parseValue(&line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading file:", err)
		return err
	}
	return nil
}

//Removes any abook entry which is either a orphan (no mail addr)
//or starts with a supplied prefix, all done
//accordingly with the supplied parameters
func (abook *Abook) Remove(orphan bool, removeWith string) {
	for k, v := range abook.DB {
		mailPresent := false
		containsRemove := false
		for kblock, vblock := range v {
			if kblock == emailKey {
				mailPresent = true
			}
			if strings.Contains(vblock, removeWith) {
				containsRemove = true
			}
		}
		if (orphan && !mailPresent) || (removeWith != "" && containsRemove) {
			delete(abook.DB, k)
		}
	}
}

//Write the abook contents into
//a supplied io.Writer
func (abook *Abook) WriteTo(w io.Writer) (n int64, err error) {
	stringReader := strings.NewReader(abook.String())
	n, err = stringReader.WriteTo(w)
	return n, err
}
