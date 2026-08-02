package ike

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

type IKEMessage struct {
	InitiatorSPI uint64
	ResponderSPI uint64
	Version      uint8
	ExchangeType uint8
	Flags        uint8
	MessageID    uint32
	Payloads     IKEPayloadContainer
}

func (ikeMessage *IKEMessage) Encode() ([]byte, error) {
	log.Println("Encoding IKE message")

	ikeMessageData := make([]byte, 28)

	binary.BigEndian.PutUint64(ikeMessageData[0:8], ikeMessage.InitiatorSPI)
	binary.BigEndian.PutUint64(ikeMessageData[8:16], ikeMessage.ResponderSPI)
	ikeMessageData[17] = ikeMessage.Version
	ikeMessageData[18] = ikeMessage.ExchangeType
	ikeMessageData[19] = ikeMessage.Flags
	binary.BigEndian.PutUint32(ikeMessageData[20:24], ikeMessage.MessageID)

	if len(ikeMessage.Payloads) > 0 {
		ikeMessageData[16] = byte(ikeMessage.Payloads[0].Type())
	} else {
		ikeMessageData[16] = NoNext
	}

	ikeMessagePayloadData, err := ikeMessage.Payloads.Encode()
	if err != nil {
		return nil, fmt.Errorf("Encode(): EncodePayload failed: %+v", err)
	}

	ikeMessageData = append(ikeMessageData, ikeMessagePayloadData...)
	binary.BigEndian.PutUint32(ikeMessageData[24:28], uint32(len(ikeMessageData)))

	return ikeMessageData, nil
}

func (ikeMessage *IKEMessage) Decode(rawData []byte) error {
	// IKE message packet format this implementation referenced is
	// defined in RFC 7296, Section 3.1
	log.Println("Decoding IKE message")

	// bounds checking
	if len(rawData) < 28 {
		return errors.New("Decode(): Received broken IKE header")
	}
	ikeMessageLength := binary.BigEndian.Uint32(rawData[24:28])
	if ikeMessageLength < 28 {
		return fmt.Errorf("Decode(): Illegal IKE message length %d < header length 20", ikeMessageLength)
	}
	// len() return int, which is 64 bit on 64-bit host and 32 bit
	// on 32-bit host, so this implementation may potentially cause
	// problem on 32-bit machine
	if len(rawData) != int(ikeMessageLength) {
		return errors.New("Decode(): The length of received message not matchs the length specified in header")
	}

	nextPayload := rawData[16]

	ikeMessage.InitiatorSPI = binary.BigEndian.Uint64(rawData[:8])
	ikeMessage.ResponderSPI = binary.BigEndian.Uint64(rawData[8:16])
	ikeMessage.Version = rawData[17]
	ikeMessage.ExchangeType = rawData[18]
	ikeMessage.Flags = rawData[19]
	ikeMessage.MessageID = binary.BigEndian.Uint32(rawData[20:24])

	err := ikeMessage.Payloads.Decode(nextPayload, rawData[28:])
	if err != nil {
		return fmt.Errorf("Decode(): DecodePayload failed: %+v", err)
	}

	return nil
}

type IKEPayloadContainer []IKEPayload

func (container *IKEPayloadContainer) Encode() ([]byte, error) {
	log.Println("Encoding IKE payloads")

	ikeMessagePayloadData := make([]byte, 0)

	for index, payload := range *container {
		payloadData := make([]byte, 4)     // IKE payload general header
		if (index + 1) < len(*container) { // if it has next payload
			payloadData[0] = uint8((*container)[index+1].Type())
		} else {
			payloadData[0] = NoNext
		}

		data, err := payload.marshal()
		if err != nil {
			return nil, fmt.Errorf("EncodePayload(): Failed to marshal payload: %+v", err)
		}

		payloadData = append(payloadData, data...)
		binary.BigEndian.PutUint16(payloadData[2:4], uint16(len(payloadData)))

		ikeMessagePayloadData = append(ikeMessagePayloadData, payloadData...)
	}

	return ikeMessagePayloadData, nil
}

