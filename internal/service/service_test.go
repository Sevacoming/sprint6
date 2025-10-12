package service

import "testing"

func TestDetectAndConvert_Empty(t *testing.T) {
 _, err := DetectAndConvert("   ")
 if err == nil {
  t.Fatal("expected error on empty input")
 }
}
func TestDetectAndConvert_TextToMorse(t *testing.T) {
 got, err := DetectAndConvert("Привет")
 if err != nil || got == "" {
  t.Fatalf("unexpected: out=%q err=%v", got, err)
 }
}
func TestDetectAndConvert_MorseToText(t *testing.T) {
 in := ".--. .-. .. ...- . -"
 got, err := DetectAndConvert(in)
 if err != nil || got == "" {
  t.Fatalf("unexpected: out=%q err=%v", got, err)
 }
}
