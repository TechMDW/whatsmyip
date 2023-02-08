package whatsmyip

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type IpFetch struct {
	Ip      string
	Website string
	Error   error
	Http    HttpInfo
}

type HttpInfo struct {
	StatusCode int
	Status     string
	Method     string
	Body       []byte
	Headers    []string
	Error      error
}

// Generic function that can be used to fetch an IP from any website that returns an IP.
//
// This uses a regex to find the first IP in the body of the response.
func IpWebScraper(ctx context.Context, url string) (IpFetch, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return IpFetch{
			Website: url,
			Error:   err,
		}, err
	}

	req = req.WithContext(ctx)

	res, err := hc.Do(req)

	if err != nil {
		return IpFetch{
			Website: url,
			Http: HttpInfo{
				Error: err,
			},
		}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return IpFetch{
			Website: url,
			Error:   err,
		}, err
	}

	if res.StatusCode != 200 {
		return IpFetch{
			Website: url,
			Http: HttpInfo{
				Error:      ErrInvalidStatusCode,
				StatusCode: res.StatusCode,
				Status:     res.Status,
				Method:     req.Method,
				Body:       body,
				Headers:    GetArrayOfAllHeaders(res),
			},
		}, ErrInvalidStatusCode
	}

	ip := ipRegex.Find(body)

	if ip == nil {
		return IpFetch{
			Website: url,
			Error:   ErrInvalidIp,
			Http: HttpInfo{
				StatusCode: res.StatusCode,
				Status:     res.Status,
				Method:     req.Method,
				Body:       body,
				Headers:    GetArrayOfAllHeaders(res),
			},
		}, ErrInvalidIp
	}

	return IpFetch{
		Ip:      string(ip),
		Website: url,
		Error:   nil,
		Http: HttpInfo{
			StatusCode: res.StatusCode,
			Status:     res.Status,
			Method:     req.Method,
			Body:       body,
			Headers:    GetArrayOfAllHeaders(res),
		},
	}, nil
}

// Converts the http.Header map to an array of strings.
//
// Seperates the key and value with "::".
func GetArrayOfAllHeaders(res *http.Response) []string {
	var headers []string

	for k, v := range res.Header {
		data := fmt.Sprintf("%s::%s", k, v[0])
		headers = append(headers, data)
	}

	return headers
}
