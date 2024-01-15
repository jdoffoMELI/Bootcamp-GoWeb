package internal

import "context"

// RepositoryTicket represents the repository interface for tickets
type RepositoryTicket interface {
	// GetAll returns all the tickets
	Get(ctx context.Context) (t map[int]TicketAttributes, err error)

	// GetTicketByDestinationCountry returns the tickets filtered by destination country
	GetTicketByDestinationCountry(ctx context.Context, country string) (t map[int]TicketAttributes, err error)
}
