package cmd

import "github.com/sirupsen/logrus"

var log = logrus.WithFields(logrus.Fields{
	"package": "cmd",
})
