package entity

import (
	"errors"
	"fmt"
	"fs/lib"
	"strings"
)

type Node struct {
	Id       string `json:"id" gorm:"id"`
	Type     int    `json:"type" gorm:"type"`
	Order    int    `json:"order" gorm:"order"`
	Parent   *Node
	Name     string `json:"name" gorm:"name"`
	Path     string `json:"path" gorm:"path"`
	Data     string `json:"data" gorm:"data"`
	Children []*Node
	Root     *Node
}

func (g *Node) TotalChildren() int {
	if g == nil {
		return 0
	}
	res := 1
	for _, x := range g.Children {
		res += x.TotalChildren()
	}
	return res
}

func (g *Node) UpdatePreOrderNumber() {
	res := []string{}
	UpdateOrderRecursive(g, &res)
}

func UpdateOrderRecursive(node *Node, res *[]string) {
	if node == nil {
		return
	}
	tmp := append(*res, fmt.Sprintf("%s", node.Name))
	*res = tmp
	node.Order = len(*res) - 1

	for _, x := range node.Children {
		UpdateOrderRecursive(x, res)
	}
	tmp2 := append(*res, ")")
	*res = tmp2
}

//TraverseByPath navigate to the destination folder in folderPath,
//
//assume that the path has the form a/b/c, 'a' is the children ItemView of the current ItemView
func (g *Node) TraverseByPath(path string) (*Node, error) {
	// check path is valid
	path = lib.CleanPathToFileOrFolder(path)
	if path == "" {
		return &Node{}, ErrorInvalidPath
	}
	if path == "/" {
		if g.IsRoot() {
			return g, nil
		}
		return g.Root, nil
	} else if lib.IsCurrentFolder(path) {
		return g, nil
	} else if lib.IsParentFolder(path) {
		if g.IsRoot() {
			return g, nil
		} else {
			return g.Parent, nil
		}
	}

	folderNames := strings.Split(path, "/")

	var curr *Node
	if strings.HasPrefix(path, "/") {
		if g.IsRoot() {
			curr = g
		} else {
			curr = g.Root
		}
		folderNames = folderNames[1:]
	} else {
		curr = g
	}

	//fmt.Printf("folder names: %v\n", folderNames)
	//fmt.Printf("curr: %v, names: %v\n", curr, folderNames)

	for i, folderName := range folderNames {
		if folderName == "." {
			//fmt.Printf("stay at current node %s\n", curr.Name)
			continue
		} else if folderName == ".." {
			if !curr.IsRoot() {
				//fmt.Printf("go to parent node %s\n", curr.Parent.Name)
				curr = curr.Parent
			}
		} else {
			noFolder := true
			for _, x := range curr.Children {
				if x.Type == 0 && x.Name == folderName {
					curr = x
					noFolder = false

					if i == len(folderNames)-1 {
						//fmt.Printf("found  at %s\n", curr.Name)
						return curr, nil
					}
					break
				}
			}

			if noFolder {
				return &Node{}, ErrorInvalidPath
			}
		}
	}

	return curr, nil
}

//InsertByPath insert a new ItemView at the specified folderPath,
//
//assume that the folderPath has the form a/b/c which included the new ItemView name
// 'a' is the children ItemView of the current ItemView
func (g *Node) InsertByPath(path string, node Node) error {
	parentNode, newNodeName, err := g.getParentNodeAndDestNodeName(path)
	if err != nil {
		return err
	}
	if len(parentNode.Children) != 0 {
		for _, x := range parentNode.Children {
			if x.Name == newNodeName {
				return ErrorPathAndItemNameNotConsistent
			}
		}
	}

	node.Path = path
	parentNode.Children = append(parentNode.Children, &node)
	node.Parent = parentNode
	return nil
}

//FindByName navigate to the destination folder in folderPath,
//
//assume that the folderPath has the form a/b/c, 'a' is the children ItemView of the current ItemView
func (g *Node) FindByName(name string) []*Node {
	res := []*Node{}

	if strings.Contains(g.Name, name) {
		res = append(res, g)
	}

	var err error
	var curr []*Node
	for _, x := range g.Children {
		curr = x.FindByName(name)
		if err == nil {
			res = append(res, curr...)
		}
	}

	return res
}

