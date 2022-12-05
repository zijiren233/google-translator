package gtranslate

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/scanner"

	"golang.org/x/text/language"

	"github.com/robertkrimen/otto"
)

var ttk, _ = otto.ToValue("0")

func translate(text, from, to, googleHost string, withVerification bool, client *http.Client) (result Translated, err error) {
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
		return result, fmt.Errorf("return err, code: %d", r.StatusCode)
	}

	return parseRawTranslated(raw), nil
}

func parseRawTranslated(data []byte) (result Translated) {
	var s scanner.Scanner
	s.Init(bytes.NewReader(data))
	var (
		coord       = []int{-1}
		textBuilder strings.Builder
	)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch tok {
		case '[':
			coord[len(coord)-1]++
			coord = append(coord, -1)
		case ']':
			coord = coord[:len(coord)-1]
		case ',':
			// no-op
		default:
			tokText := s.TokenText()
			coord[len(coord)-1]++
			if len(coord) == 4 && coord[1] == 0 && coord[3] == 0 {
				if tokText != "null" {
					textBuilder.WriteString(tokText[1 : len(tokText)-1])
				}
			}
			if len(coord) == 4 && coord[0] == 0 && coord[1] == 0 && coord[3] == 3 {
				if tokText != "null" {
					result.Pronunciation = tokText[1 : len(tokText)-1]
				}
			}
			if len(coord) == 2 && coord[0] == 0 && coord[1] == 2 {
				result.Detected.Lang = tokText[1 : len(tokText)-1]
			}
			if len(coord) == 2 && coord[0] == 0 && coord[1] == 6 {
				result.Detected.Confidence, _ = strconv.ParseFloat(s.TokenText(), 64)
			}
		}
	}
	result.Text = textBuilder.String()

	return
}
