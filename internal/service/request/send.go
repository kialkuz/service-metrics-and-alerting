package request

import (
	"kialkuz/service-metrics-and-alerting/pkg/errorclassifier"
	"net/http"
)

func Send(client *http.Client, request *http.Request) error {
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &errorclassifier.HTTPError{
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}
