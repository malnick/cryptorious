package aws

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func TestNew(t *testing.T) {
	k, err := New()
	if err != nil {
		t.Error("expected no errors calling New(), got", err)
	}

	if k.Client == nil {
		t.Error("expected *session.Session, got nil")
	}

	if k.Logger == nil {
		t.Error("expected logrus.StdLogger, got nil")
	}

	if k.Config == nil {
		t.Error("expected aws.Config, got nil")
	}
}

func TestOptionLogger(t *testing.T) {
	var kmsLog = logrus.New()

	_, err := New(OptionLogger(kmsLog))
	if err != nil {
		t.Error("expected no errors getting new KMS object, got", err)
	}
}

func TestOptionClient(t *testing.T) {
	c := session.New()

	_, err := New(OptionClient(c))
	if err != nil {
		t.Error("expected no errors getting new KMS object, got", err)
	}
}

func TestOptionConfig(t *testing.T) {
	c := aws.NewConfig().WithRegion("us-west-1")

	a, err := New(OptionConfig(c))
	if err != nil {
		t.Error("expected no errors getting new KMS object, got", err)
	}

	if *(a.Config.Region) != "us-west-1" {
		t.Error("expected region us-west-1, got ", a.Config.Region)
	}
}
