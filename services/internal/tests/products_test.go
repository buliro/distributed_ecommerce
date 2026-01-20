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

// Define a suite struct that embeds testify's suite.Suite
type ProductTestSuite struct {
	suite.Suite
	app *fiber.App
}

// SetupTest sets the app to a new instance of the app
func (suite *ProductTestSuite) SetupTest() {
	suite.app = fiber.New()
	database.ConnectDB()
	routes.SetupRoutes(suite.app)
}

// TestGetProducts tests the /products endpoint
func (suite *ProductTestSuite) TestGetProducts() {
	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	resp, err := suite.app.Test(req, -1)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), 200, resp.StatusCode)
	defer resp.Body.Close()
	var products []map[string]any
	suite.NoError(json.NewDecoder(resp.Body).Decode(&products))
}

// TestProductTestSuite runs the ProductTestSuite
func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}
