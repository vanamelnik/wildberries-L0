package server

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vanamelnik/wildberries-L0/models"
	"github.com/vanamelnik/wildberries-L0/storage"
)

var (
	//go:embed templates/index.gohtml
	indexFile string
	//go:embed templates/order.gohtml
	orderFile string
)

type Server struct {
	http.Server

	mainTpl  *template.Template
	orderTpl *template.Template
	s        storage.Storage
	router   *mux.Router
}

func New(addr string, s storage.Storage) (*Server, error) {
	mainTpl, err := template.New("index").Parse(indexFile)
	if err != nil {
		return nil, err
	}
	orderTpl, err := template.New("order").Parse(orderFile)
	if err != nil {
		return nil, err
	}
	router := mux.NewRouter()
	server := Server{
		Server: http.Server{
			Addr:    addr,
			Handler: router,
		},
		mainTpl:  mainTpl,
		orderTpl: orderTpl,
		s:        s,
	}
	router.HandleFunc("/", server.indexHandler).Methods(http.MethodGet)
	router.HandleFunc("/{uid}", server.orderHandler).Methods(http.MethodGet)
	return &server, nil
}

func (srv *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	records, err := srv.s.GetAll()
	if err != nil {
		log.Printf("server: could not get records from the storage: %s", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}
	orders := fillOrders(records)
	if err := srv.mainTpl.Execute(w, orders); err != nil {
		log.Printf("server: could not execute the main template: %s", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}
}

func (srv *Server) orderHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		log.Println("server: unreachable: no key provided")
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}
	jsonOrder, err := srv.s.Get(uid)
	if err != nil {
		log.Printf("server: could not get order %s: %s", uid, err)
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, fmt.Sprintf("Order %s not found.", uid), http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}
	var order models.Order
	if err := json.Unmarshal([]byte(jsonOrder), &order); err != nil {
		log.Printf("server: unreachable: could not unmarshal json from the database: %s", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}
	if err := srv.orderTpl.Execute(w, order); err != nil {
		log.Printf("server: could not execute the order template: %s", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}
}

func fillOrders(records []storage.OrderDB) []models.Order {
	orders := make([]models.Order, 0, len(records))
	for _, rec := range records {
		var o models.Order
		if err := json.Unmarshal([]byte(rec.JSONOrder), &o); err != nil {
			panic("unreachable: " + err.Error())
		}
		orders = append(orders, o)
	}
	return orders
}
