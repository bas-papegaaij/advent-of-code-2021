package solvers

import (
	"fmt"
	"strconv"
	"strings"
)

type Day16 struct{}

type packetHeader struct {
	ver     int
	tID     int
	decoder *packetDecoder
}

type decoderPacket interface {
	version() int
	typeID() int
	decode()
	sumVersions() int
	evaluate() int64
}

func (d Day16) Solve(input []string, part int) (int64, error) {
	binStr, err := d.inputToBinaryString(input[0])
	if err != nil {
		return -1, err
	}

	decoder := packetDecoder{
		bitStream: binStr,
		ptr:       0,
	}
	packet := decoder.decode()

	switch part {
	case 1:
		return int64(packet.sumVersions()), nil
	case 2:
		return packet.evaluate(), nil
	default:
		return -1, invalidPartError(part)
	}
}

func (d Day16) inputToBinaryString(input string) (string, error) {
	var builder strings.Builder
	for _, char := range input {
		val, err := strconv.ParseUint(string(char), 16, 4)
		if err != nil {
			return "", err
		}

		builder.WriteString(fmt.Sprintf("%04b", val))
	}
	return builder.String(), nil
}

type packetDecoder struct {
	bitStream string
	ptr       int
}

func (d *packetDecoder) decode() decoderPacket {
	return d.decodePacket()
}

func (d *packetDecoder) read(bits int) string {
	res := d.bitStream[d.ptr : d.ptr+bits]
	d.ptr += bits
	return res
}

func (d *packetDecoder) decodePacket() decoderPacket {
	ver := d.readHeader()
	tID := d.readHeader()
	header := packetHeader{ver, tID, d}
	var newPacket decoderPacket
	if tID == 4 {
		newPacket = &valuePacket{
			packetHeader: header,
		}
	} else {
		newPacket = &operatorPacket{
			packetHeader: header,
		}
	}
	newPacket.decode()
	return newPacket
}

func (d *packetDecoder) readHeader() int {
	header := d.read(3)
	val, err := strconv.ParseInt(header, 2, 32)
	if err != nil {
		panic(err)
	}
	return int(val)
}

type valuePacket struct {
	packetHeader
	value int64
}

func (p *valuePacket) typeID() int {
	return p.tID
}

func (p *valuePacket) version() int {
	return p.ver
}

func (p *valuePacket) evaluate() int64 {
	return p.value
}

func (p *valuePacket) decode() {
	bitsRead := 6
	shouldContinue := true
	var valueBits strings.Builder
	for shouldContinue {
		continueBit := p.decoder.read(1)
		if continueBit == "0" {
			shouldContinue = false
		}
		valueBits.WriteString(p.decoder.read(4))
		bitsRead += 5
	}
	val, err := strconv.ParseInt(valueBits.String(), 2, 64)
	if err != nil {
		panic(fmt.Sprintf("bad value packet %v", err))
	}
	p.value = val
}

func (p *valuePacket) sumVersions() int {
	return p.ver
}

type operatorPacket struct {
	packetHeader
	subPackets []decoderPacket
}

func (p *operatorPacket) typeID() int {
	return p.tID
}

func (p *operatorPacket) version() int {
	return p.ver
}

func (p *operatorPacket) decode() {
	mode := p.decoder.read(1)
	switch mode {
	case "0":
		p.decodeMode1()
	case "1":
		p.decodeMode2()
	default:
		panic("invalid operator packet mode")
	}
}

func (p *operatorPacket) decodeMode1() {
	lengthBits := p.decoder.read(15)
	lengthVal, err := strconv.ParseInt(lengthBits, 2, 32)
	if err != nil {
		panic(err)
	}

	startPtr := p.decoder.ptr
	for p.decoder.ptr-startPtr < int(lengthVal) {
		p.subPackets = append(p.subPackets, p.decoder.decodePacket())
	}
}

func (p *operatorPacket) decodeMode2() {
	lengthBits := p.decoder.read(11)
	lengthVal, err := strconv.ParseInt(lengthBits, 2, 32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < int(lengthVal); i++ {
		p.subPackets = append(p.subPackets, p.decoder.decodePacket())
	}
}

func (p *operatorPacket) sumVersions() int {
	sum := p.ver
	for _, p := range p.subPackets {
		sum += p.sumVersions()
	}
	return sum
}

func (p *operatorPacket) evaluate() int64 {
	switch p.tID {
	case 0:
		return p.sum()
	case 1:
		return p.product()
	case 2:
		return p.minimum()
	case 3:
		return p.maximum()
	case 5:
		return p.gt()
	case 6:
		return p.lt()
	case 7:
		return p.eq()
	default:
		panic("bad operator type")
	}
}

func (p *operatorPacket) sum() int64 {
	var sum int64
	for _, p := range p.subPackets {
		sum += p.evaluate()
	}
	return sum
}

func (p *operatorPacket) product() int64 {
	product := p.subPackets[0].evaluate()
	for _, p := range p.subPackets[1:] {
		product *= p.evaluate()
	}
	return product
}
func (p *operatorPacket) minimum() int64 {
	min := p.subPackets[0].evaluate()
	for _, p := range p.subPackets[1:] {
		val := p.evaluate()
		if val < min {
			min = val
		}
	}

	return min
}

func (p *operatorPacket) maximum() int64 {
	max := p.subPackets[0].evaluate()
	for _, p := range p.subPackets[1:] {
		val := p.evaluate()
		if val > max {
			max = val
		}
	}

	return max
}

func (p *operatorPacket) gt() int64 {
	if p.subPackets[0].evaluate() > p.subPackets[1].evaluate() {
		return 1
	}
	return 0
}

func (p *operatorPacket) lt() int64 {
	if p.subPackets[0].evaluate() < p.subPackets[1].evaluate() {
		return 1
	}
	return 0
}

func (p *operatorPacket) eq() int64 {
	if p.subPackets[0].evaluate() == p.subPackets[1].evaluate() {
		return 1
	}
	return 0
}
