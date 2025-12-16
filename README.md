# Visual Traceroute Tool
Hopings is an open-source visual traceroute tool.

It performs classic ICMP traceroute and presents the result as a clear,
hop-by-hop network path with latency visualization and country information.

![Hopings screenshot](https://net.u00.io/hopings/screen00.png)

## Advantages

- **Clear visual path**  
  Presents traceroute results as a clean, readable hop-by-hop path.  
  Latency, countries, and unreachable hops are visible at a glance.

- **Honest ICMP traceroute**  
  Uses classic ICMP traceroute without shortcuts or synthetic assumptions.  
  The output reflects exactly what the network returns.

- **Latency at a glance**  
  Round-trip times are visually categorized, making long-distance jumps and slow segments easy to spot instantly.

- **Country-aware routing**  
  Each hop is enriched with country information and flags, helping to understand geographic routing decisions.

- **Readable handling of missing replies**  
  Filtered or rate-limited hops are shown explicitly, without misleading placeholders or fake timings.

- **Single executable**  
  Distributed as a single standalone executable.  
  No installers, no runtime dependencies, no configuration files.

- **Fast and lightweight**  
  Starts instantly and runs without background services, telemetry, or unnecessary overhead.

- **Focused by design**  
  Hopings does one thing — traceroute — and does it well.  
  No maps, no animations, no clutter.

- **Open source and transparent**  
  Fully open-source. Easy to inspect, learn from, and extend.  
  Built with correctness and clarity as first principles.
