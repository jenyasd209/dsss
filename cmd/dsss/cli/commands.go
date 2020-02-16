package cli

import (
	"fmt"
	api2 "github.com/iorhachovyevhen/dsss/api"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
	"path/filepath"
)

var (
	cliPath      string
	historyFile  string
	uploadedPath string
)
var addrs = "http://localhost:8080"

var api = api2.API(addrs)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("can't get $HOME var")
	}

	cliPath = path.Join(home, ".dsss_cli")
	uploadedPath = path.Join(cliPath, "uploaded_files")

	if _, err := os.Stat(uploadedPath); err != nil {
		err = os.MkdirAll(uploadedPath, os.ModePerm)
		if err != nil {
			log.Fatalf("can't create dir by path: %v", uploadedPath)
		}
	}

	historyFile = path.Join(cliPath, "history.txt")

	if _, err := os.Stat(historyFile); err != nil {
		_, err = os.Create(historyFile)
		if err != nil {
			log.Fatalf("can't create file by path: %v", historyFile)
		}
	}
}

var RootApp = &cli.App{
	Name:        "dsss",
	Usage:       "keep your files remotely, and forget about them",
	HelpName:    "dsss",
	Version:     "0.1",
	Description: "I don't know how it will be work in future but now it's shit",
	Commands: cli.Commands{
		fileCmd,
	},
	Flags:                nil,
	CommandNotFound:      nil,
	Authors:              nil,
	Copyright:            "steal it, no one want use it",
	EnableBashCompletion: true,
}

var fileCmd = &cli.Command{
	Name:  "file",
	Usage: "work with files",
	Subcommands: cli.Commands{
		addCmd,
		getCmd,
		removeCmd,
		hystoryCmd,
	},
}

var addCmd = &cli.Command{
	Name:        "add",
	Aliases:     []string{"upload"},
	Usage:       "dsss file get <src_path>",
	Description: "download file by ID",
	Category:    "file",
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 {
			return errors.Errorf("arguments needed: 1, but got: %v", ctx.NArg())
		}

		filePath := ctx.Args().Get(0)

		filename, content, err := readFile(filePath)
		if err != nil {
			return err
		}

		id, err := api.Files().Add(filename, content)
		if err != nil {
			return err
		}

		fmt.Printf("File id: %v\n", id)

		return writeToHistoryFile(filename, id.String())
	},
}

var getCmd = &cli.Command{
	Name:        "get",
	Aliases:     []string{"download"},
	Usage:       "dsss file get <file_id> <dst_path>",
	Description: "upload a new file. If <dst_path> is empty file will be uploaded to $HOME/.dsss_cli/uploaded_files",
	Category:    "file",
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 && ctx.NArg() != 2 {
			return errors.Errorf("Needed 2 argument (second can be default '%s'), but got: %v", uploadedPath, ctx.NArg())
		}

		id := ctx.Args().Get(0)

		dst := ctx.Args().Get(1)
		if dst == "" {
			dst = uploadedPath
		}

		file, err := api.Files().Get([]byte(id))
		if err != nil {
			return err
		}

		p := path.Join(dst, file.Title())

		err = writeFile(p, file.Body())
		if err != nil {
			return err
		}

		abs, _ := filepath.Abs(p)

		fmt.Printf("File path: %s\n", abs)

		return nil
	},
}

var removeCmd = &cli.Command{
	Name:        "remove",
	Aliases:     []string{"delete"},
	Usage:       "dsss file remove <file_id>",
	Description: "remove file by ID",
	Category:    "file",
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 {
			return errors.Errorf("arguments needed: 1, but got: %v", ctx.NArg())
		}

		id := ctx.Args().Get(0)

		_, err := api.Files().Delete([]byte(id))
		if err != nil {
			return err
		}

		err = deleteFromHistoryFile(id)
		if err != nil {
			return err
		}

		fmt.Printf("ID of deleted file: %v\n", id)

		return nil
	},
}

var hystoryCmd = &cli.Command{
	Name:        "history",
	Aliases:     []string{"show"},
	Usage:       "dsss file list",
	Description: "return file uploading and deleting history",
	Category:    "file",
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 0 {
			return errors.Errorf("arguments needed: 0, but got: %v", ctx.NArg())
		}

		_, ids, err := readFile(historyFile)
		if err != nil {
			return err
		}

		list := string(ids)

		if list == "" || list == "\n" {
			return errors.New("No one uploaded file")
		}

		fmt.Println(list)

		return nil
	},
}
