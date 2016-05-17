package cli

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/malnick/cryptorious/config"
)

// Start() is a wrapper for codegansta/CLI implementation
func Start() error {
	printBanner()
	c, err := config.GetConfiguration()
	handleError(err)
	app := cli.NewApp()
	app.Version = c.Version
	app.Name = "cryptorious"
	app.Usage = "CLI-based encryption for passwords and random data"
	app.Authors = []cli.Author{
		{
			Name:  "Jeff Malnick",
			Email: "malnick@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       c.AppYamlPath,
			Usage:       "Path to YAML configuration file `/PATH/TO/CONFIG`",
			Destination: &c.AppYamlPath,
		},
		cli.StringFlag{
			Name:        "user, u",
			Value:       c.UserName,
			Usage:       "Username for vault `USERNAME`.",
			Destination: &c.UserName,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "Decrypt a value in the vault `VALUE`",
			Action: func(c *cli.Context) {
				fmt.Println("Decrypting ", c.Args().First())
			},
		},
		{
			Name:    "encrypt",
			Aliases: []string{"e"},
			Usage:   "Encrypt a vlue for the vault `VALUE`",
			Action: func(c *cli.Context) {
				fmt.Println("Encyrpting ", c.Args().First())
			},
		},
		{
			Name:    "generate-keys",
			Aliases: []string{"gk"},
			Usage:   "Generate a unique RSA public and private key pair for a user specified by user_name or with -user",
			Action: func(c *cli.Context) {
				fmt.Println("Generating new RSA public/private key pair for ", c.Args().First())
			},
		},
	}

	app.Run(os.Args)
	return nil
}

func printBanner() {
	fmt.Println(`_________                            __                   .__                        `)
	fmt.Println(`\_   ___ \ _______  ___.__.______  _/  |_   ____  _______ |__|  ____   __ __   ______`)
	fmt.Println(`/    \  \/ \_  __ \<   |  |\____ \ \   __\ /  _ \ \_  __ \|  | /  _ \ |  |  \ /  ___/`)
	fmt.Println(`\     \____ |  | \/ \___  ||  |_> > |  |  (  <_> ) |  | \/|  |(  <_> )|  |  / \___ \ `)
	fmt.Println(` \______  / |__|    / ____||   __/  |__|   \____/  |__|   |__| \____/ |____/ /____  >`)
	fmt.Println(`        \/          \/     |__|                                                   \/ `)
}

func handleError(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
