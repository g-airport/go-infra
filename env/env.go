package env

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/go_infrastructure/validate"
)

var (
	Dir        string
	RunDir     string
	LogDir     string
	LogPath    string
	ConfigPath string
	ConfigDir  string
	Pid        int
	Hostname   string
)

func init() {

	file, _ := filepath.Abs(os.Args[0])
	dir := filepath.Dir(file)

	Dir = filepath.Dir(dir + "..")

	LogDir = Dir + "/log/"
	if !IsExist(LogDir) {
		if err := os.MkdirAll(LogDir, os.ModePerm); err != nil {
			ErrExit(err)
		}
	}
	LogPath = LogDir + filepath.Base(os.Args[0]) + ".log"

	ConfigDir = Dir + "/conf/"
	ConfigPath = ConfigDir + filepath.Base(os.Args[0]) + ".conf"
	if !filepath.IsAbs(ConfigPath) {
		ConfigPath = Dir + "/" + ConfigPath
	}

	RunDir, _ = os.Getwd()
	RunDir, _ = filepath.Abs(RunDir)

	ConfigPath, _ = filepath.Abs(ConfigPath)
	ConfigDir, _ = filepath.Abs(ConfigDir)

	LogDir, _ = filepath.Abs(LogDir)
	LogPath, _ = filepath.Abs(LogPath)

	Pid = os.Getpid()

	hostname, err := os.Hostname()
	if err != nil {
		ErrExit(err)
	}

	Hostname = hostname
}

// ---------------------------------------------------------
func AbsPath(path string) string {

	if !filepath.IsAbs(path) {
		path = Dir + "/" + path
	}

	path, _ = filepath.Abs(path)

	return path
}

func ErrExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ErrExitFunc(fs ...func() error) {
	for _, f := range fs {
		ErrExit(f())
	}
}

func IsExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
	}
	return false
}
