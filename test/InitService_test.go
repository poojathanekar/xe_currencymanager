package test

import (
  "github.com/nbio/st"
  "gopkg.in/h2non/gock.v1"
	"io/ioutil"
  "net/http"
  "testing"
  "time"
)

// TestGetXECurrencyAPI test for correct result Test to PASS
func TestGetXECurrencyAPI(t *testing.T) {
xeResponse :=`{
    "terms": "http://www.xe.com/legal/dfs.php",
    "privacy": "http://www.xe.com/privacy.php",
    "from": "USD",
    "amount": 1,
    "timestamp": "2020-04-27T00:00:00Z",
    "to": [
        {
            "quotecurrency": "INR",
            "mid": 76.2834881272
        }
    ]
}`

  defer gock.Off()
  gock.New("https://xecdapi.xe.com/v1/convert_from.json").
    MatchParams(map[string]string{
			"from": "USD",
			"to":"INR",
		}).
    Reply(200).
    JSON([]byte(xeResponse))

  req, err := http.NewRequest("GET","https://xecdapi.xe.com/v1/convert_from.json/?from=USD&to=INR",nil)
  req.Header.Set("Authorization","Basic bWtjbDY3NDU1NTA4Om5kZWxpOGwzajZlZnVubzFxbGszcnE2bHRp")
  client := &http.Client{}

	client = &http.Client{
		Timeout: time.Second * time.Duration(1500),
	}
	resp, err := client.Do(req)

  st.Expect(t, err, nil)
  st.Expect(t, resp.StatusCode, 200)

  body, _ := ioutil.ReadAll(resp.Body)
  st.Expect(t, string(body), xeResponse)
  st.Expect(t, gock.IsDone(), true)
}

// TestGetXECurrencyAPI test for incorrect input
func TestGetXECurrencyAPIForIncorrectInput(t *testing.T) {
xeResponse :=`{
    "code": 7,
    "message": "No USDD found on 2020-04-27T00:00:00Z",
    "documentation_url": "https://xecdapi.xe.com/docs/v1/"
}`
  defer gock.Off()
  gock.New("https://xecdapi.xe.com/v1/convert_from.json").
    MatchParams(map[string]string{
			"from": "USSDD",
			"to":"INR",
		}).
    Reply(200).
    JSON([]byte(xeResponse))

  req, err := http.NewRequest("GET","https://xecdapi.xe.com/v1/convert_from.json/?from=USSDD&to=INR",nil)
  req.Header.Set("Authorization","Basic bWtjbDY3NDU1NTA4Om5kZWxpOGwzajZlZnVubzFxbGszcnE2bHRp")
  client := &http.Client{}

	client = &http.Client{
		Timeout: time.Second * time.Duration(1500),
	}
	resp, err := client.Do(req)

  st.Expect(t, err, nil)
  st.Expect(t, resp.StatusCode, 200)

  body, _ := ioutil.ReadAll(resp.Body)
  st.Expect(t, string(body), xeResponse)
  st.Expect(t, gock.IsDone(), true)
}


func TestGetXECurrencyAPIForCorrectInput(t *testing.T) {
xeResponse :=`{
    "terms": "http://www.xe.com/legal/dfs.php",
    "privacy": "http://www.xe.com/privacy.php",
    "from": "USD",
    "amount": 1,
    "timestamp": "2020-04-28T00:00:00Z",
    "to": [
        {
            "quotecurrency": "AFN",
            "mid": 75.8376356125
        },
        {
            "quotecurrency": "INR",
            "mid": 76.341274431
        }
    ]
}`
  defer gock.Off()
  gock.New("https://xecdapi.xe.com/v1/convert_from.json").
    MatchParams(map[string]string{
			"from": "USD",
			"to":"INR,AFN",
		}).
    Reply(200).
    JSON([]byte(xeResponse))

  req, err := http.NewRequest("GET","https://xecdapi.xe.com/v1/convert_from.json/?from=USD&to=INR,AFN",nil)
  req.Header.Set("Authorization","Basic bWtjbDY3NDU1NTA4Om5kZWxpOGwzajZlZnVubzFxbGszcnE2bHRp")
  client := &http.Client{}

	client = &http.Client{
		Timeout: time.Second * time.Duration(1500),
	}
	resp, err := client.Do(req)

  st.Expect(t, err, nil)
  st.Expect(t, resp.StatusCode, 200)

  body, _ := ioutil.ReadAll(resp.Body)
  st.Expect(t, string(body), xeResponse)
  st.Expect(t, gock.IsDone(), true)
}
