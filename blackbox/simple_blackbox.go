package blackbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"log"

	"github.com/atrariksa/awallet/configs"
	"github.com/atrariksa/awallet/constants"
	"github.com/atrariksa/awallet/drivers"
	"github.com/atrariksa/awallet/models"
)

type Blackbox struct {
	cfg *configs.Config
}

func NewBlackBox(cfg *configs.Config) *Blackbox {
	return &Blackbox{
		cfg: cfg,
	}
}

func (b *Blackbox) Run() {
	b.cleanUp()
	b.runTestAPICreateUser()
	b.runTestAPIReadBalance()
	b.runTestAPITopupBalance()
	b.runTestAPITransferBalance()
	b.runTestAPITopTransactionPerUser()
	b.runTestAPITopUsers()
	b.cleanUp()
}

func (b *Blackbox) cleanUp() {
	dbWrite := drivers.NewDBClientWrite(b.cfg)
	dbWrite.Exec("DELETE FROM mutations")
	dbWrite.Exec("ALTER TABLE mutations AUTO_INCREMENT = 1")
	dbWrite.Exec("DELETE FROM user_total_outgoings")
	dbWrite.Exec("ALTER TABLE user_total_outgoings AUTO_INCREMENT = 1")
	dbWrite.Exec("DELETE FROM users")
	dbWrite.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

func (b *Blackbox) runTestAPICreateUser() {
	var host string = b.cfg.APP.HOST + ":" + b.cfg.APP.PORT
	var path string = constants.CREATE_USER_PATH

	var tests = []struct {
		testName    string
		prepare     func() (*http.Client, *http.Request)
		expectedMet func(*http.Response) error
	}{
		{
			"CreateUser #201: ",
			func() (*http.Client, *http.Request) {
				rand.Seed(time.Now().UnixNano())
				username := "any" + fmt.Sprintf("%v", rand.Int())
				return createUserRequest(b.cfg, username)
			},
			func(resp *http.Response) error {
				if resp.StatusCode != http.StatusCreated {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusCreated)
				}
				createUser := models.CreateUserResponse{}
				err := getStruct(resp, &createUser)
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			"CreateUser #400 : ",
			func() (*http.Client, *http.Request) {
				rand.Seed(time.Now().UnixNano())
				body, _ := json.Marshal(&models.CreateUserRequest{})
				return &http.Client{Timeout: time.Second * 1},
					&http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							Scheme: "http",
							Host:   host,
							Path:   path,
						},
						Body: ioutil.NopCloser(bytes.NewReader(body)),
					}
			},
			func(resp *http.Response) error {
				if resp.StatusCode != http.StatusBadRequest {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusBadRequest)
				}
				return nil
			},
		},
		{
			"CreateUser #409: ",
			func() (*http.Client, *http.Request) {
				rand.Seed(time.Now().UnixNano())
				username := "any" + fmt.Sprintf("%v", rand.Int())
				getNewUser(b.cfg, username) // create user

				body, _ := json.Marshal(&models.CreateUserRequest{
					Username: username,
				})

				client := &http.Client{Timeout: time.Second * 1}
				req := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   host,
						Path:   path,
					},
					Body: ioutil.NopCloser(bytes.NewReader(body)),
				}

				return client, req
			},
			func(resp *http.Response) error {
				if resp.StatusCode != http.StatusConflict {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusConflict)
				}

				return nil
			},
		},
	}
	for _, v := range tests {
		client, req := v.prepare()
		resp, err := client.Do(req)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		err = v.expectedMet(resp)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		log.Println(v.testName, "PASS")
	}
}

