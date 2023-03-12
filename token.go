package gtranslate

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func sM(a, ttk string) string {
	d := strings.Split(ttk, ".")
	b, err := strconv.ParseInt(d[0], 10, 64)
	if err != nil {
		return ""
	}
	var (
		e = []int{}
	)
	runes := []rune(a)

	for g := 0; g < len(runes); g++ {
		var l = int(runes[g])
		if l < 128 {
			e = append(e, l)
		} else if l < 2048 {
			e = append(e, (l>>6)|192)
		} else if (l&64512) == 55296 && g+1 < len(runes) && (int(runes[g+1])&64512) == 56320 {
			l = 65536 + ((l & 1023) << 10) + (int(runes[g+1]) & 1023)
			e = append(e, (l>>18)|240)
			e = append(e, ((l>>12)&63)|128)
			g++
		} else {
			e = append(e, (l>>12)|224)
			e = append(e, ((l>>6)&63)|128)
			e = append(e, (l&63)|128)
		}
	}
	tmp := int(b)
	for f1 := 0; f1 < len(e); f1++ {
		tmp += e[f1]
		tmp = xr(tmp, "+-a^+6")
	}
	tmp = xr(tmp, "+-3^+b+-f")

	if len(d) >= 2 {
		i, err := strconv.ParseInt(d[1], 10, 64)
		if err != nil {
			return ""
		}
		tmp ^= int(i)
	} else {
		tmp ^= 0
	}
	if tmp < 0 {
		tmp = (tmp & 2147483647) + 2147483648
	}
	tmp %= 1e6

	return fmt.Sprintf("&tk=%d.%d", tmp, tmp^int(b))
}

func xr(a int, b string) int {
	runes := []rune(b)
	for c := 0; c < len(runes)-2; c += 3 {
		d := runes[c+2]
		if 'a' <= d {
			d = d - 87
		} else {
			d = d - '0'
		}
		var result int
		if runes[c+1] == '+' {
			result = int(uint(a) >> d)
		} else {
			result = a << d
		}

		if runes[c] == '+' {
			a = (a + result) & 4294967295
		} else {
			a = a ^ result
		}
	}
	return a
}

func updateTTK(TTK string, googleHost string, client *http.Client) (string, error) {
	t := time.Now().UnixNano() / 3600000
	now := float64(t)
	ttk, err := strconv.ParseFloat(TTK, 64)
	if err != nil {
		return "0", err
	}

	if ttk == now {
		return TTK, nil
	}

	resp, err := client.Get(fmt.Sprintf("https://translate.%s", googleHost))
	if err != nil {
		return "0", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0", err
	}
	matches := regexp.MustCompile(`tkk:\s?'(.+?)'`).FindStringSubmatch(string(body))
	if len(matches) > 0 {
		return matches[0], nil
	}

	return TTK, nil
}

func get(text, ttk, googleHost string, client *http.Client) string {
	ttk, err := updateTTK(ttk, googleHost, client)
	if err != nil {
		return ""
	}
	tk := sM(text, ttk)

	if err != nil {
		return ""
	}
	sTk := strings.Replace(tk, "&tk=", "", -1)
	return sTk
}
