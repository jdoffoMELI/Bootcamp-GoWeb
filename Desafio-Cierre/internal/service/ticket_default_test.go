package service_test

import (
	"app/internal"
	"app/internal/repository"
	"app/internal/service"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for ServiceTicketDefault.GetTotalAmountTickets
func TestServiceTicketDefault_GetTotalAmountTickets(t *testing.T) {
	t.Run("success to get total tickets", func(t *testing.T) {
		// arrange
		// - repository: mock
		rp := repository.NewRepositoryTicketMock()
		// - repository: set-up
		rp.FuncGet = func() (t map[int]internal.TicketAttributes, err error) {
			t = map[int]internal.TicketAttributes{
				1: {
					Name:    "John",
					Email:   "johndoe@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
			}
			return
		}

		// - service
		sv := service.NewServiceTicketDefault(rp)

		// act
		total, err := sv.GetTotalAmountTickets()

		// assert
		expectedTotal := 1
		require.NoError(t, err)
		require.Equal(t, expectedTotal, total)
	})
}

func TestServiceTicketDefault_GetTicketsAmountByDestinationCountry(t *testing.T) {
	t.Run("success to get total tickets by destination country", func(t *testing.T) {
		// arrange
		// - repository: mock
		rp := repository.NewRepositoryTicketMock()
		rp.FuncGetTicketsByDestinationCountry = func(destination string) (t map[int]internal.TicketAttributes, err error) {
			t = map[int]internal.TicketAttributes{
				1: {
					Name:    "John",
					Email:   "johndoe@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
				2: {
					Name:    "Jane",
					Email:   "jane@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
			}

			if destination != "USA" {
				return nil, nil
			}

			return t, nil
		}

		// - service
		sv := service.NewServiceTicketDefault(rp)

		// act
		total, err := sv.GetTicketsAmountByDestinationCountry("USA")

		// assert
		expectedTotal := 2
		require.NoError(t, err)
		require.Equal(t, expectedTotal, total)
	})

	t.Run("Should return 0", func(t *testing.T) {
		// arrange
		// - repository: mock
		rp := repository.NewRepositoryTicketMock()
		rp.FuncGetTicketsByDestinationCountry = func(destination string) (t map[int]internal.TicketAttributes, err error) {
			t = map[int]internal.TicketAttributes{
				1: {
					Name:    "John",
					Email:   "johndoe@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
				2: {
					Name:    "Jane",
					Email:   "jane@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
			}

			if destination != "USA" {
				return nil, nil
			}

			return t, nil
		}

		// - service
		sv := service.NewServiceTicketDefault(rp)

		// act
		total, err := sv.GetTicketsAmountByDestinationCountry("FOOCOUNTRY")

		// assert
		expectedTotal := 0
		require.NoError(t, err)
		require.Equal(t, expectedTotal, total)
	})

}

func TestServiceDefault_GetPercentageTicketsByDestinationCountry(t *testing.T) {
	t.Run("success to get total tickets by destination country", func(t *testing.T) {
		// arrange
		// - repository: mock
		rp := repository.NewRepositoryTicketMock()
		rp.FuncGetTicketsByDestinationCountry = func(destination string) (t map[int]internal.TicketAttributes, err error) {
			t = map[int]internal.TicketAttributes{
				1: {
					Name:    "John",
					Email:   "johndoe@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
				2: {
					Name:    "Jane",
					Email:   "jane@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
			}

			if destination != "USA" {
				return nil, nil
			}

			return t, nil
		}
		rp.FuncGet = func() (t map[int]internal.TicketAttributes, err error) {
			t = map[int]internal.TicketAttributes{
				1: {
					Name:    "John",
					Email:   "johndoe@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
				2: {
					Name:    "Jane",
					Email:   "jane@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
			}
			return
		}

		// - service
		sv := service.NewServiceTicketDefault(rp)

		// act
		total, err := sv.GetPercentageTicketsByDestinationCountry("USA")

		// assert
		expectedTotal := 100.0
		require.NoError(t, err)
		require.Equal(t, expectedTotal, total)
	})

	t.Run("Should return 0", func(t *testing.T) {
		// arrange
		// - repository: mock
		rp := repository.NewRepositoryTicketMock()
		rp.FuncGetTicketsByDestinationCountry = func(destination string) (t map[int]internal.TicketAttributes, err error) {
			t = map[int]internal.TicketAttributes{
				1: {
					Name:    "John",
					Email:   "johndoe@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
				2: {
					Name:    "Jane",
					Email:   "jane@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
			}

			if destination != "USA" {
				return nil, nil
			}

			return t, nil
		}
		rp.FuncGet = func() (t map[int]internal.TicketAttributes, err error) {
			t = map[int]internal.TicketAttributes{
				1: {
					Name:    "John",
					Email:   "johndoe@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
				2: {
					Name:    "Jane",
					Email:   "jane@gmail.com",
					Country: "USA",
					Hour:    "10:00",
					Price:   100,
				},
			}
			return
		}

		// - service
		sv := service.NewServiceTicketDefault(rp)

		// act
		total, err := sv.GetPercentageTicketsByDestinationCountry("FAKE COUNTRY")

		// assert
		expectedTotal := 0.0
		require.NoError(t, err)
		require.Equal(t, expectedTotal, total)
	})

}
