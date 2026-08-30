package dldruntime

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcHandler interface {
	Handle(context.Context, string, json.RawMessage) (any, bool, error)
}

func ServeStdio(ctx context.Context, in io.Reader, out io.Writer, h rpcHandler) error {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	encoder := json.NewEncoder(out)
	for scanner.Scan() {
		var req rpcRequest
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			_ = encoder.Encode(rpcResponse{JSONRPC: "2.0", Error: &rpcError{Code: -32700, Message: "parse error"}})
			continue
		}
		if req.JSONRPC != "2.0" || req.Method == "" || len(req.ID) == 0 {
			_ = encoder.Encode(rpcResponse{JSONRPC: "2.0", ID: req.ID, Error: &rpcError{Code: -32600, Message: "invalid request"}})
			continue
		}
		result, shutdown, err := h.Handle(ctx, req.Method, req.Params)
		resp := rpcResponse{JSONRPC: "2.0", ID: req.ID}
		if err != nil {
			var re *RPCMethodError
			if errors.As(err, &re) {
				resp.Error = &rpcError{Code: re.Code, Message: re.Message}
			} else {
				resp.Error = &rpcError{Code: -32603, Message: "internal error"}
			}
		} else {
			resp.Result = result
		}
		if err := encoder.Encode(resp); err != nil {
			return fmt.Errorf("write rpc response: %w", err)
		}
		if shutdown {
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read rpc request: %w", err)
	}
	return nil
}

type RPCMethodError struct {
	Code    int
	Message string
}

func (e *RPCMethodError) Error() string { return e.Message }

func invalidParams(message string) error   { return &RPCMethodError{Code: -32602, Message: message} }
func methodNotFound() error                { return &RPCMethodError{Code: -32601, Message: "method not found"} }
func operationFailed(message string) error { return &RPCMethodError{Code: -32000, Message: message} }

func decodeParams(raw json.RawMessage, dst any) error {
	if len(raw) == 0 || string(raw) == "null" {
		raw = []byte("{}")
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return invalidParams("invalid params")
	}
	return nil
}
