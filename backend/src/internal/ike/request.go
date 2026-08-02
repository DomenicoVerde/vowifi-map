package ike

import (
	"errors"
	"math/big"
	"math/rand"
)

var errKeyExchangeGeneration = errors.New("BuildIKESAInitRequest(): failed to generate key exchange data")

// BuildIKESAInitRequest forges a standard IKE_SA_INIT Initiator Request,
// proposing the same algorithm set as the original vulnerability scanner
// (backend/src/main.go), without the DH-downgrade retry logic since here
// we are only probing for reachability.
func BuildIKESAInitRequest() ([]byte, error) {
	ikeInitiatorSPI := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	ikeMessage := new(IKEMessage)
	ikeMessage.BuildIKEHeader(ikeInitiatorSPI, 0, IKE_SA_INIT, InitiatorBitCheck, 0)

	securityAssociation := ikeMessage.Payloads.BuildSecurityAssociation()
	proposal := securityAssociation.Proposals.BuildProposal(1, TypeIKE, nil)

	var attributeType uint16 = AttributeTypeKeyLength
	var keyLength256 uint16 = 256
	var keyLength192 uint16 = 192
	var keyLength128 uint16 = 128
	proposal.EncryptionAlgorithm.BuildTransform(TypeEncryptionAlgorithm, ENCR_AES_CBC, &attributeType, &keyLength256, nil)
	proposal.EncryptionAlgorithm.BuildTransform(TypeEncryptionAlgorithm, ENCR_AES_CBC, &attributeType, &keyLength192, nil)
	proposal.EncryptionAlgorithm.BuildTransform(TypeEncryptionAlgorithm, ENCR_AES_CBC, &attributeType, &keyLength128, nil)
	proposal.EncryptionAlgorithm.BuildTransform(TypeEncryptionAlgorithm, ENCR_AES_CTR, &attributeType, &keyLength256, nil)
	proposal.EncryptionAlgorithm.BuildTransform(TypeEncryptionAlgorithm, ENCR_AES_CTR, &attributeType, &keyLength192, nil)
	proposal.EncryptionAlgorithm.BuildTransform(TypeEncryptionAlgorithm, ENCR_AES_CTR, &attributeType, &keyLength128, nil)

	proposal.IntegrityAlgorithm.BuildTransform(TypeIntegrityAlgorithm, AUTH_HMAC_SHA1_96, nil, nil, nil)
	proposal.IntegrityAlgorithm.BuildTransform(TypeIntegrityAlgorithm, AUTH_HMAC_MD5_96, nil, nil, nil)
	proposal.IntegrityAlgorithm.BuildTransform(TypeIntegrityAlgorithm, AUTH_HMAC_SHA2_256_128, nil, nil, nil)
	proposal.IntegrityAlgorithm.BuildTransform(TypeIntegrityAlgorithm, AUTH_HMAC_SHA2_512_256, nil, nil, nil)

	proposal.PseudorandomFunction.BuildTransform(TypePseudorandomFunction, PRF_HMAC_MD5, nil, nil, nil)
	proposal.PseudorandomFunction.BuildTransform(TypePseudorandomFunction, PRF_HMAC_SHA1, nil, nil, nil)
	proposal.PseudorandomFunction.BuildTransform(TypePseudorandomFunction, PRF_HMAC_SHA2_256, nil, nil, nil)
	proposal.PseudorandomFunction.BuildTransform(TypePseudorandomFunction, PRF_HMAC_SHA2_384, nil, nil, nil)
	proposal.PseudorandomFunction.BuildTransform(TypePseudorandomFunction, PRF_HMAC_SHA2_512, nil, nil, nil)

	proposal.DiffieHellmanGroup.BuildTransform(TypeDiffieHellmanGroup, DH_768_BIT_MODP, nil, nil, nil)
	proposal.DiffieHellmanGroup.BuildTransform(TypeDiffieHellmanGroup, DH_1024_BIT_MODP, nil, nil, nil)
	proposal.DiffieHellmanGroup.BuildTransform(TypeDiffieHellmanGroup, DH_1536_BIT_MODP, nil, nil, nil)
	proposal.DiffieHellmanGroup.BuildTransform(TypeDiffieHellmanGroup, DH_2048_BIT_MODP, nil, nil, nil)
	proposal.DiffieHellmanGroup.BuildTransform(TypeDiffieHellmanGroup, DH_3072_BIT_MODP, nil, nil, nil)
	proposal.DiffieHellmanGroup.BuildTransform(TypeDiffieHellmanGroup, DH_4096_BIT_MODP, nil, nil, nil)

	generator := new(big.Int).SetUint64(Group2Generator)
	factor, ok := new(big.Int).SetString(Group2PrimeString, 16)
	if !ok {
		return nil, errKeyExchangeGeneration
	}
	secret := GenerateRandomNumber()
	localPublicKeyExchangeValue := new(big.Int).Exp(generator, secret, factor).Bytes()
	prependZero := make([]byte, len(factor.Bytes())-len(localPublicKeyExchangeValue))
	localPublicKeyExchangeValue = append(prependZero, localPublicKeyExchangeValue...)
	ikeMessage.Payloads.BUildKeyExchange(DH_1024_BIT_MODP, localPublicKeyExchangeValue)

	localNonce := GenerateRandomNumber().Bytes()
	ikeMessage.Payloads.BuildNonce(localNonce)

	return ikeMessage.Encode()
}