func (container *IKEPayloadContainer) Decode(nextPayload uint8, rawData []byte) error {
	log.Println("Decoding IKE payloads")

	for len(rawData) > 0 {
		// bounds checking
		//Trace("DecodePayload(): Decode 1 payload")
		if len(rawData) < 4 {
			return errors.New("DecodePayload(): No sufficient bytes to decode next payload")
		}
		payloadLength := binary.BigEndian.Uint16(rawData[2:4])
		if payloadLength < 4 {
			return fmt.Errorf("DecodePayload(): Illegal payload length %d < header length 4", payloadLength)
		}
		if len(rawData) < int(payloadLength) {
			return errors.New("DecodePayload(): The length of received message not matchs the length specified in header")
		}

		criticalBit := (rawData[1] & 0x80) >> 7

		var payload IKEPayload

		switch nextPayload {
		case TypeSA:
			payload = new(SecurityAssociation)
		case TypeKE:
			payload = new(KeyExchange)
		case TypeCERTreq:
			payload = new(CertificateRequest)
		case TypeNiNr:
			payload = new(Nonce)
		case TypeN:
			payload = new(Notification)
		case TypeV:
			payload = new(VendorID)
		default:
			if criticalBit == 0 {
				// Skip this payload
				nextPayload = rawData[0]
				rawData = rawData[payloadLength:]
				continue
			} else {
				// TODO: Reject this IKE message
				return fmt.Errorf("Unknown payload type: %d", nextPayload)
			}
		}

		if err := payload.unmarshal(rawData[4:payloadLength]); err != nil {
			return fmt.Errorf("DecodePayload(): Unmarshal payload failed: %+v", err)
		}

		*container = append(*container, payload)

		nextPayload = rawData[0]
		rawData = rawData[payloadLength:]
	}

	return nil
}

type IKEPayload interface {
	// Type specifies the IKE payload types
	Type() IKEPayloadType

	// Called by Encode() or Decode()
	marshal() ([]byte, error)
	unmarshal(rawData []byte) error
}

// Definition of Security Association

var _ IKEPayload = &SecurityAssociation{}

type SecurityAssociation struct {
	Proposals ProposalContainer
}

type ProposalContainer []*Proposal

type Proposal struct {
	ProposalNumber          uint8
	ProtocolID              uint8
	SPI                     []byte
	EncryptionAlgorithm     TransformContainer
	PseudorandomFunction    TransformContainer
	IntegrityAlgorithm      TransformContainer
	DiffieHellmanGroup      TransformContainer
	ExtendedSequenceNumbers TransformContainer
}

type TransformContainer []*Transform

type Transform struct {
	TransformType                uint8
	TransformID                  uint16
	AttributePresent             bool
	AttributeFormat              uint8
	AttributeType                uint16
	AttributeValue               uint16
	VariableLengthAttributeValue []byte
}

func (securityAssociation *SecurityAssociation) Type() IKEPayloadType { return TypeSA }

