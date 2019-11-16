package token

type TokenType string

const (
    // Lower command
    CMD_C = "c"
    CMD_D = "d"
    CMD_E = "e"
    CMD_F = "f"
    CMD_G = "g"
    CMD_A = "a"
    CMD_B = "b"
    CMD_V = "v"
    CMD_Q = "q"
    CMD_L = "l"
    CMD_PROGRAM = "@"
    // Upper command
    CMD_TRACK = "Track"
    CMD_CHANNEL = "Channel"
)

type Token struct {
    Type    TokenType
    Label   string
}



