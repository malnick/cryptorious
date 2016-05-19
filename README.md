# Cryptorious
1Pass for the CLI.

## Step 1: Generate keys

```
cryptorious generate 
```

Defaults to placing keys in ```$HOME/.ssh/cryptorious_privatekey``` and ```$HOME/.ssh/cryptorious_publickey```.

You can override this with ```--private-key``` and ```--public-key```:

```
cryptorious generate --private-key foo_priv --public-key foo_pub 
```


