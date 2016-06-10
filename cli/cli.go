package cli

import (
	"errors"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/malnick/cryptorious/action"
	"github.com/malnick/cryptorious/config"
)

// Start() is a wrapper for codegansta/CLI implementation
func Start() error {
	printBanner()
	config, cErr := config.GetConfiguration()
	handleError(cErr)
	app := cli.NewApp()
	app.Version = config.Version
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
			Name:        "vault-path, vp",
			Value:       config.VaultPath,
			Usage:       "Path to vault.yaml.",
			Destination: &config.VaultPath,
		},
		cli.StringFlag{
			Name:        "private-key, priv",
			Value:       config.PrivateKeyPath,
			Usage:       "Path to private key.",
			Destination: &config.PrivateKeyPath,
		},
		cli.StringFlag{
			Name:        "public-key, pub",
			Value:       config.PublicKeyPath,
			Usage:       "Path to public key.",
			Destination: &config.PublicKeyPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "Decrypt a value in the vault `VALUE`",
			Action: func(c *cli.Context) {
				handleError(action.Decrypt(c.Args().First(), config))
			},
		},
		{
			Name:            "encrypt",
			Aliases:         []string{"e"},
			Usage:           "Encrypt a value for the vault `VALUE`",
			UsageText:       "Encrypt - encrypt a password and/or note with `KEY` to cryptorious vault.",
			SkipFlagParsing: false,
			HideHelp:        false,
			HelpName:        "encrypt",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username, u",
					Usage: "Optional: `USERNAME` for encrypted value.",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Optional: `PASSWORD` for encrypted value.",
				},
				cli.StringFlag{
					Name:  "note, n",
					Usage: "Optional: `NOTE` for encrypted value.",
				},
			},
			Action: func(c *cli.Context) {
				key := c.Args().First()
				if len(c.Args()) != 1 {
					handleError(errors.New("Must pass value for key in arguments to `encrypt`: `cryptorious encrypt $KEY`"))
				} else {
					handleError(action.Encrypt(
						key,
						&action.VaultSet{
							Username:   c.String("username"),
							Password:   c.String("password"),
							SecureNote: c.String("note"),
						},
						config))
				}
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate a unique RSA public and private key pair for a user specified by user_name or with -user",
			Action: func(c *cli.Context) {
				fmt.Println("Generating new RSA public/private key pair for ", c.Args().First())
				handleError(action.GenerateKeys(config))
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
