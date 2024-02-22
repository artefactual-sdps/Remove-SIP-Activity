package remove

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	ignore "github.com/Diogenesoftoronto/go-gitignore"
)

const RemoveSIPFilesName = "remove-sip-files"

type RemoveSIPFilesActivity struct{}

func NewRemoveSIPFilesActivity() *RemoveSIPFilesActivity {
	return &RemoveSIPFilesActivity{}
}

type RemoveSIPFilesParams struct {
	DestPath   string // path to the directory where the files will be removed from.
	ConfigPath string // path to the configuration file
}

type RemoveSIPFilesResult struct{}

func (md *RemoveSIPFilesActivity) Execute(ctx context.Context, params *RemoveSIPFilesParams) (*RemoveSIPFilesResult, error) {
	obj := ignore.CompileIgnoreLines(params.ConfigPath)
	// Walk the sip path directory and remove the file if it does return a match on obj
	err := filepath.WalkDir(params.DestPath, func(path string, d fs.DirEntry, err error) error {
		if obj.MatchesPath(path) {
			err := os.RemoveAll(path)
			if err != nil {
				return err

			}
		}
		return nil

	})
	if err != nil {
		return &RemoveSIPFilesResult{}, err
	}

	return &RemoveSIPFilesResult{}, nil
}
