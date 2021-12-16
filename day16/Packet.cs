using System.Collections.Generic;

class Packet
{
    private string unparsed;
    public Packet(string binary)
    {
        unparsed = binary;

        Version = ReadInt(3);
        Type = ReadInt(3);
        Payload = ReadPayload();
    }

    public static Packet ParseHex(string hex)
    {
        return new Packet(HexToBinary(hex));
    }

    public int Version { get; init; }
    public int Type { get; init; }

    public string Unparsed {
        get { return unparsed; }
    }
    
    public Payload Payload { get; init; }

    string ReadBits(int count)
    {
        if (unparsed.Length < count) throw new Exception("Unexpected end of input");
        string bits = unparsed.Substring(0, count);
        unparsed = unparsed.Substring(count);
        return bits;
    }

    int ReadInt(int bits) 
    {
        return Convert.ToInt32(ReadBits(bits), 2);
    }

    long ReadLong(int bits) 
    {
        return Convert.ToInt64(ReadBits(bits), 2);
    }

    bool ReadBool()
    {
        return ReadBits(1) == "1";
    }

    Payload ReadPayload()
    {
        switch (Type)
        {
            case 4: return ReadLiteral();
            default: return ReadOperator();
        }
    }

    Payload ReadLiteral()
    {
        var l = new Literal();
        var more = true;
        while (more)
        {
            more = ReadBool();
            l.ApplyNibble(ReadInt(4));
        }
        return l;
    }

    Operator ReadOperator()
    {
        var op = new Operator(SelectAggregator());
        var lengthRepresentsBits = ReadBool(); 
        if (lengthRepresentsBits)
        {
            var subPacketCount = ReadInt(11);
            for (int i = 0; i < subPacketCount; i++)
            {
                var unparsedLength = unparsed.Length;
                var p = new Packet(unparsed);
                op.ApplyPacket(p);
                ReadBits(unparsedLength - p.Unparsed.Length); // consume the bits
            }
        }
        else
        {
            var contentLength = ReadInt(15);
            var subBits = ReadBits(contentLength);

            while (subBits.Length > 0)
            {
                var p = new Packet(subBits);
                op.ApplyPacket(p);
                subBits = p.Unparsed;
            }
        }

        return op;
    }

    Func<IEnumerable<long>, long> SelectAggregator()
    {
        switch (Type)
        {
            case 0: return Enumerable.Sum;
            case 1: return ints => { return ints.Aggregate((product, factor) => product * factor); };
            case 2: return Enumerable.Min;
            case 3: return Enumerable.Max;
            case 5: return ints => { return ints.ElementAt(0) > ints.ElementAt(1) ? 1 : 0; };
            case 6: return ints => { return ints.ElementAt(0) < ints.ElementAt(1) ? 1 : 0; };
            case 7: return ints => { return ints.ElementAt(0) == ints.ElementAt(1) ? 1 : 0; };
            default: return _ => { return 0; };
        }
    }

    static string HexToBinary(string hex)
    {
        return String.Join(
            String.Empty,
            hex.Select(
                c => Convert.ToString(Convert.ToInt32(c.ToString(), 16), 2).PadLeft(4, '0')
            )
        );
    }   

    public int SumOfVersions()
    {
        return Version + Payload.SumOfVersions();
    }
}

interface Payload {
    int SumOfVersions();
    long Value { get; }
}

class Literal : Payload
{
    long val = 0;
    public void ApplyNibble(int nibble)
    {
        val = val * 16 + nibble;
    }

    public int SumOfVersions()
    {
        return 0;
    }

    public long Value {
        get {
            return val;
        }
    }
}

class Operator : Payload
{
    List<Packet> packets = new List<Packet>();

    public Operator(Func<IEnumerable<long>, long> aggregator)
    {
        Aggregator = aggregator;
    }

    public Func<IEnumerable<long>, long> Aggregator { get; init; }

    public void ApplyPacket(Packet p)
    {
        packets.Add(p);
    }

    public int SumOfVersions()
    {
        return packets.Sum(p => p.SumOfVersions());
    }

    public long Value {
        get {
            return Aggregator(packets.Select(p => p.Payload.Value));
        }
    }
}