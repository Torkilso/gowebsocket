package main


//Funnet på nett
/*
Første byte inneholder typ beskrivelse
Andre byte inneholder lengden på dataen fra(/til) klienten
either two or eight bytes if the length does not fit in the second byte (the second byte is then a code saying how many bytes are used for the length)
the actual (raw) data
 */
func decode (rawBytes []byte) string {
	var idxMask int
	if rawBytes[1] == 126 {
		idxMask = 4
	} else if rawBytes[1] == 127 {
		idxMask = 10
	} else {
		idxMask = 2
	}

	masks := rawBytes[idxMask:idxMask + 4]
	data := rawBytes[idxMask + 4:len(rawBytes)]
	decoded := make([]byte, len(rawBytes) - idxMask + 4)

	for i, b := range data {
		decoded[i] = b ^ masks[i % 4]
	}
	return string(decoded)
}

//Ikke testet 27.04.17
func encode (message string) (result []byte) {
	rawBytes := []byte(message)
	var idxData int

	length := byte(len(rawBytes))
	if len(rawBytes) <= 125 { //one byte to store data length
		result = make([]byte, len(rawBytes) + 2)
		result[1] = length
		idxData = 2
	} else if len(rawBytes) >= 126 && len(rawBytes) <= 65535 { //two bytes to store data length
		result = make([]byte, len(rawBytes) + 4)
		result[1] = 126 //extra storage needed
		result[2] = ( length >> 8 ) & 255
		result[3] = ( length      ) & 255
		idxData = 4
	} else {
		result = make([]byte, len(rawBytes) + 10)
		result[1] = 127
		result[2] = ( length >> 56 ) & 255
		result[3] = ( length >> 48 ) & 255
		result[4] = ( length >> 40 ) & 255
		result[5] = ( length >> 32 ) & 255
		result[6] = ( length >> 24 ) & 255
		result[7] = ( length >> 16 ) & 255
		result[8] = ( length >>  8 ) & 255
		result[9] = ( length       ) & 255
		idxData = 10
	}

	result[0] = 129 //only text is supported

	// put raw data at the correct index
	for i, b := range rawBytes {
		result[idxData + i] = b
	}
	return
}