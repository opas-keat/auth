package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"omsoft.com/auth/cmd/models"
	"omsoft.com/auth/cmd/util"
)

// type LoginController interface {
// 	// Eat(c *gin.Context)
// }

// type loginController struct {
// 	// dutyService service.DutyService
// }

type (
	MsgLogin       models.Login
	MsgToken       models.Token
	MsgTokenDetail models.TokenDetail
	MsgJWTContext  models.JWTContext
)

// func (lc *loginController) Login(c *fiber.Ctx) error {
// 	var l MsgLogin
// 	err := c.BodyParser(&l)
// 	if err != nil {
// 		// return failOnError(c, err, "cannot parse json", fiber.StatusBadRequest)
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
// 				"success": false,
// 				"message": err,
// 				"data":    nil,
// 			})
// 		}
// 	}
// 	println(l.Username)
// 	println(l.Password)
// 	return nil
// }
// type ErrorResponse struct {
//     FailedField string
//     Tag         string
//     Value       string
// }

// var validate = validator.New()
// func ValidateStruct(user models.Login) []*ErrorResponse {
//     var errors []*ErrorResponse
//     err := validate.Struct(user)
//     if err != nil {
//         for _, err := range err.(validator.ValidationErrors) {
//             var element ErrorResponse
//             element.FailedField = err.StructNamespace()
//             element.Tag = err.Tag()
//             element.Value = err.Param()
//             errors = append(errors, &element)
//         }
//     }
//     return errors
// }
func Login(c *fiber.Ctx) error {
	l := new(MsgLogin)
	if err := c.BodyParser(l); err != nil {
		println(err.Error())
		return util.FailOnError(c, err, "cannot parse json", fiber.StatusBadRequest, l)
	}
	// errors := ValidateStruct(*l)
	// if errors != nil {
	//    return c.Status(fiber.StatusBadRequest).JSON(errors)
	// }
	// tx := db.Debug()
	// 	tx = tx.Where("user_name = ?", auth.UserName)
	// 	var login models.Login
	// 	var user models.User
	// 	response := tx.First(&user)
	// 	if response.Error != nil {
	// 		println("Error occurred while retrieving user from the database: " + response.Error.Error())
	// 		login.Message = "ชื่อผู้ใช้งานหรือรหัสผ่าน ไม่ถูกต้อง"
	// 	}
	roles := []string{
		"admin",
		"report",
	}
	uid := "144479bd-fcdc-4c9f-b116-f2a08807a4c3"
	t, err := createToken(uid, "admin", "admin a", roles)
	if err != nil {
		return util.FailOnError(c, err, "StatusForbidden", fiber.StatusForbidden, t)
	}
	// t := new(MsgToken)
	return util.SuccessResponse(c, "", fiber.StatusOK, t)
}

func generateTokenBy(uid string, rdKey string, ctx interface{}, signed string, expire int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	location, _ := time.LoadLocation("Asia/Bangkok")
	claims["iss"] = "omsoft"
	claims["sub"] = uid
	claims["exp"] = expire
	claims["iat"] = time.Now().In(location).Unix()
	claims["context"] = ctx
	if ctx != nil {
		claims["access_uuid"] = rdKey
	} else {
		claims["refresh_uuid"] = rdKey
	}
	t, err := token.SignedString([]byte(signed))
	if err != nil {
		return "", err
	}
	return t, nil
}

func createToken(uid string, user string, fullName string, roles []string) (*MsgTokenDetail, error) {
	tokenDetail := &MsgTokenDetail{}
	location, _ := time.LoadLocation("Asia/Bangkok")
	at := time.Now().In(location).Add(time.Minute * 1).Unix()
	rt := time.Now().In(location).Add(time.Hour * 24 * 7).Unix()
	tokenDetail.AccessTokenExp = at
	tokenDetail.RefreshTokenExp = rt
	tokenDetail.AccessUUid = utils.UUIDv4()
	tokenDetail.RefreshUUid = utils.UUIDv4()
	var err error
	/////mock data/////
	ctAccess := MsgJWTContext{}
	ctAccess.User = user
	ctAccess.DisplayName = fullName
	ctAccess.Roles = roles
	tokenDetail.Context = ctAccess
	/////mock data/////
	tokenDetail.Token.AccessToken, err = generateTokenBy(uid, tokenDetail.AccessUUid, ctAccess, "OMSoft", tokenDetail.AccessTokenExp)
	if err != nil {
		return tokenDetail, err
	}
	tokenDetail.Token.RefreshToken, err = generateTokenBy(uid, tokenDetail.RefreshUUid, nil, "OMSoftR", tokenDetail.RefreshTokenExp)
	if err != nil {
		return tokenDetail, err
	}
	return tokenDetail, nil
}

func AuthorizationRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SuccessHandler: AuthSuccess,
		ErrorHandler:   AuthError,
		SigningKey:     []byte("OMSoft"),
		SigningMethod:  "HS256",
		TokenLookup:    "header:Authorization",
		AuthScheme:     "Bearer",
	})
}

func AuthError(c *fiber.Ctx, e error) error {
	// return util.FailOnError(c, err, "cannot parse json", fiber.StatusUnauthorized, l)
	return util.FailOnError(c, e, "Unauthorized", fiber.StatusUnauthorized, new(MsgToken))
}

func AuthSuccess(c *fiber.Ctx) error {
	c.Next()
	return nil
}