func (securityAssociation *SecurityAssociation) marshal() ([]byte, error) {
	//Info("[SecurityAssociation] marshal(): Start marshalling")

	securityAssociationData := make([]byte, 0)

	for proposalIndex, proposal := range securityAssociation.Proposals {
		proposalData := make([]byte, 8)

		if (proposalIndex + 1) < len(securityAssociation.Proposals) {
			proposalData[0] = 2
		} else {
			proposalData[0] = 0
		}

		proposalData[4] = proposal.ProposalNumber
		proposalData[5] = proposal.ProtocolID

		proposalData[6] = uint8(len(proposal.SPI))
		if len(proposal.SPI) > 0 {
			proposalData = append(proposalData, proposal.SPI...)
		}

		// combine all transforms
		var transformList []*Transform
		transformList = append(transformList, proposal.EncryptionAlgorithm...)
		transformList = append(transformList, proposal.PseudorandomFunction...)
		transformList = append(transformList, proposal.IntegrityAlgorithm...)
		transformList = append(transformList, proposal.DiffieHellmanGroup...)
		transformList = append(transformList, proposal.ExtendedSequenceNumbers...)

		if len(transformList) == 0 {
			return nil, errors.New("One proposal has no any transform")
		}
		proposalData[7] = uint8(len(transformList))

		proposalTransformData := make([]byte, 0)

		for transformIndex, transform := range transformList {
			transformData := make([]byte, 8)

			if (transformIndex + 1) < len(transformList) {
				transformData[0] = 3
			} else {
				transformData[0] = 0
			}

			transformData[4] = transform.TransformType
			binary.BigEndian.PutUint16(transformData[6:8], transform.TransformID)

			if transform.AttributePresent {
				attributeData := make([]byte, 4)

				if transform.AttributeFormat == 0 {
					// TLV
					if len(transform.VariableLengthAttributeValue) == 0 {
						return nil, errors.New("Attribute of one transform not specified")
					}
					attributeFormatAndType := ((uint16(transform.AttributeFormat) & 0x1) << 15) | transform.AttributeType
					binary.BigEndian.PutUint16(attributeData[0:2], attributeFormatAndType)
					binary.BigEndian.PutUint16(attributeData[2:4], uint16(len(transform.VariableLengthAttributeValue)))
					attributeData = append(attributeData, transform.VariableLengthAttributeValue...)
				} else {
					// TV
					attributeFormatAndType := ((uint16(transform.AttributeFormat) & 0x1) << 15) | transform.AttributeType
					binary.BigEndian.PutUint16(attributeData[0:2], attributeFormatAndType)
					binary.BigEndian.PutUint16(attributeData[2:4], transform.AttributeValue)
				}

				transformData = append(transformData, attributeData...)
			}

			binary.BigEndian.PutUint16(transformData[2:4], uint16(len(transformData)))

			proposalTransformData = append(proposalTransformData, transformData...)
		}

		proposalData = append(proposalData, proposalTransformData...)
		binary.BigEndian.PutUint16(proposalData[2:4], uint16(len(proposalData)))

		securityAssociationData = append(securityAssociationData, proposalData...)
	}

	return securityAssociationData, nil
}

func (securityAssociation *SecurityAssociation) unmarshal(rawData []byte) error {
	//Info("[SecurityAssociation] unmarshal(): Start unmarshalling received bytes")
	//Tracef("[SecurityAssociation] unmarshal(): Payload length %d bytes", len(rawData))

	for len(rawData) > 0 {
		//Trace("[SecurityAssociation] unmarshal(): Unmarshal 1 proposal")
		// bounds checking
		if len(rawData) < 8 {
			return errors.New("Proposal: No sufficient bytes to decode next proposal")
		}
		proposalLength := binary.BigEndian.Uint16(rawData[2:4])
		if proposalLength < 8 {
			return errors.New("Proposal: Illegal payload length %d < header length 8")
		}
		if len(rawData) < int(proposalLength) {
			return errors.New("Proposal: The length of received message not matchs the length specified in header")
		}

		// Log whether this proposal is the last
		if rawData[0] == 0 {
			//Trace("[SecurityAssociation] This proposal is the last")
		}
		// Log the number of transform in the proposal
		//Tracef("[SecurityAssociation] This proposal contained %d transform", rawData[7])

		proposal := new(Proposal)
		var transformData []byte

		proposal.ProposalNumber = rawData[4]
		proposal.ProtocolID = rawData[5]

		spiSize := rawData[6]
		if spiSize > 0 {
			// bounds checking
			if len(rawData) < int(8+spiSize) {
				return errors.New("Proposal: No sufficient bytes for unmarshalling SPI of proposal")
			}
			proposal.SPI = append(proposal.SPI, rawData[8:8+spiSize]...)
		}

		transformData = rawData[8+spiSize : proposalLength]

		for len(transformData) > 0 {
			// bounds checking
			//Trace("[SecurityAssociation] unmarshal(): Unmarshal 1 transform")
			if len(transformData) < 8 {
				return errors.New("Transform: No sufficient bytes to decode next transform")
			}
			transformLength := binary.BigEndian.Uint16(transformData[2:4])
			if transformLength < 8 {
				return errors.New("Transform: Illegal payload length %d < header length 8")
			}
			if len(transformData) < int(transformLength) {
				return errors.New("Transform: The length of received message not matchs the length specified in header")
			}

			// Log whether this transform is the last
			if transformData[0] == 0 {
				//Trace("[SecurityAssociation] This transform is the last")
			}

			transform := new(Transform)

			transform.TransformType = transformData[4]
			transform.TransformID = binary.BigEndian.Uint16(transformData[6:8])
			if transformLength > 8 {
				transform.AttributePresent = true
				transform.AttributeFormat = ((transformData[8] & 0x80) >> 7)
				transform.AttributeType = binary.BigEndian.Uint16(transformData[8:10]) & 0x7f

				if transform.AttributeFormat == 0 {
					attributeLength := binary.BigEndian.Uint16(transformData[10:12])
					// bounds checking
					if (12 + attributeLength) != transformLength {
						return fmt.Errorf("Illegal attribute length %d not satisfies the transform length %d",
							attributeLength, transformLength)
					}
					copy(transform.VariableLengthAttributeValue, transformData[12:12+attributeLength])
				} else {
					transform.AttributeValue = binary.BigEndian.Uint16(transformData[10:12])
				}
			}

			switch transform.TransformType {
			case TypeEncryptionAlgorithm:
				proposal.EncryptionAlgorithm = append(proposal.EncryptionAlgorithm, transform)
			case TypePseudorandomFunction:
				proposal.PseudorandomFunction = append(proposal.PseudorandomFunction, transform)
			case TypeIntegrityAlgorithm:
				proposal.IntegrityAlgorithm = append(proposal.IntegrityAlgorithm, transform)
			case TypeDiffieHellmanGroup:
				proposal.DiffieHellmanGroup = append(proposal.DiffieHellmanGroup, transform)
			case TypeExtendedSequenceNumbers:
				proposal.ExtendedSequenceNumbers = append(proposal.ExtendedSequenceNumbers, transform)
			}

			transformData = transformData[transformLength:]
		}

		securityAssociation.Proposals = append(securityAssociation.Proposals, proposal)

		rawData = rawData[proposalLength:]
	}

	return nil
}

