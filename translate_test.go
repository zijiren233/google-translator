package gtranslate

import (
	"testing"
)

func TestTranslate(t *testing.T) {
	result, err := Translate("Go是一种语言层面支持并发（Go最大的特色、天生支持并发）\n内置runtime、支持垃圾回收（GC）、静态强类型，快速编译的语言", "en", TranslationParams{From: "auto"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("V1:\nText: %s\nPronunciation: %s\nDetectedLang: %s\nConfidence: %f", result.Text, result.Pronunciation, result.Detected.Lang, result.Detected.Confidence)
}

func TestTranslateWithClienID(t *testing.T) {
	result, err := TranslateWithClienID("Go是一种语言层面支持并发（Go最大的特色、天生支持并发）\n内置runtime、支持垃圾回收（GC）、静态强类型，快速编译的语言", "en", TranslationWithClienIDParams{From: "auto"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("V2:\nText: %s\nPronunciation: %s\nDetectedLang: %s\nConfidence: %f", result.Text, result.Pronunciation, result.Detected.Lang, result.Detected.Confidence)
}
