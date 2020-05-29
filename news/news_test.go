package news

import (
	"testing"

	"github.com/thrasher-corp/gocryptotrader/communications/slack"
)

func TestCheckOtherThings(t *testing.T) {
	err := CheckOtherThings()
	if err != nil {
		t.Error(err)
	}
}

func TestSendMessage(t *testing.T) {
	err := SendMessage("HELLoooooo")
	if err != nil {
		t.Error(err)
	}
}

func TestWriteFile(t *testing.T) {
	err := WriteFile([]string{"hello"}, "checklist.json")
	if err != nil {
		t.Error(err)
	}
	words, err := ReadFile("checklist.json")
	if words[len(words)-1] != "hello" {
		t.Fatal("hello wasn't successfully added to the word list")
	}
}

func TestReadFile(t *testing.T) {
	a, err := ReadFile("checklist.json")
	t.Log(a)
	if err != nil {
		t.Error(err)
	}
}

func TestHandleMessage(t *testing.T) {
	var s slack.Message
	s.Text = "!add hi,i,love,fomo"
	err := HandleMessage(s)
	if err != nil {
		t.Error(err)
	}
}
