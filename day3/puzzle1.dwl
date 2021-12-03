%dw 2.0
import fromBinary from dw::core::Numbers
output application/json
var diagnostics = day3 splitBy '\n'
var diagnosticBits = diagnostics map (($ splitBy "") map $ as Number)
var lineCount = sizeOf(diagnostics)
var bitSums = diagnosticBits reduce (d, counts = []) -> if (isEmpty(counts)) d else (d zip counts map sum($))
var gammaEpsilonBits = unzip((bitSums map if ($ / lineCount > 0.5) [1, 0] else [0, 1]))
fun bitsToDecimal(bits) = fromBinary((bits map $ as String) joinBy "")
var rates = {
    gamma: bitsToDecimal(gammaEpsilonBits[0]),
    epsilon: bitsToDecimal(gammaEpsilonBits[1])
}
---
rates.gamma * rates.epsilon
