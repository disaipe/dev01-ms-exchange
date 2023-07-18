package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os/exec"
	"regexp"

	rpc "github.com/disaipe/dev01-rpc-base"
)

type MailBoxItem struct {
	Id             string
	DisplayName    string
	TotalItemSize  int64
	TotalItemCount int64
}

type GetMailBoxSizeRequest struct {
	rpc.Response

	Id string
}

type GetMailboxSizeResponse struct {
	rpc.ResultResponse

	Id     string
	Status bool
	Error  string
	Items  []MailBoxItem
}

var rpcAction = rpc.ActionFunction(func(rpcServer *rpc.Rpc, body io.ReadCloser, appAuth string) (rpc.Response, error) {
	var mailboxRequest GetMailBoxSizeRequest

	err := json.NewDecoder(body).Decode(&mailboxRequest)

	if err != nil {
		return nil, err
	}

	var resultStatus = true
	var resultMessage string

	if mailboxRequest.Id == "" {
		resultStatus = false
		resultMessage = "Id is requried"
	}

	if resultStatus {
		go func() {
			var resultStatus = true
			var resultError string
			var mailboxItems []MailBoxItem

			cmd := exec.Command("powershell.exe", "-nologo", "-noprofile", "-NonInteractive", "-ExecutionPolicy", "ByPass", "-OutputFormat", "Text", "-File", "./getMailboxSizes.ps1")
			out, err := cmd.CombinedOutput()

			if err != nil {
				log.Printf("Failed to start command: %v", err)

				resultStatus = false
				resultError = err.Error()
			} else {
				cleanOut := regexp.MustCompile(`([^\pL\pM\pN\pP\pS\s]|\r\n)`).ReplaceAllLiteralString(string(out), "")
				err = json.Unmarshal([]byte(cleanOut), &mailboxItems)

				if err != nil {
					log.Printf("Failed to make results: %v", err)

					resultStatus = false
					resultError = err.Error()
				}
			}

			resultData := &GetMailboxSizeResponse{
				Id:     mailboxRequest.Id,
				Status: resultStatus,
				Items:  mailboxItems,
				Error:  resultError,
			}

			rpcServer.SendResult(*resultData, appAuth)
		}()
	}

	requestAcceptedResponse := &rpc.ActionResponse{
		Status: resultStatus,
		Data:   resultMessage,
	}

	return requestAcceptedResponse, nil
})

func main() {
	flag.Parse()

	rpc.Config.SetServiceSettings(
		"dev01-exchange",
		"Dev01 Exchaing mailbox size monitor daemon",
		"The part of the Dev01 platform",
	)

	rpc.Config.SetAction("/get", &rpcAction)

	if rpc.Config.Serving() {
		rpcServer := &rpc.Rpc{}
		rpcServer.Run()
	}
}
