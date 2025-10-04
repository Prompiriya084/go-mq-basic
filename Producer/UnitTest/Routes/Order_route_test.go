package unittest_routes_order

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	eventbus "github.com/Prompiriya084/go-mq/EventBus"
	models "github.com/Prompiriya084/go-mq/Models"
	ports_repositories "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Ports/Repositories"
	services "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Services"
	unittest_eventbus "github.com/Prompiriya084/go-mq/Producer/UnitTest/MockItem/MQ"
	unittest_repositories "github.com/Prompiriya084/go-mq/Producer/UnitTest/MockItem/Repositories"
	routes "github.com/Prompiriya084/go-mq/Producer/Web/Routes"
	"github.com/stretchr/testify/assert"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// This ensures a clean slate and avoids state leakage between tests.
func createTestApp(mockRepo ports_repositories.OrderRepository, mockEventbus eventbus.EventBus[models.Order]) *fiber.App {
	app := fiber.New()
	service := services.NewOrderService(mockRepo, mockEventbus)
	routes.OrderSetupRouter(app, service) // Set up the routes for testing
	return app
}

func TestCreate(t *testing.T) {
	// ✅ Add this line before setting up routes
	// os.Setenv("RABBITMQ_URL", "amqp://testMQ:password@localhost:5672/")
	mockRepo := &unittest_repositories.MockOrderRepo{
		MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{},
	}
	mockEventbus := &unittest_eventbus.MockEventbus[models.Order]{}
	app := createTestApp(mockRepo, mockEventbus)
	// app.Post("/orders",  ) // Replace with your actual handler
	testcase := []struct {
		description string
		requestBody models.Order
		expection   int
	}{
		{
			description: "Valid input",
			requestBody: models.Order{
				ID:        uuid.Nil,
				ProductID: "testProducID",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeleteAt:  time.Time{},
			},
			expection: fiber.StatusOK,
		},
		{
			description: "Invalid Input",
			requestBody: models.Order{
				ID:        uuid.Nil,
				ProductID: "",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeleteAt:  time.Time{},
			},
			expection: fiber.StatusBadRequest,
		},
	}
	for _, tc := range testcase {
		t.Run(tc.description, func(t *testing.T) {
			reqbody, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err, "Failed to marshal request body") // Check marshal error

			req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(reqbody))
			req.Header.Set("Content-type", "application/json")

			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expection, resp.StatusCode)

		})
	}
}
func TestGetAll(t *testing.T) {
	// ✅ Add this line before setting up routes
	// os.Setenv("RABBITMQ_URL", "amqp://testMQ:password@localhost:5672/")

	// app.Post("/orders",  ) // Replace with your actual handler
	testcase := []struct {
		description       string
		mockReturnedData  []*models.Order
		mockReturnedError error
		expection         int
	}{
		{
			description: "Returns 200 OK with orders",
			mockReturnedData: []*models.Order{
				{ID: uuid.New(), ProductID: "product-1"},
				{ID: uuid.New(), ProductID: "product-2"},
			},
			mockReturnedError: nil,
			expection:         fiber.StatusOK,
		},
		{
			description:       "Returns 404 when no orders found",
			mockReturnedData:  []*models.Order{},
			mockReturnedError: nil,
			expection:         fiber.StatusNotFound,
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
			mockEventbus := &unittest_eventbus.MockEventbus[models.Order]{}
			app := createTestApp(mockRepo, mockEventbus)

			req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
			req.Header.Set("Content-type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expection, resp.StatusCode)

		})
	}
}
func TestGet(t *testing.T) {
	// ✅ Add this line before setting up routes
	// os.Setenv("RABBITMQ_URL", "amqp://testMQ:password@localhost:5672/")

	// app.Post("/orders",  ) // Replace with your actual handler
	orderId1 := uuid.New()
	// orderId2 := uuid.New()
	testcase := []struct {
		description string
		queryString string
		mockData    *models.Order
		mockError   error
		expection   int
	}{
		{
			description: "Returns 200 OK with orders",
			queryString: orderId1.String(),
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-1",
			},
			mockError: nil,
			expection: fiber.StatusOK,
		},
		{
			description: "Returns 404 when orders not found",
			queryString: orderId1.String(),
			mockData:    nil,
			mockError:   nil,
			expection:   fiber.StatusNotFound,
		},
		{
			description: "Returns 400 Bad request when sending invalid params.",
			queryString: "%20",
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-1",
			},
			mockError: nil,
			expection: fiber.StatusBadRequest,
		},
	}
	for _, tc := range testcase {
		t.Run(tc.description, func(t *testing.T) {
			// reqbody, err := json.Marshal(tc.reqBody)
			mockRepo := &unittest_repositories.MockOrderRepo{
				MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{
					GetFn: func(filters *models.Order, preload []string) (*models.Order, error) {
						return tc.mockData, tc.mockError
					},
				},
			}
			mockEventbus := &unittest_eventbus.MockEventbus[models.Order]{}
			app := createTestApp(mockRepo, mockEventbus)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/orders/%s", tc.queryString), nil)
			req.Header.Set("Content-type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expection, resp.StatusCode)

		})
	}
}

