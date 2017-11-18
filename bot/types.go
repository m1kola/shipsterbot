package bot

// Available operations
import (
	"github.com/m1kola/telegram_shipsterbot/storage"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TODO: Rename
type BotApp struct {
	Bot     *tgbotapi.BotAPI
	Storage storage.StorageInterface
}
