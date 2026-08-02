package ike

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"
)

var (
	randomNumberMaximum big.Int
	randomNumberMinimum big.Int
)

func init() {
	randomNumberMaximum.SetString(strings.Repeat("F", 512), 16)
	randomNumberMinimum.SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
}

// GenerateRandomNumber returns a large random number suitable for use as an
// IKE nonce or Diffie-Hellman secret.
func GenerateRandomNumber() *big.Int {
	var number *big.Int
	var err error
	for {
		number, err = rand.Int(rand.Reader, &randomNumberMaximum)
		if err != nil {
			log.Panicf("Error occurs when generating random number: %+v", err)
			return nil
		}
		if number.Cmp(&randomNumberMinimum) == 1 {
			break
		}
	}
	return number
}

// Diffie-Hellman group 2 (RFC 2409, 1024-bit MODP), used for the IKE_SA_INIT
// key exchange payload.
const (
	Group2PrimeString string = "FFFFFFFFFFFFFFFFC90FDAA22168C234" +
		"C4C6628B80DC1CD129024E088A67CC74" +
		"020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F1437" +
		"4FE1356D6D51C245E485B576625E7EC6" +
		"F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE6" +
		"49286651ECE65381FFFFFFFFFFFFFFFF"
	Group2Generator = 2
)
