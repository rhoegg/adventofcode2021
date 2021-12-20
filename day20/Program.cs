using System.Linq;

var inputData = File.ReadAllText("input.txt").Split("\n\n");
var enhancement = inputData[0].Select(c => c == '#').ToList();
var enhancementZero = enhancement[0];

var inputImage = inputData[1].Split("\n").Select(line => line.Select(c => c == '#').ToList()).ToList(); 
var enhanced = inputImage;

for (var i = 0; i < 25; i++)
{
    // 1
    inputImage = FrameImage(enhanced, false);
    enhanced = Enhance(enhancement, InputSquares(inputImage, false));
    Console.WriteLine(enhanced.SelectMany(row => row).Count( pixel => pixel));

    // 2
    inputImage = FrameImage(enhanced, enhancementZero);
    enhanced = Enhance(enhancement, InputSquares(inputImage, enhancementZero));
    Console.WriteLine(enhanced.SelectMany(row => row).Count( pixel => pixel));
}

Console.WriteLine(Visualize(enhanced));

List<List<bool>> FrameImage(List<List<bool>> image, bool frameValue)
{
    image = image.Select(row => row.Prepend(frameValue).Append(frameValue).ToList()).ToList();
    var frameRow = Enumerable.Repeat(frameValue, image[0].Count).ToList();
    return image.Prepend(frameRow).Append(frameRow).ToList();
}

List<List<List<bool>>> InputSquares(List<List<bool>> image, bool surroundValue)
{
    var framed = FrameImage(image, surroundValue);
    return image.Select((row, y) => 
        row.Select((pixel, x) =>
            framed.GetRange(y, 3).SelectMany(row => row.GetRange(x, 3)).ToList()
        ).ToList()
    ).ToList();
}

List<List<bool>> Enhance(List<bool> enhancement, List<List<List<bool>>> squares)
{
    return squares.Select(row =>
        row.Select(square => {
            var binary = string.Join("", square.Select(bit => bit ? "1" : "0"));
            //Console.WriteLine(binary + " " + Convert.ToInt32(binary, 2));
            return enhancement[Convert.ToInt32(binary, 2)];
        }).ToList()
    ).ToList();
}

string Visualize(List<List<bool>> image)
{
    return String.Join("\n", image.Select(row =>
        String.Join("", row.Select(bit => bit ? "#": "."))
    ));
}
