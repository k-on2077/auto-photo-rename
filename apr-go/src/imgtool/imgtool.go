package imgtool

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type RenameRecord struct {
	OldName string `json:"OldName"`
	NewName string `json:"NewName"`
	Error   string `json:"Error"`
}

func RenameImgByTime(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	extens := []string{"jpg", "jpeg", "png", "heic"}
	records := make([]*RenameRecord, 0)

	var fileDir string
	if runtime.GOOS == "windows" {
		fileDir = dir + "\\"
	} else {
		fileDir = dir + "/"
	}

	existNames := make(map[string]int, 0)
	for _, file := range files {
		existNames[fileDir+file.Name()] = 1
		existNames[fileDir+strings.ToLower(file.Name())] = 1
	}

	for _, file := range files {
		lowerFileName := strings.ToLower(file.Name())
		exts := strings.Split(lowerFileName, ".")
		if len(exts) < 2 {
			continue
		}

		fileExten := exts[len(exts)-1]
		if !strings.HasSuffix(lowerFileName, fileExten) || !contains(extens, fileExten) {
			continue
		}

		oldName := fileDir + file.Name()
		originalTime, err := getFileExifTime(oldName)
		if err != nil {
			continue
		}

		newName := fileDir + "IMG_" + originalTime.Format("20060102_150405") + "." + fileExten
		if newName == oldName {
			continue
		}

		// check if the new name already exists
		if _, ok := existNames[newName]; ok {
			for i := 1; i < 10000; i++ {
				newName = fileDir + "IMG_" + originalTime.Format("20060102_150405") + fmt.Sprintf("_%v", i) + "." + fileExten
				if _, ok := existNames[newName]; !ok {
					break
				}
			}
		}
		existNames[newName] = 1

		err = os.Rename(oldName, newName)
		if err != nil {
			record := &RenameRecord{
				OldName: oldName,
				NewName: newName,
				Error:   err.Error(),
			}
			records = append(records, record)
		} else {
			record := &RenameRecord{
				OldName: oldName,
				NewName: newName,
				Error:   "",
			}
			records = append(records, record)
		}
	}

	logRenameRecord(records, dir)
	return nil
}

func logRenameRecord(records []*RenameRecord, dir string) {
	if len(records) == 0 {
		return
	}

	var fileDir string
	if runtime.GOOS == "windows" {
		fileDir = dir + "\\"
	} else {
		fileDir = dir + "/"
	}

	jsonBytes, err := json.MarshalIndent(records, "", "  ")
	if err == nil {
		logFileName := fileDir + "log_" + time.Now().Format("20060102_150405") + ".json"
		if _, err := os.Stat(logFileName); err == nil {
			os.Remove(logFileName)
		}
		if file, err := os.Create(logFileName); err == nil {
			defer file.Close()
			file.Write(jsonBytes)
		}
	}
}

func contains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}
