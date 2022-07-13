package entity

import "errors"

var (
	ErrorNameNotValid                    = errors.New("item name is not valid")
	ErrorInvalidPath                     = errors.New("error invalid path")
	ErrorAnItemWithGivenNameAlreadyExist = errors.New("a file/folder with given name already exist")
	ErrorPathAndItemNameNotConsistent    = errors.New("path and item name is not consistent")
	ErrorNoFileOrFolderFound             = errors.New("no file or folder found")
	ErrorCantDeleteRoot                  = errors.New("root directory can't be deleted")
	ErrorCantDeleteCurrentNodeOrParent   = errors.New("can't delete current node or parent, can only delete child node")
	ErrorIdMustNotEmpty                  = errors.New("id must be non-empty")
	ErrorFileCantContainFileAndFolder    = errors.New("files cannot contain file or folder")
	ErrorFolderCantContainData           = errors.New("folders cannot contain data themselves")
	ErrorNoNodeWithGivenNameFound        = errors.New("no file or folder under this item")
)
