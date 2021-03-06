package main 

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
)

type NewTokenStruct struct {
	 FirstName string
	 LastName string
	 Email string
	 Phone string
}


func signToken(tokenStruct NewTokenStruct, key interface{}) (string, error) {

	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "firstName": tokenStruct.FirstName,
    "lastName": tokenStruct.LastName,
	})

	if out, err := token.SignedString(key); err == nil {
		fmt.Println(out)
		return out, nil
	} else {
		return "", fmt.Errorf("Error signing token: %v", err)
	}

}


func validateToken(tokenString string) error {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	    // Don't forget to validate the alg is what you expect:
	    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	    }

	    // hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	    return []byte("EB32ODSKJN234KJNDSKJSODF89N"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	    fmt.Println(claims["foo"], claims["nbf"])
	} else {
	    fmt.Println(err)
	}

	return nil
}
