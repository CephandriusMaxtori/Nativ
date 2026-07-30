//go:build ignore

package vulkan

import "encoding/binary"

type spvBuilder struct {
	words []uint32
	bound uint32
	nid   uint32
}

func newSPV() *spvBuilder { return &spvBuilder{nid: 1} }

func (b *spvBuilder) id() uint32 {
	b.nid++
	if b.nid > b.bound {
		b.bound = b.nid
	}
	return b.nid
}

func (b *spvBuilder) op(opcode uint32, args ...uint32) {
	wc := uint32(1 + len(args))
	b.words = append(b.words, (wc<<16)|opcode)
	b.words = append(b.words, args...)
}

func (b *spvBuilder) str(str string) []uint32 {
	data := append([]byte(str), 0)
	for len(data)%4 != 0 {
		data = append(data, 0)
	}
	out := make([]uint32, len(data)/4)
	for i := range out {
		out[i] = binary.LittleEndian.Uint32(data[i*4:])
	}
	return out
}

func (b *spvBuilder) opStr(opcode uint32, target uint32, str string) {
	s := b.str(str)
	wordCount := uint32(1 + 1 + len(s))
	b.words = append(b.words, (wordCount<<16)|opcode, target)
	b.words = append(b.words, s...)
}

func (b *spvBuilder) opEntry(execModel, funcID uint32, name string, ifaces ...uint32) {
	s := b.str(name)
	wc := uint32(1 + 1 + 1 + len(s) + uint32(len(ifaces)))
	b.words = append(b.words, (wc<<16)|15, execModel, funcID)
	b.words = append(b.words, s...)
	b.words = append(b.words, ifaces...)
}

func (b *spvBuilder) build() []byte {
	if b.bound == 0 {
		b.bound = b.nid
	}
	n := uint32(5 + uint32(len(b.words)))
	buf := make([]byte, n*4)
	binary.LittleEndian.PutUint32(buf[0:], 0x07230203)
	binary.LittleEndian.PutUint32(buf[4:], 0x00010000)
	binary.LittleEndian.PutUint32(buf[8:], 0)
	binary.LittleEndian.PutUint32(buf[12:], b.bound)
	binary.LittleEndian.PutUint32(buf[16:], 0)
	for i, w := range b.words {
		binary.LittleEndian.PutUint32(buf[20+i*4:], w)
	}
	return buf
}

func vertexShaderSPIRV() []byte {
	b := newSPV()

	voidT := b.id()
	floatT := b.id()
	vec2T := b.id()
	vec4T := b.id()
	voidFnT := b.id()
	inV2P := b.id()
	inV4P := b.id()
	outV4P := b.id()

	posVar := b.id()
	colVar := b.id()
	fragVar := b.id()
	glPosVar := b.id()

	c0 := b.id()
	c1 := b.id()

	mainFn := b.id()
	entry := b.id()
	pv := b.id()
	cv := b.id()
	px := b.id()
	py := b.id()
	gp := b.id()

	b.op(17, 1)
	b.op(14, 0, 0)
	b.opEntry(0, mainFn, "main", posVar, colVar, fragVar, glPosVar)
	b.op(19, 4, 0)
	b.opStr(5, mainFn, "main")
	b.opStr(5, posVar, "position")
	b.opStr(5, colVar, "inColor")
	b.opStr(5, fragVar, "fragColor")
	b.opStr(5, glPosVar, "gl_Position")
	b.op(71, posVar, 30, 0)
	b.op(71, colVar, 30, 1)
	b.op(71, fragVar, 30, 0)
	b.op(71, glPosVar, 11, 0)

	b.op(19, voidT)
	b.op(33, voidFnT, voidT)
	b.op(22, floatT, 32)
	b.op(23, vec2T, floatT, 2)
	b.op(23, vec4T, floatT, 4)
	b.op(32, inV2P, 1, vec2T)
	b.op(32, inV4P, 1, vec4T)
	b.op(32, outV4P, 3, vec4T)

	b.op(59, posVar, inV2P, 1)
	b.op(59, colVar, inV4P, 1)
	b.op(59, fragVar, outV4P, 3)
	b.op(59, glPosVar, outV4P, 3)

	b.op(43, c0, floatT, 0)
	b.op(43, c1, floatT, 1065353216)

	b.op(54, mainFn, voidT, 0, voidFnT)
	b.op(247, entry)

	b.op(61, pv, vec2T, posVar)
	b.op(61, cv, vec4T, colVar)
	b.op(222, px, floatT, pv, 0)
	b.op(222, py, floatT, pv, 1)
	b.op(44, gp, vec4T, px, py, c0, c1)
	b.op(62, glPosVar, gp)
	b.op(62, fragVar, cv)
	b.op(253)
	b.op(56)

	return b.build()
}

func fragmentShaderSPIRV() []byte {
	b := newSPV()

	voidT := b.id()
	floatT := b.id()
	vec4T := b.id()
	voidFnT := b.id()
	inV4P := b.id()
	outV4P := b.id()

	fragVar := b.id()
	outVar := b.id()
	mainFn := b.id()
	entry := b.id()
	color := b.id()

	b.op(17, 1)
	b.op(14, 0, 0)
	b.opEntry(4, mainFn, "main", fragVar, outVar)
	b.op(16, mainFn, 0)
	b.op(19, 4, 0)
	b.opStr(5, mainFn, "main")
	b.opStr(5, fragVar, "fragColor")
	b.opStr(5, outVar, "outColor")
	b.op(71, fragVar, 30, 0)
	b.op(71, outVar, 30, 0)

	b.op(19, voidT)
	b.op(33, voidFnT, voidT)
	b.op(22, floatT, 32)
	b.op(23, vec4T, floatT, 4)
	b.op(32, inV4P, 1, vec4T)
	b.op(32, outV4P, 3, vec4T)

	b.op(59, fragVar, inV4P, 1)
	b.op(59, outVar, outV4P, 3)

	b.op(54, mainFn, voidT, 0, voidFnT)
	b.op(247, entry)

	b.op(61, color, vec4T, fragVar)
	b.op(62, outVar, color)

	b.op(253)
	b.op(56)

	return b.build()
}
