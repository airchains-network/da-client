package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/airchains-network/da-client/models"
	"github.com/airchains-network/da-client/types"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "🌟 Modular Code in Action!"})
}

// Define other handlers here

func CelestiaController(c *gin.Context) {

	var bodyData types.CelestiaData

	if err := c.BindJSON(&bodyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      400,
			"success":     false,
			"message":     "Invalid JSON format",
			"description": "Please check the JSON format and try again.",
		})
		return
	}

	jsonBodyData, err := json.Marshal(bodyData)
	if err != nil {
		log.Fatal(err)
	}
	encodedData := base64.StdEncoding.EncodeToString(jsonBodyData)
	fmt.Println(encodedData)
	//  !Change this to env variable
	rpcAUTH := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.0ohOAkKt_044L7oUXUMtGV27hoTJ0hR1fBH6p6fDhX0"
	daCelRPC := "http://34.131.171.247/celestia/"

	//! Data Value For Encording
	//encodedString := base64.StdEncoding.EncodeToString([]byte(dataValue))

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
			map[string]interface{}{"Fee": 7980, "GasLimit": 79796},
		},
	}

	//* Marshal the payload struct to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		// return false, "Error in Payload JSON"
	}

	//* Create a new HTTP client
	client := &http.Client{}

	//* Create a new POST request with headers and JSON payload
	req, err := http.NewRequest("POST", daCelRPC, bytes.NewBuffer(payloadJSON))
	if err != nil {
		// return false, "Error creating request"
	}

	//! Set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", rpcAUTH)
	// Send the request
	response, err := client.Do(req)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		// Handle the error
		fmt.Println("Error reading body:", err)
		return
	}

	responseJsonCelestia := []byte(string(body))

	var celestiaOutput models.APIResponseCelestia
	errorOfOutput := json.Unmarshal(responseJsonCelestia, &celestiaOutput)
	if errorOfOutput != nil {
		// Handle the error
		fmt.Println("Error unmarshalling JSON:", errorOfOutput)
		return
	}

	if err := json.Unmarshal([]byte(responseJsonCelestia), &celestiaOutput); err == nil {
		fmt.Println("Successfully parsed json Celestia")
		fmt.Println(celestiaOutput.Result)
	} else {
		fmt.Println("Error parsing json Celestia")
	}

	if err != nil {
		// return false, "Error sending request
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"success": false,
			"message": " Data hasn't been submitted to Celestia DA!",
			"daHash":  "daHash",
		})
	}

	defer response.Body.Close()
	fmt.Println(celestiaOutput.Result)
	c.JSON(http.StatusOK, gin.H{
		"status":           200,
		"success":          true,
		"message":          " Data has been submitted to Celestia DA!",
		"daIncludedHeight": celestiaOutput.Result,
	})
}

func AvailController(c *gin.Context) {

}
