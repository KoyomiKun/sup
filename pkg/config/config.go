package config

import (
	"os"
	"os/exec"

	"github.com/pelletier/go-toml"

	"github.com/KoyomiKun/sup/utils/log"
)

// How to react when the supervised process went down.
type ProcessRestartStrategy string

const (
	RestartStrategyAlways    ProcessRestartStrategy = "always"
	RestartStrategyOnFailure ProcessRestartStrategy = "on-failure"
	RestartStrategyNone      ProcessRestartStrategy = "none"
)

type Config struct {
	SupConfig     Sup     `toml:"sup" comment:"Config related with Sup."`
	ProgramConfig Program `toml:"program" comment:"Config related with the supervised process."`
}

type Sup struct {
	Socket string `toml:"socket" comment:"Path to an unix socket, to which Sup daemon will be listening." default:"./sup.sock"`
}

type Program struct {
	Process Process `toml:"process" comment:"Config related with process."`
	Log     Log     `toml:"log" comment:"Config related with log."`
}

type Process struct {
	Path            string                 `toml:"path" comment:"Path to an executable, which would spawn the supervised process."`
	Args            []string               `toml:"args" comment:"Arguments to the supervised process."`
	Envs            map[string]string      `toml:"envs" comment:"Environment variables to the supervised process."`
	WorkDir         string                 `toml:"workDir" comment:"Working directory of the supervised process. Current directory by default." default:"./"`
	AutoStart       bool                   `toml:"autoStart" comment:"Start the process as Sup goes up. False by default." default:"false"`
	StartSeconds    int                    `toml:"startSeconds" comment:"Sup waits 'startSeconds' after each start to avoid the process restarts too rapidly." default:"5"`
	RestartStrategy ProcessRestartStrategy `toml:"restartStrategy" comment:"How to react when the supervised process went down. One of 'on-failure', 'always', 'none'. 'on-failure' by default." default:"on-failure"`
}

type Log struct {
	Path            string `toml:"path" comment:"Path where to save the current un-rotated log. Using basename of the supervised process by default."`
	MaxSize         int    `toml:"maxSize" comment:"Maximum size in MiB of the log file before it gets rotated. 128 MiB by default." default:"134217728"`
	MaxDays         int    `toml:"maxDays" comment:"Maximum days to retain old log files based on the UTC time encoded in their filename. unlimited by default." default:"0"`
	MaxBackups      int    `toml:"maxBackups" comment:"Maximum number of old log files to retain. Retaining all old log files by default. 32 by default." default:"32"`
	Compress        bool   `toml:"compress" comment:"Whether the rotated log files should be compressed with gzip, no compression by default." default:"false"`
	MergeCompressed bool   `toml:"mergeCompressed" comment:"Whether the gzipped backups should be merged, no by default." default:"false"`
}

func NewConfig(filepath string) *Config {

	config := &Config{}

	tf, err := toml.LoadFile(filepath)
	if err != nil {
		log.Fatalf("read config %s failed: %v", filepath, err)
	}

	if err := tf.Unmarshal(config); err != nil {
		log.Fatalf("unmarshal config %s failed: %v", filepath, err)
	}

	if len(config.ProgramConfig.Process.RestartStrategy) == 0 {
		config.ProgramConfig.Process.RestartStrategy = RestartStrategyOnFailure
	}

	if len(config.ProgramConfig.Process.WorkDir) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("get working dir failed: %v", err)
		}
		config.ProgramConfig.Process.WorkDir = wd
	}

	path, err := exec.LookPath(config.ProgramConfig.Process.Path)
	if err != nil {
		log.Fatalf("lookup executable bin in %s failed: %v", config.ProgramConfig.Process.Path, err)
	}
	config.ProgramConfig.Process.Path = path

	return config
}
