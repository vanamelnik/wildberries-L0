package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
)

// Validation regexs
var (
	phoneNumRegex *regexp.Regexp
	emailRegex    *regexp.Regexp
)

func init() {
	phoneNumRegex = regexp.MustCompile(`^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s\./0-9]*$`)
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
}

type (
	// Order struct  is used only to validate received orders, because, according to the assignment,
	// the only thing known about the organization of the data is that the data is static.
	Order struct {
		OrderUID          string    `json:"order_uid"`
		TrackNumber       string    `json:"track_number"`
		Entry             string    `json:"entry"`
		Delivery          Delivery  `json:"delivery"`
		Payment           Payment   `json:"payment"`
		Items             []Item    `json:"items"`
		Locale            string    `json:"locale"`
		InternalSignature string    `json:"internal_signature"`
		Shardkey          string    `json:"shardkey"`
		SmID              int       `json:"sm_id"`
		DateCreated       time.Time `json:"date_created"`
		OOFShard          string    `json:"oof_shard"`
	}

	Delivery struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		ZIP     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	}

	Payment struct {
		Transaction  string `json:"transaction"`
		RequestID    string `json:"request_id"`
		Currency     string `json:"currency"`
		Provider     string `json:"provider"`
		Amount       int    `json:"amount"`
		PaymentDt    uint64 `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	}

	Item struct {
		ChrtID      uint64 `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		RID         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int    `json:"total_price"`
		NmID        uint64 `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int    `json:"status"`
	}
)

func (o Order) Validate() error {
	var err error
	if o.OrderUID == "" {
		err = multierror.Append(err, errors.New("empty order UID"))
	}
	if o.Payment.Transaction == "" {
		err = multierror.Append(err, errors.New("empty payment transaction field"))
	}
	if !emailRegex.MatchString(o.Delivery.Email) {
		err = multierror.Append(err, errors.New("incorrect delivery email"))
	}
	if !phoneNumRegex.MatchString(o.Delivery.Phone) {
		err = multierror.Append(err, errors.New("incorrect delivery phone number"))
	}
	// TODO: add other validation for other fields

	return err
}
