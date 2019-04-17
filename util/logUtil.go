package util

import (
	"github.com/lexkong/log"
	"gosignaler-cluster/signalerconst"
)

func InitLogCfg(){

	passLagerCfg := log.PassLagerCfg{
		Writers:       signalerconst.WRITERS,
		LoggerLevel:    signalerconst.LOGGER_LEVEL,
		LoggerFile:     signalerconst.LOGGER_DIR,
		LogFormatText:  signalerconst.LOG_FORMAT_TEXT,
		RollingPolicy:  signalerconst.ROLLING_POLICY,
		LogRotateDate:  signalerconst.LOG_ROTATE_DATE,
		LogRotateSize:  signalerconst.LOG_ROTATE_SIZE,
		LogBackupCount: signalerconst.LOG_ROTATE_COUNT,
	}
	log.InitWithConfig(&passLagerCfg)

}