//DeleteByPath delete the ItemView at path, as long as it could be navigated from current ItemView
//
//assume that the path has the form a/b/c, 'a' is the children ItemView of the current ItemView
//
// currently support alias: '.', '..', '/'
func (g *Node) DeleteByPath(path string) error { // TODO: refactor 'view' name
	if lib.IsCurrentFolder(path) || lib.IsParentFolder(path) {
		return ErrorCantDeleteCurrentNodeOrParent
	} else if path == "/" {
		return ErrorCantDeleteRoot
	}
	parentNode, destNodeName, err := g.getParentNodeAndDestNodeName(path)
	if err != nil {
		return err
	}

	n := len(parentNode.Children)
	for i, x := range parentNode.Children {
		if x.Name == destNodeName {
			parentNode.Children[i], parentNode.Children[n-1] = parentNode.Children[n-1], parentNode.Children[i]
			parentNode.Children = parentNode.Children[:n-1]
			return nil
		}
	}

	return ErrorNoFileOrFolderFound
}

func (g *Node) GetChildWithName(name string) (*Node, bool) {
	if g == nil {
		return &Node{}, false
	}
	for _, x := range g.Children {
		if x.Name == name {
			return x, true
		}
	}
	return &Node{}, false
}

// =========== HELPER FUNCTIONS ==================

//ParsePathToItemViewNameSlice check if path has the form a/b or /a/b/c
//
// and each folder name is valid, then convert to a slice of ItemViews
func ParsePathToItemViewNameSlice(path string) ([]string, error) {
	if path == "/" {
		return []string{}, nil
	}
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	} else if strings.HasPrefix(path, "./") {
		path = path[2:]
	} else if strings.HasPrefix(path, "../") {
		path = path[3:]
	} else if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	parts := strings.Split(path, "/")
	fmt.Printf("PARTS: %d\n", len(parts))

	for _, x := range parts {
		if lib.CheckName(x) == false {
			return []string{}, ErrorNameNotValid
		}
	}

	return parts, nil
}

func (g *Node) getParentNodeAndDestNodeName(path string) (*Node, string, error) {
	path = lib.CleanPathToFileOrFolder(path)
	if path == "" {
		return &Node{}, "", ErrorInvalidPath
	}
	parentPath := lib.GetParentPath(path)
	parentNode, err := g.TraverseByPath(parentPath)
	if err != nil {
		return &Node{}, "", ErrorInvalidPath
	}
	nodeNames := strings.Split(path, "/")
	destNodeName := nodeNames[len(nodeNames)-1]
	if !lib.CheckName(destNodeName) { // file name could be '.' or '..'
		return &Node{}, "", ErrorInvalidPath
	}
	return parentNode, destNodeName, nil
}

func (g *Node) Print(indent int, tmp string) {
	tmp = GetIndentString(indent)
	tmp += "+--"
	tmp += g.Name
	tmp += "/"
	tmp += "\n"

	fmt.Printf("%s", tmp)

	if g.Type == 0 {
		for _, x := range g.Children {
			x.Print(indent+1, tmp)
		}
	}
}

func (g *Node) PrintOrder(indent int, tmp string) {
	tmp = GetIndentString(indent)
	tmp += "+--"
	tmp += g.Name
	tmp += "/"
	tmp += fmt.Sprintf("(%d)", g.Order)
	tmp += "\n"

	fmt.Printf("%s ", tmp)

	if g.Type == 0 {
		for _, x := range g.Children {
			x.PrintOrder(indent+1, tmp)
		}
	}
}

func GetIndentString(indent int) string {
	res := ""
	for i := 0; i < indent; i++ {
		res += "|  "
	}
	return res
}

func (g *Node) ToItem() (Item, error) {
	if g.CheckValid() != nil {
		return Item{}, errors.New("item invalid")
	}
	parentId := ""
	if g.Parent != nil {
		parentId = g.Parent.Id
	}
	return Item{
		Id:       g.Id,
		ItemType: g.Type,
		OrderNum: g.Order,
		ParentId: parentId,
		Name:     g.Name,
		Data:     g.Data,
	}, nil
}

func (g *Node) ToString() string {
	if g.Parent != nil {
		return fmt.Sprintf("Id: %s, Name: %s, ParentId: %s, Order: %d\n", g.Id, g.Name, g.Parent.Id, g.Order)
	}
	return fmt.Sprintf("Id: %s, Name: %s, ParentId: %s, Order: %d\n", g.Id, g.Name, "", g.Order)
}

func (g *Node) IsRoot() bool {
	if g.Parent == nil {
		return true
	}
	return false
}

func (g *Node) CheckValid() error {
	if g == nil {
		return nil
	}
	return nil
}
