package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/airchains-network/da-client/models"
	"github.com/airchains-network/da-client/modules"
	celestiaTypes "github.com/airchains-network/da-client/types"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func CelestiaController(c *gin.Context) {
	rpcAUTH := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.0ohOAkKt_044L7oUXUMtGV27hoTJ0hR1fBH6p6fDhX0"
	daCelRPC := "http://34.131.171.247/celestia/"
	var bodyData celestiaTypes.CelestiaData
	if err := c.BindJSON(&bodyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   "Invalid JSON format",
			"daKeyHash": "nil",
		})
		return
	}

	jsonBodyData, err := json.Marshal(bodyData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   " Invalid JSON format",
			"daKeyHash": "nil",
		})
		return
	}
	encodedData := base64.StdEncoding.EncodeToString(jsonBodyData)

	payload := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "blob.Submit",
		"params": []interface{}{
			[]interface{}{map[string]interface{}{
				"namespace":     "AAAAAAAAAAAAAAAAAAAAAAAAAAECAwQFBgcICRA=",
				"data":          encodedData,
				"share_version": 0,
			}},
			map[string]interface{}{"Fee": 17980, "GasLimit": 179796},
		},
	}

	//* Marshal the payload struct to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":    502,
			"success":   false,
			"message":   " Error when marshalling payload",
			"daKeyHash": "nil",
		})
		return
	}

	//* Create a new HTTP client
	client := &http.Client{}

	//* Create a new POST request with headers and JSON payload
	req, err := http.NewRequest("POST", daCelRPC, bytes.NewBuffer(payloadJSON))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":    502,
			"success":   false,
			"message":   " Error when creating new request to Celestia DA",
			"daKeyHash": "nil",
		})
		return
	}

	//! Set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", rpcAUTH)
	// Send the request
	response, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   " Error when sending request to Celestia DA",
			"daKeyHash": "nil",
		})
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		// Handle the error
		fmt.Println("Error reading body:", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"status":    502,
			"success":   false,
			"message":   " Error reading body from Celestia DA",
			"daKeyHash": "nil",
		})
		return
	}

	responseJsonCelestia := []byte(string(body))

	var test map[string]interface{}
	err = json.Unmarshal(responseJsonCelestia, &test)
	fmt.Println(test)
	if test["error"] != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   "Panic Error sending request to Celestia DA",
			"daKeyHash": "nil",
		})
		return
	}

	var celestiaOutput models.APIResponseCelestia
	errorOfOutput := json.Unmarshal(responseJsonCelestia, &celestiaOutput)
	if errorOfOutput != nil {
		// Handle the error
		fmt.Println("Error unmarshalling JSON:", errorOfOutput)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   " Error unmarshalling JSON from Celestia DA",
			"daKeyHash": "nil",
		})
		return
	}

	if err := json.Unmarshal([]byte(responseJsonCelestia), &celestiaOutput); err == nil {

	} else {
		fmt.Println("Error parsing json Celestia")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   " Error parsing JSON from Celestia DA",
			"daKeyHash": "nil",
		})
		return
	}

	if err != nil {
		// return false, "Error sending request
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   " Data hasn't been submitted to Celestia DA!",
			"daKeyHash": "nil",
		})
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	stringOutputDaHashKey := strconv.Itoa(celestiaOutput.Result)
	fmt.Println(stringOutputDaHashKey)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"success":   true,
		"message":   " Data has been submitted to Celestia DA!",
		"daKeyHash": stringOutputDaHashKey,
	})

	return
}

func AvailController(c *gin.Context) {
	var bodyData celestiaTypes.CelestiaData
	if err := c.BindJSON(&bodyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   "Invalid JSON format",
			"daKeyHash": "nil",
		})
		return
	}

	jsonBodyData, err := json.Marshal(bodyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   " Invalid JSON format",
			"daKeyHash": "nil",
		})
		return
	}

	encodedData := base64.StdEncoding.EncodeToString(jsonBodyData)

	statusCheck, DaHash := modules.AvailModule(encodedData)
	fmt.Println(statusCheck)
	if statusCheck {
		c.JSON(http.StatusOK, gin.H{
			"status":    200,
			"success":   true,
			"message":   " Data has been submitted to Avail DA!",
			"daKeyHash": DaHash,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    400,
			"success":   false,
			"message":   " Data hasn't been submitted to Avail DA!",
			"daKeyHash": "nil",
		})
		return
	}
}
