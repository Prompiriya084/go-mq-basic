package unittest_services

import (
	"errors"
	"fmt"
	"testing"

	models "github.com/Prompiriya084/go-mq/Models"
	services "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Services"
	unittest_mq "github.com/Prompiriya084/go-mq/Producer/UnitTest/MockItem/MQ"
	unittest_repositories "github.com/Prompiriya084/go-mq/Producer/UnitTest/MockItem/Repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	orderId1 := uuid.New()
	orderId2 := uuid.New()
	mockReturnedData := []*models.Order{
		{ID: orderId1, ProductID: "Product-1"},
		{ID: orderId2, ProductID: "Product-2"},
	}
	testcase := []struct {
		description       string
		sendingFilter     *models.Order
		mockReturnedData  []*models.Order
		mockReturnedError error
		expection         int //Orders length
	}{
		{
			description:       "[OK]Can return orders without filters.",
			sendingFilter:     nil,
			mockReturnedData:  mockReturnedData,
			mockReturnedError: nil,
			expection:         2, //If orders found (lenght of returned data)
		},
		{
			description:       "[Error] Orders not found.",
			sendingFilter:     nil,
			mockReturnedData:  []*models.Order{},
			mockReturnedError: nil,
			expection:         0, // If orders not found (lenght of returned data)
		},
	}
	for _, tc := range testcase {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := &unittest_repositories.MockOrderRepo{
				MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{
					GetAllFn: func(filters *models.Order, preload []string) ([]*models.Order, error) {
						return tc.mockReturnedData, tc.mockReturnedError
					},
				},
			}
			mockMqProducer := &unittest_mq.MockProducer{}
			services := services.NewOrderService(mockRepo, mockMqProducer)
			response, err := services.GetAll(tc.sendingFilter, nil)

			assert.NoError(t, err)
			assert.Equal(t, tc.expection, len(response))
		})
	}
}

func TestGet(t *testing.T) {
	returnedOrderId := uuid.New()
	testcase := []struct {
		description       string
		mockReturnedData  *models.Order
		mockReturnedError error
		expection         *models.Order //Orders length
	}{
		{
			description:       "[OK]Can return order.",
			mockReturnedData:  &models.Order{ID: returnedOrderId, ProductID: "Product-1"},
			mockReturnedError: nil,
			expection:         &models.Order{ID: returnedOrderId, ProductID: "Product-1"}, //If orders found
		},
		{
			description:       "[Error] Orders not found.",
			mockReturnedData:  nil,
			mockReturnedError: nil,
			expection:         nil, // If orders not found
		},
	}
	for _, tc := range testcase {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := &unittest_repositories.MockOrderRepo{
				MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{
					GetFn: func(filters *models.Order, preload []string) (*models.Order, error) {
						return tc.mockReturnedData, tc.mockReturnedError
					},
				},
			}
			mockMqProducer := &unittest_mq.MockProducer{}
			services := services.NewOrderService(mockRepo, mockMqProducer)
			response, err := services.Get(nil, nil)
			fmt.Println("Response: ", response)
			fmt.Println("Expection: ", tc.expection)
			assert.NoError(t, err)
			assert.Equal(t, tc.expection, response)
		})
	}
}

