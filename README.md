# [![CircleCI](https://circleci.com/gh/malnick/cryptorious.svg?style=svg)](https://circleci.com/gh/malnick/cryptorious)

Like 1Password but for the CLI.

## Manpage
### Main Menu
```
 _________                            __                   .__
 \_   ___ \ _______  ___.__.______  _/  |_   ____  _______ |__|  ____   __ __   ______
 /    \  \/ \_  __ \<   |  |\____ \ \   __\ /  _ \ \_  __ \|  | /  _ \ |  |  \ /  ___/
 \     \____ |  | \/ \___  ||  |_> > |  |  (  <_> ) |  | \/|  |(  <_> )|  |  / \___ \
  \______  / |__|    / ____||   __/  |__|   \____/  |__|   |__| \____/ |____/ /____  >
         \/          \/     |__|                                                   \/
 - CLI-based encryption for passwords and random data

USAGE:
   cryptorious [global options] command [command options] [arguments...]

AUTHOR(S):
   Jeff Malnick <malnick@gmail.com>

COMMANDS:
    rename      Rename an entry in the vault
    delete      Remove an entry from the cryptorious vault
    decrypt     Decrypt a value in the vault `VALUE`
    encrypt     Encrypt a value for the vault `VALUE`
    generate    Generate a RSA keys or a secure password.

GLOBAL OPTIONS:
   --vault-path, --vp "/Users/malnick/.cryptorious/vault.yaml"  Path to vault.yaml.
   --debug                                                      Debug/Verbose log output.
   --help, -h                                                   show help
   --version, -v                                                print the version
```

### Encrypt Sub Menu
```
NAME:
   encrypt - Encrypt a value for the vault `VALUE`

USAGE:
   encrypt [command options] [arguments...]

OPTIONS:
   --key-arn    KMS key ARN
```

### Decrypt Sub Menu
```   
NAME:
   cryptorious decrypt - Decrypt a value in the vault `VALUE`

USAGE:
   cryptorious decrypt [command options] [arguments...]

OPTIONS:
   --copy, -c           Copy decrypted password to clipboard automatically
   --goto, -g           Open your default browser to https://<key_name> and login automatically
   --timeout, -t "10"   Timeout in seconds for the decrypt session window to expire
```   

### Rename Sub Menu
```
NAME:
   cryptorious rename - Rename an entry in the vault

USAGE:
   cryptorious rename [command options] [arguments...]

OPTIONS:
   --old, -o    Name of old entry name [key] in vault
   --new, -n    Name of new entry name [key] in vault
```

### Generate Sub Menu
```
NAME:
 generate - 	Generate a RSA keys or a secure password 

USAGE:
  generate command [command options] [arguments...]

COMMANDS:
   password	[--[l]ength] Generate a random password

OPTIONS:
   --help, -h	show help

```

## Step 0: Build && Alias

Build it and install: `make install`

Add to your `.[bash | zsh | whatever]rc`: `alias cpt=cryptorious`

## Step 1: Add KMS keys to AWS
NOTE: will add cmd for this soon

In your own AWS account, add a KMS key and grant your IAM user access.

## Step 2: Encrypt
NOTE: will add flag for AWS profile soon

Use your AWS profile and encrypt some data:
```
AWS_PROFILE=personal cryptorious encrypt --key-arn=<my_kms_key_arn> github.com
```

Will open a ncurses window and prompt you for username, password and a secure note. All input is optional. 

## Step 3: Decrypt 

```
AWS_PROFILE=personal cryptorious decrypt thing
```

Will open a ncurses window with the decrypted vault entry. 

Forgo the the ncurses window and copy the decrypted password stright to the system clipboard? 
```
cryptorious decrypt -[c]opy thing
```
No printing, just a message that your decrypted password is now available in the paste buffer for your user. 

If you've saved your vault entries with the URI of the site they belong to (i.e., ran `cryptorious encrypt github.com`...) then you can use the `-[g]oto` flag to open your default browser to this URI. Pair it with `-[c]opy` and the shorthand for `[d]ecrypt` and you'll have a fast way of navigating directly to your desired, secure website (let's also assume you've aliased `cpt=cryptorious`):
```
cpt d -g -c github.com
```

## Step 5: Generate Secure Password
The `generate` command also lets you generate random, secure passwords of `n` length:
```
cryptorious generate password --length 20
(yZkj,GX`w7T4x&TaYyw
```

This defaults to a length of 15 if you don't pass --[l]ength.
