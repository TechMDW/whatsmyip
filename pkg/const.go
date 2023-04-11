package whatsmyip

import (
	"errors"
	"net/http"
	"regexp"
	"time"
)

var hc http.Client = http.Client{Timeout: time.Second * 5}

var (
	// ErrInvalidStatusCode is returned when the status code is not 200
	ErrInvalidStatusCode = errors.New("invalid status code")

	// ErrInvalidIp is returned when the ip is not valid
	ErrInvalidIp = errors.New("invalid ip")
)

var ipv4Regex = regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)

// TODO: Add IPv6 support
// var ipv6Regex = regexp.MustCompile(`\b(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}\b`)

// URL of websites thats return your ip
var Websites []string = []string{
	"https://checkip.amazonaws.com",
	"https://api.ipify.org/?format=plain",
	"https://myexternalip.com/raw",
	"https://ipecho.net/plain",
	"https://ident.me",
	"https://icanhazip.com",
	"https://www.trackip.net/ip",
	"https://ifconfig.me/ip",
	"https://ip.42.pl/raw",
	"https://wtfismyip.com/text",
	"https://ipinfo.io/ip",
	"https://ipecho.net/plain",
	"https://myip.dnsomatic.com",
}
