package fs

import (
	"fmt"
	"fs/app/entity"
	"fs/app/repositories"
	"fs/app/usecase"
	"fs/lib"
	"log"
	"sort"
	"strconv"
	"strings"
)

type treeUseCaseImpl struct {
	Node        *entity.Node
	Total       int
	CurrentNode *entity.Node
	ItemRepo    repositories.ItemRepository
}

func NewTreeUseCase(node *entity.Node, repo repositories.ItemRepository) usecase.TreeUseCase {
	items, _ := repo.GetAll()
	return &treeUseCaseImpl{node, len(items), node, repo}
}

// Import directly from db into the current tree
func (g *treeUseCaseImpl) Import() error {
	items, err := g.ItemRepo.GetAll()
	if err != nil {
		return ErrorFailToGetItems
	}

	rootItem, err := g.ItemRepo.GetByOrderNum(0)
	if err != nil {
		return ErrorFailToGetItems
	}
	g.Node.Id = rootItem.Id
	g.Node.Order = rootItem.OrderNum
	g.Node.Name = rootItem.Name
	g.Node.Type = rootItem.ItemType
	g.Node.Children = []*entity.Node{}
	g.Node.Path = "/"

	root := g.Node

	memo := map[string]*entity.Node{}

	for _, x := range items {
		memo[x.Id] = &entity.Node{Id: x.Id}
		//fmt.Printf("(id: %s, %s) ", x.Id, x.ParentId)
	}
	memo[root.Id] = g.Node

	sort.SliceStable(items, func(i, j int) bool {
		a, _ := strconv.Atoi(items[i].ParentId)
		b, _ := strconv.Atoi(items[j].ParentId)
		return a < b
	})

	for _, item := range items[1:] {
		parent, ok := memo[item.ParentId]
		if !ok {
			fmt.Printf("parent: %s\n", item.ParentId)
			return ErrorParentChildOrdering
		}

		// temporary ignore path
		node := memo[item.Id]
		node.Type = item.ItemType
		node.Parent = parent
		node.Name = item.Name
		node.Data = item.Data
		node.Children = []*entity.Node{}
		node.Root = g.Node

		parent.Children = append(parent.Children, node)
	}

	g.Node.UpdatePreOrderNumber()

	return nil
}

// Export current tree to db
func (g *treeUseCaseImpl) Export() error {
	queue := []*entity.Node{}
	for _, x := range g.Node.Children {
		queue = append(queue, x)
	}

	for len(queue) > 0 {
		node := queue[0]
		item, err := node.ToItem()
		if err != nil {
			return ErrorConvertNodeToItem
		}

		//log.Printf("node: %v, item: %v\n", node, item)
		log.Printf("node: %v, item: %v\n", node.ToString(), item.ToString())

		err = g.ItemRepo.Update(item)
		if err != nil {
			return err
		}
		queue = queue[1:]

		for _, x := range node.Children {
			queue = append(queue, x)
		}
	}
	return nil
}

func (g *treeUseCaseImpl) Navigate(folderPath string) (*entity.Node, error) {
	res, err := g.Node.TraverseByPath(folderPath)
	if err == nil {
		g.CurrentNode = res
	}

	return res, err
}

func (g *treeUseCaseImpl) Create(path string, data string) error {
	// check if path is valid
	path = lib.CleanPathToFileOrFolder(path)
	if path == "" {
		return ErrorInvalidPath
	}
	parentPath := lib.GetParentPath(path)
	parentNode, err := g.Navigate(parentPath)
	if err != nil {
		return ErrorInvalidPath
	}
	nodeNames := strings.Split(path, "/")
	if !lib.CheckName(nodeNames[len(nodeNames)-1]) { // file name could be '.' or '..'
		return ErrorInvalidFileOrFolderName
	}
	newNodeName := nodeNames[len(nodeNames)-1]

	var node entity.Node
	if len(data) == 0 {
		node = entity.Node{Parent: parentNode, Name: newNodeName, Type: 0}
	} else {
		node = entity.Node{Parent: parentNode, Name: newNodeName, Type: 1}
		node.Data = data
	}

	item, err := node.ToItem()
	if err != nil {
		return ErrorConvertNodeToItem
	}

	err = g.ItemRepo.InsertNewItem(item)
	if err != nil {
		return ErrorFailToSaveItemToDB
	}

	err = g.Node.InsertByPath(path, node)
	if err != nil {
		return ErrorFailToInsertNodeToTree
	}

	return nil
}

