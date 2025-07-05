package unittest_order

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	models "github.com/Prompiriya084/go-mq/Models"
	unittest_repositories "github.com/Prompiriya084/go-mq/Producer/UnitTest/MockItem"
	routes "github.com/Prompiriya084/go-mq/Producer/Web/Routes"
	"github.com/stretchr/testify/assert"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MockOrderRepo struct {
	*unittest_repositories.MockRepositoryImpl[models.Order]
}

type MockProducer struct{}

func (m *MockProducer) PublishMessage(queue string, body []byte) error {
	return nil // simulate successful MQ send
}

// This ensures a clean slate and avoids state leakage between tests.
func createTestApp() *fiber.App {
	app := fiber.New()
	mockRepo := &MockOrderRepo{
		MockRepositoryImpl: &unittest_repositories.MockRepositoryImpl[models.Order]{},
	}
	mockProducer := &MockProducer{}
	routes.OrderSetupRouter(app, mockRepo, mockProducer) // Set up the routes for testing
	return app
}

func TestRoutes(t *testing.T) {
	// ✅ Add this line before setting up routes
	// os.Setenv("RABBITMQ_URL", "amqp://testMQ:password@localhost:5672/")
	app := createTestApp()
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
			// Use Fiber v3-compatible handler
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expection, resp.StatusCode)

		})
	}
}

// func TestAdd(t *testing.T) {
// 	testcase = []struct {
// 		name      string
// 		params    models.Order
// 		expection string
// 	}{

// 	}
// }
