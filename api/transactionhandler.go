package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vechain/thor/api/utils/httpx"
	"github.com/vechain/thor/api/utils/types"
	"github.com/vechain/thor/thor"
	"net/http"
)

//TransactionHTTPPathPrefix http path prefix
const TransactionHTTPPathPrefix = "/transactions"

//NewTransactionHTTPRouter add path to router
func NewTransactionHTTPRouter(router *mux.Router, ti *TransactionInterface) {
	sub := router.PathPrefix(TransactionHTTPPathPrefix).Subrouter()

	sub.Path("").Methods("POST").HandlerFunc(httpx.WrapHandlerFunc(ti.handleSendTransaction))
	sub.Path("/{id}").Methods("GET").HandlerFunc(httpx.WrapHandlerFunc(ti.handleGetTransactionByID))
	sub.Path("/{id}/receipts").Methods("GET").HandlerFunc(httpx.WrapHandlerFunc(ti.handleGetTransactionReceiptByID))
}

func (ti *TransactionInterface) handleGetTransactionByID(w http.ResponseWriter, req *http.Request) error {

	query := mux.Vars(req)

	if len(query) == 0 {
		return httpx.Error("No Params!", 400)
	}
	id, ok := query["id"]
	if !ok {
		return httpx.Error("Invalid Params!", 400)
	}
	txID, err := thor.ParseHash(id)
	if err != nil {
		return httpx.Error("Invalid hash!", 400)
	}
	tx, err := ti.GetTransactionByID(txID)
	if err != nil {
		return httpx.Error("Get transaction failed!", 400)
	}
	txData, err := json.Marshal(tx)
	if err != nil {
		return httpx.Error("System Error!", 400)
	}
	w.Write(txData)
	return nil
}

func (ti *TransactionInterface) handleGetTransactionReceiptByID(w http.ResponseWriter, req *http.Request) error {

	query := mux.Vars(req)
	if len(query) == 0 {
		return httpx.Error("No Params!", 400)
	}
	id, ok := query["id"]
	if !ok {
		return httpx.Error("Invalid Params!", 400)
	}
	txID, err := thor.ParseHash(id)
	if err != nil {
		return httpx.Error("Invalid hash!", 400)
	}
	receipt, err := ti.GetTransactionReceiptByID(txID)
	if err != nil {
		fmt.Println(err)
		return httpx.Error("Get transaction receipt failed!", 400)
	}
	receiptData, err := json.Marshal(receipt)
	if err != nil {
		return httpx.Error("System Error!", 400)
	}
	w.Write(receiptData)
	return nil
}

func (ti *TransactionInterface) handleSendTransaction(w http.ResponseWriter, req *http.Request) error {
	raw := []byte(req.FormValue("rawTransaction"))
	rawTransaction := new(types.RawTransaction)
	if err := json.Unmarshal(raw, &rawTransaction); err != nil {
		return err
	}
	txID, err := ti.SendRawTransaction(rawTransaction)
	if err != nil {
		return err
	}
	txIDMap := map[string]string{
		"txID": txID.String(),
	}
	data, err := json.Marshal(txIDMap)
	if err != nil {
		return httpx.Error(" System Error! ", 400)
	}
	w.Write(data)
	return nil
}
