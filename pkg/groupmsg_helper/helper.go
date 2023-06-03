package groupmsghelper

import (
	"github.com/micmonay/keybd_event"
	"golang.design/x/clipboard"
)

func SendText(msg string) error {
	clipboard.Write(clipboard.FmtText, []byte(msg))
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}
	kb.SetKeys(keybd_event.VK_V)
	kb.HasCTRL(true)
	err = kb.Launching()
	if err != nil {
		return err
	}
	kb.SetKeys(keybd_event.VK_ENTER)
	kb.HasCTRL(false)
	err = kb.Launching()
	if err != nil {
		return err
	}
	return nil
}