//ShowData shows the content of the file at filePath
func (g *treeUseCaseImpl) ShowData(filePath string) (string, error) {
	filePath = lib.CleanPathToFileOrFolder(filePath)
	if filePath == "" {
		return "", ErrorInvalidPath
	}
	parentPath := lib.GetParentPath(filePath)
	nodeNames := strings.Split(filePath, "/")
	if !lib.CheckName(nodeNames[len(nodeNames)-1]) { // file name could be '.' or '..'
		return "", ErrorInvalidFileOrFolderName
	}
	newNodeName := nodeNames[len(nodeNames)-1]
	parentNode, err := g.Node.TraverseByPath(parentPath)

	if err != nil {
		return "", err
	} else {
		for _, x := range parentNode.Children {
			if x.Type == 1 && x.Name == newNodeName {
				return x.Data, nil
			}
		}

		return "", ErrorNoFileFoundAtGivenPath
	}
}

//List out all items directly under a folder
//
// if folderPath is an empty string, list out all item under current folder
func (g *treeUseCaseImpl) List(folderPath string) ([]*entity.Node, error) {
	if len(folderPath) == 0 {
		return g.Node.Children, nil
	}

	folderPath = lib.CleanPathToFileOrFolder(folderPath)
	destFolder, err := g.Node.TraverseByPath(folderPath)

	if err != nil {
		return []*entity.Node{}, err
	} else {
		return destFolder.Children, nil
	}
}

//FindByName return entity that has the given name
//
// If folderPath is not empty, look for items inside folder at folderPath only
func (g *treeUseCaseImpl) FindByName(name, folderPath string) ([]*entity.Node, error) {
	if len(folderPath) > 0 {
		folderPath = lib.CleanPathToFileOrFolder(folderPath)
		destItem, err := g.Navigate(folderPath)
		if err != nil {
			return []*entity.Node{}, err
		}

		return destItem.FindByName(name), nil
	} else {
		return g.CurrentNode.FindByName(name), nil
	}
}

//Update the node at path to have new name (optionally new data)
func (g *treeUseCaseImpl) Update(path, newName, data string) error {
	folderPath, targetItemName, err := filterPathAndLastItemName(path)
	if !strings.HasPrefix(folderPath, "/") {
		folderPath = "/" + folderPath
	}
	log.Printf("error %s, %s, %v\n", folderPath, targetItemName, err)
	if err != nil {
		return err
	}

	destFolder := g.Node
	// check if folder is root
	if len(folderPath) != 0 {
		destFolder, err = g.Navigate(folderPath)
		if err != nil {
			return err
		}
	}

	fmt.Printf("destFolder: %s\n", destFolder.Name)
	var targetNode entity.Node

	_, ok := destFolder.GetChildWithName(newName)
	if ok {
		return ErrorItemsCantHaveSameName
	}

	for _, x := range destFolder.Children {
		if x.Name == targetItemName {
			fmt.Printf("equal %s\n", x.Name)
			x.Name = newName
			if x.Type == 1 {
				x.Data = data
			}
			targetNode = *x
			break
		}
	}

	targetItem, err := targetNode.ToItem()
	if err != nil {
		return err
	}

	log.Printf("targetItem: %v\n", targetItem)

	err = g.ItemRepo.Update(targetItem)
	if err != nil {
		return ErrorFailToSaveItemToDB
	}
	return nil
}

