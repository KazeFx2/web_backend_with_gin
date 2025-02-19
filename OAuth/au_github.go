package OAuth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"main/AccessControl"
	"main/Logger"
	"main/Net"
	"main/ParamTools"
	"main/QuickRes"
	"main/Router"
	"main/Vars"
	"net/http"
)

type tokenData struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

func AuGithubCallback(r *gin.Engine) {
	r.GET(Router.FatherPath+Router.AuGithub, func(c *gin.Context) {
		acc, _, _, _, _, _ := AccessControl.AccControl(c)
		if !acc {
			return
		}
		// DO CALLBACK FROM GITHUB

		// Get the type && Get the code
		_type, exist1 := ParamTools.GetParam(c, "type")
		Logger.Log(Logger.LogDebug, fmt.Sprintf("%v, %v", _type, exist1))
		if exist1 && _type == "github" {
			_code, exist2 := ParamTools.GetParam(c, "code")
			Logger.Log(Logger.LogDebug, fmt.Sprintf("%v, %v", _code, exist2))
			if exist2 {
				// Query FIN
				data := tokenData{
					ClientId:     Vars.ClientId,
					ClientSecret: Vars.ClientSecrets,
					Code:         _code,
				}
				payload, err := json.Marshal(data)
				response, err := Net.Request("POST", Vars.GithubTokenUrl, map[string]string{"Content-Type": "application/json"}, bytes.NewBuffer(payload))
				body, err := io.ReadAll(response.Body)
				Logger.Log(Logger.LogDebug, fmt.Sprintf("%v", err))
				Logger.Log(Logger.LogDebug, fmt.Sprintf("%s", body))
				token := ParamTools.ParseTokenBody(body)
				Logger.Log(Logger.LogDebug, fmt.Sprintf("%s, %s, %s", token.AccessToken, token.TokenType, token.Scope))

				headers := map[string]string{
					"Accept":        "application/vnd.github+json",
					"Authorization": "Bearer " + token.AccessToken,
				}
				response, err = Net.Get(Vars.GithubUserUrl, headers, map[string]string{})
				body, err = io.ReadAll(response.Body)
				Logger.Log(Logger.LogDebug, fmt.Sprintf("%s", body))
				c.JSON(http.StatusOK, gin.H{
					"token": token.AccessToken,
				})
				return
			}
		}
		// Query Failed
		QuickRes.BadRequest(c)
	})
}
