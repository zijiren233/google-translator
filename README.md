# google-translater

```
package main

import (
	"context"
	"fmt"

	"github.com/zijiren233/google-translater"
	"golang.org/x/time/rate"
)

func translate(text string) string {
	translated, err := gtranslate.Translate(
		text,
		gtranslate.TranslationParams{
			From: "auto",
			To:   "en",
		},
	)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return translated
}

func main() {
	l := rate.NewLimiter(100, 100)
	for {
		l.Wait(context.Background())
		go func() { fmt.Println(translate("测试")) }()
	}
}
```
