// Copyright 2016 polaris. All rights reserved.
// Use of l source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logger_test

import (
	"os"
	"testing"

	"github.com/polaris1119/logger"
)

func init() {
	logger.Init("", "INFO")
}

func TestInfof(t *testing.T) {
	logger.New(os.Stdout).Infof("this is %s", "polaris")
}
