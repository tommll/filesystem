package fs

import (
	"fs/app/database"
	entity2 "fs/app/entity"
	"fs/app/repositories"
	impl "fs/module/fs/repositories"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var root entity2.Node
var repo repositories.ItemRepository
var data []entity2.Node

//+--~/ (id=0)
//|  +--tung/ (id=1)
//|  |  +--file1/ (id=2)
//|  |  +--file2/ (id=3)
//|  |  +--folder1/ (id=4)
//|  |  |  +--file3/ (id=5)
//|  +--usr/ (id=6)
//|  |  +--file4/ (id=7)

func initTestData() {
	database.Init()
	db := database.GetDB()
	repo = impl.NewItemRepo(db)
	data = make([]entity2.Node, 0)
	repo.DeleteAll()

	root = entity2.Node{
		Id:       "0",
		Type:     0,
		Name:     "~",
		Data:     "",
		Path:     "/",
		Children: []*entity2.Node{},
	}
	data = append(data, root)

	x1 := entity2.Node{
		Id:       "1",
		Type:     0,
		Parent:   &root,
		Name:     "tung",
		Data:     "",
		Path:     "/tung",
		Children: []*entity2.Node{},
		Root:     &root,
	}
	root.Children = append(root.Children, &x1)
	data = append(data, x1)

	x2a := entity2.Node{
		Id:       "2",
		Type:     1,
		Parent:   &x1,
		Name:     "file1",
		Path:     "/tung/file1",
		Data:     "hello world",
		Children: []*entity2.Node{},
		Root:     &root,
	}
	x1.Children = append(x1.Children, &x2a)
	data = append(data, x2a)

	x2b := entity2.Node{
		Id:       "3",
		Type:     1,
		Parent:   &x1,
		Name:     "file2",
		Path:     "/tung/file2",
		Data:     "as98d9asd",
		Children: []*entity2.Node{},
		Root:     &root,
	}
	x1.Children = append(x1.Children, &x2b)
	data = append(data, x2b)

	x2c := entity2.Node{
		Id:       "4",
		Type:     0,
		Parent:   &x1,
		Name:     "folder1",
		Path:     "/tung/folder1",
		Data:     "",
		Children: []*entity2.Node{},
		Root:     &root,
	}
	x1.Children = append(x1.Children, &x2c)
	data = append(data, x2c)

	x2d := entity2.Node{
		Id:       "5",
		Type:     1,
		Parent:   &x2c,
		Name:     "file3",
		Path:     "/tung/folder1/file3",
		Data:     "as98d9asd",
		Children: []*entity2.Node{},
		Root:     &root,
	}
	x2c.Children = append(x2c.Children, &x2d)
	data = append(data, x2d)

	x3 := entity2.Node{
		Id:       "6",
		Type:     0,
		Parent:   &root,
		Name:     "usr",
		Path:     "/usr",
		Data:     "",
		Children: []*entity2.Node{},
		Root:     &root,
	}
	root.Children = append(root.Children, &x3)
	data = append(data, x3)

	x4 := entity2.Node{
		Id:       "7",
		Type:     1,
		Parent:   &x3,
		Name:     "file4",
		Path:     "/usr/file4",
		Data:     "hello tung",
		Children: []*entity2.Node{},
		Root:     &root,
	}
	x3.Children = append(x3.Children, &x4)
	data = append(data, x4)

	for _, x := range data[1:] {
		item, _ := x.ToItem()
		repo.InsertNewItem(item)
	}

}

func TestBasicExport(t *testing.T) {
	initTestData()

}

func TestNavigate(t *testing.T) {
	initTestData()

	uc := NewTreeUseCase(&root, repo)

	var gotFolder *entity2.Node
	var gotErr error

	gotFolder, gotErr = uc.Navigate("/tung/../usr")
	assert.NoError(t, gotErr)
	assert.Equal(t, "usr", gotFolder.Name)

	gotFolder, gotErr = uc.Navigate("./tung")
	assert.NoError(t, gotErr)
	assert.Equal(t, "tung", gotFolder.Name)

	gotFolder, gotErr = uc.Navigate("../")
	assert.NoError(t, gotErr)
	assert.Equal(t, "~", gotFolder.Name)

	gotFolder, gotErr = uc.Navigate(".")
	assert.NoError(t, gotErr)
	assert.Equal(t, "~", gotFolder.Name)

}

func TestCreate(t *testing.T) {
	initTestData()

	uc := NewTreeUseCase(&root, repo)
	var gotErr error

	t.Run("normal case", func(t *testing.T) {
		root.Print(0, "")
		//newNode1 := entity2.Node{Id: "10", Type: 0, Name: "new folder here", Parent: &data[1]}
		gotErr = uc.Create("/tung/new folder here", "")
		assert.NoError(t, gotErr)

		gotNodes, gotErr := uc.Navigate("/tung/new folder here")
		assert.NoError(t, gotErr)
		assert.Equal(t, "new folder here", gotNodes.Name)
		root.Print(0, "")
	})

	t.Run("normal case 2", func(t *testing.T) {
		root.Print(0, "")
		//newNode1 := entity2.Node{Id: "10", Type: 0, Name: "new folder here", Parent: &data[1]}
		gotErr = uc.Create("/tung/new2/", "")
		assert.NoError(t, gotErr)

		gotNodes, gotErr := uc.Navigate("/tung/new2")
		assert.NoError(t, gotErr)
		assert.Equal(t, "new2", gotNodes.Name)
		root.Print(0, "")
	})

	t.Run("wrong case: insert under file", func(t *testing.T) {
		root.Print(0, "")
		//newNode2 := entity2.Node{Id: "11", Type: 1, Name: "new file here", Parent: &data[1]}
		gotErr = uc.Create("/tung/file1.txt/wrong_file.bin", "asd1223")
		assert.Equal(t, ErrorInvalidPath, gotErr)
	})

	t.Run("wrong case: invalid file name", func(t *testing.T) {
		root.Print(0, "")
		gotErr = uc.Create("/tung/.", "asd1223")
		assert.Equal(t, ErrorInvalidPath, gotErr)

		root.Print(0, "")
		gotErr = uc.Create("/tung/..", "asd1223")
		assert.Equal(t, ErrorInvalidFileOrFolderName, gotErr)
	})
}

func TestShowData(t *testing.T) {
	initTestData()

	uc := NewTreeUseCase(&root, repo)

	gotData, gotErr := uc.ShowData("/tung/file1/")
	assert.NoError(t, gotErr)
	assert.Equal(t, "hello world", gotData)

	gotData, gotErr = uc.ShowData("/tung/file1/")
	assert.NoError(t, gotErr)
	assert.Equal(t, "hello world", gotData)

	// wrong path case
	gotData, gotErr = uc.ShowData("/tu g/file1")
	assert.Equal(t, ErrorInvalidPath, gotErr)

	// no filename case
	gotData, gotErr = uc.ShowData("/tung/file_wrong")
	assert.Equal(t, ErrorNoFileFoundAtGivenPath, gotErr)

	// invalid path case
	gotData, gotErr = uc.ShowData("/tung/.")
	assert.Equal(t, ErrorInvalidPath, gotErr)

	// invalid file name case
	gotData, gotErr = uc.ShowData("/tung/..")
	assert.Equal(t, ErrorInvalidFileOrFolderName, gotErr)
}

func TestList(t *testing.T) {
	initTestData()
	root.Print(0, "")

	database.Init()
	db := database.GetDB()
	repo := impl.NewItemRepo(db)

	uc := NewTreeUseCase(&root, repo)

	gotItems, gotErr := uc.List("/tung")
	assert.NoError(t, gotErr)
	assert.Equal(t, 3, len(gotItems))

	gotItems, gotErr = uc.List("/tung/")
	assert.NoError(t, gotErr)
	assert.Equal(t, 3, len(gotItems))

	gotItems, gotErr = uc.List("/tung/folder1")
	assert.NoError(t, gotErr)
	assert.Equal(t, 1, len(gotItems))

	// wrong folder path
	gotItems, gotErr = uc.List("/tung/folder_wrong")
	assert.Equal(t, ErrorInvalidPath, gotErr)

	gotItems, gotErr = uc.List("/tung/.")
	assert.Equal(t, ErrorInvalidPath, gotErr)

	gotItems, gotErr = uc.List("/./")
	assert.Equal(t, ErrorInvalidPath, gotErr)

	// wrong case: list file
	gotItems, gotErr = uc.List("/tung/file1")
	assert.Equal(t, ErrorInvalidPath, gotErr)
}

func TestFind(t *testing.T) {
	initTestData()
	root.Print(0, "")

	uc := NewTreeUseCase(&root, repo)

	gotItem, gotErr := uc.FindByName("file3", "/tung/folder1")
	assert.NoError(t, gotErr)
	assert.Equal(t, 1, len(gotItem))

	// fail case
	gotItem, gotErr = uc.FindByName("file_wrong", "/")
	assert.Equal(t, 0, len(gotItem))

	// find substring case
	gotItem, gotErr = uc.FindByName("file", "/tung")
	assert.Equal(t, 3, len(gotItem))

	// find from current folder
	gotItem, gotErr = uc.FindByName("file4", ".")
	assert.Equal(t, 1, len(gotItem))

	// find from parent folder
	_, gotErr = uc.Navigate("/usr")
	gotItem, gotErr = uc.FindByName("file4", ".")
	assert.NoError(t, gotErr)
	assert.Equal(t, 1, len(gotItem))

	gotItem, gotErr = uc.FindByName("file4", "..")
	assert.NoError(t, gotErr)
	assert.Equal(t, 1, len(gotItem))

	// find from a separated folder
	gotItem, gotErr = uc.FindByName("file4", "../tung")
	assert.NoError(t, gotErr)
	assert.Equal(t, 0, len(gotItem))
}

func TestUpdate(t *testing.T) {
	initTestData()
	root.Print(0, "")

	uc := NewTreeUseCase(&root, repo)
	var gotErr error
	var gotItem entity2.Item
	log.Println(gotItem)

	t.Run("a", func(t *testing.T) {
		////// file case
		newData := ""
		gotErr = uc.Update("/tung/folder1", "folder1_new", newData)
		assert.NoError(t, gotErr)

		root.Print(0, "")
		gotItem, gotErr = repo.GetById("4")
		assert.NoError(t, gotErr)
		assert.Equal(t, "folder1_new", gotItem.Name)
	})
	var path string
	var newName string

	//// wrong case: update duplicate name
	path = "/tung/file2"
	newName = "file1"
	gotErr = uc.Update(path, newName, "data of fake file1")
	assert.Equal(t, ErrorItemsCantHaveSameName, gotErr)
	root.Print(0, "")

	//// wrong case: invalid name
	path = "/tung/file1"
	newName = "file% +2"
	gotErr = uc.Update(path, newName, "data of invalid file")
	assert.NotEqual(t, nil, gotErr)
	root.Print(0, "")

	//t.Run("edge case: using alias", func(t *testing.T) {
	//	// 1st
	//	gotErr = uc.Update("./tung", "tung1", "")
	//	assert.NoError(t, gotErr)
	//	root.Print(0, "")
	//
	//	gotItem, gotErr = repo.GetById("1")
	//	assert.NoError(t, gotErr)
	//	assert.Equal(t, "tung1", gotItem.Name)
	//
	//	// 2nd
	//	gotErr = uc.Update("../tung1", "tung2", "")
	//	assert.NoError(t, gotErr)
	//	root.Print(0, "")
	//
	//	gotItem, gotErr = repo.GetById("1")
	//	assert.NoError(t, gotErr)
	//	assert.Equal(t, "tung2", gotItem.Name)
	//
	//	// 3rd
	//	uc.Navigate("/tung2/folder1")
	//	gotErr = uc.Update("../../usr/file4", "file4_new", "mock data")
	//	assert.NoError(t, gotErr)
	//	root.Print(0, "")
	//
	//	gotItem, gotErr = repo.GetById("7")
	//	assert.NoError(t, gotErr)
	//	assert.Equal(t, "file4_new", gotItem.Name)
	//})
}

func TestMove(t *testing.T) {
	initTestData()

	uc := NewTreeUseCase(&root, repo)
	var gotErr error

	t.Run("move folder", func(t *testing.T) {
		root.Print(0, "")
		gotErr = uc.Move("/usr", "/tung")
		assert.NoError(t, gotErr)

		gotNode, gotErr := uc.Navigate("/tung/usr")
		assert.NoError(t, gotErr)
		assert.Equal(t, "usr", gotNode.Name)
		root.Print(0, "")
	})

	t.Run("move folder to root", func(t *testing.T) {
		root.Print(0, "")
		gotErr = uc.Move("/tung/usr", "/")
		assert.NoError(t, gotErr)

		gotNode, gotErr := uc.Navigate("/usr")
		assert.NoError(t, gotErr)
		assert.Equal(t, "usr", gotNode.Name)
		root.Print(0, "")
	})

	t.Run("move file to folder", func(t *testing.T) {
		root.Print(0, "")
		gotErr = uc.Move("/usr/file4", "/tung")
		assert.NoError(t, gotErr)

		root.Print(0, "")
	})

	t.Run("move folder to its child's folder", func(t *testing.T) {
		gotErr = uc.Move("/tung", "/tung/folder1")
		assert.Equal(t, ErrorDestinationIsSubPathOfSource, gotErr)

		root.Print(0, "")
	})

	t.Run("move folder/file under root", func(t *testing.T) {
		gotErr = uc.Move("/tung/file1", "/")
		assert.Equal(t, nil, gotErr)

		nodes, gotErr := uc.List("/")
		hasNewFile := false
		for _, x := range nodes {
			if x.Name == "file1" {
				hasNewFile = true
			}
		}
		assert.Equal(t, true, hasNewFile)
		assert.NoError(t, gotErr)
		root.Print(0, "")
	})

}

func TestRemove(t *testing.T) {
	initTestData()

	uc := NewTreeUseCase(&root, repo)
	var gotErr error
	var paths []string
	var gotNodes []*entity2.Node

	root.Print(0, "")
	paths = []string{"/tung/file1"}
	gotErr = uc.Remove(paths)
	assert.NoError(t, gotErr)
	gotNodes, gotErr = uc.FindByName("file1", "/tung")
	assert.Equal(t, 0, len(gotNodes))
	root.Print(0, "")

	// folder
	paths = []string{"/usr"}
	gotErr = uc.Remove(paths)
	assert.NoError(t, gotErr)
	gotNodes, gotErr = uc.FindByName("usr", "/")
	assert.Equal(t, 0, len(gotNodes))
	root.Print(0, "")

	//// fault
	paths = []string{"/ur"}
	gotErr = uc.Remove(paths)
	assert.Equal(t, ErrorFailToDeleteNodeInTree, gotErr)

	root.Print(0, "")

	// delete root case
	paths = []string{"/"}
	gotErr = uc.Remove(paths)
	assert.Equal(t, ErrorFailToDeleteNodeInTree, gotErr)

	//root.Print(0, "")

	// delete current node or higher
	paths = []string{"./"}
	gotErr = uc.Remove(paths)
	assert.Equal(t, ErrorFailToDeleteNodeInTree, gotErr)

	paths = []string{"."}
	gotErr = uc.Remove(paths)
	assert.Equal(t, ErrorFailToDeleteNodeInTree, gotErr)

	paths = []string{"../"}
	gotErr = uc.Remove(paths)
	assert.Equal(t, ErrorFailToDeleteNodeInTree, gotErr)

	paths = []string{".."}
	gotErr = uc.Remove(paths)
	assert.Equal(t, ErrorFailToDeleteNodeInTree, gotErr)

	paths = []string{"/usr"}
	uc.Navigate("/tung/folder1") // at folder1
	gotErr = uc.Remove(paths)
	assert.Equal(t, ErrorFailToDeleteNodeInTree, gotErr)
}

func TestTreeUseCaseImpl_Import(t *testing.T) {
	newRoot := entity2.Node{
		Id:       "0",
		Type:     0,
		Order:    0,
		Parent:   nil,
		Name:     "~",
		Path:     "/",
		Data:     "",
		Children: []*entity2.Node{},
		Root:     nil,
	}

	database.Init()
	db := database.GetDB()
	repo = impl.NewItemRepo(db)
	uc := NewTreeUseCase(&newRoot, repo)

	gotErr := uc.Import()
	assert.NoError(t, gotErr)

	newRoot.Print(0, "")
}

func TestTreeUseCaseImpl_Export(t *testing.T) {
	newRoot := entity2.Node{}

	initTestData()
	uc := NewTreeUseCase(&newRoot, repo)

	gotErr := uc.Import()
	assert.NoError(t, gotErr)

	log.Println("BEFORE MANIPULATION")
	newRoot.PrintOrder(0, "")

	gotErr = uc.Move("/tung/file1", "/usr")
	assert.NoError(t, gotErr)

	log.Println("AFTER MANIPULATION, now order should change")
	newRoot.PrintOrder(0, "")

	gotErr = uc.Export()
	assert.NoError(t, gotErr)
}
