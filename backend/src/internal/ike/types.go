package ike

// IKE types
type IKEPayloadType uint8

// Payload type numbers, as assigned by RFC 7296. Only the payload types
// that can actually appear in an IKE_SA_INIT exchange are kept: this
// scanner only ever sends/receives that one exchange, never completing the
// handshake, so payload types exclusive to IKE_AUTH/CREATE_CHILD_SA/
// INFORMATIONAL (encapsulated in an encrypted SK payload we can't decrypt
// anyway) don't need a representation here. Kept as explicit values
// (rather than iota) so removing the unused ones doesn't shift the wire
// value of the ones that remain.
const (
	NoNext      = 0
	TypeSA      = 33
	TypeKE      = 34
	TypeCERTreq = 38
	TypeNiNr    = 40
	TypeN       = 41
	TypeV       = 43
)

// used for SecurityAssociation-Proposal-Transform TransformType
const (
	TypeEncryptionAlgorithm = iota + 1
	TypePseudorandomFunction
	TypeIntegrityAlgorithm
	TypeDiffieHellmanGroup
	TypeExtendedSequenceNumbers
)

// used for SecurityAssociation-Proposal-Transform AttributeFormat
const (
	AttributeFormatUseTLV = iota
	AttributeFormatUseTV
)

// used for SecurityAssociation-Proposal-Trandform AttributeType
const (
	AttributeTypeKeyLength = 14
)

// used for SecurityAssociation-Proposal-Transform TransformID
const (
	ENCR_DES_IV64 = 1
	ENCR_DES      = 2
	ENCR_3DES     = 3
	ENCR_RC5      = 4
	ENCR_IDEA     = 5
	ENCR_CAST     = 6
	ENCR_BLOWFISH = 7
	ENCR_3IDEA    = 8
	ENCR_DES_IV32 = 9
	ENCR_NULL     = 11
	ENCR_AES_CBC  = 12
	ENCR_AES_CTR  = 13
)

const (
	PRF_HMAC_MD5 = iota + 1
	PRF_HMAC_SHA1
	PRF_HMAC_TIGER
	PRF_AES128_XCBC
	PRF_HMAC_SHA2_256
	PRF_HMAC_SHA2_384
	PRF_HMAC_SHA2_512
	PRF_AES128_CMAC
	PRF_HMAC_STREEBOG_512
)

const (
	AUTH_NONE = iota
	AUTH_HMAC_MD5_96
	AUTH_HMAC_SHA1_96
	AUTH_DES_MAC
	AUTH_KPDK_MD5
	AUTH_AES_XCBC_96
	AUTH_HMAC_MD5_128
	AUTH_HMAC_SHA1_160
	AUTH_AES_CMAC_96
	AUTH_AES_128_GMAC
	AUTH_AES_192_GMAC
	AUTH_AES_256_GMAC
	AUTH_HMAC_SHA2_256_128
	AUTH_HMAC_SHA2_384_192
	AUTH_HMAC_SHA2_512_256
)

const (
	DH_NONE          = 0
	DH_768_BIT_MODP  = 1
	DH_1024_BIT_MODP = 2
	DH_1536_BIT_MODP = 5
	DH_2048_BIT_MODP = iota + 10
	DH_3072_BIT_MODP
	DH_4096_BIT_MODP
	DH_6144_BIT_MODP
	DH_8192_BIT_MODP
)

// Exchange Type
const (
	IKE_SA_INIT = iota + 34
	IKE_AUTH
	CREATE_CHILD_SA
	INFORMATIONAL
)

// Protocol ID
const (
	TypeNone = iota
	TypeIKE
	TypeAH
	TypeESP
)

// Flags
const (
	ResponseBitCheck  = 0x20
	VersionBitCheck   = 0x10
	InitiatorBitCheck = 0x08
)
