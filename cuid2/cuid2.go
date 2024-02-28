package cuid2

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

var primeNumbers = []int{
	109717,
	109721,
	109741,
	109751,
	109789,
	109793,
	109807,
	109819,
	109829,
}

var alphabet_array = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

var machineFingerprint string

func init() {
	pid := strconv.Itoa(os.Getpid())
	hostname, _ := os.Hostname()
	acc := len(hostname) + 36
	for i := 0; i < len(hostname); i = i + 1 {
		acc = acc + int(hostname[i])
	}
	idBlock := padWithZero(pid, 2)
	nameBlock := padWithZero(strconv.Itoa(acc), 2)
	machineFingerprint = idBlock + nameBlock
}

var counter = math.MaxInt
var random_buffer_size = 4096
var random_buffer = make([]byte, random_buffer_size)
var random_buffer_index = random_buffer_size

func nextIntValue() int {
	if random_buffer_index == random_buffer_size {
		_, err := rand.Read(random_buffer)
		if err != nil {
			panic(err)
		}
		random_buffer_index = 0
	}
	result := int(random_buffer[random_buffer_index])<<24 | int(random_buffer[random_buffer_index+1]&0xff)<<16 | int(random_buffer[random_buffer_index+2]&0xff)<<8 | int(random_buffer[random_buffer_index+3]&0xff)
	return result
}

func nextCounterValue() int {
	if counter < math.MaxInt {
		counter++
	} else {
		counter = nextIntValue()
	}
	return counter
}

func createEntropy(length int) string {
	entropy := ""
	var primeNumber int
	for len(entropy) < length {
		primeNumber = primeNumbers[rand.Intn(len(primeNumbers))]
		entropy = fmt.Sprintf("%s%s", entropy, strconv.FormatUint(uint64(float64(primeNumber)*rand.Float64()), 36))
	}
	return entropy
}

func computeHash(content string, saltLength int) string {
	salt := createEntropy(saltLength)
	bytes := []byte(sha3.New256().Sum([]byte(content + salt)))
	bi := big.Int{}
	return bi.SetBytes(bytes).Text(36)
}

func padWithZero(str string, size int) string {
	paddedString := "000000000" + str
	return paddedString[:len(paddedString)-size]
}

type CUID2 string

func (c CUID2) String() string {
	return string(c)
}

// New returns a new CUID2 identifier with length of l
func New(l int) CUID2 {
	t := strconv.FormatInt(time.Now().UnixMilli(), 36)
	firstLetter := alphabet_array[int(math.Abs(float64(nextIntValue()%len(alphabet_array))))]
	content := fmt.Sprintf("%s%s%d%s", t, createEntropy(l), nextCounterValue(), machineFingerprint)
	hash := computeHash(content, l)
	return CUID2(string(firstLetter) + hash[1:l])
}
