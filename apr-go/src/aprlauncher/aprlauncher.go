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
	fmt.Printf("Usage:\n")
	fmt.Printf("- input [help] to get command instruction;\n")
	fmt.Printf("- input [exit] or [quit] to stop;\n")
	fmt.Printf("- input [name] [folder path] to rename all photos in the specific folder;\n\n")
	fmt.Printf("please input command:\n")
}

func processRename(dir string) {
	err := imgtool.RenameImgByTime(dir)
	if err != nil {
		fmt.Printf("rename error: %v\n", err)
	} else {
		fmt.Printf("rename success\n")
	}

	fmt.Printf("\n")
	fmt.Println("please input command:")
}
