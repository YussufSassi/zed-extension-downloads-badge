package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Extension struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	DownloadCount int    `json:"download_count"`
}

type Response struct {
	Data []Extension `json:"data"`
}

type DownloadResponse struct {
	Downloads int `json:"download_count"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://zed.dev/api/extensions?max_schema_version=1&include_native=true", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := Response{}

	json.Unmarshal([]byte(bodyText), &response)

	extensionId := r.URL.Query().Get("extensionId")

	if extensionId == "" {
		fmt.Fprintln(w, "No extensionId provided")
	}

	for _, element := range response.Data {
		if element.Id == extensionId {
			resData := DownloadResponse{}
			resData.Downloads = element.DownloadCount
			json_msg, err := json.Marshal(resData)

			if err != nil {
				panic(err)
			}

			fmt.Fprintf(w, "%s", json_msg)

		}
	}

}
