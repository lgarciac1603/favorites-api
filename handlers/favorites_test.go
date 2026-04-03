// handlers/favorites_test.go
package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FavoritesHandlerTestSuite struct {
	suite.Suite
	handler *FavoritesHandler
	db      sqlmock.Sqlmock
	mockDB  *sql.DB
}

// before test setup
func (suite *FavoritesHandlerTestSuite) SetupTest() {
	var err error
	suite.mockDB, suite.db, err = sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.handler = NewFavoritesHandler(suite.mockDB)
}

// Close mock connection
func (suite *FavoritesHandlerTestSuite) TeardownTest() {
	suite.mockDB.Close()
}

// ==================== GET /favorites ====================

func (suite *FavoritesHandlerTestSuite) TestGetFavorites_Success() {
	userID := 1
	rows := sqlmock.NewRows([]string{"id", "user_id", "crypto_id", "crypto_name", "created_at"}).
		AddRow(1, 1, "bitcoin", "Bitcoin", "2026-04-02T10:30:00Z").
		AddRow(2, 1, "ethereum", "Ethereum", "2026-04-02T10:35:00Z")

	suite.db.ExpectQuery(`SELECT id, user_id, crypto_id, crypto_name, created_at FROM user_favorites WHERE user_id = \$1 ORDER BY created_at DESC`).
		WithArgs(userID).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID)

	suite.handler.GetFavorites(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(2), response["total"])
	assert.NotNil(suite.T(), response["data"])
	assert.NoError(suite.T(), suite.db.ExpectationsWereMet())
}

func (suite *FavoritesHandlerTestSuite) TestGetFavorites_Empty() {
	userID := 2
	rows := sqlmock.NewRows([]string{"id", "user_id", "crypto_id", "crypto_name", "created_at"})

	suite.db.ExpectQuery(`SELECT id, user_id, crypto_id, crypto_name, created_at FROM user_favorites WHERE user_id = \$1 ORDER BY created_at DESC`).
		WithArgs(userID).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID)

	suite.handler.GetFavorites(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(0), response["total"])
	assert.NoError(suite.T(), suite.db.ExpectationsWereMet())
}

func (suite *FavoritesHandlerTestSuite) TestGetFavorites_NoAuth() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	suite.handler.GetFavorites(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Usuario no autenticado", response["error"])
}

func (suite *FavoritesHandlerTestSuite) TestGetFavorites_DBError() {
	userID := 1

	suite.db.ExpectQuery(`SELECT id, user_id, crypto_id, crypto_name, created_at FROM user_favorites WHERE user_id = \$1 ORDER BY created_at DESC`).
		WithArgs(userID).
		WillReturnError(sql.ErrConnDone)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID)

	suite.handler.GetFavorites(c)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Error consultando BD", response["error"])
}

// ==================== POST /favorites ====================

func (suite *FavoritesHandlerTestSuite) TestPostFavorite_Success() {
	userID := 1
	cryptoID := "solana"
	cryptoName := "Solana"

	suite.db.ExpectQuery(`SELECT id FROM user_favorites WHERE user_id = \$1 AND crypto_id = \$2`).
		WithArgs(userID, cryptoID).
		WillReturnError(sql.ErrNoRows)

	suite.db.ExpectQuery(`INSERT INTO user_favorites`).
		WithArgs(userID, cryptoID, cryptoName, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "crypto_id", "crypto_name", "created_at"}).
			AddRow(3, 1, "solana", "Solana", "2026-04-02T11:00:00Z"))

	body := map[string]string{
		"cryptoId":   cryptoID,
		"cryptoName": cryptoName,
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID)

	suite.handler.PostFavorite(c)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Crypto añadida a favoritos", response["message"])
	assert.NoError(suite.T(), suite.db.ExpectationsWereMet())
}

