package action

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/malnick/cryptorious/config"
)

var (
	tmpPriv = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQD25aZJq71JXxeRdkm6sN3nJSu9rNxDRtJftnMR//j0V9tMLHMD
yg/jQQDBy8XVtFigd8YFr92mtwhzzp7QfCPp1ZUuHqoK/sb84MIPtNL1fcdutKlP
YoouM1XAvicd4YHpgft391QmjCDxiSXusqblZToSGEZP81JU5JpIyIn13wIDAQAB
AoGBALLUYh6oW0FCtEJzKDImG4RpwwXup9e++2/CKhTWkA8Dd97zrxcGi31yPscf
/pqstyj7uB20ZVp05pVjCls+H4+6IWzxEB6Eu2e0IIIvDmGq4+FS1rXAZ2LsH14k
9qb44hA6X9F81r2UnVE0w3eIABcZVZOFnSDZGRgWqrbwpVORAkEA/ffn465i8LYB
ByGPRsXAA7NVXbYVAUH6lgLZ70+iQsHlH8kDsWrfOfQlN9Vz6vqor7ELrnPTXDm5
o9MDJqR+VwJBAPjfQzqxBQqvcvvdrZk68WZPa7a3/PxNifbjnD+AoQtrL0sTuyZP
/cT558ihdvAJN5rlR0NZuWqkWUQChuV+/7kCQB26/7Jvn7V+GPC0xQkL7UaBn+Sw
hBT5nFQjUU/qipw2BpSJ+5yxXiByrEi0/DTt0wF+QFfTx1Jsj4bWFPBZIVECQQCS
w1bqLapDeuPcFAZj7padNwjWX/oY78EEj4V9DWXrTfI93AzpsxJ8LsO5VT7Gnyjj
d4Jm/WqSsQjTSooynIpBAkEA3ecVC4I3hVITueZ8KxoLBerIslecpvfdB/FehhZ5
E+N+L/j5fHPlqfK9GdXuFjDRed8zBsopkgsxw8cSRUlVig==
-----END RSA PRIVATE KEY-----`

	tmpPub = `-----BEGIN RSA PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD25aZJq71JXxeRdkm6sN3nJSu9
rNxDRtJftnMR//j0V9tMLHMDyg/jQQDBy8XVtFigd8YFr92mtwhzzp7QfCPp1ZUu
HqoK/sb84MIPtNL1fcdutKlPYoouM1XAvicd4YHpgft391QmjCDxiSXusqblZToS
GEZP81JU5JpIyIn13wIDAQAB
-----END RSA PUBLIC KEY-----`

	badPrivHead = `-----FOO-----
MIICXgIBAAKBgQD25aZJq71JXxeRdkm6sN3nJSu9rNxDRtJftnMR//j0V9tMLHMD
yg/jQQDBy8XVtFigd8YFr92mtwhzzp7QfCPp1ZUuHqoK/sb84MIPtNL1fcdutKlP
YoouM1XAvicd4YHpgft391QmjCDxiSXusqblZToSGEZP81JU5JpIyIn13wIDAQAB
AoGBALLUYh6oW0FCtEJzKDImG4RpwwXup9e++2/CKhTWkA8Dd97zrxcGi31yPscf
/pqstyj7uB20ZVp05pVjCls+H4+6IWzxEB6Eu2e0IIIvDmGq4+FS1rXAZ2LsH14k
9qb44hA6X9F81r2UnVE0w3eIABcZVZOFnSDZGRgWqrbwpVORAkEA/ffn465i8LYB
ByGPRsXAA7NVXbYVAUH6lgLZ70+iQsHlH8kDsWrfOfQlN9Vz6vqor7ELrnPTXDm5
o9MDJqR+VwJBAPjfQzqxBQqvcvvdrZk68WZPa7a3/PxNifbjnD+AoQtrL0sTuyZP
/cT558ihdvAJN5rlR0NZuWqkWUQChuV+/7kCQB26/7Jvn7V+GPC0xQkL7UaBn+Sw
hBT5nFQjUU/qipw2BpSJ+5yxXiByrEi0/DTt0wF+QFfTx1Jsj4bWFPBZIVECQQCS
w1bqLapDeuPcFAZj7padNwjWX/oY78EEj4V9DWXrTfI93AzpsxJ8LsO5VT7Gnyjj
d4Jm/WqSsQjTSooynIpBAkEA3ecVC4I3hVITueZ8KxoLBerIslecpvfdB/FehhZ5
E+N+L/j5fHPlqfK9GdXuFjDRed8zBsopkgsxw8cSRUlVig==
-----END RSA PRIVATE KEY-----`

	badPubHead = `-----FOO-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD25aZJq71JXxeRdkm6sN3nJSu9
rNxDRtJftnMR//j0V9tMLHMDyg/jQQDBy8XVtFigd8YFr92mtwhzzp7QfCPp1ZUu
HqoK/sb84MIPtNL1fcdutKlPYoouM1XAvicd4YHpgft391QmjCDxiSXusqblZToS
GEZP81JU5JpIyIn13wIDAQAB
-----END RSA PUBLIC KEY-----`
)

func TestGenerateKeysPath(t *testing.T) {
	// Generate temp pub and priv keys
	privFile, err := ioutil.TempFile(os.TempDir(), "privkey")
	defer os.Remove(privFile.Name())
	if err != nil {
		t.Error("Error writing temp file for private key")
	}
	pubFile, err := ioutil.TempFile(os.TempDir(), "pubkey")
	defer os.Remove(pubFile.Name())
	if err != nil {
		t.Error("Error writing temp file for public key")
	}
	privFile.Write([]byte(tmpPriv))
	pubFile.Write([]byte(tmpPub))

	// Config setup
	config := config.Config{
		PrivateKeyPath: privFile.Name(),
		PublicKeyPath:  pubFile.Name(),
	}

	genErr := GenerateKeys(config)
	if genErr != nil {
		t.Error("Expected no errors generating keys, got ", genErr.Error)
	}

	config.PrivateKeyPath = "/foo"
	genErr = GenerateKeys(config)
	if genErr == nil {
		t.Error("Expected errors with bad priv key path, got ", genErr.Error)
	}

	config.PublicKeyPath = "/foo"
	config.PrivateKeyPath = privFile.Name()
	genErr = GenerateKeys(config)
	if genErr == nil {
		t.Error("Expected errors with bad pub key path, got ", genErr.Error)
	}
	config.PublicKeyPath = pubFile.Name()

	privFile.Write([]byte(badPrivHead))
	genErr = GenerateKeys(config)
	if genErr != nil {
		t.Error("Expected errors with bad private key header, got", genErr)
	}
	privFile.Write([]byte(tmpPriv))

	pubFile.Write([]byte(badPubHead))
	genErr = GenerateKeys(config)
	if genErr != nil {
		t.Error("Expected errors with bad private key header, got", genErr)
	}
}
