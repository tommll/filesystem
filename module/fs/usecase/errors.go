package fs

import "errors"

var (
	ErrorNoFileFoundAtGivenPath       = errors.New("no file found at the given path")
	ErrorDestinationIsSubPathOfSource = errors.New("destination path is sub-path of source path")
	ErrorItemsCantHaveSameName        = errors.New("can't have 2 items with the same name ")
	ErrorAtLeastOnePathIsFaulty       = errors.New("at least one path is faulty while deleting")
	ErrorFailToGetItems               = errors.New("fail to get items")
	ErrorParentChildOrdering          = errors.New("fail to attach child node to its parent node")
	ErrorConvertNodeToItem            = errors.New("fail to convert node to item")
	ErrorFailToSaveItemToDB           = errors.New("fail to save item to db")
	ErrorFailToDeleteItemFromDB       = errors.New("fail to delete item from db")
	ErrorFailToDeleteNodeInTree       = errors.New("fail to delete item in tree")
	ErrorFailToInsertNodeToTree       = errors.New("fail to insert node to tree")
	ErrorInvalidPath                  = errors.New("error invalid path")
	ErrorCantDividePrePathAndDestName = errors.New("can't divide a pre-path and a destination name")
	ErrorInvalidFileOrFolderName      = errors.New("invalid file name")
)
