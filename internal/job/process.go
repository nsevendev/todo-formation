package job

import (
	"context"
	"encoding/json"

	"github.com/nsevenpack/logger/v2/logger"
)

func ProcessJob(ctx context.Context, job Job) {
	go func() {
		data, err := json.Marshal(job)
		if err != nil {
			logger.Ef("❌ Erreur lors du Marshal du job: %v", err)
			return
		}

		if ClientRedis == nil {
			logger.Ef("❌ ClientRedis n'est pas initialisé")
			return
		}
		result := ClientRedis.LPush(ctx, "job:queue", data)
		if err := result.Err(); err != nil {
			logger.Ef("❌ Erreur lors de l'enregistrement du job '%s' : %v", job.Name, err)
			return
		}
		logger.Sf("📥 Job '%s' enregistré dans Redis avec payload: %v", job.Name, job.Payload)
	}()
}
