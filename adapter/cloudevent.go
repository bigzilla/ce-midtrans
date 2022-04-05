package adapter

import (
	"errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var (
	ErrUnknownType   = errors.New("unknown type")
	ErrUnknownSource = errors.New("unknown source")
)

const (
	typePrefix = "com.midtrans"

	creditCard    = "credit_card"
	bankTransfer  = "bank_transfer"
	echannel      = "echannel"
	bcaKlikpay    = "bca_klikpay"
	bcaKlikbca    = "bca_klikbca"
	briEpay       = "bri_epay"
	cimbClicks    = "cimb_clicks"
	danamonOnline = "danamon_online"
	qris          = "qris"
	gopay         = "gopay"
	shopeepay     = "shopeepay"
	cstore        = "cstore"
	akulaku       = "akulaku"

	bankPermata = "permata"
	bankBCA     = "bca"
	bankBNI     = "bni"
	bankBRI     = "bri"

	storeIndomaret = "indomaret"
	storeAlfamart  = "alfamart"
)

func toCloudEvent(p Payload, data []byte) (cloudevents.Event, error) {
	e := cloudevents.NewEvent()

	if _, ok := types[p.PaymentType]; !ok {
		return e, ErrUnknownType
	}

	var source string
	switch p.PaymentType {
	case creditCard, echannel, bcaKlikpay, bcaKlikbca, briEpay, cimbClicks, danamonOnline, qris, gopay, shopeepay, akulaku:
		source = sources[p.PaymentType]
	case bankTransfer:
		var bank string
		if p.PermataVANumber != "" {
			bank = bankPermata
		}
		if len(p.VANumbers) != 0 {
			bank = p.VANumbers[0].Bank
		}

		s, ok := sources[bankTransfer+bank]
		if !ok {
			return e, ErrUnknownSource
		}
		source = s
	case cstore:
		s, ok := sources[p.PaymentType+p.Store]
		if !ok {
			return e, ErrUnknownSource
		}
		source = s
	default:
		return e, ErrUnknownSource
	}

	e.SetType(fmt.Sprintf("%s.%s", typePrefix, p.PaymentType))
	e.SetSource(source)
	e.SetExtension("transactionstatus", p.TransactionStatus)
	e.SetExtension("fraudstatus", p.FraudStatus)

	if err := e.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return e, err
	}

	return e, nil
}

var types = map[string]bool{
	creditCard:    true,
	bankTransfer:  true,
	echannel:      true,
	bcaKlikpay:    true,
	bcaKlikbca:    true,
	briEpay:       true,
	cimbClicks:    true,
	danamonOnline: true,
	qris:          true,
	gopay:         true,
	shopeepay:     true,
	cstore:        true,
	akulaku:       true,
}

var sources = map[string]string{
	creditCard:                 "#card-payment",
	bankTransfer + bankPermata: "#permata-virtual-account",
	bankTransfer + bankBCA:     "#bca-virtual-account",
	bankTransfer + bankBNI:     "#bni-virtual-account",
	bankTransfer + bankBRI:     "#bri-virtual-account",
	echannel:                   "#mandiri-bill-payment",
	bcaKlikpay:                 "#bca-klikpay",
	bcaKlikbca:                 "#klikbca",
	briEpay:                    "#brimo",
	cimbClicks:                 "#cimb-clicks",
	danamonOnline:              "#danamon-online-banking",
	qris:                       "#qris",
	gopay:                      "#gopay",
	shopeepay:                  "#shopeepay",
	cstore + storeIndomaret:    "#indomaret",
	cstore + storeAlfamart:     "#alfamart",
	akulaku:                    "#akulaku",
}
