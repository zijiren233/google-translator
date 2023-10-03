package gtranslator

import (
	"net/http"
	"time"
)

type Translated struct {
	Detected      Detected `json:"detected"`
	Text          string   `json:"text"`          // translated text
	Pronunciation string   `json:"pronunciation"` // pronunciation of translated text
}

// Detected represents language detection result
type Detected struct {
	Lang       string  `json:"lang"`       // detected language
	Confidence float64 `json:"confidence"` // the confidence of detection result (0.00 to 1.00)
}

var googleHostList = []string{"google.com"}

// come from http://m.news.xixik.com/content/25df0c00ac300bdc/
var googleHostAsia = []string{"google.com.hk", "google.mn", "google.co.kr", "google.co.jp", "google.com.vn", "google.la", "google.com.kh", "google.co.th", "google.com.my", "google.com.sg", "google.com.bn", "google.com.ph", "google.co.id", "google.kz", "google.kg", "google.com.tj", "google.co.uz", "google.tm", "google.com.af", "google.com.pk", "google.com.np", "google.co.in", "google.com.bd", "google.lk", "google.mv", "google.com.kw", "google.com.sa", "google.com.bh", "google.ae", "google.com.om", "google.jo", "google.co.il", "google.com.lb", "google.com.tr", "google.az", "google.am", "google.co.ls"}
var googleHostEurope = []string{"google.is", "google.dk", "google.no", "google.se", "google.fi", "google.ee", "google.lv", "google.lt", "google.ie", "google.co.uk", "google.gg", "google.je", "google.im", "google.fr", "google.nl", "google.be", "google.lu", "google.de", "google.at", "google.ch", "google.li", "google.pt", "google.es", "google.com.gi", "google.ad", "google.it", "google.com.mt", "google.sm", "google.gr", "google.ru", "google.com.by", "google.com.ua", "google.pl", "google.cz", "google.sk", "google.hu", "google.si", "google.hr", "google.ba", "google.me", "google.rs", "google.mk", "google.bg", "google.ro", "google.md"}
var googleHostAfrica = []string{"google.com.eg", "google.com.ly", "google.dz", "google.co.ma", "google.sn", "google.gm", "google.ml", "google.bf", "google.com.sl", "google.ci", "google.com.gh", "google.tg", "google.bj", "google.ne", "google.com.ng", "google.sh", "google.cm", "google.td", "google.cf", "google.ga", "google.cg", "google.cd", "google.it.ao", "google.com.et", "google.dj", "google.co.ke", "google.co.ug", "google.co.tz", "google.rw", "google.bi", "google.mw", "google.co.mz", "google.mg", "google.sc", "google.mu", "google.co.zm", "google.co.zw", "google.co.bw", "google.com.na", "google.co.za"}
var googleHostAtlantic = []string{"google.com.au", "google.com.nf", "google.co.nz", "google.com.sb", "google.com.fj", "google.fm", "google.ki", "google.nr", "google.tk", "google.ws", "google.as", "google.to", "google.nu", "google.co.ck", "google.com.do", "google.tt", "google.com.co", "google.com.ec", "google.co.ve", "google.gy", "google.com.pe", "google.com.bo", "google.com.py", "google.com.br", "google.com.uy", "google.com.ar", "google.cl"}
var googleHostAmerica = []string{"google.gl", "google.com.mx", "google.com.gt", "google.com.bz", "google.com.sv", "google.hn", "google.com.ni", "google.co.cr", "google.com.pa", "google.bs", "google.com.cu", "google.com.jm", "google.ht"}
var sw chan string

func init() {
	googleHostList = append(googleHostList, googleHostAsia...)
	googleHostList = append(googleHostList, googleHostEurope...)
	googleHostList = append(googleHostList, googleHostAfrica...)
	googleHostList = append(googleHostList, googleHostAtlantic...)
	googleHostList = append(googleHostList, googleHostAmerica...)
	sw = make(chan string, len(googleHostList))
	for _, v := range googleHostList {
		sw <- v
	}
}

// TranslationParams is a util struct to pass as parameter to indicate how to translate
type TranslationParams struct {
	From             string
	Retry            int
	RetryDelay       time.Duration
	LangVerification bool
	GoogleHost       string
	Client           *http.Client
}

type TranslationWithClienIDParams struct {
	From             string
	Retry            int
	RetryDelay       time.Duration
	LangVerification bool
	ClientID         int
	Client           *http.Client
}

const (
	defaultNumberOfRetries = 2
)

// TranslateWithParams translate a text with simple params as string
func Translate(text, To string, params TranslationParams) (translated *Translated, err error) {
	if params.Retry == 0 {
		params.Retry = defaultNumberOfRetries
	}

	for params.Retry > 0 {
		if params.GoogleHost == "" {
			params.GoogleHost = <-sw
			defer func() { sw <- params.GoogleHost }()
		}
		translated, err = translate(text, To, &params)
		if err == nil {
			return
		}
		params.Retry--
		time.Sleep(params.RetryDelay)
	}

	return
}

func TranslateWithClienID(text, To string, params TranslationWithClienIDParams) (translated *Translated, err error) {
	if params.Retry == 0 {
		params.Retry = defaultNumberOfRetries
	}

	for params.Retry > 0 {
		translated, err = translateWithClienID(text, To, &params)
		if err == nil {
			return
		}
		params.Retry--
		time.Sleep(params.RetryDelay)
	}

	return
}
