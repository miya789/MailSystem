package redmine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type (
	Redmine struct {
		client   http.Client
		UseProxy bool
	}

	// RedmineResponse represents response.
	//
	// https://tikasan.hatenablog.com/entry/2017/04/26/110854
	// from https://mholt.github.io/json-to-go/
	RedmineResponse struct {
		Issues     []Issue `json:"issues"`
		TotalCount int     `json:"total_count"`
		Offset     int     `json:"offset"`
		Limit      int     `json:"limit"`
	}

	Issue struct {
		ID      int `json:"id"`
		Project struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		Tracker struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tracker"`
		Status struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"status"`
		Priority struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"priority"`
		Author struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"author"`
		AssignedTo struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"assigned_to,omitempty"`
		Parent struct {
			ID int `json:"id"`
		} `json:"parent,omitempty"`
		Subject        string      `json:"subject"`
		Description    string      `json:"description"`
		StartDate      string      `json:"start_date"`
		DueDate        interface{} `json:"due_date"`
		DoneRatio      int         `json:"done_ratio"`
		IsPrivate      bool        `json:"is_private"`
		EstimatedHours interface{} `json:"estimated_hours"`
		CreatedOn      time.Time   `json:"created_on"`
		UpdatedOn      time.Time   `json:"updated_on"`
		ClosedOn       interface{} `json:"closed_on"`
		CustomFields   []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"custom_fields,omitempty"`
	}
)

const keyPage = "page"

var (
	redmineAPIKey string
	proxyURL      string
)

func init() {
	if err := godotenv.Load("../../config/.env"); err != nil {
		log.Println(fmt.Errorf("Failed to getIssues(): failed to read \".env\""))
		return
	}
	redmineAPIKey = os.Getenv("redmineAPIKey")
	proxyURL = os.Getenv("proxyURL")
	return
}

func (r *Redmine) GetIssues(redmineURL string) ([]Issue, error) {
	req, _ := http.NewRequest(http.MethodGet, redmineURL, nil)
	req.Header.Set("X-Redmine-API-Key", redmineAPIKey)
	if r.UseProxy {
		proxyURL, _ := url.Parse(proxyURL)
		r.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to getIssues(): %w", err)
	}
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	rResp := new(RedmineResponse)
	if err := json.Unmarshal([]byte(byteArray), rResp); err != nil {
		return nil, fmt.Errorf("Failed to getIssues(): %w", err)
	}

	log.Printf("TotalCount: \t%d\n", rResp.TotalCount)
	log.Printf("Limit: \t\t%d\n", rResp.Limit)
	log.Printf("Offset: \t\t%d\n", rResp.Offset)

	var respIssues []Issue
	rest := rResp.TotalCount
	for i := 1; rest > 0; i++ {
		// ページ移動
		q := req.URL.Query()
		q.Set(keyPage, strconv.Itoa(i))
		req.URL.RawQuery = q.Encode()

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Failed to getIssues(): %w", err)
		}
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		rResp := new(RedmineResponse)
		if err := json.Unmarshal([]byte(byteArray), rResp); err != nil {
			return nil, fmt.Errorf("Failed to getIssues(): %w", err)
		}

		respIssues = append(respIssues, rResp.Issues...)

		rest -= rResp.Limit
	}

	return respIssues, nil
}