func (b *Blackbox) runTestAPIReadBalance() {

	var tests = []struct {
		testName    string
		prepare     func() (*http.Client, *http.Request)
		expectedMet func(*http.Response) error
	}{
		{
			"ReadBalance #200: ",
			func() (*http.Client, *http.Request) {
				rand.Seed(time.Now().UnixNano())
				username := "any" + fmt.Sprintf("%v", rand.Int())
				newUser := getNewUser(b.cfg, username)

				header := http.Header{}
				header.Add("Authorization", newUser.Token)

				return readBalanceRequest(b.cfg, newUser.Token)
			},
			func(resp *http.Response) error {
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusOK)
				}
				readBalanceResp := models.ReadBalanceResponse{}
				err := getStruct(resp, &readBalanceResp)
				if err != nil {
					return err
				}
				if readBalanceResp.Balance != 0 {
					return fmt.Errorf("Got %v, Want %v", readBalanceResp.Balance, 0)
				}
				return nil
			},
		},
		{
			"ReadBalance #401 : ",
			func() (*http.Client, *http.Request) {
				return readBalanceRequest(b.cfg, "invalid token")
			},
			func(resp *http.Response) error {
				if resp.StatusCode != http.StatusUnauthorized {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusUnauthorized)
				}
				return nil
			},
		},
	}
	for _, v := range tests {
		client, req := v.prepare()
		resp, err := client.Do(req)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		err = v.expectedMet(resp)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		log.Println(v.testName, "PASS")
	}
}

func (b *Blackbox) runTestAPITopupBalance() {

	var tests = []struct {
		testName    string
		prepare     func() (*http.Client, *http.Request, models.CreateUserResponse)
		expectedMet func(*http.Response, models.CreateUserResponse) error
	}{
		{
			"BalanceTopup #204: ",
			func() (*http.Client, *http.Request, models.CreateUserResponse) {
				rand.Seed(time.Now().UnixNano())
				username := "any" + fmt.Sprintf("%v", rand.Int())
				newUser := getNewUser(b.cfg, username)
				client, req := topupBalanceRequest(b.cfg, newUser.Token, uint32(20000))
				return client, req, newUser
			},
			func(resp *http.Response, newUser models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusNoContent {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusNoContent)
				}
				newUserBalanceResp := getBalance(b.cfg, newUser.Token)
				if newUserBalanceResp.Balance != 20000 {
					return fmt.Errorf("Got %v, Want %v", newUserBalanceResp.Balance, 20000)
				}
				return nil
			},
		},
		{
			"BalanceTopup #400 : ",
			func() (*http.Client, *http.Request, models.CreateUserResponse) {
				rand.Seed(time.Now().UnixNano())
				username := "any" + fmt.Sprintf("%v", rand.Int())
				newUser := getNewUser(b.cfg, username)
				client, req := topupBalanceRequest(b.cfg, newUser.Token, uint32(0))
				return client, req, newUser
			},
			func(resp *http.Response, newUser models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusBadRequest {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusBadRequest)
				}
				return nil
			},
		},
		{
			"BalanceTopup #401 : ",
			func() (*http.Client, *http.Request, models.CreateUserResponse) {
				client, req := topupBalanceRequest(b.cfg, "newUserToken", uint32(0))
				return client, req, models.CreateUserResponse{}
			},
			func(resp *http.Response, newUser models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusUnauthorized {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusUnauthorized)
				}
				return nil
			},
		},
	}
	for _, v := range tests {
		client, req, newUser := v.prepare()
		resp, err := client.Do(req)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		err = v.expectedMet(resp, newUser)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		log.Println(v.testName, "PASS")
	}
}

