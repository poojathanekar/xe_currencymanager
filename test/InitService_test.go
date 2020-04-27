package test

import (
  "github.com/nbio/st"
  "gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"xe_currencymanager/service"
  "net/http"
  "testing"
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
  req.Header.Set("Authorization","Basic am9zaDU4ODE3NTI1Mjo0dGNmbHQ4a2swMWlhcHVqMDR1aHA3MnVvZw==")
  res, err := service.RequestForAPI(req)
  st.Expect(t, err, nil)
  st.Expect(t, res.StatusCode, 200)

  body, _ := ioutil.ReadAll(res.Body)
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
  req.Header.Set("Authorization","Basic am9zaDU4ODE3NTI1Mjo0dGNmbHQ4a2swMWlhcHVqMDR1aHA3MnVvZw==")
  res, err := service.RequestForAPI(req)
  st.Expect(t, err, nil)
  st.Expect(t, res.StatusCode, 200)

  body, _ := ioutil.ReadAll(res.Body)
  st.Expect(t, string(body), xeResponse)
  st.Expect(t, gock.IsDone(), true)
}

