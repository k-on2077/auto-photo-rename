package common

import (
	"encoding/json"
	"os"
	"runtime"
	"time"
)

func Contains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func LogRenameRecord(records []*RenameRecord, dir string) {
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