func (suite *FavoritesHandlerTestSuite) TestPostFavorite_AlreadyExists() {
	userID := 1
	cryptoID := "bitcoin"

	suite.db.ExpectQuery(`SELECT id FROM user_favorites WHERE user_id = \$1 AND crypto_id = \$2`).
		WithArgs(userID, cryptoID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	body := map[string]string{
		"cryptoId":   cryptoID,
		"cryptoName": "Bitcoin",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID)

	suite.handler.PostFavorite(c)

	assert.Equal(suite.T(), http.StatusConflict, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Esta crypto ya está en favoritos", response["error"])
}

func (suite *FavoritesHandlerTestSuite) TestPostFavorite_MissingCryptoId() {
	userID := 1
	body := map[string]string{
		"cryptoName": "Bitcoin",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID)

	suite.handler.PostFavorite(c)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(suite.T(), response["error"].(string), "requerido")
}

func (suite *FavoritesHandlerTestSuite) TestPostFavorite_MissingCryptoName() {
	userID := 1
	body := map[string]string{
		"cryptoId": "bitcoin",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID)

	suite.handler.PostFavorite(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *FavoritesHandlerTestSuite) TestPostFavorite_InvalidJSON() {
	userID := 1

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/favorites", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID)

	suite.handler.PostFavorite(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *FavoritesHandlerTestSuite) TestPostFavorite_NoAuth() {
	body := map[string]string{
		"cryptoId":   "bitcoin",
		"cryptoName": "Bitcoin",
	}
	bodyBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.handler.PostFavorite(c)
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *FavoritesHandlerTestSuite) TestPostFavorite_DBInsertError() {
	userID := 1
	cryptoID := "cardano"

	suite.db.ExpectQuery(`SELECT id FROM user_favorites WHERE user_id = \$1 AND crypto_id = \$2`).
		WithArgs(userID, cryptoID).
		WillReturnError(sql.ErrNoRows)

	suite.db.ExpectQuery(`INSERT INTO user_favorites`).
		WithArgs(userID, cryptoID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	body := map[string]string{
		"cryptoId":   cryptoID,
		"cryptoName": "Cardano",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/favorites", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID)

	suite.handler.PostFavorite(c)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

// ==================== DELETE /favorites/{cryptoId} ====================
func (suite *FavoritesHandlerTestSuite) TestDeleteFavorite_Success() {
	userID := 1
	cryptoID := "bitcoin"

	suite.db.ExpectExec(`DELETE FROM user_favorites WHERE user_id = \$1 AND crypto_id = \$2`).
		WithArgs(userID, cryptoID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "cryptoId", Value: cryptoID})
	c.Set("userID", userID)

	suite.handler.DeleteFavorite(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Crypto eliminada de favoritos", response["message"])
	assert.NoError(suite.T(), suite.db.ExpectationsWereMet())
}

func (suite *FavoritesHandlerTestSuite) TestDeleteFavorite_NotFound() {
	userID := 1
	cryptoID := "nonexistent"

	suite.db.ExpectExec(`DELETE FROM user_favorites WHERE user_id = \$1 AND crypto_id = \$2`).
		WithArgs(userID, cryptoID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "cryptoId", Value: cryptoID})
	c.Set("userID", userID)

	suite.handler.DeleteFavorite(c)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Crypto no encontrada en favoritos", response["error"])
}

func (suite *FavoritesHandlerTestSuite) TestDeleteFavorite_MissingCryptoId() {
	userID := 1

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID)

	suite.handler.DeleteFavorite(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "CryptoId es requerido", response["error"])
}

func (suite *FavoritesHandlerTestSuite) TestDeleteFavorite_NoAuth() {
	cryptoID := "bitcoin"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "cryptoId", Value: cryptoID})

	suite.handler.DeleteFavorite(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *FavoritesHandlerTestSuite) TestDeleteFavorite_DBError() {
	userID := 1
	cryptoID := "ethereum"

	suite.db.ExpectExec(`DELETE FROM user_favorites WHERE user_id = \$1 AND crypto_id = \$2`).
		WithArgs(userID, cryptoID).
		WillReturnError(sql.ErrConnDone)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "cryptoId", Value: cryptoID})
	c.Set("userID", userID)

	suite.handler.DeleteFavorite(c)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Error eliminando de BD", response["error"])
}

func TestFavoritesHandlerSuite(t *testing.T) {
	suite.Run(t, new(FavoritesHandlerTestSuite))
}
