package vault

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	expectedUsername   = "test-username"
	expectedPassword   = "test-password"
	expectedSecureNote = "test-secure-note"

	testVaultSet = Set{
		Username: EncryptedEntry{
			Ciphertext: "AAAAAAAAAAAAAAAAAAAAALthIXEjHFp2BzZFaM2NzZI=",
			Key:        "AQIBAHjjCeNKbXvpFm2HyIXoKqPvjyyCDoSbmOkO8uSNKqBm0wE/eSqj6B0NnrczvwcLSlzrAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMrzjLCH0Te9Egi1UsAgEQgDtiaSFOxBU3jpEU1+OBCg0fOMTCin1IO7uWMtPB5Rq+MArZiFDNNIl4V7P7WTGZk5p9nEVXZibi0evnKA==",
			IV:         "4VvVxXdouSf+ANN+PJ5+tg==",
		},

		Password: EncryptedEntry{
			Ciphertext: "AAAAAAAAAAAAAAAAAAAAAGbAqmydcY5wALQOAz0/m28=",
			IV:         "fbAcI8NkQbFYWHCtOkCm/Q==",
			Key:        "AQIBAHjjCeNKbXvpFm2HyIXoKqPvjyyCDoSbmOkO8uSNKqBm0wFX/eO80576OpuMIJRk6ZvlAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMfYDlWalA7BKoXxJEAgEQgDvLKQigz+X1DDNGf0KAhWLcBXKh1DItzyZsPodwArUoU0wAj7tb3MWhQXJN4XcDFLfh4sln0ozr4vxtYA==",
		},

		SecureNote: EncryptedEntry{
			Ciphertext: "AAAAAAAAAAAAAAAAAAAAAD+/u65AUMH576+xZgGZexMmaWcbhxwSkHm0SMt1Xh4Z",
			Key:        "AQIBAHjjCeNKbXvpFm2HyIXoKqPvjyyCDoSbmOkO8uSNKqBm0wG7UlVRaAOaTAbr9fqcochxAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMCupa+YgNV/EQizynAgEQgDuQ14WOf0gVBLYvjfXMAREMzygPGVfmNxRsx7SExpyHTRl0tQIE//7LztBvIFCjN7PNuDWACQuiHcW1Yw==`",
			IV:         "DRE6jVdAjY/cIsVsl1/lug==",
		},
	}
)

func TestNew(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "vaultTest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	v, err := New(tmpfile.Name())
	if err != nil {
		t.Error(err)
	}

	if v.Path != tmpfile.Name() {
		t.Errorf("Expected %s, got %s", tmpfile.Name(), v.Path)
	}
}

func TestWriteAndLoad(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "vaultTest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tv := Vault{
		Path: tmpfile.Name(),
		Data: map[string]*Set{
			"test-key": &testVaultSet,
		},
	}

	err = tv.Write()
	if err != nil {
		t.Error(err)
	}

	err = tv.Load()
	if err != nil {
		t.Error(err)
	}

	if _, ok := tv.Data["test-key"]; !ok {
		t.Error("expected 'testKey' to be present, is not.")
	}
}

func TestAdd(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "vaultTest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tv := Vault{Path: tmpfile.Name()}

	err = tv.Write()
	if err != nil {
		t.Error(err)
	}

	err = tv.Load()
	if err != nil {
		t.Error(err)
	}

	updateWithSet := &Set{
		Username: EncryptedEntry{
			Ciphertext: "phony-cipher",
			Key:        "phony-key",
			IV:         "phony-iv",
		},
		Password: EncryptedEntry{
			Ciphertext: "phony-cipher-pw",
			Key:        "phony-key-pw",
			IV:         "phon-iv-pw",
		},
	}

	err = tv.Add("anotherUser", updateWithSet)
	if err != nil {
		t.Error(err)
	}

	if _, ok := tv.Data["anotherUser"]; !ok {
		t.Error("expected 'anotherUser' to be present, is not.")
	}

	if tv.Data["anotherUser"].Username.Ciphertext != "phony-cipher" {
		t.Error("expected username to be 'phony-cipher', got, ", tv.Data["anotherUser"].Username.Ciphertext)
	}

	if tv.Data["anotherUser"].Password.Ciphertext != "phony-cipher-pw" {
		t.Error("expected passwrod to be 'phony-cipher-pw', got, ", tv.Data["anotherUser"].Password.Ciphertext)
	}
}
