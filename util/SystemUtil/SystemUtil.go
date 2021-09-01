package SystemUtil

import (
	"helloworld-blockchain-go/util/LogUtil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func ErrorExit(message string, exception interface{}) {
	LogUtil.Error("system error occurred, and exited, please check the error！"+message, exception)
	os.Exit(1)
}
func SystemRootDirectory() string {
	_, b, _, _ := runtime.Caller(0)
	systemRootDirectory := filepath.Join(filepath.Dir(b), "../..")
	return systemRootDirectory
}
func CallDefaultBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	exec.Command(cmd, args...).Start()
}
