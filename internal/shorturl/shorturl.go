package shorturl

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func Encode(last int) string {
	short := make([]byte, 5)
	for step := 0; step < 5; step++ {
		short[4-step] = alphabet[last%len(alphabet)]
		last = last / len(alphabet)
	}
	return string(short)
}
