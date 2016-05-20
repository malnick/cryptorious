# [![CircleCI](https://circleci.com/gh/malnick/cryptorious.svg?style=svg)](https://circleci.com/gh/malnick/cryptorious)

```
_________                            __                   .__
\_   ___ \ _______  ___.__.______  _/  |_   ____  _______ |__|  ____   __ __   ______
/    \  \/ \_  __ \<   |  |\____ \ \   __\ /  _ \ \_  __ \|  | /  _ \ |  |  \ /  ___/
\     \____ |  | \/ \___  ||  |_> > |  |  (  <_> ) |  | \/|  |(  <_> )|  |  / \___ \
 \______  / |__|    / ____||   __/  |__|   \____/  |__|   |__| \____/ |____/ /____  >
        \/          \/     |__|                                                   \/
```
Like 1Password but for the CLI. Stores your encrypted data in eyaml using generic SSH keys as the basis for encryption/decryption so you never have to type a password to get your passwords ever again.

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
cryptorious encrypt github p@$$
```

Will encrypt ```p@$$``` using your private key and dump out a PKCS1 hash, writing it to ```$HOME/.cryptorious/vault.yaml``` with the key `github`.


## Step 3: Decrypt 

```
cryptorious decrypt github
```

Will write the decrypted password stored for the key `github` to stdout. 

## YAML?
Yes, YAML. 

Why? Because it's easy to digest visually and provides the easiest way to store the encrypted hash with a lookup key. 
