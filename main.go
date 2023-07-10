package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	__NAME__           = "cn"
	__USAGE__          = "replace text of file name ,folder name and file text of a folder"
	__VERSION__        = "0.0.1"
	__DESCRIPTION__    = "replace text of file name ,folder name and file text of a folder"
	__DEFAULT_FOLDER__ = "."
	TEXT_CHARS         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]}\\|;:'\",<.>/? \t\n\r"
)

func main() {
	app := &cli.App{
		Name:        __NAME__,
		Version:     __VERSION__,
		Usage:       __USAGE__,
		Description: __DESCRIPTION__,
		Commands:    nil,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "open_folder",
				Aliases: []string{"o"},
				Usage:   "base folder",
			},
			&cli.BoolFlag{
				Name:    "replace_folder_name",
				Aliases: []string{"d"},
				Usage:   "whether replace folder name",
			},
			&cli.BoolFlag{
				Name:    "replace_file_name",
				Aliases: []string{"f"},
				Usage:   "whether replace file name",
			},
			&cli.BoolFlag{
				Name:    "replace_file_text",
				Aliases: []string{"t"},
				Usage:   "whether replace file text",
			},
			&cli.BoolFlag{
				Name:    "recursion",
				Aliases: []string{"r"},
				Usage:   "whether to operate recursively",
			},
		},
		Action: func(cCtx *cli.Context) error {
			var replacedStr = cCtx.Args().Get(0)
			if replacedStr == "" {
				return fmt.Errorf("the replaced text must be given")
			}
			var replacingStr = cCtx.Args().Get(1)
			if replacingStr == "" {
				return fmt.Errorf("the replacing text must be given")
			}
			var openFolder = __DEFAULT_FOLDER__
			if cCtx.String("open_folder") != "" {
				openFolder = cCtx.String("open_folder")
			}
			return cn(replacedStr, replacingStr,
				openFolder,
				cCtx.Bool("replace_folder_name"),
				cCtx.Bool("replace_file_name"),
				cCtx.Bool("replace_file_text"),
				cCtx.Bool("recursion"))
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cn(replacedStr string,
	replacingStr string,
	openFolder string,
	replaceFolderName bool,
	replaceFileName bool,
	replaceFileText bool,
	recursion bool) error {
	openFolderInfo, err := os.Stat(openFolder)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("path is not exist")
		} else {
			panic(err)
		}
	}
	if !openFolderInfo.IsDir() {
		panic(fmt.Errorf("this path is not dir"))
	}

	operatorDir, err := filepath.Abs(openFolder)
	if err != nil {
		panic(err)
	}
	fmt.Printf("operate based on %s \n", operatorDir)
	if replaceFolderName {
		changeFolderName(replacedStr, replacingStr, operatorDir, recursion)
	}

	if replaceFileName {
		changeFileName(replacedStr, replacingStr, operatorDir, recursion)
	}

	if replaceFileText {
		changeFileText(replacedStr, replacingStr, operatorDir, recursion)
	}

	return nil
}

func changeFolderName(replacedStr string,
	replacingStr string,
	operatorDir string,
	recursion bool) error {
	dirList := make([]string, 0)
	if recursion {
		err := filepath.Walk(operatorDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if path != operatorDir && strings.Contains(info.Name(), replacedStr) {
					dirList = append(dirList, path)
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	} else {
		dirEntries, err := os.ReadDir(operatorDir)
		if err != nil {
			panic(err)
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				if strings.Contains(entry.Name(), replacedStr) {
					path := filepath.Join(operatorDir, entry.Name())
					dirList = append(dirList, path)
				}
			}
		}
	}

	if len(dirList) == 0 {
		fmt.Println("no folder name changed")
		return nil
	}
	fmt.Printf("replace %s to %s for these folders \n", replacedStr, replacingStr)
	for i := len(dirList) - 1; i >= 0; i-- {
		treeDirPath := dirList[i]
		parentDir := filepath.Dir(treeDirPath)
		oldFolderBaseName := filepath.Base(treeDirPath)
		newFolderPath := filepath.Join(parentDir, strings.Replace(oldFolderBaseName, replacedStr, replacingStr, 1))
		fmt.Printf("%s -> %s \n", treeDirPath, newFolderPath)
		err := os.Rename(treeDirPath, newFolderPath)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func changeFileName(replacedStr string,
	replacingStr string,
	operatorDir string,
	recursion bool) error {
	fileList := make([]string, 0)
	if recursion {
		err := filepath.Walk(operatorDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if path != operatorDir && strings.Contains(info.Name(), replacedStr) {
					fileList = append(fileList, path)
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	} else {
		dirEntries, err := os.ReadDir(operatorDir)
		if err != nil {
			panic(err)
		}
		for _, entry := range dirEntries {
			if !entry.IsDir() {
				if strings.Contains(entry.Name(), replacedStr) {
					path := filepath.Join(operatorDir, entry.Name())
					fileList = append(fileList, path)
				}
			}
		}
	}

	if len(fileList) == 0 {
		fmt.Println("no file name changed")
		return nil
	}
	fmt.Printf("replace %s to %s for these file \n", replacedStr, replacingStr)
	for i := len(fileList) - 1; i >= 0; i-- {
		filePath := fileList[i]
		parentDir := filepath.Dir(filePath)
		oldFileBaseName := filepath.Base(filePath)
		newFilePath := filepath.Join(parentDir, strings.Replace(oldFileBaseName, replacedStr, replacingStr, 1))
		fmt.Printf("%s -> %s \n", filePath, newFilePath)
		err := os.Rename(filePath, newFilePath)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// todo use rg to change file text
func changeFileText(replacedStr string,
	replacingStr string,
	operatorDir string,
	recursion bool) error {
	println("begin replace text")
	if recursion {
		err := filepath.Walk(operatorDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				replaceText(replacedStr, replacingStr, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	} else {
		dirEntries, err := os.ReadDir(operatorDir)
		if err != nil {
			panic(err)
		}
		for _, entry := range dirEntries {
			if !entry.IsDir() {
				path := filepath.Join(operatorDir, entry.Name())
				replaceText(replacedStr, replacingStr, path)
			}
		}
	}

	println("end replace text")
	return nil
}

func replaceText(replacedStr string,
	replacingStr string,
	filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	scanner := bufio.NewScanner(file)
	tmpFile, err := os.Create(filePath + ".tmp")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	writer := bufio.NewWriter(tmpFile)
	for scanner.Scan() {
		line := scanner.Text()
		newLine := strings.Replace(line, replacedStr, replacingStr, -1)
		_, err := fmt.Fprintln(writer, newLine)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Error:", scanner.Err())
		return
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = tmpFile.Close()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = file.Close()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = os.Rename(tmpFile.Name(), filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
