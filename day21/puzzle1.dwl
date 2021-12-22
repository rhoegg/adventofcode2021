%dw 2.0
import divideBy from dw::core::Arrays
import mergeWith from dw::core::Objects
output application/json
var starts = [8, 3]
var winningScore = 1000
var diceRolls = (0 to 999) as Array map (($ mod 100) + 1) divideBy 3
var board = (1 to 10) as Array
var turns = diceRolls map (rolls, i) -> {
    rolls: rolls,
    i: i mod 2,
    roll: i + 1
}
var gameStart = starts map (position, i) -> {player: i, position: position, score: 0, diceRolls: 0}

fun playOneTurn(gameState, rolls) = do
    {
        var newPosition = ((gameState.position - 1 + sum(rolls)) mod 10) + 1
        var newScore = gameState.score + board[newPosition - 1]
        ---
        if (gameState.score >= winningScore) gameState else {
            position: newPosition,
            score: newScore,
            diceRolls: gameState.diceRolls + sizeOf(rolls)
        }
    }

var gameOver = turns reduce ( (turn, gameState = gameStart) ->
    if (max(gameState map $.score) >= winningScore)
        gameState
    else [
        gameState[turn.i] mergeWith playOneTurn(gameState[turn.i], turn.rolls),
        gameState[turn.i + 1 mod 2]
    ] orderBy $.player
)


---
{
    gameOver: gameOver,
    rolls: sum(gameOver map $.diceRolls),
    lowScore: min(gameOver map $.score)
}


