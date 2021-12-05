%dw 2.0
import mergeWith from dw::core::Objects
output application/json

var lines = day5
    splitBy  '\n'
        map (
            $ splitBy /\s*\->\s*/
            map ($ splitBy ',')
            map {x: $[0], y: $[1]}
        )
var puzzle1Lines = lines 
    filter (
        ($[0].x == $[1].x) or ($[0].y == $[1].y)
    ) 

fun pointsInLine(line) = 
    (line[0].x to line[1].x) 
        flatMap (x) ->
            (line[0].y to line[1].y)
                map (y) -> {x: x, y: y}

fun ventMapValues(line, ventMap) = {(
    pointsInLine(line) map (point) -> do {
        var k = (point.x ++ ":" ++ point.y)
        ---
        {
           k : 1 + (ventMap[k] default 0)
        }
    }
)}
---

sizeOf(
    (puzzle1Lines flatMap pointsInLine($))
    map ($.x ++ " " ++ $.y)
    groupBy $
    pluck $
    filter sizeOf($) > 1
)

