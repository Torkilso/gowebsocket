package websocket

func decode (inputBytes []byte) string {
	mask := 2
	if inputBytes[1]-128 == 126 {
		mask = 4
	} else if inputBytes[1]-128 == 127 {
		mask = 10
	}

	masks := inputBytes[mask:mask + 4]

  dataEnd := mask + 4

  for r := mask + 4; r <= len(inputBytes); r++ {
    if inputBytes[r] == 0 {
      check := true
      for t := 1; t <= 10; t++ {
        if inputBytes[r+t] != 0 {
          check = false
          break
        }
      }
      if check {
        dataEnd = r
        break
      }
    }
  }

  data := inputBytes[mask + 4:dataEnd]
	decoded := make([]byte, len(data))

	for i, b := range data {
		decoded[i] = b ^ masks[i % 4]
	}
	return string(decoded)
}

func encode (message string) (result []byte) {

	input := []byte(message)
	var dataIndex int

	length := byte(len(input))

	if len(input) <= 125 { //one byte to store data length
		result = make([]byte, len(input)+2)
		result[1] = length
		dataIndex = 2
	} else if len(input) >= 126 && len(input) <= 65535 { //two bytes to store data length
		result = make([]byte, len(input)+4)

		result[1] = 126 //extra storage needed
		result[2] = byte(len(input) >> 8)
		result[3] = length

		dataIndex = 4
	} else {
		result = make([]byte, len(input)+10)
		result[1] = 127
		result[2] = byte(len(input) >> 56)
		result[3] = byte(len(input) >> 48)
		result[4] = byte(len(input) >> 40)
		result[5] = byte(len(input) >> 32)
		result[6] = byte(len(input) >> 24)
		result[7] = byte(len(input) >> 16)
		result[8] = byte(len(input) >> 8)
		result[9] = length
		dataIndex = 10
	}

	result[0] = 129

	// put data at the correct index
	for i, b := range input {
		result[dataIndex+i] = b
	}
	return result
}
