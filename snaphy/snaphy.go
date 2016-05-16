package main

import (
	"os"
	"github.com/codegangsta/cli"
	//_ "snaphyAuth/models"
	"fmt"
	"github.com/astaxie/beego"
	"path/filepath"
	"runtime"
)

func init()  {
	//Current path of calling file..
	_, file, _, _ := runtime.Caller(1)
	//adding app path..first..
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "../")))
	fmt.Println(appPath)
	beego.AppPath = appPath
}



func main() {
	app := cli.NewApp()
	app.Name = "snaphy"
	app.Usage = "snaphy admin cli control!"
	/*var user string
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "lang",
			Value:       "english",
			Usage:       "username for creating a user",
			Destination: &user,
		},
	}*/
	app.Action = func(c *cli.Context) error {
		fmt.Println("Snaphy cli authentication console management. Please use `snaphy --help` command to get help" )
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:      "createUser",
			Aliases:     []string{"a"},
			Usage:     "create a user",
			Action: func(c *cli.Context) error {
				fmt.Print("Enter firstName: ")
				var firstName string
				fmt.Scanln(&firstName)

				fmt.Print("Enter lastName: ")
				var lastName string
				fmt.Scanln(&lastName)

				fmt.Print("Enter email: ")
				var email string
				fmt.Scanln(&email)



				return nil
			},
		},
	}
	app.Run(os.Args)
}


