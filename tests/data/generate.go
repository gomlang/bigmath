package main

import (
	"fmt"
	"math/big"
	"math/rand"
)

func randomInteger(random *rand.Rand, size int) *big.Int {
	data := make([]byte, size)
	_, _ = random.Read(data)
	data[0] |= 1
	value := new(big.Int).SetBytes(data)
	if random.Intn(2) == 0 {
		value.Neg(value)
	}
	return value
}

func randomRat(random *rand.Rand) *big.Rat {
	numerator := randomInteger(random, 16)
	denominator := new(big.Int).Abs(randomInteger(random, 10))
	return new(big.Rat).SetFrac(numerator, denominator)
}

func emit(op string, precision int, mode string, left, right, expected *big.Rat, inexact bool) {
	fmt.Printf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%t\n",
		op, precision, mode,
		left.Num(), left.Denom(), right.Num(), right.Denom(),
		expected.Num(), expected.Denom(), inexact)
}

func main() {
	random := rand.New(rand.NewSource(20260923))
	for sample := 0; sample < 64; sample++ {
		left := randomRat(random)
		right := randomRat(random)
		emit("radd", 0, "none", left, right, new(big.Rat).Add(left, right), false)
		emit("rsub", 0, "none", left, right, new(big.Rat).Sub(left, right), false)
		emit("rmul", 0, "none", left, right, new(big.Rat).Mul(left, right), false)
		emit("rdiv", 0, "none", left, right, new(big.Rat).Quo(left, right), false)
	}
	modes := []struct {
		name string
		mode big.RoundingMode
	}{
		{"half_even", big.ToNearestEven},
		{"half_away", big.ToNearestAway},
		{"toward_zero", big.ToZero},
		{"away_zero", big.AwayFromZero},
		{"floor", big.ToNegativeInf},
		{"ceiling", big.ToPositiveInf},
	}
	operations := []string{"round", "add", "sub", "mul", "div", "sqrt"}
	for _, mode := range modes {
		for _, precision := range []uint{3, 7, 16, 53, 113} {
			for sample := 0; sample < 8; sample++ {
				left := randomRat(random)
				right := randomRat(random)
				operation := operations[sample%len(operations)]
				if operation == "sqrt" {
					left = new(big.Rat).Abs(left)
				}
				leftFloat := new(big.Float).SetPrec(precision).SetMode(mode.mode).SetRat(left)
				rightFloat := new(big.Float).SetPrec(precision).SetMode(mode.mode).SetRat(right)
				result := new(big.Float).SetPrec(precision).SetMode(mode.mode)
				switch operation {
				case "round":
					result.SetRat(left)
				case "add":
					result.Add(leftFloat, rightFloat)
				case "sub":
					result.Sub(leftFloat, rightFloat)
				case "mul":
					result.Mul(leftFloat, rightFloat)
				case "div":
					result.Quo(leftFloat, rightFloat)
				case "sqrt":
					result.Sqrt(leftFloat)
				}
				exact, _ := result.Rat(nil)
				inexact := result.Acc() != big.Exact
				if operation == "sqrt" {
					operand, _ := leftFloat.Rat(nil)
					inexact = new(big.Rat).Mul(exact, exact).Cmp(operand) != 0
				}
				emit(operation, int(precision), mode.name, left, right, exact, inexact)
			}
		}
	}
}
