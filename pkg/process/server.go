package process

import (
	"net/rpc"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/KoyomiKun/sup/pkg/config"
	"github.com/KoyomiKun/sup/utils/log"
)

type Server struct {
	server *rpc.Server
}

func NewServer(
	config config.Process,
	logConfig config.Log) {

	// init envs
	envs := os.Environ()
	newEnvs := make([]string, 0, len(envs)+len(config.Envs))

	for _, env := range envs {
		eqi := strings.Index(env, "=")
		if eqi == -1 {
			log.Warnf("invalid env %s", env)
			continue
		}
		newEnvs = append(newEnvs, env)
	}

	for k, v := range config.Envs {
		newEnvs = append(newEnvs, k+"="+v)
	}

	cmd := exec.Command(config.Path, config.Args...)
	cmd.Env = newEnvs
	cmd.Dir = config.WorkDir

	// init logger
	logger, err := NewFileWriter(
		WithFilename(logConfig.Path),
		WithMaxBytes(int64(logConfig.MaxSize)*1024*1024),
		WithMaxBackups(logConfig.MaxBackups),
		WithCompress(logConfig.Compress),
		WithMergeCompressedBackups(logConfig.MergeCompressed),
		WithMaxAge(time.Hour*24*time.Duration(logConfig.MaxDays)),
	)
	if err != nil {
		log.Fatalf("init rotate logger failed: %v", err)
	}

	// init controller
	controller := &Controller{}
}
