package main

import (
	"bufio"
	"fmt"
	"github.com/tonnytg/talk2me/pkg/web"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("FriendAI talk to OpenAI")
	fmt.Println("type some question!!!")
	fmt.Println("---------------------")

	for {
		// get question
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		// send question to OpenAI
		web2talk.SendWithArgs(text)
	}
}
