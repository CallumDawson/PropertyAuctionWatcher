package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const PropertyUrl = "https://www.iamsold.co.uk/property/..."
const DiscordWebhook = "https://discord.com/api/webhooks/..."

func getBidCount() int {
	resp, err := http.Get(PropertyUrl)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err == nil {
				bodyStr := string(bodyBytes)
				bidStart := strings.Index(bodyStr, "bid_count")
				if bidStart != -1 {
					bidEnd := strings.Index(bodyStr[bidStart:], "<")
					if bidEnd != -1 {
						i, _ := strconv.Atoi(bodyStr[bidStart+11 : bidStart+bidEnd])
						return i
					}
				}
			}
		}
	}

	return -1
}

func sendDiscord(text string) {
	buf := []byte(`{"content": "` + text + `"}`)
	_, _ = http.Post(DiscordWebhook, "application/json", bytes.NewBuffer(buf))
}

func main() {
	sleepDuration := time.Minute * 15
	bids := getBidCount()
	for {
		time.Sleep(sleepDuration)
		currentBids := getBidCount()
		if currentBids > bids {
			sendDiscord("New bids detected! Total: " + strconv.Itoa(currentBids))
			bids = currentBids
		}
		fmt.Println("Checked: " + time.Now().Format(time.RFC822))
	}
}
