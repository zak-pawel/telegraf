# StatsD Input Plugin

This service plugin gathers metrics from a [Statsd][statsd] server.

⭐ Telegraf v0.2.0
🏷️ applications
💻 all

[statsd]: https://github.com/statsd/statsd

## Service Input <!-- @/docs/includes/service_input.md -->

This plugin is a service input. Normal plugins gather metrics determined by the
interval setting. Service plugins start a service to listen and wait for
metrics or events to occur. Service plugins have two key differences from
normal plugins:

1. The global or plugin specific `interval` setting may not apply
2. The CLI options of `--test`, `--test-wait`, and `--once` may not produce
   output for this plugin

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md#plugins

## Configuration

```toml @sample.conf
# Statsd Server
[[inputs.statsd]]
  ## Protocol, must be "tcp", "udp4", "udp6" or "udp" (default=udp)
  protocol = "udp"

  ## MaxTCPConnection - applicable when protocol is set to tcp (default=250)
  max_tcp_connections = 250

  ## Enable TCP keep alive probes (default=false)
  tcp_keep_alive = false

  ## Specifies the keep-alive period for an active network connection.
  ## Only applies to TCP sockets and will be ignored if tcp_keep_alive is false.
  ## Defaults to the OS configuration.
  # tcp_keep_alive_period = "2h"

  ## Address and port to host UDP listener on
  service_address = ":8125"

  ## The following configuration options control when telegraf clears it's cache
  ## of previous values. If set to false, then telegraf will only clear it's
  ## cache when the daemon is restarted.
  ## Reset gauges every interval (default=true)
  delete_gauges = true
  ## Reset counters every interval (default=true)
  delete_counters = true
  ## Reset sets every interval (default=true)
  delete_sets = true
  ## Reset timings & histograms every interval (default=true)
  delete_timings = true

  ## Enable aggregation temporality adds temporality=delta or temporality=commulative tag, and
  ## start_time field, which adds the start time of the metric accumulation.
  ## You should use this when using OpenTelemetry output.
  # enable_aggregation_temporality = false

  ## Percentiles to calculate for timing & histogram stats.
  percentiles = [50.0, 90.0, 99.0, 99.9, 99.95, 100.0]

  ## separator to use between elements of a statsd metric
  metric_separator = "_"

  ## Parses extensions to statsd in the datadog statsd format
  ## currently supports metrics and datadog tags.
  ## http://docs.datadoghq.com/guides/dogstatsd/
  datadog_extensions = false

  ## Parses distributions metric as specified in the datadog statsd format
  ## https://docs.datadoghq.com/developers/metrics/types/?tab=distribution#definition
  datadog_distributions = false

  ## Keep or drop the container id as tag. Included as optional field
  ## in DogStatsD protocol v1.2 if source is running in Kubernetes
  ## https://docs.datadoghq.com/developers/dogstatsd/datagram_shell/?tab=metrics#dogstatsd-protocol-v12
  datadog_keep_container_tag = false

  ## Statsd data translation templates, more info can be read here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/TEMPLATE_PATTERN.md
  # templates = [
  #     "cpu.* measurement*"
  # ]

  ## Number of UDP messages allowed to queue up, once filled,
  ## the statsd server will start dropping packets
  allowed_pending_messages = 10000

  ## Number of worker threads used to parse the incoming messages.
  # number_workers_threads = 5

  ## Number of timing/histogram values to track per-measurement in the
  ## calculation of percentiles. Raising this limit increases the accuracy
  ## of percentiles but also increases the memory usage and cpu time.
  percentile_limit = 1000

  ## Maximum socket buffer size in bytes, once the buffer fills up, metrics
  ## will start dropping.  Defaults to the OS default.
  # read_buffer_size = 65535

  ## Max duration (TTL) for each metric to stay cached/reported without being updated.
  # max_ttl = "10h"

  ## Sanitize name method
  ## By default, telegraf will pass names directly as they are received.
  ## However, upstream statsd now does sanitization of names which can be
  ## enabled by using the "upstream" method option. This option will a) replace
  ## white space with '_', replace '/' with '-', and remove characters not
  ## matching 'a-zA-Z_\-0-9\.;='.
  #sanitize_name_method = ""

  ## Replace dots (.) with underscore (_) and dashes (-) with
  ## double underscore (__) in metric names.
  # convert_names = false

  ## Convert all numeric counters to float
  ## Enabling this would ensure that both counters and guages are both emitted
  ## as floats.
  # float_counters = false

  ## Emit timings `metric_<name>_count` field as float, the same as all other
  ## histogram fields
  # float_timings = false

  ## Emit sets as float
  # float_sets = false
```

