package folder

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
)

// pre defined variables for error handling
var (
	ErrInvalidOrg       = errors.New("error: OrgId does not exist or is invalid")
	ErrInvalidFolder    = errors.New("error: Folder does not exist")
	ErrFolderInvalidOrg = errors.New("error: Folder does not exist in the given org")
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

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

// helper function to check if a specified name exists in a list of folders.
// returns pointer to the folder and true if found, or nil and false if not found.
func includesFolder(folders []Folder, folderName string) (*Folder, bool) {
	for _, folder := range folders {
		if folder.Name == folderName {
			return &folder, true
		}
	}
	return nil, false
}

// retrieves all child folders associated with given orgId and specific folder name.
// returns a slice of folders containing matching children, or error if it doesn't exist.
func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// error check: invalid orgID
	if orgID == uuid.Nil {
		return []Folder{}, ErrInvalidOrg
	}

	all_folders := f.folders
	folder, found := includesFolder(all_folders, name)

	// routes to appropriate error checks
	if !found {
		return []Folder{}, ErrInvalidFolder
	} else if folder.OrgId != orgID {
		return []Folder{}, ErrFolderInvalidOrg
	}

	// retrieve all folders with corresponding orgID
	folders_matching_orgID := f.GetFoldersByOrgID(orgID)
	res := []Folder{}

	for _, folder := range folders_matching_orgID {
		if strings.Contains(folder.Paths, (name + ".")) {
			res = append(res, folder)
		}
	}

	return res, nil
}
