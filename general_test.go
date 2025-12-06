package listmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Prop struct {
	UID         string `yaml:"uid"`
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	CSS         string `yaml:"css,omitzero,omitempty"`
}

type Metadata struct {
	Layout string `yaml:"layout"`
	Props  []Prop `yaml:"props"`
}

func TestUnmarshalAndMarshal(t *testing.T) {
	content, err := os.ReadFile("./data/test.md")
	if err != nil {
		t.Error(err)
	}

	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))

	var meta Metadata
	lists, err := Unmarshal(content, &meta)
	if err != nil {
		t.Error(err)
	}

	json.Marshal()

	_, err = json.MarshalIndent(lists, " ", " ")
	if err != nil {
		t.Error(err)
	}

	data, err := Marshal(lists)
	if err != nil {
		t.Error(err)
	}

	// log.Println(string(data))
	if diff := cmp.Diff(content, data); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

}
