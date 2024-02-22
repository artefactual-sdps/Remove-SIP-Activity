package remove

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	ignore "github.com/Diogenesoftoronto/go-gitignore"
)

const RemoveFilesName = "remove-files-activity"

type RemoveFilesActivity struct{}

func NewRemoveFilesActivity() *RemoveFilesActivity {
	return &RemoveFilesActivity{}
}

type RemoveFilesParams struct {
	RemovePath string // The remove path is the path that files will be removed from.
	IgnorePath string // The ignore path is the path that contains the ignore/succumb file that configures which files will be removed.
}

type RemoveFilesResult struct {
	Removed []string
}

func (a *RemoveFilesActivity) Execute(ctx context.Context, params *RemoveFilesParams) (*RemoveFilesResult, error) {
	ig, err := ignore.CompileIgnoreFile(params.IgnorePath)
	if err != nil {
		return nil, err
	}

	// Walk the target directory and remove the file if it does return a match
	// one of the ignore patterns.
	var deleted []string
	var prevPath string
	err = filepath.WalkDir(params.RemovePath, func(path string, d fs.DirEntry, err error) error {
		if ig.MatchesPath(path) {
			err := os.RemoveAll(path)
			if err != nil {
				return err
			}

			// Deleted may have duplicates if the removed path is
			// a directory, for every item in the directory there
			// will be a duplicate. To avoid this you could check
			// the previous item in the list to see if it is the
			// same and only add it to the list if it is not.
			if path != prevPath {
				fmt.Println(path)
				deleted = append(deleted, path)
			}
			prevPath = path
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &RemoveFilesResult{Removed: deleted}, nil
}
