package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/go-playground/validator/v10"
	ipdata "github.com/ipdata/go"
	"github.com/mssola/useragent"
	"github.com/tomasen/realip"
)

func (app *application) formatValidationErr(err error) (map[string]string, error) {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{}, err
	}

	errs := map[string]string{}

	for _, e := range validationErrs {
		errs[strings.ToLower(e.Field())] = e.Translate(app.translator)
	}

	return errs, nil
}

func (app *application) newToken() string {
	return uniuri.NewLenChars(6, []byte("01234567890"))
}

func (app *application) userIP(r *http.Request) string {
	ip := realip.FromRequest(r)

	if ip == "127.0.0.1" || ip == "::1" || ip == "192.0.2.1" {
		ip = "86.44.17.109"
	}

	return ip
}

func (app *application) userLocation(ip string) (string, error) {
	client, err := ipdata.NewClient(app.config.IPKey)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := client.LookupWithContext(ctx, ip)
	if err != nil {
		return "", err
	}

	location := fmt.Sprintf("%s, %s", data.City, data.CountryName)

	return location, nil
}

func (app *application) getUserAgent(r *http.Request) string {
	ua := r.UserAgent()

	valid := strings.HasPrefix(ua, "Mozilla")

	if !valid {
		ua = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
	}

	return ua
}

func (app *application) getUserDeviceInfo(ua string) string {
	agent := useragent.New(ua)

	browser, _ := agent.Browser()

	operatingSystem := strings.Split(agent.OS(), " ")
	os := operatingSystem[0]

	deviceInfo := fmt.Sprintf("%s on %s", browser, os)
	return deviceInfo
}
