package lspserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"sylmark/data"
	"sylmark/lsp"
	"sylmark/server"
	"time"

	"github.com/sourcegraph/jsonrpc2"
	"github.com/tj/go-naturaldate"
)

func (h *LangHandler) handleWorkspaceExecuteCommand(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {

	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.ExecuteCommandParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}
	// slog.Info(fmt.Sprintf("Execute the command %s %v ", params.Command, params.Arguments))

	switch params.Command {
	case "show":
		{
			arg := "today"
			// slog.Info(fmt.Sprintf("Arg is bef %s %d", arg, len(params.Arguments)))
			if len(params.Arguments) > 0 && len(params.Arguments[0]) > 0 {
				arg = params.Arguments[0]
			}
			// slog.Info("Arg is " + arg)

			date, err := naturaldate.Parse(arg, time.Now())
			if err != nil {
				slog.Error("Date is wrong")
				return nil, nil
			}
			fileName := h.Store.Config.GetDateString(date) + ".md"
			uri, err := h.Store.Config.GetFileURI(fileName, "journal/")
			h.ShowDocument(uri, false, lsp.Range{})
		}
	case "create":
		{
			if len(params.Arguments) > 0 && len(params.Arguments[0]) > 0 {
				filePath := params.Arguments[0]
				_, err := os.Create(filePath)
				if err != nil {
					slog.Error("Could not create error is " + err.Error())
				}
				uri, err := data.UriFromPath(filePath)
				if err != nil {
					slog.Error("Failed to get uri err " + err.Error())
					return nil, nil
				}
				// update data
				h.onDocCreated(uri, "")
			}
		}
	case "append":
		{
			if len(params.Arguments) > 1 && len(params.Arguments[0]) > 0 {
				filePath := params.Arguments[0]
				if heading := params.Arguments[1]; len(heading) > 0 {
					f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend)
					if err != nil {
						slog.Error("Could not open file error is " + err.Error())
						break
					}
					_, err = f.WriteString(heading)
					if err != nil {
						slog.Error("Could not write to file error is " + err.Error())
					}
					f.Close()
					h.loadDocData(filePath)
				}
			}
		}
	case "graph":
		{
			server := server.NewServer(&h.Store, &h.Store.Config, h.ShowDocument)
			go server.StartAndListen()
		}
	}

	return nil, nil
}
