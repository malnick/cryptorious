package cli

import (
	"errors"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/malnick/cryptorious/action"
	vaultConfig "github.com/malnick/cryptorious/config"
)

// Start is a wrapper for codegansta/CLI implementation
func Start() error {
	printBanner()
	config, cErr := vaultConfig.GetConfiguration()
	handleError(cErr)
	app := cli.NewApp()
	app.Version = config.Version
	app.Name = printBanner()
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
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Debug/Verbose log output.",
			Destination: &config.DebugMode,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "rename",
			Usage: "Rename an entry in the vault",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "old, o",
					Usage: "Name of old entry name [key] in vault",
				},
				cli.StringFlag{
					Name:  "new, n",
					Usage: "Name of new entry name [key] in vault",
				},
			},
			Action: func(c *cli.Context) {
				setLogger(config.DebugMode)
				handleError(action.RenameVaultEntry(c.String("old"), c.String("new"), config.VaultPath))
			},
		},
		{
			Name:  "rotate",
			Usage: "Rotate your cryptorious vault",
			Action: func(c *cli.Context) {
				setLogger(config.DebugMode)
				handleError(action.RotateVault(config))
			},
		},
		{
			Name:  "delete",
			Usage: "Remove an entry from the cryptorious vault",
			Action: func(c *cli.Context) {
				setLogger(config.DebugMode)
				handleError(action.DeleteVaultEntry(c.Args().First(), config.VaultPath))
			},
		},
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "Decrypt a value in the vault `VALUE`",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "copy, c",
					Usage:       "Copy decrypted password to clipboard automatically",
					Destination: &config.Clipboard,
				},
				cli.BoolFlag{
					Name:        "goto, g",
					Usage:       "Open your default browser to https://<key_name> and login automatically.",
					Destination: &config.Goto,
				},
				cli.IntFlag{
					Name:        "timeout, t",
					Usage:       "Timeout in seconds for the decrypt session window to expire.",
					Value:       10,
					Destination: &config.DecryptSessionTimeout,
				},
			},
			Action: func(c *cli.Context) {
				setLogger(config.DebugMode)
				handleError(action.Decrypt(c.Args().First(), config))
			},
			Subcommands: []cli.Command{
				{
					Name:  "all",
					Usage: "Decrypt the entire vault",
					Action: func(c *cli.Context) {
						setLogger(config.DebugMode)
						handleError(action.PrintAll(config))
					},
				},
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
					Name:        "key-arn",
					Usage:       "KMS key ARN",
					Value:       "",
					Destination: &config.KMSKeyARN,
				},
			},
			Action: func(c *cli.Context) {
				setLogger(config.DebugMode)
				key := c.Args().First()
				if len(c.Args()) != 1 {
					handleError(errors.New("Must pass value for key in arguments to `encrypt`: `cryptorious encrypt $KEY`"))
				}

				handleError(action.Encrypt(key, config))
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate a RSA keys or a secure password.",
			Subcommands: []cli.Command{
				{
					Name:  "keys",
					Usage: "Generate KMS key for cryptorious",
					Action: func(c *cli.Context) {
						setLogger(config.DebugMode)
						fmt.Println("Generating new KMS key pair for ", c.Args().First())
						handleError(action.GenerateKeys(config))
					},
				},
				{
					Name:  "password",
					Usage: "Generate a random password",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "length, l",
							Usage: "Length of the password to generate",
							Value: 15,
						},
					},
					Action: func(c *cli.Context) {
						setLogger(config.DebugMode)
						handleError(action.NewPassword(c.Int("length")))
					},
				},
			},
		},
	}

	app.Run(os.Args)
	return nil
}

func printBanner() string {
	banner := `
 _________                            __                   .__                        
 \_   ___ \ _______  ___.__.______  _/  |_   ____  _______ |__|  ____   __ __   ______
 /    \  \/ \_  __ \<   |  |\____ \ \   __\ /  _ \ \_  __ \|  | /  _ \ |  |  \ /  ___/
 \     \____ |  | \/ \___  ||  |_> > |  |  (  <_> ) |  | \/|  |(  <_> )|  |  / \___ \ 
  \______  / |__|    / ____||   __/  |__|   \____/  |__|   |__| \____/ |____/ /____  >
         \/          \/     |__|                                                   \/ 
`
	return banner
}

func handleError(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func setLogger(debug bool) {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}