//Move an item from source to dest,
//
//return an error if any path is invalid or destPath is a sub-path of sourcePath
func (g *treeUseCaseImpl) Move(sourcePath, destPath string) error {
	// find source item
	var sourceNode *entity.Node
	sourceNode, err := g.Navigate(sourcePath)
	if err != nil {
		parts := strings.Split(sourcePath, "/")
		sourceNode, err = g.Navigate(strings.Join(parts[:len(parts)-1], "/"))
		if err != nil {
			return err
		}
		for _, x := range sourceNode.Children {
			if x.Name == parts[len(parts)-1] {
				sourceNode = x
				break
			}
		}
	}

	// find dest folder
	destNode, err := g.Node.TraverseByPath(destPath)
	if err != nil {
		return err
	}

	// detect if dest is child folder of source
	if strings.HasPrefix(removeSlash(destPath), removeSlash(sourcePath)) {
		return ErrorDestinationIsSubPathOfSource
	}

	// remove source from source's parent
	sourceParent := sourceNode.Parent
	n := len(sourceParent.Children)
	for i, x := range sourceParent.Children {
		if x.Name == sourceNode.Name {
			sourceParent.Children[i], sourceParent.Children[n-1] = sourceParent.Children[n-1], sourceParent.Children[i]
			sourceParent.Children = sourceParent.Children[:n-1]
			break
		}
	}

	sourceNode.Parent = destNode
	// scan if dest already has an item with source's name
	for _, x := range destNode.Children {
		if x.Name == sourceNode.Name {
			return ErrorItemsCantHaveSameName
		}
	}

	// TODO: update new path
	sourceItem, err := sourceNode.ToItem()
	if err != nil {
		return err
	}

	err = g.ItemRepo.Update(sourceItem)
	if err != nil {
		return err
	}

	destNode.Children = append(destNode.Children, sourceNode)
	g.Node.UpdatePreOrderNumber()
	return nil
}

func (g *treeUseCaseImpl) Remove(paths []string) error {
	atLeastOneError := false

	for _, path := range paths {
		path = lib.CleanPathToFileOrFolder(path)
		if path == "" {
			atLeastOneError = true
			return ErrorInvalidPath
		}

		parentPath := lib.GetParentPath(path)
		nodeNames := strings.Split(path, "/")
		if !lib.CheckName(nodeNames[len(nodeNames)-1]) { // file name could be '.' or '..'
			atLeastOneError = true
			return ErrorFailToDeleteNodeInTree
		}
		parentNode, err := g.Node.TraverseByPath(parentPath)

		err = parentNode.DeleteByPath(path)
		if err != nil {
			atLeastOneError = true
			return ErrorFailToDeleteNodeInTree
		} else {
			parentPath := lib.GetParentPath(path)
			parentNode, err := g.Navigate(parentPath)
			if err != nil {
				atLeastOneError = true
				return ErrorInvalidPath
			}
			nodeNames := strings.Split(path, "/")
			if !lib.CheckName(nodeNames[len(nodeNames)-1]) { // file name could be '.' or '..'
				atLeastOneError = true
				return ErrorInvalidFileOrFolderName
			}
			deleteNodeName := nodeNames[len(nodeNames)-1]

			for _, x := range parentNode.Children {
				if x.Name == deleteNodeName {
					err = g.ItemRepo.Delete(x.Id)
					if err != nil {
						return ErrorFailToDeleteItemFromDB
					}
					break
				}
			}
		}

	}
	if atLeastOneError {
		return ErrorAtLeastOnePathIsFaulty
	}

	return nil
}

// filterPathAndLastItemName return path before an item and the item name
//
// In the case the file is under root, 1st string will return empty
func filterPathAndLastItemName(path string) (string, string, error) {
	if path == "/" || path == "~" {
		return "/", "", ErrorCantDividePrePathAndDestName
	}
	path = removeSlash(path)
	parts, err := entity.ParsePathToItemViewNameSlice(path)
	if err != nil {
		return "", "", err
	}
	folderPath := strings.Join(parts[:len(parts)-1], "/")
	itemName := parts[len(parts)-1]
	return folderPath, itemName, nil
}

func removeSlash(x string) string {
	if strings.HasSuffix(x, "/") {
		x = x[:len(x)-1]
	}
	return x
}
