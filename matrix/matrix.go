package matrix

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
	"time"
)

func GenerateMatrix(seed int64, size int) [][]float64 {
	m := make([][]float64, size)
	for r := 0; r < size; r++ {
		m[r] = make([]float64, size)
		for n := 0; n < size; n++ {
			var e = int64(0x4b72e682d)
			var a = int64(0x2675dcd22)
			m[r][n] = float64((e*seed + a) % 1000)
			seed = int64(m[r][n])
		}
	}
	return m
}

func CalculateHash(matrix [][]float64, size int) int64 {
	var sb strings.Builder
	for n := 0; n < size; n++ {
		for r := 0; r < size; r++ {
			sb.WriteString(fmt.Sprintf("%.0f", matrix[n][r]))
		}
	}
	hash := sha256.Sum256([]byte(sb.String()))
	bi := new(big.Int).SetBytes(hash[:])
	mod := new(big.Int).SetInt64(1e7)
	bi.Mod(bi, mod)
	return bi.Int64()
}

func Multiple(a, b [][]float64, size int) [][]float64 {
	result := make([][]float64, size)
	for i := 0; i < size; i++ {
		result[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			var sum float64
			for k := 0; k < size; k++ {
				sum += a[i][k] * b[k][j]
			}
			result[i][j] = sum
		}
	}
	return result
}

// CalculateResult CPU运算
func CalculateResult(seed1, seed2 int64, seedSize int) (float64, float64) {
	start := time.Now()
	matrix1 := GenerateMatrix(seed1, seedSize)
	matrix2 := GenerateMatrix(seed2, seedSize)
	res := Multiple(matrix1, matrix2, seedSize)
	hash := CalculateHash(res, seedSize)
	end := time.Now()
	result1 := float64(start.UnixMilli()) / float64(hash)
	result2 := float64(hash) / float64(end.Sub(start).Milliseconds())
	return result1, result2
}

// 统一入口：自动选择Go/Go GPU/Python Metal
func AutoCalculateResult(seed1, seed2 int64, seedSize int, gpuEnabled bool) (float64, float64, error) {
	// 默认CPU
	res1, res2 := CalculateResult(seed1, seed2, seedSize)
	return res1, res2, nil
}
