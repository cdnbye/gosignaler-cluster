package signalerconst

/**
  log config const
 */
const(

	WRITERS string = "file"
	LOGGER_LEVEL string = "WARN"
	LOGGER_DIR string = "/home/ubuntu/gopub/log/gosignaler.log"
	LOG_FORMAT_TEXT bool =true
	ROLLING_POLICY string = "size"
	LOG_ROTATE_DATE int =1
	LOG_ROTATE_SIZE int =1
	LOG_ROTATE_COUNT int =30
)
