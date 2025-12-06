package listmd

type ListMd struct {
	Metadata any       `json:"metadata"`
	Lists    [][]*Node `json:"lists,omitempty"`
}

type anyMeta struct{}

type Node struct {
	Id              int `json:"id"`
	parentId        int
	level           int
	rawContents     []string
	collectionIndex int
	Markdown        string  `json:"markdown"`
	Children        []*Node `json:"children,omitempty"`
	UID             string  `json:"uid"`
}

type NodePosition struct {
	CollectionIndex int
	Path            []int
}
