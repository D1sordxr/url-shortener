package params

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/notification/vo"
	"time"

	"github.com/google/uuid"
)

type UpdateNotificationStatusParams struct {
	ID       uuid.UUID  `json:"id"`
	Status   vo.Status  `json:"status"`
	Attempts int16      `json:"attempts"`
	SentAt   *time.Time `json:"sent_at"`
}
