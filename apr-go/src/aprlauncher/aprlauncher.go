package main

import (
	"apr-go/src/common"
	"apr-go/src/imgtool"
	"apr-go/src/videotool"
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	CmdExit = "exit"
	CmdQuit = "quit"
	CmdHelp = "help"
	CmdName = "rename"
)

var imgExtens = []string{"jpg", "jpeg", "png", "heic"}
var videoExtens = []string{"mov", "mp4", "mkv"}

func main() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "-rename" {
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Printf("get current path error: %v\n", err)
				fmt.Printf("please input a correct path\n")
			} else {
				fmt.Printf("rename files in current path: %v\n", currentDir)
				processRename(currentDir)
				fmt.Printf("bye...\n")
				os.Exit(0)
			}
		}
	}

	printHelp()

	for {
		fmt.Println("please input command:")
		reader := bufio.NewReader(os.Stdin)
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("read command error: %v\n", err)
			continue
		}

		input := string(line)
		fields := strings.Fields(input)
		if len(fields) == 0 {
			fmt.Printf("please input\n")
			continue
		}

		cmd := strings.ToLower(fields[0])

		if cmd == CmdExit || cmd == CmdQuit {
			fmt.Printf("bye...\n")
			break
		} else if cmd == CmdHelp {
			printHelp()
		} else if cmd == CmdName {
			if len(fields) == 2 {
				processRename(fields[1])
			}
		}
	}

	os.Exit(0)
}

func printHelp() {
	fmt.Printf("Usage:\n")
	fmt.Printf("- execute with option [-rename] to rename all files in the current folder;\n")
	fmt.Printf("- input [%v] to get command instruction;\n", CmdHelp)
	fmt.Printf("- input [%v] or [%v] to stop;\n", CmdExit, CmdQuit)
	fmt.Printf("- input [%v] [folder path] to rename all photos in the specific folder;\n\n", CmdName)
}

func processRename(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("read dir error: %v\n", err)
	} else {
		fileDir := ""
		if runtime.GOOS == "windows" {
			fileDir = dir + "\\"
		} else {
			fileDir = dir + "/"
		}

		records := make([]*common.RenameRecord, 0)

		existNames := make(map[string]int, 0)
		for _, file := range files {
			existNames[fileDir+file.Name()] = 1
			existNames[fileDir+strings.ToLower(file.Name())] = 1
		}

		for _, file := range files {
			lower := strings.ToLower(file.Name())
			exts := strings.Split(lower, ".")
			if len(exts) < 2 {
				continue
			}

			fileExten := exts[len(exts)-1]
			if !strings.HasSuffix(lower, fileExten) {
				continue
			}

			oldName := fileDir + file.Name()
			newName := oldName
			filePrefix := ""
			var fileTime time.Time
			var err error

			if common.Contains(imgExtens, fileExten) {
				filePrefix = "IMG_"
				fileTime, err = imgtool.GetImgCreationTime(oldName)
			} else if common.Contains(videoExtens, fileExten) {
				filePrefix = "IMG_"
				fileTime, err = videotool.GetVideoCreationTime(oldName)
			} else {
				err = fmt.Errorf("unsupported file type")
			}

			if err == nil {
				newName = fileDir + filePrefix + fileTime.Format("20060102_150405") + "." + fileExten
			}

			if newName == oldName {
				continue
			}

			// check if the new name already exists
			if _, ok := existNames[newName]; ok {
				for i := 1; i < 10000; i++ {
					newName = fileDir + filePrefix + fileTime.Format("20060102_150405") + fmt.Sprintf("_%v", i) + "." + fileExten
					if _, ok := existNames[newName]; !ok {
						break
					}
				}
			}
			existNames[newName] = 1

			err = os.Rename(oldName, newName)
			if err != nil {
				record := &common.RenameRecord{
					OldName: oldName,
					NewName: newName,
					Error:   err.Error(),
				}
				records = append(records, record)
			} else {
				record := &common.RenameRecord{
					OldName: oldName,
					NewName: newName,
					Error:   "",
				}
				records = append(records, record)
			}
		}

		common.LogRenameRecord(records, dir)
		fmt.Printf("rename done\n")
	}
}