func TestUpdate(t *testing.T) {
	// ✅ Add this line before setting up routes
	// os.Setenv("RABBITMQ_URL", "amqp://testMQ:password@localhost:5672/")

	// app.Post("/orders",  ) // Replace with your actual handler
	orderId1 := uuid.New()
	// orderId2 := uuid.New()
	testcase := []struct {
		description string
		queryString string
		requestBody models.Order
		mockData    *models.Order
		mockError   error
		expection   int
	}{
		{
			description: "Returns 200 OK with orders",
			queryString: orderId1.String(),
			requestBody: models.Order{
				ProductID: "product-5",
			},
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-5",
			},
			mockError: nil,
			expection: fiber.StatusOK,
		},
		{
			description: "Returns 400 when not sending query string",
			queryString: "%20",
			requestBody: models.Order{
				ID:        orderId1,
				ProductID: "",
			},
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-1",
			},
			mockError: nil,
			expection: fiber.StatusBadRequest,
		},
		{
			description: "Returns 400 when sending request body incorrect",
			queryString: orderId1.String(),
			requestBody: models.Order{
				ID: orderId1,
			},
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-1",
			},
			mockError: nil,
			expection: fiber.StatusBadRequest,
		},
		{
			description: "Returns 400 when not sending anything.",
			queryString: "%20",
			requestBody: models.Order{},
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-1",
			},
			mockError: nil,
			expection: fiber.StatusBadRequest,
		},
		{
			description: "Returns 500 when no orders found",
			queryString: orderId1.String(),
			requestBody: models.Order{
				// ID:        orderId1,
				ProductID: "product-5",
			},
			mockData:  nil,
			mockError: nil,
			expection: fiber.StatusInternalServerError,
		},
		{
			description: "Returns 500 when someting went wrong.",
			queryString: orderId1.String(),
			requestBody: models.Order{
				ID:        orderId1,
				ProductID: "product-5",
			},
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-1",
			},
			mockError: errors.New("Return error!!!"),
			expection: fiber.StatusInternalServerError,
		},
	}
	for _, tc := range testcase {
		t.Run(tc.description, func(t *testing.T) {
			reqbody, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err, "Failed to marshal request body") // Check marshal error
			// reqbody, err := json.Marshal(tc.reqBody)
			mockRepo := &unittest_repositories.MockOrderRepo{
				MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{
					GetFn: func(filters *models.Order, preload []string) (*models.Order, error) {
						return tc.mockData, tc.mockError
					},
				},
			}
			mockEventbus := &unittest_eventbus.MockEventbus[models.Order]{}
			app := createTestApp(mockRepo, mockEventbus)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/orders/%s", tc.queryString), bytes.NewReader(reqbody))
			req.Header.Set("Content-type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expection, resp.StatusCode)

		})
	}
}

func TestDelete(t *testing.T) {
	// ✅ Add this line before setting up routes
	// os.Setenv("RABBITMQ_URL", "amqp://testMQ:password@localhost:5672/")

	// app.Post("/orders",  ) // Replace with your actual handler
	orderId1 := uuid.New()
	testcase := []struct {
		description string
		queryString string
		mockData    *models.Order
		mockError   error
		expection   int
	}{
		{
			description: "Returns 200 OK with deleting orders complete",
			queryString: orderId1.String(),
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-1",
			},
			mockError: nil,
			expection: fiber.StatusOK,
		},
		{
			description: "Returns 400 Sending incorrect query string.",
			queryString: "testttt",
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-1",
			},
			mockError: nil,
			expection: fiber.StatusBadRequest,
		},
		{
			description: "Returns 400 Sending space query string.",
			queryString: "%20",
			mockData: &models.Order{
				ID:        orderId1,
				ProductID: "product-1",
			},
			mockError: nil,
			expection: fiber.StatusBadRequest,
		},
		{
			description: "Returns 500 when no orders found",
			queryString: orderId1.String(),
			mockData:    nil,
			mockError:   nil,
			expection:   fiber.StatusInternalServerError,
		},
		{
			description: "Returns 500 something went wrong",
			queryString: orderId1.String(),
			mockData:    nil,
			mockError:   errors.New("Test error"),
			expection:   fiber.StatusInternalServerError,
		},
	}
	for _, tc := range testcase {
		t.Run(tc.description, func(t *testing.T) {
			// reqbody, err := json.Marshal(tc.reqBody)
			mockRepo := &unittest_repositories.MockOrderRepo{
				MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{
					GetFn: func(filters *models.Order, preload []string) (*models.Order, error) {
						return tc.mockData, tc.mockError
					},
				},
			}
			mockEventbus := &unittest_eventbus.MockEventbus[models.Order]{}
			app := createTestApp(mockRepo, mockEventbus)
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/orders/%s", tc.queryString), nil)
			req.Header.Set("Content-type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expection, resp.StatusCode)

		})
	}
}
