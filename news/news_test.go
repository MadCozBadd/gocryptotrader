package news

import "testing"

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
