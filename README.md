# `list.md` — Parse Markdown Lists + Frontmatter Into Structured JSON (Go Library)

`list.md` is a lightweight Go library for **parsing Markdown list documents**—including multi-level nested lists, checkboxes, fenced code blocks, and YAML frontmatter—into a clean, structured Go type (`ListMd`).
It also supports the reverse operation: **convert structured JSON or `ListMd` back into markdown list format**.

Its goal is to make Markdown list documents behave like a **typed data structure**, enabling workflows such as:

* Markdown → JSON (for APIs, rendering engines, note-taking apps, static sites)
* JSON → Markdown (for editors, generators, automation)
* Frontmatter + list blocks → structured trees of nodes with stable UIDs
* Round-trippable parsing for building tools **like Notion, Obsidian, workflow editors, or checklist systems**.

---

## ✨ Features

* ✅ Parse multi-section Markdown separated by `---`
* ✅ Parse YAML frontmatter (`metadata`)
* ✅ Parse nested list structures (`- item`, children, grandchildren, etc.)
* ✅ Preserve fenced code blocks inside list items
* ✅ Extract checkboxes (`- [ ]`, `- [x]`)
* ✅ Auto-generate stable UIDs for each list node
* ✅ Convert JSON → markdown and markdown → JSON
* ✅ Simple Go API: `Marshal`, `Unmarshal`, `MarshalJSON`, `UnmarshalJSON`

---

## Installation

```bash
go get github.com/raitucarp/list.md
```

---

## 🚀 Quick Example

Given a markdown file like:

````md
---
layout: vertical
props:
- uid: item_0_1_1
  label: Custom Label
  description: This is actually custom label
  css: |
    .item {
      etawer
    }
---
- test
- this is test
  - # what do you want to test
    ```yml
    ewarwaer: ewrarar
    ```
````

You can parse it to a structured `ListMd`:

```go
data, _ := os.ReadFile("example.md")

lists, err := listmd.Unmarshal(data, nil)
if err != nil {
    panic(err)
}

fmt.Println(lists.Metadata)
fmt.Println(lists.Lists[0][1].Children[0].Markdown)
```

You can also convert JSON → Markdown:

```go
jsonData := []byte(`{"metadata": {...}, "lists": [...]}`)
md, err := listmd.MarshalJSON(jsonData)
```

---

## 📘 Output Structure (`ListMd`)

`list.md` parses into the following structure:

```go
type ListMd struct {
	Metadata any       `json:"metadata"`
	Lists    [][]*Node `json:"lists,omitempty"`
}

type Node struct {
	Id       int      `json:"id"`
	Markdown string   `json:"markdown"`
	Children []*Node  `json:"children,omitempty"`
	UID      string   `json:"uid"`
}
```

* `Metadata` → your YAML frontmatter.
* `Lists` → a slice of list-blocks (`---` separated groups).
* `Node` represents a single markdown list item.
* `Markdown` preserves multi-line content, code fences, inline formatting, etc.
* `UID` is auto-generated or picked up from metadata if provided.

---

## API Overview

### Parse JSON → Markdown

```go
func MarshalJSON(v []byte) ([]byte, error)
```

### Parse `ListMd` → Markdown

```go
func Marshal(v ListMd) ([]byte, error)
```

---

### Parse Markdown → `ListMd`

```go
func Unmarshal(data []byte, metaType any) (ListMd, error)
```

### Parse JSON → `ListMd`

```go
func UnmarshalJSON(v []byte, metaType any) (ListMd, error)
```

---

## 📚 Full Example

Markdown input:

```md
---
layout: vertical
props:
  - uid: item_0_1
    label: well
---
- this is
- [ ] true
- [x] something
  [x] continued
- for you test
  - nested item
```

Parsed result (excerpt):

```json
{
  "metadata": {
    "layout": "vertical",
    "props": [...]
  },
  "lists": [
    [
      {
        "id": 0,
        "markdown": "this is",
        "uid": "item_0_0"
      },
      {
        "id": 1,
        "markdown": "[ ] true",
        "uid": "item_0_1"
      },
      {
        "id": 2,
        "markdown": "[x] something\n[x] continued",
        "uid": "item_0_2"
      }
    ]
  ]
}
```
---

## Design Philosophy

* Keep markdown human-friendly.
* Make list structures **typed**, predictable, and machine-readable.
* Preserve every character the user writes—especially code blocks.
* Never lose information when converting back and forth.

---

## Roadmap

* [ ] Add options for UID strategy
* [ ] Render trees with custom templates
* [ ] Plugin system for metadata extraction
* [ ] Optional strict mode for malformed markdown

---

## Contributing

PRs welcome!
If you find unexpected markdown structures, please open an issue with examples.

---

## 📄 License

MIT
