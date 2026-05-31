package token

  import "time"

  // Maker 是 Token 创建的接口。支持 JWT 和 PASETO 两种实现
  type Maker interface {
        CreateToken(username string, role string, duration time.Duration, tokenType TokenType) (string, *Payload, error)
        VerifyToken(token string, tokenType TokenType) (*Payload, error)
  }  