// Definition of Key Exchange

var _ IKEPayload = &KeyExchange{}

type KeyExchange struct {
	DiffieHellmanGroup uint16
	KeyExchangeData    []byte
}

func (keyExchange *KeyExchange) Type() IKEPayloadType { return TypeKE }

func (keyExchange *KeyExchange) marshal() ([]byte, error) {
	//Info("[KeyExchange] marshal(): Start marshalling")

	keyExchangeData := make([]byte, 4)

	binary.BigEndian.PutUint16(keyExchangeData[0:2], keyExchange.DiffieHellmanGroup)
	keyExchangeData = append(keyExchangeData, keyExchange.KeyExchangeData...)

	return keyExchangeData, nil
}

func (keyExchange *KeyExchange) unmarshal(rawData []byte) error {
	//Info("[KeyExchange] unmarshal(): Start unmarshalling received bytes")
	//Tracef("[KeyExchange] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		//Trace("[KeyExchange] unmarshal(): Unmarshal 1 key exchange data")
		// bounds checking
		if len(rawData) <= 4 {
			return errors.New("KeyExchange: No sufficient bytes to decode next key exchange data")
		}

		keyExchange.DiffieHellmanGroup = binary.BigEndian.Uint16(rawData[0:2])
		keyExchange.KeyExchangeData = append(keyExchange.KeyExchangeData, rawData[4:]...)
	}

	return nil
}

// Definition of Certificate Request

var _ IKEPayload = &CertificateRequest{}

type CertificateRequest struct {
	CertificateEncoding    uint8
	CertificationAuthority []byte
}

func (certificateRequest *CertificateRequest) Type() IKEPayloadType { return TypeCERTreq }

func (certificateRequest *CertificateRequest) marshal() ([]byte, error) {
	//Info("[CertificateRequest] marshal(): Start marshalling")

	certificateRequestData := make([]byte, 1)

	certificateRequestData[0] = certificateRequest.CertificateEncoding
	certificateRequestData = append(certificateRequestData, certificateRequest.CertificationAuthority...)

	return certificateRequestData, nil
}

