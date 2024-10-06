package folder

import (
	"errors"

	"github.com/gofrs/uuid"
)

// Pre defined variables for error handling
var (
	// get_folder errors
	ErrInvalidOrg       = errors.New("error: OrgId does not exist or is invalid")
	ErrInvalidFolder    = errors.New("error: Folder does not exist")
	ErrFolderInvalidOrg = errors.New("error: Folder does not exist in the given org")
	// move_folder errors
	ErrInvalidSrcFolder = errors.New("error: Source folder does not exist")
	ErrInvalidDstFolder = errors.New("error: Destination folder does not exist")
	ErrMoveToSame       = errors.New("error: Cannot move a folder to itself")
	ErrMoveToChildOf    = errors.New("error: Cannot move a folder to a child of itself")
	ErrMoveToDiffOrg    = errors.New("error: Cannot move a folder to a different organization")
)

// Fetches data from sample.json
func GetAllFolders() []Folder {
	return GetSampleData()
}

// Filters and returns all folders belonging to a specific orgId
func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	folders := f.folders

	res := []Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, f)
		}
	}

	return res

}

// Filters finds a folder of type Folder given a folder name. Returns nil and false if not found.
func includesFolder(folders []Folder, folderName string) (*Folder, bool) {
	for _, folder := range folders {
		if folder.Name == folderName {
			return &folder, true
		}
	}
	return nil, false
}