func (b *Blackbox) runTestAPITransferBalance() {

	var tests = []struct {
		testName    string
		prepare     func() (*http.Client, *http.Request, models.CreateUserResponse, models.CreateUserResponse)
		expectedMet func(*http.Response, models.CreateUserResponse, models.CreateUserResponse) error
	}{
		{
			"Transfer #204: ",
			func() (*http.Client, *http.Request, models.CreateUserResponse, models.CreateUserResponse) {
				rand.Seed(time.Now().UnixNano())

				senderUsername := "sender" + fmt.Sprintf("%v", rand.Int())
				sender := getNewUser(b.cfg, senderUsername)

				topupBalance(b.cfg, sender.Token, 50000)

				destUsername := "dest" + fmt.Sprintf("%v", rand.Int())
				destUser := getNewUser(b.cfg, destUsername)

				client, req := transferRequest(b.cfg, sender.Token, 20000, destUser.UserDetails.Username)
				return client, req, sender, destUser
			},
			func(resp *http.Response, sender, destUser models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusNoContent {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusNoContent)
				}

				senderBalanceResp := getBalance(b.cfg, sender.Token)
				if senderBalanceResp.Balance != 30000 {
					return fmt.Errorf("Got %v, Want %v", senderBalanceResp.Balance, 30000)
				}

				destBalanceResp := getBalance(b.cfg, destUser.Token)
				if destBalanceResp.Balance != 20000 {
					return fmt.Errorf("Got %v, Want %v", destBalanceResp.Balance, 20000)
				}

				return nil
			},
		},
		{
			"Transfer #400: ",
			func() (*http.Client, *http.Request, models.CreateUserResponse, models.CreateUserResponse) {
				rand.Seed(time.Now().UnixNano())

				senderUsername := "sender" + fmt.Sprintf("%v", rand.Int())
				sender := getNewUser(b.cfg, senderUsername)

				destUsername := "dest" + fmt.Sprintf("%v", rand.Int())
				destUser := getNewUser(b.cfg, destUsername)

				client, req := transferRequest(b.cfg, sender.Token, 20000, destUser.UserDetails.Username)
				return client, req, sender, destUser
			},
			func(resp *http.Response, sender, destUser models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusBadRequest {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusBadRequest)
				}

				return nil
			},
		},
		{
			"Transfer #401: ",
			func() (*http.Client, *http.Request, models.CreateUserResponse, models.CreateUserResponse) {
				rand.Seed(time.Now().UnixNano())

				senderUsername := "sender" + fmt.Sprintf("%v", rand.Int())
				sender := getNewUser(b.cfg, senderUsername)

				destUsername := "dest" + fmt.Sprintf("%v", rand.Int())
				destUser := getNewUser(b.cfg, destUsername)

				client, req := transferRequest(b.cfg, "sender.Token", 20000, destUser.UserDetails.Username)
				return client, req, sender, destUser
			},
			func(resp *http.Response, sender, destUser models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusUnauthorized {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusUnauthorized)
				}

				return nil
			},
		},
		{
			"Transfer #404: ",
			func() (*http.Client, *http.Request, models.CreateUserResponse, models.CreateUserResponse) {
				rand.Seed(time.Now().UnixNano())

				senderUsername := "sender" + fmt.Sprintf("%v", rand.Int())
				sender := getNewUser(b.cfg, senderUsername)

				topupBalance(b.cfg, sender.Token, 50000)

				destUsername := "dest" + fmt.Sprintf("%v", rand.Int())

				client, req := transferRequest(b.cfg, sender.Token, 20000, destUsername)
				return client, req, sender, models.CreateUserResponse{}
			},
			func(resp *http.Response, sender, destUser models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusNotFound {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusNotFound)
				}

				return nil
			},
		},
	}
	for _, v := range tests {
		client, req, sender, destUser := v.prepare()
		resp, err := client.Do(req)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		err = v.expectedMet(resp, sender, destUser)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		log.Println(v.testName, "PASS")
	}
}

