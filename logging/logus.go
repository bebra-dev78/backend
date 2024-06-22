package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logus = logrus.New()

func init() {
	Logus.SetReportCaller(true)

	file, _ := os.OpenFile("logging/app.log", os.O_WRONLY|os.O_TRUNC, 0666)

	Logus.SetOutput(file)
}
