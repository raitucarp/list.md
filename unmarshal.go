package listmd

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/goccy/go-yaml"
)

var breakString = "---"
var bulletPrefix = "- "
var spaceIndentLength = 2

func countLevel(s string) int {
	return (len(s) - len(strings.TrimSpace(s))) / spaceIndentLength
}

// UnmarshalJSON returns ListMd structure from JSON bytes.
func UnmarshalJSON(v []byte, metaType any) (lists ListMd, err error) {
	data, err := MarshalJSON(v)
	if err != nil {
		return
	}

	return UnmarshalJSON(data, metaType)
}

// Unmarshal returns ListMd structure from list markdown bytes .
func Unmarshal(data []byte, metaType any) (lists ListMd, err error) {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)

	var meta *anyMeta
	metaContent := ""
	metaMode := false
	listColIndex := 0
	listIndex := 0
	listsRaw := []Node{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "&nbsp;", " ")
		line = strings.ReplaceAll(line, "\t", strings.Repeat(" ", spaceIndentLength))
		if line == breakString && !metaMode && meta == nil {
			metaMode = true
			meta = &anyMeta{}
			continue
		}

		if line == breakString && metaMode && meta != nil {
			metaMode = false
			continue
		}

		if line == breakString && !metaMode && meta != nil {
			listColIndex++
			continue
		}

		if metaMode {
			metaContent += line + "\n"
			continue
		}

		lineTrimmed := strings.TrimSpace(line)
		if strings.HasPrefix("\\-", line) {
			lineTrimmed = strings.ReplaceAll(lineTrimmed, "\\-", "-")
		}

		if strings.HasPrefix(lineTrimmed, bulletPrefix) {
			level := countLevel(line)
			newList := Node{
				rawContents:     []string{strings.Trim(lineTrimmed, bulletPrefix)},
				level:           level,
				collectionIndex: listColIndex,
				parentId:        -1,
				Id:              listIndex,
			}

			if listIndex == 0 {
				listsRaw = append(listsRaw, newList)
				listIndex++
				continue
			}

			lastRaw := listsRaw[len(listsRaw)-1]
			if level > lastRaw.level {
				newList.parentId = lastRaw.Id
			}

			if level == lastRaw.level {
				newList.parentId = lastRaw.parentId
			}

			if level < lastRaw.level {
				lastRawSameLevelIndex := 0
				for listIndex, list := range listsRaw {
					if list.level == level {
						lastRawSameLevelIndex = listIndex
					}
				}
				newList.parentId = listsRaw[lastRawSameLevelIndex].parentId
			}

			listsRaw = append(listsRaw, newList)
			listIndex++
			continue
		}

		raw := ""
		if strings.TrimSpace(line) == "" {
			raw = ""
		} else {
			raw = strings.TrimSpace(line)
		}

		if len(listsRaw) > 0 {
			listsRaw[len(listsRaw)-1].rawContents = append(listsRaw[len(listsRaw)-1].rawContents, raw)
		}
	}

	if errScan := scanner.Err(); errScan != nil {
		return lists, errScan
	}

	for listIndex, list := range listsRaw {
		list.Markdown = strings.Join(list.rawContents, "\n")
		listsRaw[listIndex] = list
	}

	listTree := buildTree(listsRaw)

	lists.Lists = listTree
	err = yaml.Unmarshal([]byte(metaContent), metaType)
	if err != nil {
		return
	}
	lists.Metadata = metaType

	return
}
