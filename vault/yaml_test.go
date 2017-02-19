package vault

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
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
		Data: map[string]*VaultSet{
			"testKey": &VaultSet{
				Username: "foo",
				Password: "bar",
			},
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

	if _, ok := tv.Data["testKey"]; !ok {
		t.Error("expected 'testKey' to be present, is not.")
	}

	if tv.Data["testKey"].Username != "foo" {
		t.Error("expected username to be 'foo', got, ", tv.Data["testKey"].Username)
	}

	if tv.Data["testKey"].Password != "bar" {
		t.Error("expected password to be 'bar', got, ", tv.Data["testKey"].Password)
	}
}

func TestUpdate(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "vaultTest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tv := Vault{
		Path: tmpfile.Name(),
		Data: map[string]*VaultSet{
			"testKey": &VaultSet{
				Username: "foo",
				Password: "bar",
			},
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

	updateWithSet := &VaultSet{
		Username: "newUser",
		Password: "newPass",
	}

	err = tv.Update("anotherUser", updateWithSet)
	if err != nil {
		t.Error(err)
	}

	if _, ok := tv.Data["testKey"]; !ok {
		t.Error("expected 'testKey' to be present, is not.")
	}

	if tv.Data["testKey"].Username != "foo" {
		t.Error("expected username to be 'foo', got, ", tv.Data["testKey"].Username)
	}

	if tv.Data["testKey"].Password != "bar" {
		t.Error("expected password to be 'bar', got, ", tv.Data["testKey"].Password)
	}

	if _, ok := tv.Data["anotherUser"]; !ok {
		t.Error("expected 'anotherUser' to be present, is not.")
	}

	if tv.Data["anotherUser"].Username != "newUser" {
		t.Error("expected username to be 'newUser', got, ", tv.Data["anotherUser"].Username)
	}

	if tv.Data["anotherUser"].Password != "newPass" {
		t.Error("expected passwrod to be 'newPass', got, ", tv.Data["anotherUser"].Password)
	}
}
