using System;
using System.IO;
using System.Collections.Generic;

namespace Day14
{
    // needs to give back character counts instead of the actual polymer for a pair and steps
    class InefficientProgram
    {
        static void Main(string[] args)
        {
            Polymerization polymerization;

            using (var reader = new StreamReader("example1.txt")) 
            {
                var polymerTemplate = reader.ReadLine();
                Console.WriteLine("Polymer Template: " + polymerTemplate);
                polymerization = new Polymerization(polymerTemplate);

                string s = "";
                while ((s = reader.ReadLine()) != null)
                {
                    if (s.Length > 0)
                    {
                        var injectionSpec = s.Split(" -> ");
                        polymerization.AddInsertion(injectionSpec[0], injectionSpec[1]);
                    }
                }
                reader.Close();
            }

            string expanded = polymerization.Expand(10);
            Dictionary<char, int> counts = new Dictionary<char, int>();
            foreach (char c in expanded)
            {
                if (counts.ContainsKey(c))
                {
                    counts[c] += 1;
                }
                else
                {
                    counts[c] = 1;
                }
            }
            foreach (var entry in counts)
            {
                Console.WriteLine(entry.Key + " " + entry.Value);
            }
        }
    }

    class Polymerization {
        readonly string polymerTemplate;
        Dictionary<string, string> insertions = new Dictionary<string, string>();
        Dictionary<int, Dictionary<string, string>> cache = 
            new Dictionary<int, Dictionary<string, string>>();
        
        public Polymerization(string template) 
        {
            this.polymerTemplate = template;
        }

        public void AddInsertion(string template, string insertion)
        {
            this.insertions[template] = insertion;
        }
        
        public void Remember(string template, int step, string expanded)
        {
            if (step <= 20) 
            {
                GetStepCache(step)[template] = expanded;
                return;
            } 
            // store in file
            File.WriteAllText("cache-" + template + "-" + step + ".txt", expanded);
        }

        public string GetCached(string template, int step)
        {
            if (step <= 20)
            {
                if (! GetStepCache(step).ContainsKey(template))
                {
                    return null;
                }
                return GetStepCache(step)[template];
            }
            // retrieve from file
            string filename = "cache-" + template + "-" + step + ".txt";
            if (File.Exists(filename))
            {
                return File.ReadAllText(filename);
            }
            return null;
        }

        public string Translate(string pair)
        {
            var insertion = "";
            if (insertions.ContainsKey(pair))
            {
                insertion = insertions[pair];
            }
            return pair[0] + insertion + pair[1];
        }

        private Dictionary<string, string> GetStepCache(int step)
        {
            if (! cache.ContainsKey(step))
            {
                cache[step] = new Dictionary<string, string>();
            }
            return cache[step];
        }

        public string Expand(int steps)
        {
            return this.Expand(polymerTemplate, steps);
        }
        public string Expand(string template, int steps)
        {
            // first letter then expand each pair after
            if (0 == steps) return template;

            string expanded = template[0].ToString();
            var pairs = new List<string>();
            for (int i = 1; i < template.Length; i++)
            {
                pairs.Add(template.Substring(i - 1, 2));
            }

            foreach (var pair in pairs)
            {
                string expandedPair = null;
                // int knownSteps = steps;
                // while( knownSteps > 0 && expandedPair == null)
                // {
                //     expandedPair = GetCached(pair, knownSteps);
                //     if (expandedPair == null) knownSteps--;
                // }
                // now we just have steps - knownSteps left
                expandedPair = Expand(Translate(pair), steps - 1);
                Remember(pair, steps, expandedPair);
                    
                expanded += expandedPair.Substring(1);
            }
            return expanded;
        }
    }
}
