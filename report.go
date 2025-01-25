package security

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"suiet_server/llm"
	"suiet_server/resolver/meta"
	"suiet_server/resolver/rpc_client"
	"suiet_server/resolver/rpc_client/rpc_object"
	"suiet_server/schema/model"
	"suiet_server/utils/config"
	"suiet_server/utils/fn"
	"time"
)

var (
	reportHandler *ReportHandler
)

type ReportHandler struct {
	spreadsheetID string
	sheetService  *sheets.Service
}

func init() {
	if config.GetServerType() != "api" {
		return
	}
	// Get base64 encoded credentials from environment variable
	credentialsB64 := os.Getenv("GOOGLE_CREDENTIALS_BASE64")
	if credentialsB64 == "" {
		log.Fatal("GOOGLE_CREDENTIALS_BASE64 environment variable is not set")
	}

	// Decode base64 credentials
	credentialsJSON, err := base64.StdEncoding.DecodeString(credentialsB64)
	if err != nil {
		log.Fatalf("Failed to decode credentials: %v", err)
	}

	// Get spreadsheet ID from environment variable
	spreadsheetID := os.Getenv("REPORT_SPREADSHEET_ID")
	if spreadsheetID == "" {
		log.Fatal("REPORT_SPREADSHEET_ID environment variable is not set")
	}

	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		log.Fatalf("Failed to initialize report handler: %v", err)
	}

	reportHandler = &ReportHandler{
		spreadsheetID: spreadsheetID,
		sheetService:  srv,
	}
}

func HandleReport(c *gin.Context) {
	if reportHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Report handler not initialized",
		})
		return
	}

	var req meta.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	// get Object Info
	obj, err := rpc_object.GetObject(req.ObjectID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get object",
		})
		return
	}

	nfts := rpc_client.ObjectsToNFTs([]model.Object{obj}, true)
	if len(nfts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": "NFT not found",
		})
		return
	}
	go func() {
		isSafe, invalidateCategories, err := llm.ValidateNFTImage(nfts[0].Name, fn.NilToEmpty(nfts[0].Description), nfts[0].URL)

		var detectResult string
		if err != nil {
			detectResult = "error: " + err.Error()
		}

		if isSafe {
			detectResult = "safe"
		} else {
			detectResult = "unsafe" + fmt.Sprintf(" (%v)", invalidateCategories)
		}

		// Prepare the data for Google Sheets
		values := []interface{}{
			func() string {
				if isSafe {
					return "pending"
				} else {
					return "scam"
				}
			}(),
			time.Now().Format(time.RFC3339), //date
			req.Type,
			req.ObjectID,
			req.ObjectType,
			req.Submitter,
			nfts[0].Name,        // NFT Name
			nfts[0].Description, // description
			fmt.Sprintf("=IMAGE(\"%s\")", nfts[0].URL), // url
			nfts[0].Kiosk, // Kiosk
			detectResult,  // detect result
		}

		valueRange := &sheets.ValueRange{
			Values: [][]interface{}{values},
		}

		// Append the data to the spreadsheet
		_, err = reportHandler.sheetService.Spreadsheets.Values.Append(
			reportHandler.spreadsheetID,
			"A1", // Assuming we start from A1, adjust as needed
			valueRange,
		).ValueInputOption("USER_ENTERED").Do()

		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{
			//	"error": "Failed to save report",
			//})
			log.Printf("Failed to save report: %v", err)
			return
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Report submitted successfully",
	})
}
