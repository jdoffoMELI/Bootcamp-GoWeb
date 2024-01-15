package service

import (
	"app/internal"
	"context"
)

// ServiceTicketDefault represents the default service of the tickets
type ServiceTicketDefault struct {
	// rp represents the repository of the tickets
	rp internal.RepositoryTicket
}

// NewServiceTicketDefault creates a new default service of the tickets
func NewServiceTicketDefault(rp internal.RepositoryTicket) *ServiceTicketDefault {
	return &ServiceTicketDefault{
		rp: rp,
	}
}

// GetTotalTickets returns the total number of tickets
func (s *ServiceTicketDefault) GetTotalAmountTickets() (total int, err error) {
	tickets, err := s.rp.Get(context.Background())
	if err != nil {
		return
	}
	total = len(tickets)
	return
}

// GetTicketsAmountByDestinationCountry returns the total number of tickets by destination country
func (s *ServiceTicketDefault) GetTicketsAmountByDestinationCountry(destination string) (total int, err error) {
	tickets, err := s.rp.GetTicketByDestinationCountry(context.Background(), destination)
	if err != nil {
		return
	}
	total = len(tickets)
	return
}

// GetPercentageTicketsByDestinationCountry returns the percentage of tickets by destination country
func (s *ServiceTicketDefault) GetPercentageTicketsByDestinationCountry(country string) (percentage float64, err error) {
	total, err := s.GetTotalAmountTickets()
	if err != nil {
		return
	}

	ticketsByCountry, err := s.GetTicketsAmountByDestinationCountry(country)
	if err != nil {
		return
	}

	return float64(ticketsByCountry) / float64(total) * 100.0, nil

}
