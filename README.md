# bigmath

`ecosystem::bigmath` provides immutable, bounded arbitrary-precision rational
numbers and binary floating-point values in pure GoML. It depends on
`ecosystem::bigint = "0.1.0"`; decimal arithmetic is a separate base-10 model and
is not used here. The `examples/basic` example exercises the
public API and frozen `math/big` reference vectors.

## Rational

`Rational::new(BigInt, BigInt)` accepts a signed numerator and denominator. It
rejects zero denominators, moves any sign into the numerator, divides both parts
by their greatest common divisor, and stores zero as `0/1`. `from_parts` accepts
an unsigned denominator. `add`, `sub`, `mul`, `div`, `reciprocal`, `cmp`, `neg`,
`abs`, and exact integer conversion operate without floating-point
intermediates. `to_string` uses `n/d`, or `n` when the denominator is one.
`round_to_bigint(rounding)` converts directly to an integer using any of the six
`Rounding` modes below, without an intermediate binary float or loss of large
integer digits. HalfEven resolves exact halves to an even integer; HalfAway
resolves them away from zero. TowardZero truncates; Floor and Ceiling round
toward negative and positive infinity. `to_bigint_exact` continues to reject
fractions with a nonzero remainder.

Stored numerator and denominator magnitudes have at most 8192 bits each.
Arithmetic may use an intermediate of at most 16385 bits so cancellation and
reduction can produce an in-range result. Inputs beyond that intermediate
budget or results beyond the stored budget return `Error::LimitExceeded`.
Multiplication and division cross-cancel numerator/denominator factors before
multiplying. Since input rationals are canonical, the resulting products are
already reduced; this avoids constructing large factors that would immediately
cancel. Addition and subtraction first remove common denominator factors before
forming cross products, then cancel only against the original denominator GCD.
Equal denominators need only a numerator addition; coprime denominators produce
an already reduced result. Zero remains `0/1`, and final numerator/denominator
limits still apply. `reciprocal` swaps the already reduced magnitudes and
preserves the numerator's sign without repeating reduction. Zero still returns
`Error::DivisionByZero`. Comparison continues to use full cross products.

## Binary Float

`Context::new(precision, rounding)` selects 1–4096 significant **binary** bits
and one of `HalfEven`, `HalfAway`, `TowardZero`, `AwayFromZero`, `Floor`, and
`Ceiling`. `Context::standard()` uses 53 bits and `HalfEven`. Construction from
an integer or rational, `round`, `add`, `sub`, `mul`, `div`, and `sqrt` return an
`Outcome { value, inexact }`; each result is rounded once from exact rational
arithmetic. `inexact` is true only when the stored value differs from the exact
operation result. There is no ambient precision or rounding state.

Each nonzero `Float` stores a signed integer significand and a power-of-two
exponent, representing `significand × 2^exponent`. The significand has exactly
the context precision in bits. Its normalized most-significant-bit exponent is
in `-4096..4096`; a result outside that interval returns
`Error::ExponentOutOfRange` rather than silently overflowing, underflowing, or
creating infinity or a subnormal. Zero is canonical and unsigned.
An exact input just below the minimum exponent can still succeed when rounding
carries its significand into exponent -4096; the lower limit applies to the
rounded result. Results that remain below the limit still return an error.
`to_rational` recovers the exact stored value; `to_string` displays the binary representation
as `significand*2^exponent` when the exponent is nonzero. `sqrt` rejects negative
inputs, and division rejects zero. NaN, infinities, decimal parsing/formatting,
operator overloading, and compiler-enforced constant time are not provided.

The rational 8192-bit limit bounds stored Float operands and their exact
arithmetic inputs. Binary rounding may shift an intermediate by at most 8192
bits; square-root rounding may use about 20480-bit scratch integers. A
recoverable resource error may occur before a mathematical exponent error when
extreme operands create an exact intermediate outside that budget. Operations
return new values and leave their inputs unchanged.

## Verification

From the repository root, run `(cd ../verification && just ecosystem-test bigmath)`. Library tests
check reduction, signs, exact operations, midpoint and directional rounding,
square roots, limits, and error paths. The example uses independently
generated Go `math/big` reference vectors for exact rational results and
correctly rounded binary results.

## Development and examples

Requires GoML 0.1.56 or newer. The `examples/basic/` example shares the root manifest. From the library root, run:

```sh
goml run --example basic
goml test
goml verify --timeout 300s
```

`goml test` builds the example and runs its tests. `goml verify` repeats the example checks as an independent module against an isolated registry snapshot. `(cd ../verification && just ecosystem-test bigmath)` also retains the library-specific smoke and compatibility checks.