func (b *Blackbox) runTestAPITopTransactionPerUser() {

	var tests = []struct {
		testName    string
		prepare     func() (*http.Client, *http.Request, models.CreateUserResponse, models.CreateUserResponse)
		expectedMet func(*http.Response, models.CreateUserResponse, models.CreateUserResponse) error
	}{
		{
			"TopTransactionPerUser #200: ",
			func() (*http.Client, *http.Request, models.CreateUserResponse, models.CreateUserResponse) {
				rand.Seed(time.Now().UnixNano())

				senderUsername := "sender" + fmt.Sprintf("%v", rand.Int())
				sender := getNewUser(b.cfg, senderUsername)

				destUsername := "dest" + fmt.Sprintf("%v", rand.Int())
				destUser := getNewUser(b.cfg, destUsername)

				topupBalance(b.cfg, sender.Token, uint32(5000000))
				topupBalance(b.cfg, destUser.Token, uint32(2000000))

				for i := 1; i < 11; i++ {
					transfer(b.cfg, sender.Token, uint32(i*25000), destUsername)
					transfer(b.cfg, destUser.Token, uint32(i*17000), senderUsername)
				}

				client, req := topTransactionPerUserRequest(b.cfg, sender.Token)
				return client, req, sender, destUser
			},
			func(resp *http.Response, sender, destUser models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusOK)
				}

				topTransactionPerUser := []models.TopTransactionsPerUser{}
				err := getStruct(resp, &topTransactionPerUser)
				if err != nil {
					return err
				}

				expectedResp := `[
					{
					  "username": "%v",
					  "amount": -250000
					},
					{
					  "username": "%v",
					  "amount": -225000
					},
					{
					  "username": "%v",
					  "amount": -200000
					},
					{
					  "username": "%v",
					  "amount": -175000
					},
					{
					  "username": "%v",
					  "amount": 170000
					},
					{
					  "username": "%v",
					  "amount": 153000
					},
					{
					  "username": "%v",
					  "amount": -150000
					},
					{
					  "username": "%v",
					  "amount": 136000
					},
					{
					  "username": "%v",
					  "amount": -125000
					},
					{
					  "username": "%v",
					  "amount": 119000
					}
				  ]`
				expectedResp = fmt.Sprintf(expectedResp,
					destUser.UserDetails.Username,
					destUser.UserDetails.Username,
					destUser.UserDetails.Username,
					destUser.UserDetails.Username,
					destUser.UserDetails.Username,
					destUser.UserDetails.Username,
					destUser.UserDetails.Username,
					destUser.UserDetails.Username,
					destUser.UserDetails.Username,
					destUser.UserDetails.Username,
				)

				expectedTopTrf := []models.TopTransactionsPerUser{}
				err = json.Unmarshal([]byte(expectedResp), &expectedTopTrf)
				if err != nil {
					return err
				}

				if len(expectedTopTrf) != len(topTransactionPerUser) {
					return fmt.Errorf("Got %v, Want %v", len(topTransactionPerUser), len(expectedTopTrf))
				}

				for k, v := range expectedTopTrf {
					if v.Username != topTransactionPerUser[k].Username {
						return fmt.Errorf("Got %v, Want %v", topTransactionPerUser[k].Username, v.Username)
					}
					if v.Amount != topTransactionPerUser[k].Amount {
						return fmt.Errorf("Got %v, Want %v", topTransactionPerUser[k].Amount, v.Amount)
					}
				}
				return nil
			},
		},
		{
			"TopTransactionPerUser #401: ",
			func() (*http.Client, *http.Request, models.CreateUserResponse, models.CreateUserResponse) {
				client, req := topTransactionPerUserRequest(b.cfg, "")
				return client, req, models.CreateUserResponse{}, models.CreateUserResponse{}
			},
			func(resp *http.Response, sender, destUser models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusUnauthorized {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusUnauthorized)
				}

				return nil
			},
		},
	}
	for _, v := range tests {
		client, req, sender, destUser := v.prepare()
		resp, err := client.Do(req)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		err = v.expectedMet(resp, sender, destUser)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		log.Println(v.testName, "PASS")
	}
}

