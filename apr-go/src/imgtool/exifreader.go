package imgtool

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/dsoprea/go-exif/v3"
)

func getFileExifTime(filePath string) (time.Time, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return time.Now(), err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return time.Now(), err
	}

	rawExif, err := exif.SearchAndExtractExif(data)
	if err != nil {
		return time.Now(), err
	}

	entries, _, err := exif.GetFlatExifDataUniversalSearch(rawExif, nil, true)
	if err != nil {
		return time.Now(), err
	}

	for _, entry := range entries {
		if entry.TagName == "DateTimeOriginal" {
			t, err := time.Parse("2006:01:02 15:04:05", entry.Value.(string))
			if err != nil {
				return time.Now(), err
			}
			return t, nil
		}
		//fmt.Printf("IFD-PATH=[%s] ID=(0x%04x) NAME=[%s] COUNT=(%d) TYPE=[%s] VALUE=[%s]\n", entry.IfdPath, entry.TagId, entry.TagName, entry.UnitCount, entry.TagTypeName, entry.Formatted)
	}

	return time.Now(), errors.New("no data")
}
