%dw 2.0
import fromBinary from dw::core::Numbers
output application/json
var originalDiagnostics = day3 splitBy '\n'
var diagnosticBits = originalDiagnostics map (($ splitBy "") map $ as Number)
var lineCount = sizeOf(originalDiagnostics)
var bitSums = diagnosticBits reduce (d, counts = []) -> if (isEmpty(counts)) d else (d zip counts map sum($))
var gammaEpsilonBits = unzip((bitSums map if ($ / lineCount > 0.5) [1, 0] else [0, 1]))
fun bitsToDecimal(bits) = fromBinary((bits map $ as String) joinBy "")
var rates = {
    gamma: bitsToDecimal(gammaEpsilonBits[0]),
    epsilon: bitsToDecimal(gammaEpsilonBits[1])
}
fun bitAverage(diagnostics, position) = 
    sum(diagnostics map $[position] as Number) / sizeOf(diagnostics)
fun oxygenBit(diagnostics, position) = if(bitAverage(diagnostics, position) >= 0.5) 1 else 0
fun co2Bit(diagnostics, position) =    if(bitAverage(diagnostics, position) >= 0.5) 0 else 1
fun oxygen(diagnostics, position) = do {
    var bit = oxygenBit(diagnostics, position) as String
    var remaining = diagnostics filter ($[position] == bit)
    ---
    if (sizeOf(remaining) == 1) remaining[0] else oxygen(remaining, position + 1)
}
fun co2(diagnostics, position) = do {
    var bit = co2Bit(diagnostics, position) as String
    var remaining = diagnostics filter ($[position] == bit)
    ---
    if (sizeOf(remaining) == 1) remaining[0] else co2(remaining, position + 1)
}
var ratings = {
    oxygen: fromBinary(oxygen(originalDiagnostics, 0)),
    co2: fromBinary(co2(originalDiagnostics, 0))
}  
---
ratings.oxygen * ratings.co2


