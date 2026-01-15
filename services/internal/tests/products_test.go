package tests

import (
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
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
}

// TestGetProducts tests the /products endpoint
func (suite *ProductTestSuite) TestGetProducts() {
	req, _ := http.NewRequest("GET", "/products", nil)
	resp, _ := suite.app.Test(req)

	suite.Equal(200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	suite.Equal("OK", string(body))
}

// TestProductTestSuite runs the ProductTestSuite
func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}
