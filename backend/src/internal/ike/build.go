package ike

func (ikeMessage *IKEMessage) BuildIKEHeader(
	initiatorSPI uint64,
	responsorSPI uint64,
	exchangeType uint8,
	flags uint8,
	messageID uint32) {
	ikeMessage.InitiatorSPI = initiatorSPI
	ikeMessage.ResponderSPI = responsorSPI
	ikeMessage.Version = 0x20
	ikeMessage.ExchangeType = exchangeType
	ikeMessage.Flags = flags
	ikeMessage.MessageID = messageID
}

func (container *IKEPayloadContainer) BUildKeyExchange(diffiehellmanGroup uint16, keyExchangeData []byte) {
	keyExchange := new(KeyExchange)
	keyExchange.DiffieHellmanGroup = diffiehellmanGroup
	keyExchange.KeyExchangeData = append(keyExchange.KeyExchangeData, keyExchangeData...)
	*container = append(*container, keyExchange)
}

func (container *IKEPayloadContainer) BuildNonce(nonceData []byte) {
	nonce := new(Nonce)
	nonce.NonceData = append(nonce.NonceData, nonceData...)
	*container = append(*container, nonce)
}

func (container *IKEPayloadContainer) BuildSecurityAssociation() *SecurityAssociation {
	securityAssociation := new(SecurityAssociation)
	*container = append(*container, securityAssociation)
	return securityAssociation
}

func (container *ProposalContainer) BuildProposal(proposalNumber uint8, protocolID uint8, spi []byte) *Proposal {
	proposal := new(Proposal)
	proposal.ProposalNumber = proposalNumber
	proposal.ProtocolID = protocolID
	proposal.SPI = append(proposal.SPI, spi...)
	*container = append(*container, proposal)
	return proposal
}

func (container *TransformContainer) BuildTransform(
	transformType uint8,
	transformID uint16,
	attributeType *uint16,
	attributeValue *uint16,
	variableLengthAttributeValue []byte) {
	transform := new(Transform)
	transform.TransformType = transformType
	transform.TransformID = transformID
	if attributeType != nil {
		transform.AttributePresent = true
		transform.AttributeType = *attributeType
		if attributeValue != nil {
			transform.AttributeFormat = AttributeFormatUseTV
			transform.AttributeValue = *attributeValue
		} else if len(variableLengthAttributeValue) != 0 {
			transform.AttributeFormat = AttributeFormatUseTLV
			transform.VariableLengthAttributeValue = append(transform.VariableLengthAttributeValue,
				variableLengthAttributeValue...)
		} else {
			return
		}
	} else {
		transform.AttributePresent = false
	}
	*container = append(*container, transform)
}
