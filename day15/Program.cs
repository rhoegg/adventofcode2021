using System.Drawing;

var chitonRoute = new ChitonRoute(File.ReadAllLines("example1.txt"));
chitonRoute.Start = new Point(0, 0);
chitonRoute.End = new Point(9, 9);
// Console.WriteLine("Lowest total risk");
int risk = chitonRoute.LowestTotalRisk();
Console.WriteLine("Lowest total risk " + risk);
// Console.WriteLine(chitonRoute.PathsToEnd().Count());