package main

import (
	"net/http"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/injustease/ce-midtrans/adapter"
	log "github.com/sirupsen/logrus"
)

var cfg = struct {
	K_SINK     string
	SERVER_KEY string
	PORT       string
}{}

func init() {
	cfg.K_SINK = os.Getenv("K_SINK")
	cfg.SERVER_KEY = os.Getenv("SERVER_KEY")
	cfg.PORT = os.Getenv("PORT")
}

func main() {
	client, err := cloudevents.NewClientHTTP(cloudevents.WithTarget(cfg.K_SINK))
	if err != nil {
		log.Fatal(err)
	}

	handler := adapter.NewHandler(client, cfg.SERVER_KEY)

	r := http.NewServeMux()
	r.HandleFunc("/payment", handler.PaymentNotification())
	r.HandleFunc("/recurring", handler.RecurringNotification())
	r.HandleFunc("/pay-account", handler.PayAccountNotification())

	srv := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: r,
	}

	log.Infof("Server is ready to handle request at %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
