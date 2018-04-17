package aws

import (
	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// DefaultRegion is the default AWS region
const DefaultRegion = "us-east-1"

// AWS is the AWS implementation
type AWS struct {
	Client *session.Session
	Config *aws.Config
	Logger *logrus.Logger
}

// New returns an initialized AWS
func New(opts ...Option) (*AWS, error) {
	a := &AWS{
		Client: session.New(),
		Config: aws.NewConfig().WithRegion(DefaultRegion),
		Logger: logrus.New(),
	}

	for _, o := range opts {
		if err := o(a); err != nil {
			return a, err
		}
	}

	return a, nil
}

// Option defines a interface for passing functional options to the fog/aws package.
type Option func(*AWS) error

// OptionLogger implements a functional option for logging facility
func OptionLogger(l *logrus.Logger) Option {
	return func(a *AWS) error {
		a.Logger = l
		return nil
	}
}

// OptionClient allows overriding the default aws client
func OptionClient(s *session.Session) Option {
	return func(a *AWS) error {
		a.Client = s
		return nil
	}
}

// OptionConfig allows overriding the default config.
// Handy for setting region other than DefaultRegion
func OptionConfig(c *aws.Config) Option {
	return func(a *AWS) error {
		a.Config = c
		return nil
	}
}
