package listmd

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

func collectNode(nodes []*Node) (newNodes []*Node) {
	for _, node := range nodes {
		uidSlice := strings.Split(node.UID, "_")
		node.collectionIndex, _ = strconv.Atoi(uidSlice[1])

		newNodes = append(newNodes, node)
		if len(node.Children) > 0 {
			newNodes = append(newNodes, collectNode(node.Children)...)
		}
	}
	return
}

func flattenJSONToNode(nodeCollections [][]*Node) (flatNodes []*Node) {
	for _, nodeGroup := range nodeCollections {
		flatNodes = append(flatNodes, collectNode(nodeGroup)...)
	}
	return
}

// MarshalJSON returns list markdown bytes from JSON v.
func MarshalJSON(v []byte) (data []byte, err error) {
	var list ListMd
	err = json.Unmarshal(v, list)
	if err != nil {
		return
	}

	return Marshal(list)
}

// Marshal returns list markdown bytes of v.
func Marshal(v ListMd) (data []byte, err error) {
	flatList := flattenJSONToNode(v.Lists)
	builder := strings.Builder{}

	if v.Metadata != nil {
		builder.WriteString(breakString)
		builder.WriteString("\n")
		metadataString, errYaml := yaml.Marshal(v.Metadata)
		if errYaml != nil {
			return data, errYaml
		}
		builder.WriteString(string(metadataString))

		builder.WriteString(breakString)
		builder.WriteString("\n")
	}

	sortNodesByID(flatList)

	collectionIndex := 0
	for _, list := range flatList {
		if list.collectionIndex > collectionIndex {
			builder.WriteString(breakString)
			builder.WriteString("\n")
		}

		contentIndent := (list.level + 1) * len(bulletPrefix)
		if list.level > 0 {
			builder.WriteString(strings.Repeat(" ", list.level*2))
		}
		builder.WriteString(bulletPrefix)

		splitMd := strings.Split(list.Markdown, "\n")

		for mdIndex, md := range splitMd {
			if mdIndex > 0 {
				builder.WriteString(strings.Repeat(" ", contentIndent))
			}

			builder.WriteString(md)
			builder.WriteString("\n")
		}

		collectionIndex = list.collectionIndex
	}

	data = []byte(builder.String())
	return
}
