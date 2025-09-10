package imgtool

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type ExifMetadata struct {
	DateTime            string `json:"Exif.Image.DateTime"`
	DateTimeOriginal    string `json:"Exif.Photo.DateTimeOriginal"`
	DateTimeDigitized   string `json:"Exif.Photo.DateTimeDigitized"`
	OffsetTime          string `json:"Exif.Photo.OffsetTime"`
	OffsetTimeOriginal  string `json:"Exif.Photo.OffsetTimeOriginal"`
	OffsetTimeDigitized string `json:"Exif.Photo.OffsetTimeDigitized"`
}

func GetCreationTimeByExiv2(filePath string) (time.Time, error) {
	cmdStr := fmt.Sprintf("exiv2 -pe %v", filePath)
	cmdArgs := strings.Fields(cmdStr)
	output, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).CombinedOutput()
	if err != nil {
		fmt.Printf("call exiv2 error: %v\n", err)
	}

	var metadata ExifMetadata
	re := regexp.MustCompile(`\s+`)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "Exif.Image.DateTime") {
			parts := re.Split(line, 4)
			if len(parts) == 4 {
				for i := range parts {
					parts[i] = strings.TrimSpace(parts[i])
				}
				metadata.DateTime = parts[3]
			}
		} else if strings.Contains(line, "Exif.Photo.DateTimeOriginal") {
			parts := re.Split(line, 4)
			if len(parts) == 4 {
				for i := range parts {
					parts[i] = strings.TrimSpace(parts[i])
				}
				metadata.DateTimeOriginal = parts[3]
			}
		} else if strings.Contains(line, "Exif.Photo.DateTimeDigitized") {
			parts := re.Split(line, 4)
			if len(parts) == 4 {
				for i := range parts {
					parts[i] = strings.TrimSpace(parts[i])
				}
				metadata.DateTimeDigitized = parts[3]
			}
		} else if strings.Contains(line, "Exif.Photo.OffsetTime") {
			parts := re.Split(line, 4)
			if len(parts) == 4 {
				for i := range parts {
					parts[i] = strings.TrimSpace(parts[i])
				}
				metadata.OffsetTime = parts[3]
			}
		} else if strings.Contains(line, "Exif.Photo.OffsetTimeOriginal") {
			parts := re.Split(line, 4)
			if len(parts) == 4 {
				for i := range parts {
					parts[i] = strings.TrimSpace(parts[i])
				}
				metadata.OffsetTimeOriginal = parts[3]
			}
		} else if strings.Contains(line, "Exif.Photo.OffsetTimeDigitized") {
			parts := re.Split(line, 4)
			if len(parts) == 4 {
				for i := range parts {
					parts[i] = strings.TrimSpace(parts[i])
				}
				metadata.OffsetTimeDigitized = parts[3]
			}
		}
	}

	if metadata.DateTime != "" && metadata.OffsetTime == "+08:00" {
		t, err := time.Parse("2006:01:02 15:04:05", metadata.DateTime)
		if err == nil {
			return t, nil
		}
	} else if metadata.DateTimeOriginal != "" && metadata.OffsetTimeOriginal == "+08:00" {
		t, err := time.Parse("2006:01:02 15:04:05", metadata.DateTimeOriginal)
		if err == nil {
			return t, nil
		}
	} else if metadata.DateTimeDigitized != "" && metadata.OffsetTimeDigitized == "+08:00" {
		t, err := time.Parse("2006:01:02 15:04:05", metadata.DateTimeDigitized)
		if err == nil {
			return t, nil
		}
	}

	return time.Now(), fmt.Errorf("failed to get creation time from exif data")
}
