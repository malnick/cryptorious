package cli

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// Start() is a wrapper for codegansta/CLI implementation
func Start() {
	app := cli.NewApp()
	app.Name = "cryptorious"
	app.Usage = "CLI-based encryption for passwords and random data"

	app.Commands = []cli.Command{
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "Decrypt a value in the vault <$value>",
			Action: func(c *cli.Context) {
				fmt.Println("Decrypting ", c.Args().First())
			},
		},
		{
			Name:    "encrypt",
			Aliases: []string{"e"},
			Usage:   "Encrypt a vlue for the vault <$value>",
			Action: func(c *cli.Context) {
				fmt.Println("Encyrpting ", c.Args().First())
			},
		},
		{
			Name:    "generate-keys",
			Aliases: []string{"gk"},
			Usage:   "Generate a unique RSA public and private key pair for a user <$user>",
			Action: func(c *cli.Context) {
				fmt.Println("Generating new RSA public/private key pair for ", c.Args().First())
			},
		},
	}
}
