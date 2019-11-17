package utils

import (
    "fmt"
)

func BytesToHex(b []byte) string {
    s := ""
    for _, v := range b {
        s += fmt.Sprintf("%02x", v)
    }
    return s
}

