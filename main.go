package main

import (
	"ibtools_server/cmd"
	_ "ibtools_server/docs"
	"log"
	"os"

	"github.com/urfave/cli"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

var (
	cliApp        *cli.App
	configBackend string
)

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "IBTools"
	cliApp.Usage = "IBTools 1.0 Server"
	cliApp.Author = "lizhe"
	cliApp.Email = "XXXX@qq.com"
	cliApp.Version = "1.0.0"
	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "configBackend",
			Value:       "etcd",
			Destination: &configBackend,
		},
	}
}

// @title IBtools Server API document
// @version 1.0
// @description 描述了server与client的交互接口
// @termsOfService https://www.XXX.com

// @contact.name lizhe
// @contact.url https://www.XXX.com
// @contact.email XXXXX@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes https
// @host apidoc.XX.com
// @BasePath /v1
func main() {
	cliApp.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "run migrations",
			Action: func(c *cli.Context) error {
				return cmd.Migrate(configBackend)
			},
		},
		{
			Name:  "loaddata",
			Usage: "load data from fixture",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "file",
					Usage:    "文件路径及名称",
					Required: true,
				},
				cli.StringFlag{
					Name:     "type",
					Usage:    "文件类型，可填card、relation",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				return cmd.LoadData(c.String("file"), c.String("type"), configBackend)
			},
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) error {
				return cmd.RunServer(configBackend)
			},
		},
	}

	// Run the CLI app
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
