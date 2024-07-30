package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetSubGraphOutputFileName(partitionId int) string {
	return OUTPUT_FILES_DIR + "SubGraph-" + fmt.Sprint(partitionId) + JSON_FILE_EXTENSION
}

func GetSuperStepSubGraphOutputFileName(partitionId int, superStep int) string {
	return OUTPUT_FILES_DIR + "SuperStep-" + fmt.Sprint(superStep) + "-SubGraph-" + fmt.Sprint(partitionId) + JSON_FILE_EXTENSION
}

func GetCheckpointFileNamesForSuperstep(superStep int) ([]string, error) {
	return GetFilesWithPrefix(OUTPUT_FILES_DIR, fmt.Sprintf("SuperStep-%d-", superStep))
}

func GetFilesWithPrefix(dir, prefix string) ([]string, error) {
	var filesWithPrefix []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), prefix) {
			filesWithPrefix = append(filesWithPrefix, dir+info.Name())
		}
		return nil
	})
	return filesWithPrefix, err
}
