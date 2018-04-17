package kms

import (
	"crypto/aes"
	"encoding/base64"
	"os"
	"testing"

	"github.com/malnick/cryptorious/aws"
)

func TestNew(t *testing.T) {
	a, err := aws.New()
	if err != nil {
		t.Error(err)
	}

	k := New(a)

	if k.Client == nil {
		t.Error("expected *kms.KMS, got nil")
	}

	if k.Logger == nil {
		t.Error("expected logrus.Logger, got nil")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	os.Setenv("AWS_PROFILE", "personal")

	a, err := aws.New()
	if err != nil {
		t.Error(err)
	}

	k := New(a)

	expected := "foo"

	ciphertext, key, iv, err := k.Encrypt([]byte(expected), "arn:aws:kms:us-east-1:720066741276:key/3d8b504d-761a-4401-a985-b4ff6673e69e")
	if err != nil {
		t.Error(err)
	}

	if len(ciphertext) == 0 {
		t.Error("expected ciphertext length greater than 0")
	}

	if len(key) == 0 {
		t.Error("expected key length greater than 0")
	}

	decodediv, err := base64.StdEncoding.DecodeString(string(iv))
	if err != nil {
		t.Error(err)
	}

	if len(decodediv) != aes.BlockSize {
		t.Error("expected iv length of ", aes.BlockSize, " got ", len(decodediv))
	}

	cleartext, err := k.Decrypt(ciphertext, key, iv)
	if err != nil {
		t.Error(err)
	}

	if cleartext != string(expected) {
		t.Errorf("ciphertext and cleartext do not match: %s != %s", ciphertext, cleartext)
	}
}

func TestGetRandom(t *testing.T) {
	os.Setenv("AWS_PROFILE", "personal")

	a, err := aws.New()
	if err != nil {
		t.Error(err)
	}

	k := New(a)

	b, err := k.GetRandom(10)
	if err != nil {
		t.Error(err)
	}

	if len(b) == 0 {
		t.Error("expected random byte array greater than 0")
	}
}
