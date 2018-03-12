package nct

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ndphu/music-downloader/provider"
	iohelper "github.com/ndphu/music-downloader/utils/io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

var (
	AUTH_DIR    = path.Join(iohelper.GetAuthDir(), "nct")
	COOKIE_PATH = path.Join(AUTH_DIR, "cookie")
)

func (p *NCTProvider) Login(c *provider.LoginContext) error {
	loginUrl := "https://sso.nct.vn/auth/login?method=xlogin"
	client := http.Client{}

	form := getLoginForm(c)
	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	fmt.Println("Login...")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if err = os.MkdirAll(AUTH_DIR, 0777); err != nil {
		panic(err)
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "NCT_AUTH" {
			data, err := json.Marshal(cookie)
			if err != nil {
				panic(err)
			}
			if err = ioutil.WriteFile(COOKIE_PATH, data, 0777); err != nil {
				panic(err)
			}
			fmt.Println("Login successfully!")
			return nil
		}
	}

	fmt.Println("Login fail!")

	return errors.New("Login fail!")
}

func getAuthCookie() *http.Cookie {
	data, err := ioutil.ReadFile(COOKIE_PATH)
	if err != nil {
		return nil
	}
	cookie := http.Cookie{}

	if err = json.Unmarshal(data, &cookie); err != nil {
		panic(err)
	}
	return &cookie
}

func getLoginForm(c *provider.LoginContext) url.Values {
	return url.Values{
		"uname":     {c.UserName},
		"password":  {c.Password},
		"appName":   {"nhaccuatui"},
		"su":        {"https://www.nhaccuatui.com/ajax/user?type=login&"},
		"fu":        {"https://www.nhaccuatui.com/ajax/user?type=login&"},
		"ps":        {""},
		"pf":        {""},
		"checkbox1": {"on"},
	}
}
