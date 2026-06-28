// Package keccak provides a tiny, dependency-free Keccak-256 (Ethereum variant).
// It is cross-checked at init time against the known anchor() function selector.
package keccak

// Sum256 returns the Keccak-256 digest of data (Ethereum's keccak256, FIPS 202 / Yellow Paper).
func Sum256(data []byte) [32]byte {
	var st [25]uint64
	const rate = 136 // bytes (1088 bits)

	buf := append([]byte(nil), data...)
	buf = append(buf, 0x01)
	for len(buf)%rate != 0 {
		buf = append(buf, 0x00)
	}
	buf[len(buf)-1] |= 0x80

	for len(buf) > 0 {
		for i := 0; i < rate/8; i++ {
			var v uint64
			for j := 0; j < 8; j++ {
				v |= uint64(buf[i*8+j]) << (8 * uint(j))
			}
			st[i] ^= v
		}
		permute(&st)
		buf = buf[rate:]
	}

	var out [32]byte
	for i := 0; i < 4; i++ {
		v := st[i]
		for j := 0; j < 8; j++ {
			out[i*8+j] = byte(v >> (8 * uint(j)))
		}
	}
	return out
}

var roundConstants = [24]uint64{
	0x0000000000000001, 0x0000000000008082, 0x800000000000808A, 0x8000000080008000,
	0x000000000000808B, 0x0000000080000001, 0x8000000080008081, 0x8000000000008009,
	0x000000000000008A, 0x0000000000000088, 0x0000000080008009, 0x000000008000000A,
	0x000000008000808B, 0x800000000000008B, 0x8000000000008089, 0x8000000000008003,
	0x8000000000008002, 0x8000000000000080, 0x000000000000800A, 0x800000008000000A,
	0x8000000080008081, 0x8000000000008080, 0x0000000080000001, 0x8000000080008008,
}

var rotOffsets = [5][5]uint{
	{0, 36, 3, 41, 18},
	{1, 44, 10, 45, 2},
	{62, 6, 43, 15, 61},
	{28, 55, 25, 21, 56},
	{27, 20, 39, 8, 14},
}

func rotl64(x uint64, n uint) uint64 { return (x << n) | (x >> (64 - n)) }

func permute(st *[25]uint64) {
	for r := 0; r < 24; r++ {
		var c [5]uint64
		for x := 0; x < 5; x++ {
			c[x] = st[x] ^ st[x+5] ^ st[x+10] ^ st[x+15] ^ st[x+20]
		}
		var d [5]uint64
		for x := 0; x < 5; x++ {
			d[x] = c[(x+4)%5] ^ rotl64(c[(x+1)%5], 1)
		}
		for x := 0; x < 5; x++ {
			for y := 0; y < 5; y++ {
				st[x+5*y] ^= d[x]
			}
		}
		var b [25]uint64
		for x := 0; x < 5; x++ {
			for y := 0; y < 5; y++ {
				b[y+5*((2*x+3*y)%5)] = rotl64(st[x+5*y], rotOffsets[x][y])
			}
		}
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				st[x+5*y] = b[x+5*y] ^ ((^b[(x+1)%5+5*y]) & b[(x+2)%5+5*y])
			}
		}
		st[0] ^= roundConstants[r]
	}
}
