package videotool

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type TagsData struct {
	CreationTime string `json:"creation_time"`
}

type FormatData struct {
	Format struct {
		Tags TagsData `json:"tags"`
	} `json:"format"`
}

func GetVideoCreationTime(filePath string) (time.Time, error) {
	cmdStr := fmt.Sprintf("ffprobe -v error -show_entries format_tags=creation_time -print_format json -i %v", filePath)
	cmdArgs := strings.Fields(cmdStr)
	output, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).CombinedOutput()
	if err != nil {
		fmt.Printf("get video creation time error: %v\n", err)
	}

	var data FormatData
	jsonError := json.Unmarshal(output, &data)
	if jsonError != nil {
		return time.Now(), jsonError
	}

	creationTime, timeError := time.Parse(time.RFC3339, data.Format.Tags.CreationTime)
	if timeError != nil {
		return time.Now(), timeError
	}
	eightZone, zoneError := time.LoadLocation("Asia/Shanghai")
	if zoneError != nil {
		return time.Now(), zoneError
	}
	return creationTime.UTC().In(eightZone), nil
}
