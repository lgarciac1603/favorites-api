// models/favorite_test.go
package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FavoriteTestSuite struct {
	suite.Suite
}

func (suite *FavoriteTestSuite) TestFavorite_JSONMarshaling() {
	fav := Favorite{
		ID:         1,
		UserID:     10,
		CryptoID:   "bitcoin",
		CryptoName: "Bitcoin",
		CreatedAt:  "2026-04-02T10:30:00Z",
	}

	jsonBytes, err := json.Marshal(fav)

	assert.NoError(suite.T(), err)

	var decoded map[string]interface{}
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), float64(1), decoded["id"])
	assert.Equal(suite.T(), float64(10), decoded["userId"])
	assert.Equal(suite.T(), "bitcoin", decoded["cryptoId"])
	assert.Equal(suite.T(), "Bitcoin", decoded["cryptoName"])
	assert.Equal(suite.T(), "2026-04-02T10:30:00Z", decoded["createdAt"])
}

func (suite *FavoriteTestSuite) TestFavorite_JSONUnmarshaling() {
	jsonStr := `{
		"id": 5,
		"userId": 20,
		"cryptoId": "ethereum",
		"cryptoName": "Ethereum",
		"createdAt": "2026-04-02T11:00:00Z"
	}`

	var fav Favorite
	err := json.Unmarshal([]byte(jsonStr), &fav)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 5, fav.ID)
	assert.Equal(suite.T(), 20, fav.UserID)
	assert.Equal(suite.T(), "ethereum", fav.CryptoID)
	assert.Equal(suite.T(), "Ethereum", fav.CryptoName)
	assert.Equal(suite.T(), "2026-04-02T11:00:00Z", fav.CreatedAt)
}

func (suite *FavoriteTestSuite) TestFavorite_TagsMapping() {
	fav := Favorite{
		ID:         1,
		UserID:     1,
		CryptoID:   "bitcoin",
		CryptoName: "Bitcoin",
		CreatedAt:  "2026-04-02T10:30:00Z",
	}

	jsonBytes, _ := json.Marshal(fav)
	var result map[string]interface{}
	json.Unmarshal(jsonBytes, &result)

	assert.NotNil(suite.T(), result["id"])
	assert.NotNil(suite.T(), result["userId"])
	assert.NotNil(suite.T(), result["cryptoId"])
	assert.NotNil(suite.T(), result["cryptoName"])
	assert.NotNil(suite.T(), result["createdAt"])

	assert.Nil(suite.T(), result["user_id"])
	assert.Nil(suite.T(), result["crypto_id"])
}

func (suite *FavoriteTestSuite) TestFavorite_FieldTypes() {
	fav := Favorite{
		ID:         42,
		UserID:     99,
		CryptoID:   "test",
		CryptoName: "Test Coin",
		CreatedAt:  "2026-04-02T10:30:00Z",
	}

	assert.IsType(suite.T(), int(0), fav.ID)
	assert.IsType(suite.T(), int(0), fav.UserID)
	assert.IsType(suite.T(), "", fav.CryptoID)
	assert.IsType(suite.T(), "", fav.CryptoName)
	assert.IsType(suite.T(), "", fav.CreatedAt)
}

func (suite *FavoriteTestSuite) TestFavorite_EmptyStruct() {
	fav := Favorite{}

	jsonBytes, _ := json.Marshal(fav)

	assert.NotNil(suite.T(), jsonBytes)

	var result map[string]interface{}
	json.Unmarshal(jsonBytes, &result)
	assert.Equal(suite.T(), float64(0), result["id"])
	assert.Equal(suite.T(), float64(0), result["userId"])
	assert.Equal(suite.T(), "", result["cryptoId"])
}

func TestFavoriteTestSuite(t *testing.T) {
	suite.Run(t, new(FavoriteTestSuite))
}
