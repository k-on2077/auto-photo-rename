package main

import (
	"apr-go/src/imgtool"
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	CmdExit = "exit"
	CmdQuit = "quit"
	CmdHelp = "help"
	CmdName = "name"
)

func main() {
	printHelp()

	for {
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
	fmt.Println("1. input help to get command instruction;")
	fmt.Println("2. input exit or quit to stop;")
	fmt.Println("3. input name [folder path] to rename all photos in the specific folder;")
}

func processRename(dir string) {
	imgtool.RenameImgByTime(dir)
}
