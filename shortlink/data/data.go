package data

import "shortlink/config"

var ConfigureReaderList = config.NewReaderList(DbConfigureReader, RedisConfigureReader)
