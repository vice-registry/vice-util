package communication

import "github.com/vice-registry/vice-util/models"

// GetRuntimeStats returns statistics about import, export and store queue size and consumers
func GetRuntimeStats() *models.RuntimeStats {
	importQueue, err := rabbitmqCredentials.Channel.QueueInspect("import")
	importMessages := int64(0)
	importConsumers := int64(0)
	if err == nil {
		importMessages = int64(importQueue.Messages)
		importConsumers = int64(importQueue.Consumers)
	}

	exportQueue, err := rabbitmqCredentials.Channel.QueueInspect("export")
	exportMessages := int64(0)
	exportConsumers := int64(0)
	if err == nil {
		exportMessages = int64(exportQueue.Messages)
		exportConsumers = int64(exportQueue.Consumers)
	}

	storeQueue, err := rabbitmqCredentials.Channel.QueueInspect("store")
	storeMessages := int64(0)
	storeConsumers := int64(0)
	if err == nil {
		storeMessages = int64(storeQueue.Messages)
		storeConsumers = int64(storeQueue.Consumers)
	}

	stats := models.RuntimeStats{
		ExportsPending: exportMessages,
		ExportWorker:   exportConsumers,
		ImportsPending: importMessages,
		ImportWorker:   importConsumers,
		StorePending:   storeMessages,
		StoreWorker:    storeConsumers,
	}
	return &stats
}
