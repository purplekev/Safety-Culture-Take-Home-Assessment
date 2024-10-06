package folder

import (
	"strings"

	"github.com/gofrs/uuid"
)

// Retrieves all child folders associated with given orgId and specific folder name, or error otherwise.
func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// error check: invalid orgID
	if orgID == uuid.Nil {
		return []Folder{}, ErrInvalidOrg
	}

	all_folders := f.folders
	folder, found := includesFolder(all_folders, name)

	if !found { // folder not found in data
		return []Folder{}, ErrInvalidFolder
	} else if folder.OrgId != orgID { // folder found but in different org
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
