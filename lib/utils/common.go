package utils

import (
    "math/rand"
    "time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func InitializeRandomSeed(){
    rand.Seed(time.Now().UTC().UnixNano())
}

func GenerateRandomString(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}