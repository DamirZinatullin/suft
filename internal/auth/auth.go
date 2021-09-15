package auth

type Token struct {
	AccessToken string
	RefreshToken string
}

func Authenticate(email string, password string)(Token, error){
	return Token{}, nil
}