func (certificateRequest *CertificateRequest) unmarshal(rawData []byte) error {
	//Info("[CertificateRequest] unmarshal(): Start unmarshalling received bytes")
	//Tracef("[CertificateRequest] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		//Trace("[CertificateRequest] unmarshal(): Unmarshal 1 certificate request")
		// bounds checking
		if len(rawData) <= 1 {
			return errors.New("CertificateRequest: No sufficient bytes to decode next certificate request")
		}

		certificateRequest.CertificateEncoding = rawData[0]
		certificateRequest.CertificationAuthority = append(certificateRequest.CertificationAuthority, rawData[1:]...)
	}

	return nil
}

// Definition of Nonce

var _ IKEPayload = &Nonce{}

type Nonce struct {
	NonceData []byte
}

func (nonce *Nonce) Type() IKEPayloadType { return TypeNiNr }

func (nonce *Nonce) marshal() ([]byte, error) {
	//Info("[Nonce] marshal(): Start marshalling")

	nonceData := make([]byte, 0)
	nonceData = append(nonceData, nonce.NonceData...)

	return nonceData, nil
}

func (nonce *Nonce) unmarshal(rawData []byte) error {
	//Info("[Nonce] unmarshal(): Start unmarshalling received bytes")
	//Tracef("[Nonce] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		//Trace("[Nonce] unmarshal(): Unmarshal 1 nonce")
		nonce.NonceData = append(nonce.NonceData, rawData...)
	}

	return nil
}

// Definition of Notification

var _ IKEPayload = &Notification{}

type Notification struct {
	ProtocolID        uint8
	NotifyMessageType uint16
	SPI               []byte
	NotificationData  []byte
}

func (notification *Notification) Type() IKEPayloadType { return TypeN }

func (notification *Notification) marshal() ([]byte, error) {
	//Info("[Notification] marshal(): Start marshalling")

	notificationData := make([]byte, 4)

	notificationData[0] = notification.ProtocolID
	notificationData[1] = uint8(len(notification.SPI))
	binary.BigEndian.PutUint16(notificationData[2:4], notification.NotifyMessageType)

	notificationData = append(notificationData, notification.SPI...)
	notificationData = append(notificationData, notification.NotificationData...)

	return notificationData, nil
}

func (notification *Notification) unmarshal(rawData []byte) error {
	//Info("[Notification] unmarshal(): Start unmarshalling received bytes")
	//Tracef("[Notification] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		//Trace("[Notification] unmarshal(): Unmarshal 1 notification")
		// bounds checking
		if len(rawData) < 4 {
			return errors.New("Notification: No sufficient bytes to decode next notification")
		}
		spiSize := rawData[1]
		if len(rawData) < int(4+spiSize) {
			return errors.New("Notification: No sufficient bytes to get SPI according to the length specified in header")
		}

		notification.ProtocolID = rawData[0]
		notification.NotifyMessageType = binary.BigEndian.Uint16(rawData[2:4])

		notification.SPI = append(notification.SPI, rawData[4:4+spiSize]...)
		notification.NotificationData = append(notification.NotificationData, rawData[4+spiSize:]...)
	}

	return nil
}

// Definition of Vendor ID

var _ IKEPayload = &VendorID{}

type VendorID struct {
	VendorIDData []byte
}

func (vendorID *VendorID) Type() IKEPayloadType { return TypeV }

func (vendorID *VendorID) marshal() ([]byte, error) {
	//Info("[VendorID] marshal(): Start marshalling")
	return vendorID.VendorIDData, nil
}

func (vendorID *VendorID) unmarshal(rawData []byte) error {
	//Info("[VendorID] unmarshal(): Start unmarshalling received bytes")
	//Tracef("[VendorID] unmarshal(): Payload length %d bytes", len(rawData))

	if len(rawData) > 0 {
		//Trace("[VendorID] unmarshal(): Unmarshal 1 vendor ID")
		vendorID.VendorIDData = append(vendorID.VendorIDData, rawData...)
	}

	return nil
}
