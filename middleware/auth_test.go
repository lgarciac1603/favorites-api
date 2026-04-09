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

func newAuthServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/me" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_ValidToken() {
	server := newAuthServer(200, `{"id":123}`)
	defer server.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer user-123")

	AuthMiddleware(server.URL)(c)

	userID, exists := c.Get("userID")
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), "123", userID)
	assert.False(suite.T(), c.IsAborted())
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_NoToken() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	AuthMiddleware("http://example")(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.True(suite.T(), c.IsAborted())
	assert.Contains(suite.T(), w.Body.String(), "Token not found")
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidFormat() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	AuthMiddleware("http://example")(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.True(suite.T(), c.IsAborted())
	assert.Contains(suite.T(), w.Body.String(), "Invalid token format")
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_MissingBearer() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "user-123")

	AuthMiddleware("http://example")(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid token format")
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_EmptyToken() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer ")

	AuthMiddleware("http://example")(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Invalid token")
}

func (suite *AuthMiddlewareTestSuite) TestValidateToken_ValidToken() {
	server := newAuthServer(200, `{"id":456}`)
	defer server.Close()

	userID, err := ValidateToken(server.URL, "test-token-123")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "456", userID)
}

func (suite *AuthMiddlewareTestSuite) TestValidateToken_EmptyToken() {
	userID, err := ValidateToken("http://example", "")

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), userID)
	assert.Contains(suite.T(), err.Error(), "empty token")
}

func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
