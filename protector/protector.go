package protector

import (
	"math/rand"
	"strconv"
	"unicode"
)

func Get_session_key() string {
	var result = ""
	for i := 0; i < 10; i++ {
		result += strconv.Itoa(rand.Intn(9) + 1)
	}
	return result
}

func Get_hash_str() string {
	var li = ""
	for i := 0; i < 5; i++ {
		li += strconv.Itoa(rand.Intn(6) + 1)
	}
	return li
}

type SessionProtector struct {
	hash string
}

func NewSessionProtector(hash_str string) *SessionProtector {
	sp := new(SessionProtector)
	sp.hash = hash_str
	return sp
}

func reverse(value string) string {
	data := []rune(value)
	result := []rune{}

	for i := len(data) - 1; i >= 0; i-- {
		result = append(result, data[i])
	}
	return string(result)
}

func (sp *SessionProtector) Next_session_key(session_key string) string {
	if sp.hash == "" {
		return ""
	}
	for i := 0; i < len(sp.hash); i++ {
		if !unicode.IsDigit(rune(sp.hash[i])) {
			return ""
		}
	}
	var result_int = 0
	var int_hash int
	var int_key int
	for i := 0; i < len(sp.hash); i++ {
		int_hash, _ = strconv.Atoi(string(sp.hash[i]))
		int_key, _ = strconv.Atoi(sp.calc_hash(session_key, int_hash))
		result_int += int_key
	}
	result := "0000000000" + strconv.Itoa(result_int)[0:10]
	return result[len(result)-10:]
}

func (sp *SessionProtector) calc_hash(session_key string, vall int) string {
	switch vall {
	case 1:
		var result = ""
		key_int, _ := strconv.Atoi(session_key[0:5])
		result = "00" + strconv.Itoa(key_int%97)
		return result[len(result)-2:]
	case 2:
		return reverse(session_key)
	case 3:
		return session_key[5:] + session_key[:5]
	case 4:
		var num = 0
		var int_key int
		for i := 1; i < len(session_key)-1; i++ {
			int_key, _ = strconv.Atoi(string(session_key[i]))
			num += int_key + 41
		}
		return strconv.Itoa(num)
	case 5:
		var num = 0
		for i := 0; i < len(session_key); i++ {
			num += int((session_key[i]) ^ 43)
		}
		return strconv.Itoa(num)
	}
	key_int, err := strconv.Atoi(session_key)
	if err != nil {
		return ""
	}
	return strconv.Itoa(key_int + vall)

}