package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var DriveService *drive.Service

// InitDriveService initializes the Google Drive API client using OAuth2.
func InitDriveService(credentialsFile string) error {
	ctx := context.Background()
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	// Read token.json
	tokenFile := "token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		return fmt.Errorf("token.json tidak ditemukan. Silakan jalankan 'go run setup_drive.go' terlebih dahulu: %v", err)
	}

	client := config.Client(ctx, tok)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	DriveService = srv
	return nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// UploadFileToDrive uploads a file to Google Drive and returns the file ID.
func UploadFileToDrive(localFilePath, customFileName, folderId string) (string, error) {
	if DriveService == nil {
		return "", fmt.Errorf("drive service is not initialized")
	}

	file, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("unable to open local file: %v", err)
	}
	defer file.Close()

	driveFile := &drive.File{
		Name: customFileName,
	}
	
	if folderId != "" {
		driveFile.Parents = []string{folderId}
	}

	// We assume it's an image based on extension or we can let Drive figure it out.

	res, err := DriveService.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create file in drive: %v", err)
	}

	return res.Id, nil
}
