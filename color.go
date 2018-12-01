package main

import (
	"fmt"
	"os"
)

const (
	Green  = "%s\033[1;32m%s\033[0m %s\n"
	Cyan   = "%s\033[1;36m%s\033[0m %s\n"
	Yellow = "%s\033[1;33m%s\033[0m %s\n"
	Red    = "%s\033[1;31m%s\033[0m %s\n"
	Debug  = "%s\033[0;36m%s\033[0m %s\n"

	Logo = "\033[1;33m%s\033[0m"
)

func live(host string, nl bool) {
	var prefix string

	if nl {
		prefix = "\n"
	}

	fmt.Printf(Green, prefix, "[LIVE] ", host)
}

func dead(host string, nl bool) {
	var prefix string

	if nl {
		prefix = "\n"
	}

	fmt.Printf(Red, prefix, "[DEAD] ", host)

}

func debug(err error, nl bool) {
	var prefix string

	if nl {
		prefix = "\n"
	}

	fmt.Printf(Debug, prefix, "[ERROR] ", err)
}

func info(msg string, nl bool) {
	var prefix string

	if nl {
		prefix = "\n"
	}

	fmt.Printf(Yellow, prefix, "[MSG] ", msg)
}

func errMsg(errmsg error, nl bool) {
	var prefix string

	if nl {
		prefix = "\n"
	}

	fmt.Printf(Red, prefix, "[ERR] ", errmsg)
	os.Exit(1)
}

func color(color, msg string) string {
	return fmt.Sprintf(color, msg)
}
