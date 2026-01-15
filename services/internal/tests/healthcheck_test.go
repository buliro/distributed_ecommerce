package tests

import (
	"io"
	"net/http"
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
	c := http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/status", nil)
	resp, _ := c.Do(req)

	assert.Equal(suite.T(), 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"Postgres": "OK"}`, string(body))
}

func TestHealthCheckTestSuite(t *testing.T) {
	suite.Run(t, new(HealthCheckTestSuite))
}
