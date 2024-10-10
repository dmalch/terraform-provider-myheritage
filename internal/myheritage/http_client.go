package myheritage

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func doRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response", "error", err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		slog.Error("Non-OK HTTP status", "status", res.StatusCode, "body", string(body))
		return nil, fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	return body, nil
}
