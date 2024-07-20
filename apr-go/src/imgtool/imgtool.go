package imgtool

import (
	"encoding/json"
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

	for _, file := range files {
		fileName := strings.ToLower(file.Name())
		exts := strings.Split(fileName, ".")
		if len(exts) < 2 {
			continue
		}

		fileExten := exts[len(exts)-1]
		if !strings.HasSuffix(fileName, fileExten) || !contains(extens, fileExten) {
			continue
		}

		var fileDir string
		if runtime.GOOS == "windows" {
			fileDir = dir + "\\"
		} else {
			fileDir = dir + "/"
		}

		oldName := fileDir + fileName
		originalTime, err := getFileExifTime(oldName)
		if err != nil {
			continue
		}

		newName := fileDir + "IMG_" + originalTime.Format("20060102_150405") + "." + fileExten
		if newName == oldName {
			newName = fileDir + "IMG_" + originalTime.Format("20060102_150405") + "1." + fileExten
		}

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
