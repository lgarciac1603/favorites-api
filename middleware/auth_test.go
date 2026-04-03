// middleware/auth_test.go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_ValidToken() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer user-123")

	AuthMiddleware()(c)

	userID, exists := c.Get("userID")
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), "user-123", userID)
	assert.False(suite.T(), c.IsAborted())
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_NoToken() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	AuthMiddleware()(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.True(suite.T(), c.IsAborted())
	assert.Contains(suite.T(), w.Body.String(), "Token not found")
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidFormat() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	AuthMiddleware()(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.True(suite.T(), c.IsAborted())
	assert.Contains(suite.T(), w.Body.String(), "Invalid token format")
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_MissingBearer() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "user-123")

	AuthMiddleware()(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid token format")
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_EmptyToken() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer ")

	AuthMiddleware()(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid token")
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_TokenExtraction() {
	token := "my-secret-token-12345"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	AuthMiddleware()(c)

	userID, exists := c.Get("userID")
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), token, userID)
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_ContinuesToNextHandler() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")

	AuthMiddleware()(c)

	userID, _ := c.Get("userID")
	assert.Equal(suite.T(), "valid-token", userID)
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_MultipleSpaces() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer  token-123")

	AuthMiddleware()(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *AuthMiddlewareTestSuite) TestValidateToken_ValidToken() {
	userID, err := ValidateToken("test-token-123")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-token-123", userID)
}

func (suite *AuthMiddlewareTestSuite) TestValidateToken_EmptyToken() {
	userID, err := ValidateToken("")

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), userID)
	assert.Contains(suite.T(), err.Error(), "empty token")
}

func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}