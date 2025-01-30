package shortcode

import "math/rand"

type ShortCode struct {
	length int
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (s *ShortCode) GengerateShortCode() (string, error) {
	length := len(chars)
	result := make([]byte, s.length)

	for i := 0; i < s.length; i++ {
		result[i] = chars[rand.Intn(length)]
	}
	return string(result), nil
}
