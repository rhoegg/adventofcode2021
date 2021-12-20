using System.Collections.Generic;
using System.Linq;

namespace Beacons
{
    public class Program
    {
        static void Main(string[] args)
        {
            List<List<Vector3>> scanners = new List<List<Vector3>>();
            foreach (var scannerData in File.ReadAllText("input.txt").Split("\n\n"))
            {
                var lines = scannerData.Split("\n");
                var scanner =  lines.First().Split(" ")[2];
                Console.WriteLine($"Loaded scanner {scanner}");
                scanners.Add(new List<Vector3>(lines.Skip(1).Select(line => Vector3.Parse(line))));
            }

            var allBeaconsInWorldSpace = new HashSet<Vector3>();
            var scannerPositions = new List<Vector3>();

            allBeaconsInWorldSpace.UnionWith(scanners[0]);

            var q = new Queue<int>(scanners.Skip(1).Select((s, i) => i + 1));
            while (q.Count() > 0)
            {
                var scannerIndex = q.Dequeue();
                Console.WriteLine("\nScanner " + scannerIndex);
                var scanner = scanners[scannerIndex];

                var worldVectors = BeaconVectors(allBeaconsInWorldSpace);
                // inefficient, we recalc all the old ones each time around

                var rotated = scanner;
                var rotations = "DRRRDRRRDRRRDRDDRRRDRRRDRRR";
                var matched = false;
                foreach (var direction in rotations)
                {
                    switch (direction)
                    {
                        case 'R':
                            rotated = rotated.Select(v => v.RollRight()).ToList();
                            break;
                        case 'U':
                            rotated = rotated.Select(v => v.PitchUp()).ToList();
                            break;
                        case 'D':
                            rotated = rotated.Select(v => v.PitchDown()).ToList();
                            break;
                        case 'B':
                            rotated = rotated.Select(v => v.PitchBackward()).ToList();
                            break;
                        case '_':
                            break;
                    }
                    Console.Write(direction);
                    var matchedBeacons = worldVectors.Join(
                        BeaconVectors(rotated),
                        w => w.local,
                        s => s.local,
                        (w, s) => (world: w, scanner: s));
                    if (12 <= matchedBeacons.Count())
                    {
                        var offset = matchedBeacons.Select(match => match.scanner.peer.ToLocal(match.world.peer)).First();
                        Console.WriteLine($"*\nScanner offset {offset}");
                        scannerPositions.Add(offset);
                        allBeaconsInWorldSpace.UnionWith(rotated.Select(b => b.Add(offset)));
                        matched = true;
                        break;
                    }
                    else
                    {
                        if (matchedBeacons.Count() > 0)
                        {
                            Console.Write("?");
                        }
                    }
                }
                if (! matched)
                {
                    Console.WriteLine($"\nDid not match {scannerIndex} yet");
                    q.Enqueue(scannerIndex);
                }
            }
            Console.WriteLine($"\nWorld contains {allBeaconsInWorldSpace.Count} beacons");

            foreach (int i in scannerPositions.SelectMany(p1 => scannerPositions.Select(p2 => 
                Math.Abs(p1.X - p2.X) + Math.Abs(p1.Y - p2.Y) + Math.Abs(p1.Z - p2.Z))))
                {
                    Console.WriteLine(i);
                }
        }

        public static IEnumerable<(Vector3 origin, Vector3 peer, Vector3 local)> BeaconVectors(IEnumerable<Vector3> scanner)
        {
            return scanner.SelectMany(v0 => scanner.Select(v1 => (origin: v0, peer: v1, local: v0.ToLocal(v1))))
                .Where(pair => pair.origin != pair.peer);
        }

    }
}