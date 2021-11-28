package blackbox

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/atrariksa/awallet/configs"
	"github.com/atrariksa/awallet/constants"
	"github.com/atrariksa/awallet/models"
)

func createUserRequest(cfg *configs.Config, username string) (*http.Client, *http.Request) {
	var host string = cfg.APP.HOST + ":" + cfg.APP.PORT
	body, _ := json.Marshal(&models.CreateUserRequest{
		Username: username,
	})
	return &http.Client{Timeout: time.Second * 5},
		&http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Scheme: "http",
				Host:   host,
				Path:   constants.CREATE_USER_PATH,
			},
			Body: ioutil.NopCloser(bytes.NewReader(body)),
		}
}

func getNewUser(cfg *configs.Config, username string) models.CreateUserResponse {
	c, req := createUserRequest(cfg, username)
	resp, _ := c.Do(req)
	if resp.StatusCode != http.StatusCreated {
		return models.CreateUserResponse{}
	}
	createUser := models.CreateUserResponse{}
	err := getStruct(resp, &createUser)
	if err != nil {
		return models.CreateUserResponse{}
	}
	return createUser
}

func readBalanceRequest(cfg *configs.Config, token string) (*http.Client, *http.Request) {
	var host string = cfg.APP.HOST + ":" + cfg.APP.PORT
	header := http.Header{}
	header.Add("Authorization", token)
	return &http.Client{Timeout: time.Second * 5},
		&http.Request{
			Header: header,
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: "http",
				Host:   host,
				Path:   constants.READ_BALANCE_PATH,
			},
		}
}

func getBalance(cfg *configs.Config, token string) models.ReadBalanceResponse {
	c, req := readBalanceRequest(cfg, token)
	resp, _ := c.Do(req)
	if resp.StatusCode != http.StatusOK {
		return models.ReadBalanceResponse{}
	}
	readBalanceResp := models.ReadBalanceResponse{}
	err := getStruct(resp, &readBalanceResp)
	if err != nil {
		return models.ReadBalanceResponse{}
	}
	return readBalanceResp
}

func topupBalanceRequest(cfg *configs.Config, token string, amount uint32) (*http.Client, *http.Request) {
	var host string = cfg.APP.HOST + ":" + cfg.APP.PORT
	body, _ := json.Marshal(&models.TopupBalanceRequest{
		Amount: amount,
	})
	header := http.Header{}
	header.Add("Authorization", token)
	return &http.Client{Timeout: time.Second * 5},
		&http.Request{
			Header: header,
			Method: http.MethodPost,
			URL: &url.URL{
				Scheme: "http",
				Host:   host,
				Path:   constants.TOPUP_BALANCE_PATH,
			},
			Body: ioutil.NopCloser(bytes.NewReader(body)),
		}
}

func topupBalance(cfg *configs.Config, token string, amount uint32) int {
	c, req := topupBalanceRequest(cfg, token, amount)
	resp, _ := c.Do(req)
	if resp.StatusCode != http.StatusCreated {
		return 0
	}
	return resp.StatusCode
}

func transferRequest(cfg *configs.Config, token string, amount uint32, usernameDest string) (*http.Client, *http.Request) {
	var host string = cfg.APP.HOST + ":" + cfg.APP.PORT
	body, _ := json.Marshal(&models.TransferRequest{
		Amount:     amount,
		ToUsername: usernameDest,
	})
	header := http.Header{}
	header.Add("Authorization", token)
	return &http.Client{Timeout: time.Second * 5},
		&http.Request{
			Header: header,
			Method: http.MethodPost,
			URL: &url.URL{
				Scheme: "http",
				Host:   host,
				Path:   constants.TRANSFER_PATH,
			},
			Body: ioutil.NopCloser(bytes.NewReader(body)),
		}
}

func transfer(cfg *configs.Config, token string, amount uint32, usernameDest string) int {
	c, req := transferRequest(cfg, token, amount, usernameDest)
	resp, _ := c.Do(req)
	if resp.StatusCode != http.StatusNoContent {
		return 0
	}
	return resp.StatusCode
}

func topTransactionPerUserRequest(cfg *configs.Config, token string) (*http.Client, *http.Request) {
	var host string = cfg.APP.HOST + ":" + cfg.APP.PORT
	header := http.Header{}
	header.Add("Authorization", token)
	return &http.Client{Timeout: time.Second * 5},
		&http.Request{
			Header: header,
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: "http",
				Host:   host,
				Path:   constants.TOP_TRANSACTIONS_PER_USER_PATH,
			},
		}
}

func getTopTransactionPerUser(cfg *configs.Config, token string) models.TopTransactionsPerUser {
	c, req := topTransactionPerUserRequest(cfg, token)
	resp, _ := c.Do(req)
	if resp.StatusCode != http.StatusOK {
		return models.TopTransactionsPerUser{}
	}
	topTransactipPerUserResp := models.TopTransactionsPerUser{}
	err := getStruct(resp, &topTransactipPerUserResp)
	if err != nil {
		return models.TopTransactionsPerUser{}
	}
	return topTransactipPerUserResp
}

func topUsersRequest(cfg *configs.Config, token string) (*http.Client, *http.Request) {
	var host string = cfg.APP.HOST + ":" + cfg.APP.PORT
	header := http.Header{}
	header.Add("Authorization", token)
	return &http.Client{Timeout: time.Second * 5},
		&http.Request{
			Header: header,
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: "http",
				Host:   host,
				Path:   constants.TOP_USERS_PATH,
			},
		}
}

func getTopUsers(cfg *configs.Config, token string) models.TopUserResponse {
	c, req := topUsersRequest(cfg, token)
	resp, _ := c.Do(req)
	if resp.StatusCode != http.StatusOK {
		return models.TopUserResponse{}
	}
	topUsersResp := models.TopUserResponse{}
	err := getStruct(resp, &topUsersResp)
	if err != nil {
		return models.TopUserResponse{}
	}
	return topUsersResp
}

func getStruct(resp *http.Response, respStruct interface{}) (err error) {
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyByte, &respStruct)
	if err != nil {
		return
	}
	return
}