## Description

The statsd plugin is a special type of plugin which runs a backgrounded statsd
listener service while telegraf is running.

The format of the statsd messages was based on the format described in the
original [etsy
statsd](https://github.com/etsy/statsd/blob/master/docs/metric_types.md)
implementation. In short, the telegraf statsd listener will accept:

- Gauges
  - `users.current.den001.myapp:32|g` <- standard
  - `users.current.den001.myapp:+10|g` <- additive
  - `users.current.den001.myapp:-10|g`
- Counters
  - `deploys.test.myservice:1|c` <- increments by 1
  - `deploys.test.myservice:101|c` <- increments by 101
  - `deploys.test.myservice:1|c|@0.1` <- with sample rate, increments by 10
- Sets
  - `users.unique:101|s`
  - `users.unique:101|s`
  - `users.unique:102|s` <- would result in a count of 2 for `users.unique`
- Timings & Histograms
  - `load.time:320|ms`
  - `load.time.nanoseconds:1|h`
  - `load.time:200|ms|@0.1` <- sampled 1/10 of the time
- Distributions
  - `load.time:320|d`
  - `load.time.nanoseconds:1|d`
  - `load.time:200|d|@0.1` <- sampled 1/10 of the time

It is possible to omit repetitive names and merge individual stats into a
single line by separating them with additional colons:

- `users.current.den001.myapp:32|g:+10|g:-10|g`
- `deploys.test.myservice:1|c:101|c:1|c|@0.1`
- `users.unique:101|s:101|s:102|s`
- `load.time:320|ms:200|ms|@0.1`

This also allows for mixed types in a single line:

- `foo:1|c:200|ms`

The string `foo:1|c:200|ms` is internally split into two individual metrics
`foo:1|c` and `foo:200|ms` which are added to the aggregator separately.

## Influx Statsd

In order to take advantage of InfluxDB's tagging system, we have made a couple
additions to the standard statsd protocol. First, you can specify
tags in a manner similar to the line-protocol, like this:

```shell
users.current,service=payroll,region=us-west:32|g
```

<!-- TODO Second, you can specify multiple fields within a measurement:

```
current.users,service=payroll,server=host01:west=10,east=10,central=2,south=10|g
```

-->

## Metrics

Meta:

- tags: `metric_type=<gauge|set|counter|timing|histogram>`

Outputted measurements will depend entirely on the measurements that the user
sends, but here is a brief rundown of what you can expect to find from each
metric type:

- Gauges
  - Gauges are a constant data type. They are not subject to averaging, and they
    don’t change unless you change them. That is, once you set a gauge value, it
    will be a flat line on the graph until you change it again.
- Counters
  - Counters are the most basic type. They are treated as a count of a type of
    event. They will continually increase unless you set `delete_counters=true`.
- Sets
  - Sets count the number of unique values passed to a key. For example, you
    could count the number of users accessing your system using `users:<user_id>|s`.
    No matter how many times the same user_id is sent, the count will only increase
    by 1.
- Timings & Histograms
  - Timers are meant to track how long something took. They are an invaluable
    tool for tracking application performance.
  - The following aggregate measurements are made for timers:
    - `statsd_<name>_lower`: The lower bound is the lowest value statsd saw
        for that stat during that interval.
    - `statsd_<name>_upper`: The upper bound is the highest value statsd saw
        for that stat during that interval.
    - `statsd_<name>_mean`: The mean is the average of all values statsd saw
        for that stat during that interval.
    - `statsd_<name>_median`: The median is the middle of all values statsd saw
        for that stat during that interval.
    - `statsd_<name>_stddev`: The stddev is the sample standard deviation
        of all values statsd saw for that stat during that interval.
    - `statsd_<name>_sum`: The sum is the sample sum of all values statsd saw
        for that stat during that interval.
    - `statsd_<name>_count`: The count is the number of timings statsd saw
        for that stat during that interval. It is not averaged.
    - `statsd_<name>_percentile_<P>` The `Pth` percentile is a value x such
        that `P%` of all the values statsd saw for that stat during that time
        period are below x. The most common value that people use for `P` is the
        `90`, this is a great number to try to optimize.
- Distributions
  - The Distribution metric represents the global statistical distribution of a set of values calculated across your entire distributed infrastructure in one time interval. A Distribution can be used to instrument logical objects, like services, independently from the underlying hosts.
  - Unlike the Histogram metric type, which aggregates on the Agent during a given time interval, a Distribution metric sends all the raw data during a time interval.

## Plugin arguments

- **protocol** string: Protocol used in listener - tcp or udp options
- **max_tcp_connections** []int: Maximum number of concurrent TCP connections
to allow. Used when protocol is set to tcp.
- **tcp_keep_alive** boolean: Enable TCP keep alive probes
- **tcp_keep_alive_period** duration: Specifies the keep-alive period for an active network connection
- **service_address** string: Address to listen for statsd UDP packets on
- **delete_gauges** boolean: Delete gauges on every collection interval
- **delete_counters** boolean: Delete counters on every collection interval
- **delete_sets** boolean: Delete set counters on every collection interval
- **delete_timings** boolean: Delete timings on every collection interval
- **percentiles** []int: Percentiles to calculate for timing & histogram stats
- **allowed_pending_messages** integer: Number of messages allowed to queue up
waiting to be processed. When this fills, messages will be dropped and logged.
- **percentile_limit** integer: Number of timing/histogram values to track
per-measurement in the calculation of percentiles. Raising this limit increases
the accuracy of percentiles but also increases the memory usage and cpu time.
- **templates** []string: Templates for transforming statsd buckets into influx
measurements and tags.
- **parse_data_dog_tags** boolean: Enable parsing of tags in DataDog's dogstatsd format (<http://docs.datadoghq.com/guides/dogstatsd/>)
- **datadog_extensions** boolean: Enable parsing of DataDog's extensions to dogstatsd format (<http://docs.datadoghq.com/guides/dogstatsd/>)
- **datadog_distributions** boolean: Enable parsing of the Distribution metric in DataDog's dogstatsd format (<https://docs.datadoghq.com/developers/metrics/types/?tab=distribution#definition>)
- **datadog_keep_container_tag** boolean: Keep or drop the container id as tag. Included as optional field in DogStatsD protocol v1.2 if source is running in Kubernetes.
- **max_ttl** config.Duration: Max duration (TTL) for each metric to stay cached/reported without being updated.

## Statsd bucket -> InfluxDB line-protocol Templates

The plugin supports specifying templates for transforming statsd buckets into
InfluxDB measurement names and tags. The templates have a _measurement_ keyword,
which can be used to specify parts of the bucket that are to be used in the
measurement name. Other words in the template are used as tag names. For
example, the following template:

```toml
templates = [
    "measurement.measurement.region"
]
```

would result in the following transformation:

```shell
cpu.load.us-west:100|g
=> cpu_load,region=us-west 100
```

Users can also filter the template to use based on the name of the bucket,
using glob matching, like so:

```toml
templates = [
    "cpu.* measurement.measurement.region",
    "mem.* measurement.measurement.host"
]
```

which would result in the following transformation:

```shell
cpu.load.us-west:100|g
=> cpu_load,region=us-west 100

mem.cached.localhost:256|g
=> mem_cached,host=localhost 256
```

Consult the [Template Patterns](/docs/TEMPLATE_PATTERN.md) documentation for
additional details.

## Example Output
