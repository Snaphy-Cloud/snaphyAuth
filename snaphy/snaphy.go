package main

import (
	"os"
	"github.com/codegangsta/cli"
	"fmt"
	"snaphyAuth/models"
	"github.com/asaskevich/govalidator"
	"github.com/Snaphy-Cloud/Validate"
	"github.com/ttacon/chalk"
	"strconv"
)







func main() {

	app := cli.NewApp()
	app.Name = "snaphy"
	app.Usage = "snaphy Admin CLI control!"
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
			Aliases:     []string{"cu"},
			Usage:     "create a user",
			Action: func(c *cli.Context) error {
				fmt.Print("Enter firstName: ")
				var firstName string
				fmt.Scanln(&firstName)
				if govalidator.IsNull(firstName){
					fmt.Print("FirstName is a required field")
					fmt.Print("Enter firstName: ")

					fmt.Scanln(&firstName)
				}


				fmt.Print("Enter lastName: ")
				var lastName string
				fmt.Scanln(&lastName)

				fmt.Print("Enter email: ")
				var email string
				fmt.Scanln(&email)
				if !govalidator.IsEmail(email){
					fmt.Print("Email must a a valid email")
					fmt.Print("Enter email: ")
					fmt.Scanln(&email)
				}

				authUser := new(models.AuthUser)
				authUser.FirstName = firstName
				authUser.LastName = lastName
				authUser.Email = email
				authUser.Status = "active"

				id, err := authUser.Create()
				fmt.Println(id, err)
				if err != nil{
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println("User created successfully, ID: ", id)
				fmt.Println(authUser)
				return err
			},
		},
		{
			Name:      "deleteUser",
			Aliases:     []string{"du"},
			Usage:     "delete a user. Accept a user id to delete a user from database.",
			Action: func(c *cli.Context) error {
				var (
					id string
					err error
				)
				Validate.GetInput(&id, "Enter user Id: ", func(value string)(message string, isValid bool){
					return "User Id must be a number", govalidator.IsNumeric(id)

				})

				authUser := new (models.AuthUser)
				authUser.Id, err = strconv.Atoi(id)

				_, err = authUser.Delete()
				fmt.Println(chalk.Green, "Successfully deleted User.", chalk.Reset)
				return err
			},
		},
		{
			Name:      "fetchUserById",
			Aliases:     []string{"du"},
			Usage:     "fetch a user details by . Accept user Id  as input.",
			Action: func(c *cli.Context) error {
				var (
					id string
					err error
				)
				Validate.GetInput(&id, "Enter user Id: ", func(value string)(message string, isValid bool){
					return "User Id must be a number", govalidator.IsNumeric(id)

				})

				authUser := new (models.AuthUser)
				authUser.Id, err = strconv.Atoi(id)

				err = authUser.GetUser()
				fmt.Println(chalk.Magenta, "Details of user fetched: \n", chalk.Reset)
				fmt.Println(chalk.Magenta, "ID: ", chalk.Green, authUser.Id, chalk.Reset)
				fmt.Println(chalk.Magenta, "Email: ", chalk.Green, authUser.Email, chalk.Reset)
				fmt.Println(chalk.Magenta, "First Name: ", chalk.Green, authUser.FirstName, chalk.Reset)
				fmt.Println(chalk.Magenta, "Last Name: ", chalk.Green, authUser.LastName, chalk.Reset)
				fmt.Println(chalk.Magenta, "Status: ", chalk.Green, authUser.Status, chalk.Reset)

				return err
			},
		},
		{
			Name:      "deactivateUser",
			Usage:     "Deactivated  a user. Sets its status to DEACTIVATED.",
			Action: func(c *cli.Context) error {
				var (
					id string
					err error
				)
				Validate.GetInput(&id, "Enter user Id : ", func(value string)(message string, isValid bool){
					return "User Id must be a number", govalidator.IsNumeric(id)

				})

				authUser := new (models.AuthUser)
				authUser.Id, err = strconv.Atoi(id)
				err = authUser.GetUser()
				if err != nil{
					return err
				}
				_, err = authUser.Deactivate()
				fmt.Println(chalk.Green, "Successfully deactivated User.", chalk.Reset)
				return err
			},
		},
		{
			Name:      "activateUser",
			Usage:     "activated  a user. Sets its status to ACTIVATED.",
			Action: func(c *cli.Context) error {
				var (
					id string
					err error
				)
				Validate.GetInput(&id, "Enter user Id : ", func(value string)(message string, isValid bool){
					return "User Id must be a number", govalidator.IsNumeric(id)

				})

				authUser := new (models.AuthUser)
				authUser.Id, err = strconv.Atoi(id)
				err = authUser.GetUser()
				if err != nil{
					return err
				}
				_, err = authUser.Activate()
				fmt.Println(chalk.Green, "Successfully activated a User.", chalk.Reset)
				return err
			},
		},
		{
			Name:      "fetchUserApp",
			Usage:     "Fetch Application created by User",
			Action: func(c *cli.Context) error {
				var (
					id string
					err error
				)
				Validate.GetInput(&id, "Enter user Id : ", func(value string)(message string, isValid bool){
					return "User Id must be a number", govalidator.IsNumeric(id)

				})

				authUser := new (models.AuthUser)
				authUser.Id, err = strconv.Atoi(id)
				err = authUser.GetUser()
				if err != nil{
					return err
				}

				num, err := authUser.FetchApps()
				if num == 0{
					fmt.Println(chalk.Green, "No application present for this user.", chalk.Reset)
				}else{
					fmt.Println(chalk.Magenta, "Details of Apps present for given user: \n")
					for i :=0; i< len(authUser.Apps); i++{
						var app *models.Application
						app = authUser.Apps[i]
						//Now print data...
						fmt.Println("\n")
						printAppDetail(app)

					}
				}
				return err
			},
		},


		{
			Name:      "createApp",
			Aliases:     []string{"app"},
			Usage:     "Create an application",
			Action: func(c *cli.Context) error {
				var (
					email string
					err error
					name string
				)

				Validate.GetInput(&email, "Enter user email address : ", func(value string)(message string, isValid bool){
					if govalidator.IsNull(email){
						return "Email address is required", false
					}
					return "Email address must be valid", govalidator.IsEmail(email)
				})

				authUser := new(models.AuthUser)
				authUser.Email = email

				err = authUser.GetCustomUser("Email")
				if err != nil{
					fmt.Println(chalk.Red, "Email not exist! Please provide a valid email.", chalk.Reset)
					return err
				}


				Validate.GetInput(&name, "Enter application name: ", func(value string)(message string, isValid bool){
					if govalidator.IsNull(name){
						return "You must enter an application name.", false
					}
					return "", true
				})

				app := new(models.Application)
				app.Owner = authUser
				app.Name = name
				app.Status = "active"
				err = app.Create()
				fmt.Println(chalk.Green, "Application created successfully\n", chalk.Reset)
				printAppDetail(app)
				return err
			},
		},
		{
			Name:      "deleteApp",
			Usage:     "Delete an application",
			Action: func(c *cli.Context) error {
				var (
					err error
					name string
					num int64
				)

				Validate.GetInput(&name, "Enter application name: ", func(value string)(message string, isValid bool){

					return "Application name is required", govalidator.IsNull(name)
				})


				app := new(models.Application)
				app.Name = name

				num, err = app.Delete()
				if err != nil{
					fmt.Println(chalk.Red, "Error deleting application", chalk.Reset)
				}
				if num == 0{
					fmt.Println(chalk.Green, "Application name not found\n", chalk.Reset)

				}else{
					fmt.Println(chalk.Green, "Application deleted successfully\n", chalk.Reset)
				}

				return err
			},
		},
		{
			Name:      "deactivateApp",
			Usage:     "Deactivate an application",
			Action: func(c *cli.Context) error {
				var (
					err error
					name string

				)

				Validate.GetInput(&name, "Enter application name: ", func(value string)(message string, isValid bool){

					return "Application name is required", !govalidator.IsNull(value)
				})


				app := new(models.Application)
				app.Name = name
				err = app.GetApp()

				if err != nil {
					fmt.Println(chalk.Red, "Error deactivating an application", chalk.Reset)
				}else{
					fmt.Println(chalk.Green, "Application deactivated successfully\n", chalk.Reset)
					_, err = app.Deactivate()
					printAppDetail(app)
				}



				return err
			},
		},
		{
			Name:      "activateApp",
			Usage:     "Activate an application",
			Action: func(c *cli.Context) error {
				var (
					err error
					name string
				)

				Validate.GetInput(&name, "Enter application name: ", func(value string)(message string, isValid bool){

					return "Application name is required", !govalidator.IsNull(value)
				})


				app := new(models.Application)
				app.Name = name
				err = app.GetApp()

				if err != nil{
					fmt.Println(chalk.Red, "Error activating application", chalk.Reset)
				}else{
					fmt.Println(chalk.Green, "Application activated successfully\n", chalk.Reset)
					_, err = app.Activate()
					printAppDetail(app)
				}


				return err
			},
		},
		{
			Name:      "fetchAppTokens",
			Usage:     "Fetch Tokens for Application",
			Action: func(c *cli.Context) error {
				var (
					name string
					err error
					app *models.Application
				)
				Validate.GetInput(&name, "Enter App name : ", func(value string)(message string, isValid bool){
					return "App name is required", !govalidator.IsNull(value)

				})

				app = new (models.Application)
				app.Name = name
				err = app.GetApp()

				if err != nil{
					fmt.Println(err)
					return err
				}

				fmt.Println("App details \n")
				printAppDetail(app)

				num, err := app.FetchAppTokens()

				if num == 0{
					fmt.Println(chalk.Green, "No token present for this application.", chalk.Reset)
				}else{
					fmt.Println(chalk.Magenta, "Details of Tokens present for given application: \n")
					for _, token := range app.TokenInfo{
						fmt.Println("\n")
						printTokenDetail(token)
					}

				}
				return err
			},
		},
		{
			Name:      "generateAppTokens",
			Usage:     "Generate Tokens for Application",
			Action: func(c *cli.Context) error {
				var (
					name string
					err error
					app *models.Application
					token *models.Token
				)
				Validate.GetInput(&name, "Enter App name : ", func(value string)(message string, isValid bool){
					return "App name is required", !govalidator.IsNull(value)

				})

				app = new (models.Application)

				app.Name = name
				err = app.GetApp()

				if err != nil{
					fmt.Println(err)
					return err
				}


				fmt.Println("App details \n")
				printAppDetail(app)

				token = new(models.Token)
				token.Application = app
				token.Status = "active"
				_, err = token.Create()

				if err != nil{
					fmt.Println(chalk.Red, "Error creating token for application.", chalk.Reset)
					fmt.Println(chalk.Red, err, chalk.Reset)
				}else{
					fmt.Println(chalk.Magenta, "Successfully created token for given application: \n")
					fmt.Println("\n")
					printTokenDetail(token)
				}
				return err
			},
		},
		{
			Name:      "downloadPrivateKey",
			Usage:     "Download private key of a token",
			Action: func(c *cli.Context) error {
				var (
					appId string
					err error
					token *models.Token
				)
				Validate.GetInput(&appId, "Enter token AppId: ", func(value string)(message string, isValid bool){
					return "AppId is required", !govalidator.IsNull(value)

				})

				token = new(models.Token)
				token.AppId = appId
				err = token.GetToken()


				if err != nil{
					fmt.Println(chalk.Red, "Error fetching token. wrong AppId", chalk.Reset)
					fmt.Println(chalk.Red, err, chalk.Reset)
					return err
				}else{
					fmt.Println("Token details \n")
					printTokenDetail(token)
				}

				//Now download private key ..
				err = token.DownloadPrivateKey()
				if err != nil{
					fmt.Println(chalk.Red, "Error downloading privateKey file", chalk.Reset)
					fmt.Println(chalk.Red, err, chalk.Reset)
					return err
				}else{
					fmt.Println("\n")
					fmt.Println(chalk.Green, "Successfully downloaded private key file for token.\n", chalk.Reset)
				}
				return err
			},
		},
	}
	app.Run(os.Args)
}


func printAppDetail(app *models.Application){
	fmt.Println(chalk.Magenta, "ID: ", chalk.Green, app.Id, chalk.Reset)
	fmt.Println(chalk.Magenta, "Name: ", chalk.Green, app.Name, chalk.Reset)
	fmt.Println(chalk.Magenta, "Status: ", chalk.Green, app.Status, chalk.Reset)
}


func printTokenDetail(token *models.Token){
	fmt.Println(chalk.Magenta, "ID: ", chalk.Green, token.Id, chalk.Reset)
	fmt.Println(chalk.Magenta, "Public Key: ", chalk.Green, token.PublicKey, chalk.Reset)
	fmt.Println(chalk.Magenta, "Private Key: ", chalk.Green, token.PrivateKey, chalk.Reset)
	fmt.Println(chalk.Magenta, "App Id: ", chalk.Green, token.AppId, chalk.Reset)
	fmt.Println(chalk.Magenta, "App Secret: ", chalk.Green, token.AppSecret, chalk.Reset)
	fmt.Println(chalk.Magenta, "Hash Algorithm: ", chalk.Green, token.HashType, chalk.Reset)
	fmt.Println(chalk.Magenta, "Status: ", chalk.Green, token.Status, chalk.Reset)
}




