%dw 2.0
import some, every, drop, firstWith from dw::core::Arrays
output application/json

var drawnNumbers = [26,38,2,15,36,8,12,46,88,72,32,35,64,19,5,66,20,52,74,3,59,94,45,56,0,6,67,24,97,50,92,93,84,65,71,90,96,21,87,75,58,82,14,53,95,27,49,69,16,89,37,13,1,81,60,79,51,18,48,33,42,63,39,34,62,55,47,54,23,83,77,9,70,68,85,86,91,41,4,61,78,31,22,76,40,17,30,98,44,25,80,73,11,28,7,99,29,57,43,10]

var boards = day4 
    splitBy '\n\n' 
        map ($ splitBy '\n'
            map (trim($) splitBy /\s+/
                map ($ as Number)))

var gameBoards = boards map {
    board: $,
    marked: $ map [false, false, false, false, false]
}

fun transpose(a) = a reduce (row, cols = []) -> 
    if (cols == []) 
        row map [$] 
    else 
        cols map (r, i) -> r ++ [row[i]]

fun boardWins(board) = 
    (board.marked some (row) -> (row every $))
    or
    (transpose(board.marked) some (row) -> (row every $))

fun playBingo(remainingNumbers, boards) = do {
    var marked = boards map {
        board: $.board,
        marked: mark(remainingNumbers[0], $)
    }
    var winner = (marked firstWith (b) -> boardWins(b))
    ---
    winner then {winner: winner, lastNumber: remainingNumbers[0]}
    default (
        if (isEmpty(remainingNumbers)) 
            {winner: null, lastNumber: remainingNumbers[0]}
        else 
            playBingo(remainingNumbers drop 1, marked)
    )
}

fun mark(n, board) = board.board map (row, rowIndex) -> 
    row map (cell, colIndex) -> board.marked[rowIndex][colIndex] or (cell == n)

fun unmarked(board) = 
    board.board 
        map (row, rowIndex) ->
            (
                row reduce (cell, acc = {row: [], colIndex: 0}) ->
                    {
                        row: acc.row ++ [if (board.marked[rowIndex][acc.colIndex]) 0 else cell],
                        colIndex: acc.colIndex + 1
                    }
            ).row


var bingoResult = playBingo(drawnNumbers, gameBoards)
var puzzle1Inputs = {
    unmarkedSum: sum(unmarked(bingoResult.winner) map sum($)),
    lastNumber: bingoResult.lastNumber,
    bingoResult: bingoResult,
}
---
{
    inputs: puzzle1Inputs,
    result: puzzle1Inputs.unmarkedSum * 94
}

