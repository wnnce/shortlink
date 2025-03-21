package shortlink

import "shortlink/config"

var ConfigureReaderList = config.NewReaderList(KafkaConfigureReader)
