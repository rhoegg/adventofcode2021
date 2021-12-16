%dw 2.0
output application/json
import drop from dw::core::Arrays

var puzzleInputParts = (example1 splitBy '\n\n')
var polymerTemplate = puzzleInputParts[0]
var pairInsertionRules = 
    puzzleInputParts[1] splitBy '\n' 
        map ($ splitBy ' -> ')
        map { pair: $[0], insertion: $[1] }
var insertions = {(pairInsertionRules map { ($.pair): $.insertion})}

fun pairs(template) = do {
    var elements = template splitBy ""
    ---
    elements zip (elements drop 1 ++ [""])
        map ($ joinBy "")
}

fun step(template) = 
    pairs(template)
        map {
            pair: $,
            insertion: insertions[$] default ""
        }
        map ($.pair[0] ++ $.insertion)
        joinBy ""

fun steps(template, numSteps) =
    numSteps match {
        case 0 -> template
        case 1 -> step(template)
        case n if (n > 1) -> steps(step(template), n-1)
    }
---
steps(polymerTemplate, 10) splitBy "" groupBy $ pluck sizeOf($) orderBy $
