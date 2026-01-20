package tests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/leroysb/go_kubernetes/internal/api/routes"
	"github.com/leroysb/go_kubernetes/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HealthCheckTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (suite *HealthCheckTestSuite) SetupTest() {
	suite.app = fiber.New()
	database.ConnectDB()
	routes.SetupRoutes(suite.app)
}

func (suite *HealthCheckTestSuite) TestHealthCheck() {
	req := httptest.NewRequest("GET", "/api/v1/status", nil)
	resp, err := suite.app.Test(req, -1)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), 200, resp.StatusCode)
	defer resp.Body.Close()
	var payload map[string]string
	json.NewDecoder(resp.Body).Decode(&payload)
	assert.Equal(suite.T(), "OK", payload["Postgres"])
}

func TestHealthCheckTestSuite(t *testing.T) {
	suite.Run(t, new(HealthCheckTestSuite))
}