func TestCreate(t *testing.T) {
	orderId := uuid.New()
	// orderId2 := uuid.New()
	mockOrder := &models.Order{ID: orderId, ProductID: "Product-1"}
	testcase := []struct {
		description       string
		Params            *models.Order
		mockReturnedData  *models.Order
		mockReturnedError error
		expectionErr      bool //Orders length
	}{
		{
			description:       "[OK]Create order.",
			Params:            mockOrder,
			mockReturnedData:  nil,
			mockReturnedError: nil,
			expectionErr:      false, //If there is no error
		},
		{
			description:       "[Error]The order already exists.",
			Params:            mockOrder,
			mockReturnedData:  mockOrder,
			mockReturnedError: nil,
			expectionErr:      true, // If orders not found
		},
		{
			description:       "[Error]Something went wrong.",
			Params:            mockOrder,
			mockReturnedData:  nil,
			mockReturnedError: errors.New("Something went wrong."),
			expectionErr:      true, // If orders not found
		},
	}
	for _, tc := range testcase {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := &unittest_repositories.MockOrderRepo{
				MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{
					GetFn: func(filters *models.Order, preload []string) (*models.Order, error) {
						return tc.mockReturnedData, tc.mockReturnedError
					},
				},
			}
			mockMqProducer := &unittest_mq.MockProducer{}
			services := services.NewOrderService(mockRepo, mockMqProducer)

			response := services.Create(tc.Params)

			var result bool
			if response == nil {
				result = false
			} else {
				result = true
			}

			assert.Equal(t, tc.expectionErr, result)

		})
	}
}
func TestGetUpdate(t *testing.T) {
	orderId := uuid.New()
	mockOrder := &models.Order{ID: orderId, ProductID: "Product-1"}
	testcase := []struct {
		description       string
		Params            *models.Order
		mockReturnedData  *models.Order
		mockReturnedError error
		expectionErr      bool //Orders length
	}{
		{
			description:       "[OK]Update order.",
			Params:            mockOrder,
			mockReturnedData:  mockOrder,
			mockReturnedError: nil,
			expectionErr:      false, //If there is no error
		},
		{
			description:       "[Error]The order does not exists.",
			Params:            mockOrder,
			mockReturnedData:  nil,
			mockReturnedError: nil,
			expectionErr:      true, // If orders not found
		},
		{
			description:       "[Error]Something went wrong.",
			Params:            mockOrder,
			mockReturnedData:  mockOrder,
			mockReturnedError: errors.New("Something went wrong."),
			expectionErr:      true, // If orders not found
		},
	}
	for _, tc := range testcase {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := &unittest_repositories.MockOrderRepo{
				MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{
					GetFn: func(filters *models.Order, preload []string) (*models.Order, error) {
						return tc.mockReturnedData, tc.mockReturnedError
					},
				},
			}
			mockMqProducer := &unittest_mq.MockProducer{}
			services := services.NewOrderService(mockRepo, mockMqProducer)

			response := services.Update(tc.Params)

			var result bool
			if response == nil {
				result = false
			} else {
				result = true
			}

			assert.Equal(t, tc.expectionErr, result)

		})
	}
}
func TestDelete(t *testing.T) {
	orderId := uuid.New()
	mockOrder := &models.Order{ID: orderId, ProductID: "Product-1"}
	testcase := []struct {
		description       string
		Params            *models.Order
		mockReturnedData  *models.Order
		mockReturnedError error
		expectionErr      bool //Orders length
	}{
		{
			description:       "[OK]Delete order.",
			Params:            mockOrder,
			mockReturnedData:  mockOrder,
			mockReturnedError: nil,
			expectionErr:      false, //If there is no error
		},
		{
			description:       "[Error]The order does not exists.",
			Params:            mockOrder,
			mockReturnedData:  nil,
			mockReturnedError: nil,
			expectionErr:      true, // If orders not found
		},
		{
			description:       "[Error]Something went wrong.",
			Params:            mockOrder,
			mockReturnedData:  mockOrder,
			mockReturnedError: errors.New("Something went wrong."),
			expectionErr:      true, // If orders not found
		},
	}
	for _, tc := range testcase {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := &unittest_repositories.MockOrderRepo{
				MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{
					GetFn: func(filters *models.Order, preload []string) (*models.Order, error) {
						return tc.mockReturnedData, tc.mockReturnedError
					},
				},
			}
			mockMqProducer := &unittest_mq.MockProducer{}
			services := services.NewOrderService(mockRepo, mockMqProducer)

			response := services.Delete(tc.Params)

			var result bool
			if response == nil {
				result = false
			} else {
				result = true
			}

			assert.Equal(t, tc.expectionErr, result)

		})
	}
}
