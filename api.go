package gtranslate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/text/language"

	"github.com/robertkrimen/otto"
)

var ttk, _ = otto.ToValue("0")

const (
	defaultNumberOfRetries = 2
)

func translate(text, from, to, googleHost string, withVerification bool, tries int, delay time.Duration, client *http.Client) (string, error) {
	if tries == 0 {
		tries = defaultNumberOfRetries
	}

	if withVerification {
		if _, err := language.Parse(from); err != nil && from != "auto" {
			log.Println("[WARNING], '" + from + "' is a invalid language, switching to 'auto'")
			from = "auto"
		}
		if _, err := language.Parse(to); err != nil {
			log.Println("[WARNING], '" + to + "' is a invalid language, switching to 'en'")
			to = "en"
		}
	}

	t, _ := otto.ToValue(text)

	urll := fmt.Sprintf("https://translate.%s/translate_a/single", googleHost)

	token := get(t, ttk, googleHost, client)

	data := map[string]string{
		"client": "gtx",
		"sl":     from,
		"tl":     to,
		"hl":     to,
		"ie":     "UTF-8",
		"oe":     "UTF-8",
		"otf":    "1",
		"ssel":   "0",
		"tsel":   "0",
		"kc":     "7",
		"q":      text,
	}

	u, err := url.Parse(urll)
	if err != nil {
		return "", nil
	}

	parameters := url.Values{}

	for k, v := range data {
		parameters.Add(k, v)
	}
	for _, v := range []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"} {
		parameters.Add("dt", v)
	}

	parameters.Add("tk", token)
	u.RawQuery = parameters.Encode()

	var r *http.Response

	for tries > 0 {
		if client != nil {
			r, err = client.Get(u.String())
		} else {
			r, err = http.Get(u.String())
		}

		if err != nil {
			if err == http.ErrHandlerTimeout {
				return "", errors.New("bad network, please check your internet connection")
			}
			return "", err
		}

		if r.StatusCode == http.StatusOK {
			break
		}

		if r.StatusCode == http.StatusForbidden {
			tries--
			time.Sleep(delay)
		}
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var resp []interface{}

	err = json.Unmarshal([]byte(raw), &resp)
	if err != nil {
		return "", err
	}

	responseText := ""
	for _, obj := range resp[0].([]interface{}) {
		if len(obj.([]interface{})) == 0 {
			break
		}

		t, ok := obj.([]interface{})[0].(string)
		if ok {
			responseText += t
		}
	}

	return responseText, nil
}
