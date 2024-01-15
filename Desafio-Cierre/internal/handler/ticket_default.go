package handler

import (
	"app/internal"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TicketDefaultHandler represents the default handler of the tickets
type TicketDefaultHandler struct {
	service internal.ServiceTicket
}

// NewTicketDefaultHandler creates a new default handler of the tickets
func NewTicketDefaultHandler(service internal.ServiceTicket) *TicketDefaultHandler {
	return &TicketDefaultHandler{
		service: service,
	}
}

type DefaultResponse struct {
	Message string
	Data    interface{}
}

func (t *TicketDefaultHandler) GetTickestByCountry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get destination from query string
		destination := chi.URLParam(r, "dest")
		if destination == "" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		total, err := t.service.GetTicketsAmountByDestinationCountry(destination)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DefaultResponse{
			Message: "OK",
			Data:    total,
		})
	}
}

func (t *TicketDefaultHandler) GetProportionByCountry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		destination := chi.URLParam(r, "dest")
		if destination == "" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		total, err := t.service.GetPercentageTicketsByDestinationCountry(destination)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DefaultResponse{
			Message: "OK",
			Data:    total,
		})
	}
}