func (b *Blackbox) runTestAPITopUsers() {

	var tests = []struct {
		testName    string
		prepare     func() (*http.Client, *http.Request, []models.CreateUserResponse, []models.CreateUserResponse)
		expectedMet func(*http.Response, []models.CreateUserResponse, []models.CreateUserResponse) error
	}{
		{
			"TopUsers #200: ",
			func() (*http.Client, *http.Request, []models.CreateUserResponse, []models.CreateUserResponse) {

				var groupA []models.CreateUserResponse
				var groupB []models.CreateUserResponse
				for i := 0; i < 7; i++ {
					rand.Seed(time.Now().UnixNano())
					groupAUsername := fmt.Sprintf(`groupA_%v_%v`, i, rand.Int())
					userA := getNewUser(b.cfg, groupAUsername)
					groupA = append(groupA, userA)
					topupBalance(b.cfg, userA.Token, uint32(9000000))

					rand.Seed(time.Now().UnixNano())
					groupBUsername := fmt.Sprintf(`groupB_%v_%v`, i, rand.Int())
					userB := getNewUser(b.cfg, groupBUsername)
					groupB = append(groupB, userB)
					topupBalance(b.cfg, userB.Token, uint32(9000000))
				}

				for k, v := range groupA {
					groupBUsername := groupB[k].UserDetails.Username
					for i := 0; i < 20; i++ {
						if k%2 == 0 {
							transfer(b.cfg, v.Token, uint32(i*12000), groupBUsername)
						} else {
							transfer(b.cfg, v.Token, uint32(i*9000), groupBUsername)
						}
					}

					groupAUsername := v.UserDetails.Username
					for i := 0; i < 17; i++ {
						if k%2 == 0 {
							transfer(b.cfg, groupB[k].Token, uint32(i*13000), groupAUsername)
						} else {
							transfer(b.cfg, groupB[k].Token, uint32(i*10000), groupAUsername)
						}
					}
				}

				client, req := topUsersRequest(b.cfg, groupA[0].Token)
				return client, req, groupA, groupB
			},
			func(resp *http.Response, groupA, groupB []models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusOK)
				}

				topUsers := []models.TopUserResponse{}
				err := getStruct(resp, &topUsers)
				if err != nil {
					return err
				}

				expectedResp := `[
					{
					  "username": "%v",
					  "transacted_value": 2280000
					},
					{
					  "username": "%v",
					  "transacted_value": 2280000
					},
					{
					  "username": "%v",
					  "transacted_value": 2280000
					},
					{
					  "username": "%v",
					  "transacted_value": 2280000
					},
					{
					  "username": "%v",
					  "transacted_value": 1768000
					},
					{
					  "username": "%v",
					  "transacted_value": 1768000
					},
					{
					  "username": "%v",
					  "transacted_value": 1768000
					},
					{
					  "username": "%v",
					  "transacted_value": 1768000
					},
					{
					  "username": "%v",
					  "transacted_value": 1710000
					},
					{
					  "username": "%v",
					  "transacted_value": 1710000
					}
				  ]`
				expectedResp = fmt.Sprintf(expectedResp,
					groupA[6].UserDetails.Username,
					groupA[4].UserDetails.Username,
					groupA[2].UserDetails.Username,
					groupA[0].UserDetails.Username,
					groupB[6].UserDetails.Username,
					groupB[4].UserDetails.Username,
					groupB[2].UserDetails.Username,
					groupB[0].UserDetails.Username,
					groupA[5].UserDetails.Username,
					groupA[3].UserDetails.Username,
				)

				expectedTopUsers := []models.TopUserResponse{}
				err = json.Unmarshal([]byte(expectedResp), &expectedTopUsers)
				if err != nil {
					return err
				}

				if len(expectedTopUsers) != len(topUsers) {
					return fmt.Errorf("Got %v, Want %v", len(topUsers), len(expectedTopUsers))
				}

				for k, v := range expectedTopUsers {
					if v.Username != topUsers[k].Username {
						return fmt.Errorf("Got %v, Want %v", topUsers[k].Username, v.Username)
					}
					if v.TransactedValue != topUsers[k].TransactedValue {
						return fmt.Errorf("Got %v, Want %v", topUsers[k].TransactedValue, v.TransactedValue)
					}
				}

				return nil
			},
		},
		{
			"TopUsers #401: ",
			func() (*http.Client, *http.Request, []models.CreateUserResponse, []models.CreateUserResponse) {
				client, req := topUsersRequest(b.cfg, "")
				return client, req, []models.CreateUserResponse{}, []models.CreateUserResponse{}
			},
			func(resp *http.Response, groupA, groupB []models.CreateUserResponse) error {
				if resp.StatusCode != http.StatusUnauthorized {
					return fmt.Errorf("Got %v, Want %v", resp.StatusCode, http.StatusUnauthorized)
				}

				return nil
			},
		},
	}
	for _, v := range tests {
		client, req, sender, destUser := v.prepare()
		resp, err := client.Do(req)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		err = v.expectedMet(resp, sender, destUser)
		if err != nil {
			log.Println(v.testName, "FAILED", err)
			continue
		}
		log.Println(v.testName, "PASS")
	}
}
