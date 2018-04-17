package action

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/malnick/cryptorious/vault"
)

func TestDeleteVaultEntry(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testVault")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	tv, err := vault.New(tmpfile.Name())
	if err != nil {
		t.Error(err)
	}

	err = tv.Add("testKey", &testVaultSet)
	if err != nil {
		t.Error(err)
	}

	err = tv.Delete("testKey")
	if err != nil {
		t.Error(err)
	}

	if _, ok := tv.Data["testKey"]; ok {
		t.Error("Key for 'testKey' still exists after .Delete(key) called, expected to not exist.")
	}

}
