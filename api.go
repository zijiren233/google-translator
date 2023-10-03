package gtranslator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/text/language"
)

func translate(text, to string, params *TranslationParams) (result *Translated, err error) {
	if params.Client == nil {
		params.Client = http.DefaultClient
	}

	if params.LangVerification {
		if _, err := language.Parse(params.From); err != nil {
			params.From = "auto"
		}
		if _, err := language.Parse(to); err != nil {
			to = "en"
		}
	}

	uP, err := url.Parse(fmt.Sprintf("https://translate.%s/translate_a/single", params.GoogleHost))
	if err != nil {
		return result, err
	}
	parameters := uP.Query()

	data := map[string]string{
		"client": "gtx",
		"sl":     params.From,
		"tl":     to,
		"q":      text,
	}

	for k, v := range data {
		parameters.Add(k, v)
	}
	parameters.Add("dt", "t")

	parameters.Add("dt", "rm")

	uP.RawQuery = parameters.Encode()

	r, err := params.Client.Get(uP.String())
	if err != nil {
		if err == http.ErrHandlerTimeout {
			return result, errors.New("bad network, please check your internet connection")
		}
		return result, err
	}
	defer r.Body.Close()

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

// id must between 1 and 5
func translateWithClienID(text, to string, params *TranslationWithClienIDParams) (result *Translated, err error) {
	if params.Client == nil {
		params.Client = http.DefaultClient
	}
	if params.ClientID < 1 || params.ClientID > 5 {
		params.ClientID = 5
	}

	if params.LangVerification {
		if _, err := language.Parse(params.From); err != nil {
			params.From = "auto"
		}
		if _, err := language.Parse(to); err != nil {
			to = "en"
		}
	}

	uP, err := url.Parse(fmt.Sprintf("https://clients%d.google.com/translate_a/single", params.ClientID))
	if err != nil {
		return result, err
	}
	parameters := uP.Query()

	data := map[string]string{
		"client": "dict-chrome-ex",
		"sl":     params.From,
		"tl":     to,
		"q":      text,
	}

	for k, v := range data {
		parameters.Add(k, v)
	}
	parameters.Add("dt", "t")

	parameters.Add("dt", "rm")
	uP.RawQuery = parameters.Encode()

	r, err := params.Client.Get(uP.String())
	if err != nil {
		if err == http.ErrHandlerTimeout {
			return result, errors.New("bad network, please check your internet connection")
		}
		return result, err
	}
	defer r.Body.Close()

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
	for _, obj := range d[0].([]interface{}) {
		lObg := len(obj.([]interface{}))
		if lObg == 0 {
			break
		}

		if t, ok := obj.([]interface{})[0].(string); ok {
			resp.Text += t
		} else if t, ok := obj.([]interface{})[lObg-1].(string); ok {
			resp.Pronunciation = t
			break
		}

	}

	resp.Detected.Lang = d[2].(string)
	resp.Detected.Confidence = d[6].(float64)

	return resp, nil
}
