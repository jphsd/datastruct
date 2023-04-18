package datastruct

// https://medium.com/@val_deleplace/7-ways-to-implement-a-bit-set-in-go-91650229b386 has a good
// breakdown of the various ways of implementing this.

// Bits contains a compact representation of booleans in a uint64.
type Bits []uint64

// NewBits allocates a new Bits with the given length rounded up to the nearest whole uint64
func NewBits(n int) Bits {
	return make(Bits, (n+63)/64)
}

// Get returns the state of bit i
func (b Bits) Get(i int) bool {
	p := i / 64
	j := uint(i % 64)
	return (b[p] & (uint64(1) << j)) != 0
}

// Set sets bit i to true
func (b Bits) Set(i int) {
	p := i / 64
	j := uint(i % 64)
	b[p] |= (uint64(1) << j)
}

// Reset clears bit i
func (b Bits) Clear(i int) {
	p := i / 64
	j := uint(i % 64)
	b[p] &= ^(uint64(1) << j)
}

// Slice returns the bit array as a slice of bool
func (b Bits) Slice() []bool {
	lw := len(b)
	res := make([]bool, lw*64)
	ip := 0
	for i := 0; i < lw; i++ {
		word := b[i]
		mask := uint64(1)
		for j := 0; j < 64; j++ {
			if (word & mask) != 0 {
				res[ip] = true
			}
			ip++
			mask <<= 1
		}
	}
	return res
}

// BitsFromSlice returns Bits initialized with the contents of the boolean slice
func BitsFromSlice(in []bool) Bits {
	lb := len(in)
	b := NewBits(lb)
	// Do it in chunks of 64
	n := lb / 64
	ip := 0
	for i := 0; i < n; i++ {
		var word uint64
		mask := uint64(1)
		for j := 0; j < 64; j++ {
			if in[ip] {
				word |= mask
			}
			ip++
			mask <<= 1
		}
		b[i] = word
	}
	// Remainder
	var word uint64
	mask := uint64(1)
	for j := 0; ip < lb; ip, j = ip+1, j+1 {
		if in[ip] {
			word |= mask
		}
		mask <<= 1
	}
	b[n] = word

	return b
}
