# [![CircleCI](https://circleci.com/gh/malnick/cryptorious.svg?style=svg)](https://circleci.com/gh/malnick/cryptorious)

Like 1Password but for the CLI. Stores your encrypted data in eyaml using generic SSH keys as the basis for encryption/decryption so you never have to type a password to get your passwords ever again.

## Manpage
### Main Menu
```
_________                            __                   .__                        
\_   ___ \ _______  ___.__.______  _/  |_   ____  _______ |__|  ____   __ __   ______
/    \  \/ \_  __ \<   |  |\____ \ \   __\ /  _ \ \_  __ \|  | /  _ \ |  |  \ /  ___/
\     \____ |  | \/ \___  ||  |_> > |  |  (  <_> ) |  | \/|  |(  <_> )|  |  / \___ \ 
 \______  / |__|    / ____||   __/  |__|   \____/  |__|   |__| \____/ |____/ /____  >
        \/          \/     |__|                                                   \/ 
NAME:
   cryptorious - CLI-based encryption for passwords and random data

USAGE:
   cryptorious [global options] command [command options] [arguments...]
   
VERSION:
   1.0.0
   
AUTHOR(S):
   Jeff Malnick <malnick@gmail.com> 
   
COMMANDS:
    rotate      Rotate your cryptorious SSH keys and vault automatically
    delete      Remove an entry from the cryptorious vault
    decrypt     Decrypt a value in the vault `VALUE`
    encrypt     Encrypt a value for the vault `VALUE`
    generate    Generate a unique RSA public and private key pair for a user specified by user_name or with -user

GLOBAL OPTIONS:
   --vault-path, --vp "/home/malnick/.cryptorious/vault.yaml"           Path to vault.yaml.
   --private-key, --priv "/home/malnick/.ssh/cryptorious_privatekey"    Path to private key.
   --public-key, --pub "/home/malnick/.ssh/cryptorious_publickey"       Path to public key.
   --debug                                                              Debug/Verbose log output.
   --help, -h                                                           show help
   --version, -v                                                        print the version
```
### Decrypt Sub Menu
```   
_________                            __                   .__                        
\_   ___ \ _______  ___.__.______  _/  |_   ____  _______ |__|  ____   __ __   ______
/    \  \/ \_  __ \<   |  |\____ \ \   __\ /  _ \ \_  __ \|  | /  _ \ |  |  \ /  ___/
\     \____ |  | \/ \___  ||  |_> > |  |  (  <_> ) |  | \/|  |(  <_> )|  |  / \___ \ 
 \______  / |__|    / ____||   __/  |__|   \____/  |__|   |__| \____/ |____/ /____  >
        \/          \/     |__|                                                   \/ 
NAME:
   cryptorious decrypt - Decrypt a value in the vault `VALUE`

USAGE:
   cryptorious decrypt [command options] [arguments...]

OPTIONS:
   --copy, -c           Copy decrypted password to clipboard automatically
   --goto, -g           Open your default browser to https://<key_name> and login automatically.
   --timeout, -t "10"   Timeout in seconds for the decrypt session window to expire.

```   


## Step 0: Build && Alias

Build it and install: `make install`

Add to your `.[bash | zsh | whatever]rc`: `alias cpt=cryptorious`

## Step 1: Generate keys

```
cryptorious generate 
```

Defaults to placing keys in ```$HOME/.ssh/cryptorious_privatekey``` and ```$HOME/.ssh/cryptorious_publickey```.

You can override this with ```--private-key``` and ```--public-key```:

```
cryptorious generate --private-key foo_priv --public-key foo_pub 
```

## Step 2: Encrypt

```
cryptorious encrypt github  
```

Will open a ncurses window and prompt you for username, password and a secure note. All input is optional. 


## Step 3: Decrypt 

```
cryptorious decrypt thing
```

Will open a ncurses window with the decrypted vault entry. 

Forgo the the ncurses window and copy the decrypted password stright to the system clipboard? 
```
cryptorious decrypt -[c]opy thing
```
No printing, just a message that your decrypted password is now available in the paste buffer for your user. 

## Step 4: Rotate Keys & Vault
Compromised your keys? Not a problem. 

```
cryptorious rotate
```

1. Backs up your old keys to `keyPath.bak`
1. Backs up your old vault to `vaultPath.bak`
1. Generates new keys to `keyPath`
1. Decrypts vault using `privateKey.bak` and encrypts vault in place with new `privateKey`
1. Writes the vault back to disk at `vaultPath`
