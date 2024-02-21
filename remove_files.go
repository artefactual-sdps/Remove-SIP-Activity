package main

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
	SipPath     string
	SuccumbPath string
}

type RemoveSIPFilesResult struct{}

func (md *RemoveSIPFilesActivity) Execute(ctx context.Context, params *RemoveSIPFilesParams) (*RemoveSIPFilesResult, error) {
	obj := ignore.CompileIgnoreLines(params.SuccumbPath)
	// Walk the sip path directory and remove the file if it does return a match on obj
	alreadyDeletedPaths := make(map[string]any)

	err := filepath.WalkDir(params.SipPath, func(path string, d fs.DirEntry, err error) error {
		if obj.MatchesPath(path) {
			err := os.RemoveAll(path)
			if err != nil {
				return err

			}
			alreadyDeletedPaths[path] = nil
		}
		return nil

	})
	if err != nil {
		return &RemoveSIPFilesResult{}, err
	}

	return &RemoveSIPFilesResult{}, nil
}
