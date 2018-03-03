package telegram

import (
	"fmt"
	"strings"
)

const callbackQueryDataSeparator = ":"

func splitCallbackQueryData(data string) (command, payload string, err error) {
	dataPieces := strings.Split(data, callbackQueryDataSeparator)
	if len(dataPieces) != 2 || dataPieces[0] == "" || dataPieces[1] == "" {
		return "", "", fmt.Errorf(
			"Wrong data format for %#v: expected format is \"command:payload\"",
			data,
		)
	}

	return dataPieces[0], dataPieces[1], nil
}

func joinCallbackQueryData(command, payload string) string {
	return fmt.Sprintf("%s%s%s", command, callbackQueryDataSeparator, payload)
}
