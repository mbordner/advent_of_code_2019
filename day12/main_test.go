package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PrimeFactorization(t *testing.T) {
	factors := primeFactors(uint64(12))
	assert.Equal(t,[]uint64{2,2,3},factors)

	factors = primeFactors(uint64(49))
	assert.Equal(t,[]uint64{7,7},factors)
}

func Test_LCM(t *testing.T) {
	assert.Equal(t,uint64(12),lcm([]uint64{2,3,4}))
	assert.Equal(t,uint64(12),lcm1([]uint64{2,3,4}))
}