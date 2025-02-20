package token

import "time"

type Maker interface {
	//创建token
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)

	//校验token是否有效
	VerifyToken(token string) (*Payload, error)
}
