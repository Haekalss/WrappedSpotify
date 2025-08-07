package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func CallSpotifyWithRefresh(requestURL, accessToken, refreshToken string) ([]byte, string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, accessToken, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		// refresh access token
		newAccessToken, err := RefreshAccessToken(refreshToken)
		if err != nil {
			return nil, "", fmt.Errorf("refresh failed: %v", err)
		}

		// try again with new token
		req, _ := http.NewRequest("GET", requestURL, nil)
		req.Header.Set("Authorization", "Bearer "+newAccessToken)
		resp, err = client.Do(req)
		if err != nil {
			return nil, "", err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return nil, "", fmt.Errorf("retry failed: %s", body)
		}
		return body, newAccessToken, nil
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("status %d: %s", resp.StatusCode, body)
	}
	return body, accessToken, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("spotify refresh failed: %s", body)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result.AccessToken, nil
}
