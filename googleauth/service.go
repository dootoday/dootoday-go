package service

import (
	"net/http"

	"google.golang.org/api/oauth2/v2"
)

var httpClient = &http.Client{}

// GoogleAuthService :
type GoogleAuthService struct{}

// NewGoogleAuthService :
func NewGoogleAuthService() *GoogleAuthService {
	return &GoogleAuthService{}
}

// VerifyIDToken : function to verify id token recieved from client
func (g *GoogleAuthService) VerifyIDToken(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}
