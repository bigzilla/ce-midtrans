package adapter

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	log "github.com/sirupsen/logrus"
)

var ErrUnknownPayload = errors.New("unknown payload")

type Payload struct {
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`

	// #permata-virtual-account
	PermataVANumber string `json:"permata_va_number"`

	// com.midtrans.bank_transfer
	VANumbers []struct {
		VANumber string `json:"va_number"`
		Bank     string `json:"bank"`
	} `json:"va_numbers"`

	// com.midtrans.cstore
	Store string `json:"store"`
}

func (p *Payload) IsValid(serverKey string) bool {
	hash := sha512.Sum512([]byte(p.OrderID + p.StatusCode + p.GrossAmount + serverKey))
	base16hash := fmt.Sprintf("%x", hash[:])
	return p.SignatureKey == base16hash
}

type Handler struct {
	client    cloudevents.Client
	serverKey string
}

func NewHandler(client cloudevents.Client, serverKey string) *Handler {
	return &Handler{
		client:    client,
		serverKey: serverKey,
	}
}

func (h *Handler) PaymentNotification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("%v: %v", ErrUnknownPayload, err)
			http.Error(w, ErrUnknownPayload.Error(), http.StatusBadRequest)
			return
		}

		var payload Payload
		if err := json.Unmarshal(data, &payload); err != nil {
			log.Errorf("%v: %v", ErrUnknownPayload, err)
			http.Error(w, ErrUnknownPayload.Error(), http.StatusBadRequest)
			return
		}

		if !payload.IsValid(h.serverKey) {
			log.Warnf("%v: signature key is not valid", ErrUnknownPayload)
			http.Error(w, ErrUnknownPayload.Error(), http.StatusBadRequest)
			return
		}

		e, err := toCloudEvent(payload, data)
		if err != nil {
			log.Errorf("cannot create CloudEvent: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if res := h.client.Send(context.Background(), e); cloudevents.IsUndelivered(res) {
			log.Errorf("failed to send event to the sink: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h *Handler) RecurringNotification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (h *Handler) PayAccountNotification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}
