package job

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"todof/internal/mailer"

	"github.com/nsevenpack/logger/v2/logger"
)

var onceWorker sync.Once

func StartWorker() {
	onceWorker.Do(func() {
		go func() {
			for {
				data, err := ClientRedis.RPop(context.Background(), "job:queue").Result()
				if err != nil {
					time.Sleep(2 * time.Second)
					continue
				}

				var job Job
				if err := json.Unmarshal([]byte(data), &job); err != nil {
					logger.Ef("Erreur de décodage job : %v", err)
					continue
				}

				if err := routeJob(job); err != nil {
					job.Retry++
					if job.Retry >= job.MaxRetry {
						saveFailedJob(job)
						continue
					}
					retryJob(job)
				}
			}
		}()
		logger.S("🛠️ Worker Redis démarré")
	})
}

func routeJob(job Job) error {
	switch job.Name {
	case "user:welcome":
		return handleSendWelcomeEmail(job)
	default:
		return fmt.Errorf("job inconnu : %s", job.Name)
	}
}

func handleSendWelcomeEmail(job Job) error {
	to := job.Payload["email"]
	username := job.Payload["username"]

	msg := fmt.Sprintf("Bienvenue %s 👋\n\nMerci pour ton inscription sur TodoF !", username)

	return mailer.SendMail(mailer.MailData{
		To:      to,
		Subject: "Bienvenue sur TodoF 🎉",
		Body:    msg,
	})
}

func retryJob(job Job) {
	data, _ := json.Marshal(job)
	ClientRedis.LPush(context.Background(), "job:queue", data)
}

func saveFailedJob(job Job) {
	key := fmt.Sprintf("job:failed:%s:%d", job.Name, time.Now().Unix())
	data, _ := json.Marshal(job)
	ClientRedis.Set(context.Background(), key, data, 0)
}
