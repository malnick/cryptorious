package kms

import (
	"crypto/aes"
	"encoding/base64"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	thisAWS "github.com/malnick/cryptorious/aws"
	"github.com/malnick/cryptorious/crypto"
)

const (
	keyLength = 32
)

// Impl implements KMS needs
type Impl struct {
	Client *kms.KMS
	Logger *logrus.Logger
}

// New returns an initialized KMS Impl
func New(a *thisAWS.AWS) Impl {
	return Impl{
		Client: kms.New(a.Client, a.Config),
		Logger: a.Logger,
	}
}

// GetRandom queries KMS for a random bytes array and returns it with an error
// if applicable
func (k *Impl) GetRandom(size int64) ([]byte, error) {
	k.Logger.Debugf("generating randomness")
	params := &kms.GenerateRandomInput{
		NumberOfBytes: aws.Int64(size), // 16 is the AES blocksize, for reference
	}
	resp, err := k.Client.GenerateRandom(params)
	k.Logger.Debugf("respose from GenerateRandomInput(): %s", resp.Plaintext)

	if err != nil {
		return nil, err
	}

	k.Logger.Debugf("success")
	return resp.Plaintext, nil
}

// Encrypt returns a bytes encoded payload{} initiliazed with the ciphertext of
// the plaintext parameter and associated key generated via KMS
func (k *Impl) Encrypt(plaintext []byte, keyID string) ([]byte, []byte, []byte, error) {
	k.Logger.Debugf("encrypting %s via KMS", plaintext)

	// Get the initialization vector from KMS GenerateRandom()
	iv, err := k.GetRandom(aes.BlockSize)
	if err != nil {
		return nil, nil, nil, err
	}

	// Generate data key
	kmsDataKeyInput := &kms.GenerateDataKeyInput{
		KeyId:         aws.String(keyID),
		NumberOfBytes: aws.Int64(keyLength),
	}

	rsp, err := k.Client.GenerateDataKey(kmsDataKeyInput)
	if err != nil {
		return nil, nil, nil, err
	}
	k.Logger.Debugf("response from KMS %+v", rsp)

	key := rsp.CiphertextBlob
	k.Logger.Debugf("ciphertext blob:\n%s", key)

	// Create key
	sessKey := make([]byte, keyLength)
	copy(sessKey[:], rsp.Plaintext)
	k.Logger.Debugf("created key %s", sessKey)

	// Encrypt message
	ciphertext, err := crypto.AES256EncryptWithCBC(sessKey, iv, plaintext)
	if err != nil {
		return nil, nil, nil, err
	}
	k.Logger.Debugf("ciphertext: %s", ciphertext)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), []byte(base64.StdEncoding.EncodeToString(key)), []byte(base64.StdEncoding.EncodeToString(iv)), nil
}

// Decrypt returns the plaintext of a ciphertext that is assumed to be a bytes
// encoded payload type.
func (k *Impl) Decrypt(ciphertext, key, iv []byte) (string, error) {
	k.Logger.Debugf("decrypting via KMS:\n%s", ciphertext)
	k.Logger.Debugf("using key %s", key)

	decodedKey, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if _, iserr := err.(base64.CorruptInputError); iserr {
		k.Logger.Debugf("ciphertext input does not appear to be base64 encoded, falling back to original ciphertext input")
		decodedCiphertext = ciphertext
	}

	decodedIv, err := base64.StdEncoding.DecodeString(string(iv))
	if _, iserr := err.(base64.CorruptInputError); iserr {
		k.Logger.Debugf("IV input does not appear to be base64 encoded, falling back to original IV input")
		decodedIv = iv
	}

	kmsDecryptInput := &kms.DecryptInput{
		CiphertextBlob: decodedKey,
	}

	decryptRsp, err := k.Client.Decrypt(kmsDecryptInput)
	if err != nil {
		return "", err
	}

	sessKey := make([]byte, keyLength)
	copy(sessKey[:], decryptRsp.Plaintext)
	k.Logger.Debugf("session decryption key: %s", sessKey)

	// Decrypt message
	plaintext, err := crypto.AES256DecryptWithCBC(sessKey, decodedIv, string(decodedCiphertext))
	if err != nil {
		return "", err
	}
	k.Logger.Debugf("decrypted ciphertext:\n%s", plaintext)
	k.Logger.Debugf("successfully decrypted payload")

	return plaintext, nil
}
