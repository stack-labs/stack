package logrus

import (
	"fmt"
	"io"

	"github.com/stack-labs/stack/logger"
	ls "github.com/stack-labs/stack/plugin/logger/logrus/logrus"
	"github.com/stack-labs/stack/plugin/logger/logrus/lumberjack.v2"
)

func prepareLevelHooks(opts logger.PersistenceOptions, l ls.Level) ls.LevelHooks {
	hooks := make(ls.LevelHooks)

	for _, level := range ls.AllLevels {
		if level <= l {
			fileName := fmt.Sprintf("%s%s%s.log", opts.Dir, pathSeparator, level.String())
			logger.Infof("level %s logs to file: %s", level.String(), fileName)
			// todo default options?
			maxBackups := 14
			if opts.MaxFileSize != 0 {
				maxBackups = opts.MaxBackupSize / opts.MaxFileSize
			}

			hook := &PersistenceLevelHook{
				Writer: &lumberjack.Logger{
					Filename:   fileName,
					MaxSize:    opts.MaxFileSize,
					MaxBackups: maxBackups,
					MaxAge:     opts.MaxBackupKeepDays,
					Compress:   true,
					BackupDir:  opts.BackupDir,
				},
				Fired:  true,
				levels: []ls.Level{level},
			}

			hooks[level] = []ls.Hook{hook}
		}
	}

	return hooks
}

type PersistenceLevelHook struct {
	Writer io.Writer
	Fired  bool
	levels []ls.Level
}

func (hook *PersistenceLevelHook) Levels() []ls.Level {
	return hook.levels
}

func (hook *PersistenceLevelHook) Fire(entry *ls.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}
