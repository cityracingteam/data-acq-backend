package jwt

import (
	"errors"
	"time"

	"github.com/cityracingteam/data-acq-backend/models"
	"github.com/cityracingteam/data-acq-backend/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func generateRelativeExpiryTime(offset ExpiryOffset) time.Time {
	time := time.Now()
	time = time.AddDate(offset.years, offset.months, offset.days)
	return time
}

func newJwt(user models.User, offset ExpiryOffset, extraClaims ExtraClaims) (str string, err error) {
	claims := Claims{
		jwt.RegisteredClaims{
			Issuer:    "fosscat",
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(generateRelativeExpiryTime(offset)),
		},
		extraClaims, // Extra claims not in jwt.RegisteredClaims
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES384, claims)

	str, err = token.SignedString(key)

	return // returns str, err

}

func NewAccessJwt(user models.User) (str string, err error) {
	extraClaims := ExtraClaims{
		KeyID: "access",
	}

	str, err = newJwt(user, ExpiryOffset{days: 1}, extraClaims)
	return // str, err
}

// NOTE: We don't need to store the refresh token in the database, instead we check to see that we issued it and that it has not expired.
func NewRefreshJwt(user models.User) (str string, err error) {
	extraClaims := ExtraClaims{
		KeyID: "refresh",
	}

	// Create the refresh token
	str, err = newJwt(user, ExpiryOffset{months: 12}, extraClaims)
	return // str, err
}

// Returns token, claims and nil error if jwt can be verified
func VerifyJwt(jwtStr string) (token *jwt.Token, claims *Claims, err error) {
	token, err = jwt.ParseWithClaims(jwtStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return &key.PublicKey, nil
	})

	if token == nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return token, claims, nil
	} else {
		return nil, nil, err
	}
}

func GetUserForJwt(db *gorm.DB, token *jwt.Token) (*models.User, error) {
	uid, err := token.Claims.GetSubject()
	if err == nil {
		uuid, err2 := uuid.Parse(uid)
		if err2 == nil {
			// uid is a valid UUID

			user, err3 := util.GetObjectById[models.User](db, uuid)

			if err3 == nil {
				// User found sucessfully.
				return user, nil
			} else {
				// user does not exist, or other error fetching
				return nil, errors.New("unable to fetch user from database by user id")
			}
		} else {
			// uid is not a valid UUID
			return nil, errors.New("subject was not a valid user id")
		}
	} else {
		// subject was not present in jwt
		return nil, errors.New("subject not present in jwt")
	}
}
