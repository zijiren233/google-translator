package gtranslate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/text/language"

	"github.com/robertkrimen/otto"
)

var ttk, _ = otto.ToValue("0")

func translate(text, from, to, googleHost string, withVerification bool, client *http.Client) (result *Translated, err error) {
	if client == nil {
		client = http.DefaultClient
	}

	if withVerification {
		if _, err := language.Parse(from); err != nil && from != "auto" {
			from = "auto"
		}
		if _, err := language.Parse(to); err != nil {
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
		return result, err
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

	if client != nil {
		r, err = client.Get(u.String())
	} else {
		r, err = http.Get(u.String())
	}

	if err != nil {
		if err == http.ErrHandlerTimeout {
			return result, errors.New("bad network, please check your internet connection")
		}
		return result, err
	}

	if r.StatusCode != http.StatusOK {
		return result, fmt.Errorf("return err, code: %d", r.StatusCode)
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return result, err
	}

	if http.DetectContentType(raw) != `text/plain; charset=utf-8` {
		return result, fmt.Errorf("return err, code: %d\n raw: %s", r.StatusCode, raw)
	}

	return parseRawTranslated(raw)
}

func parseRawTranslated(data []byte) (*Translated, error) {

	var d []interface{}

	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	resp := &Translated{}
	l := len(d[0].([]interface{}))
	for k, obj := range d[0].([]interface{}) {
		lObg := len(obj.([]interface{}))
		if lObg == 0 {
			break
		}

		if t, ok := obj.([]interface{})[0].(string); ok {
			resp.Text += t
		} else if t, ok := obj.([]interface{})[lObg-1].(string); ok {
			if k == l-1 {
				resp.Pronunciation = t
				break
			}
		}

	}

	resp.Detected.Lang = d[2].(string)
	resp.Detected.Confidence = d[6].(float64)

	return resp, nil
}
