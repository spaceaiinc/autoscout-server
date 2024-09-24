package utility

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OneSignalを使ったWebPush通知
func WebPush(
	appID,
	apiKey,
	firebaseID,
	title,
	contents,
	topic,
	redirectURL string,
) error {
	url := "https://onesignal.com/api/v1/notifications"

	p := fmt.Sprintf(
		"{\"app_id\":\"%s\",\"include_external_user_ids\":[\"%s\"],\"headings\":{\"en\":\"%s\"},\"contents\":{\"en\":\"%s\"},\"web_push_topic\":\"%s\",\"url\":\"%s\"}",
		appID, firebaseID, title, contents, topic, redirectURL,
	)

	payload := strings.NewReader(p)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", apiKey))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	return nil
}

func WebPushMulti(
	appID,
	apiKey,
	title,
	contents,
	topic,
	redirectURL string,
	firebaseIDList []string,
) error {
	url := "https://onesignal.com/api/v1/notifications"

	p := fmt.Sprintf(
		"{\"app_id\":\"%s\",\"include_external_user_ids\":[\"%s\"],\"headings\":{\"en\":\"%s\"},\"contents\":{\"en\":\"%s\"},\"web_push_topic\":\"%s\",\"url\":\"%s\"}",
		appID, strings.Join(firebaseIDList, ", "), title, contents, topic, redirectURL,
	)

	payload := strings.NewReader(p)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", apiKey))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	return nil
}
