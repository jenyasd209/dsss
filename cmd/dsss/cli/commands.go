package cli

import (
	"fmt"
	"log"
	"os"
	"path"

	api2 "github.com/iorhachovyevhen/dsss/api"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	cliPath   string
	idFile    string
	filesPath string
)
var addrs = "http://localhost:8080"

var api = api2.API(addrs)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("can't get $HOME var")
	}

	cliPath = path.Join(home, ".dsss_cli")
	filesPath = path.Join(cliPath, "uploaded_files")

	if _, err := os.Stat(filesPath); err != nil {
		err = os.Mkdir(filesPath, os.ModePerm)
		if err != nil {
			log.Fatalf("can't create dir by path: %v", filesPath)
		}
	}

	idFile = path.Join(cliPath, "id_list.txt")
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
	Flags:           nil,
	CommandNotFound: nil,
	Authors:         nil,
	Copyright:       "steal it, no one want use it",
	BashComplete: func(ctx *cli.Context) {
		fmt.Println("file")
	},
}

var fileTasks = []string{"add", "get", "remove", "list"}

var fileCmd = &cli.Command{
	Name:  "file",
	Usage: "work with files",
	Subcommands: cli.Commands{
		addCmd,
		getCmd,
		removeCmd,
		listCmd,
	},
	BashComplete: func(ctx *cli.Context) {
		if ctx.NArg() > 0 {
			return
		}
		for _, t := range fileTasks {
			fmt.Println(t)
		}
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

		return newFileId(string(id))
	},
}

var getCmd = &cli.Command{
	Name:        "get",
	Aliases:     []string{"download"},
	Usage:       "dsss file get <file_id> <dst_path>",
	Description: "upload a new file",
	Category:    "file",
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 || ctx.NArg() != 2 {
			return errors.Errorf("at least 1 argument needed, but got: %v", ctx.NArg())
		}

		id := ctx.Args().Get(0)

		dst := "/"
		if ctx.Args().Get(1) != "" {
			dst = filesPath
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

		fmt.Printf("File path: %v\n", dst)

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

		k, err := api.Files().Delete([]byte(id))
		if err != nil {
			return err
		}

		err = deleteFromFile(string(k))
		if err != nil {
			return err
		}

		fmt.Printf("ID of deleted file: %v\n", string(k))

		return nil
	},
}

var listCmd = &cli.Command{
	Name:        "list",
	Aliases:     []string{"show"},
	Usage:       "dsss file list",
	Description: "return ids of uploaded files",
	Category:    "file",
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 0 {
			return errors.Errorf("arguments needed: 0, but got: %v", ctx.NArg())
		}

		_, ids, err := readFile(idFile)
		if err != nil {
			return err
		}

		fmt.Println(string(ids))

		return nil
	},
}
