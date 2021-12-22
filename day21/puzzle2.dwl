// timed out on dwlang.fun, which is why I switched to go
%dw 2.0
import divideBy from dw::core::Arrays
import mergeWith from dw::core::Objects
output application/json
var starts = [8, 3]
var winningScore = 21
var dieFaces = (1 to 3)

var rollFrequency = (dieFaces flatMap (d1) ->
    dieFaces map [$, d1])
    flatMap ((prev) -> dieFaces map ((d3) -> prev ++ [d3]))
    map sum($)
    groupBy $
    mapObject({($$): sizeOf($)})

// (((roll + last.position - 1) mod 10) + 1)
fun movePawn(position, roll) = ((position - 1 + roll) mod 10) + 1
fun possibilities(states) = states flatMap (state) ->
    (rollFrequency pluck (count, roll) -> do {
        var newPosition = movePawn(state.positions[0], roll)
        ---
        {
            positions: [state.positions[1], newPosition],
            scores: [state.scores[1], state.scores[0] + newPosition],
            universes: state.universes + count
        }
    }
    )
fun winCounts(states) = do {
    var outcomes = possibilities(states)
    var indeterminate = outcomes filter (max($.scores) < winningScore)
    var nextMove = if (isEmpty(indeterminate)) [0, 0] else winCounts(indeterminate)
    ---
    [
        sum(outcomes filter $.scores[0] >= winningScore map $.universes) + nextMove[1],
        sum(outcomes filter $.scores[1] >= winningScore map $.universes) + nextMove[0]
    ]
}
---
winCounts([{
    positions: starts,
    scores: [10,9],
    universes: 1
}])