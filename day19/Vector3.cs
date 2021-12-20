using System.Linq;

namespace Beacons
{
    public class Vector3
    {
        public Vector3(int x, int y, int z)
        {
            X = x;
            Y = y;
            Z = z;
        }

        public int X { get; init; }
        public int Y { get; init; }
        public int Z { get; init; }

        // aka SubtractFrom
        public Vector3 ToLocal(Vector3 peer)
        {
            // peer - this
            return new Vector3(peer.X - X, peer.Y - Y, peer.Z - Z);
        }

        public Vector3 Add(Vector3 v)
        {
            return new Vector3(v.X + X, v.Y + Y, v.Z + Z);
        }

        public Vector3 PitchBackward()
        {
            // rotation y1 = y cos angle - z sin angle
            // rotation z1 = y sin angle + z cos angle
            // 180 degrees
            // y = -y
            // z = -z;
            return new Vector3(X, -1 * Z, Y);
        }

        public Vector3 RollRight()
        {
            // rotation x1 = x cos angle - y sin angle
            // rotation y1 = x sin angle + y cos angle
            // 90 degrees, points move counterclockwise
            // x = -y
            // y = x
            return new Vector3(-1 * Y, X, Z);
        }

        public Vector3 PitchUp()
        {
            // rotation y1 = y cos angle - z sin angle
            // rotation z1 = y sin angle + z cos angle
            // 90 degrees
            // y = -z
            // z = y;
            return new Vector3(X, -1 * Z, Y);
        }

        public Vector3 PitchDown()
        {
            // rotation y1 = y cos angle - z sin angle
            // rotation z1 = y sin angle + z cos angle
            // 90 degrees
            // y = z
            // z = -y
            return new Vector3(X, Z, -1 * Y);
        }
        
        public override string ToString()
        {
            return $"({X},{Y},{Z})";
        }

        public override bool Equals(object obj)
        {
            if (! obj.GetType().IsAssignableFrom(typeof(Vector3))) return false;
            var other = obj as Vector3;
            return other.X == X && other.Y == Y && other.Z == Z;
        }

        public override int GetHashCode()
        {
            return X.GetHashCode() + 2 * Y.GetHashCode() + 3 * Z.GetHashCode();
        }

        public static Vector3 Parse(string coords)
        {
            var readings = coords.Split(",").Select(text => int.Parse(text)).ToArray();
            return new Vector3(readings[0], readings[1], readings[2]);
        }
    }
}