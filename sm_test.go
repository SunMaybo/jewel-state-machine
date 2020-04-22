package jewel_state_machine

import (
	"testing"
	"fmt"
	"log"
)

func TestStateMachine(t *testing.T) {

	sm := New("close",
		[]EventDesc{EventDesc{
			Name: "open",
			Src:  []string{"close"},
		}},
		map[CallBackType]Callback{
			CallBackEnterEvent: func(event *Event) error {
				fmt.Println("-------close------")
				return nil
			},
		},
	)

	err := sm.Transaction("open")
	if err != nil {
		log.Println(err)
	}
	log.Println(sm.Current())

}
