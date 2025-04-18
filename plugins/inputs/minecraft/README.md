# Minecraft Input Plugin

This plugin collects score metrics from a [Minecraft][minecraft] server using
the RCON protocol.

> [!NOTE]
> This plugin supports Minecraft Java Edition versions 1.11 - 1.14. When using
> a version earlier than 1.13, be aware that the values for some criteria has
> changed and need to be modified.

⭐ Telegraf v1.4.0
🏷️ server
💻 all

[minecraft]: https://www.minecraft.net/

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md#plugins

## Configuration

```toml @sample.conf
# Collects scores from a Minecraft server's scoreboard using the RCON protocol
[[inputs.minecraft]]
  ## Address of the Minecraft server.
  # server = "localhost"

  ## Server RCON Port.
  # port = "25575"

  ## Server RCON Password.
  password = ""

  ## Uncomment to remove deprecated metric components.
  # tagdrop = ["server"]
```

### Server Setup

Enable [RCON][rcon] on the Minecraft server and add the following to your
[`server.properties`][propfile] file:

```conf
enable-rcon=true
rcon.password=<your password>
rcon.port=<1-65535>
```

Scoreboard [objectives][objectives] must be added using the server console for
the plugin to collect. These can be added in game by players with op status,
from the server console, or over an RCON connection.

When getting started pick an easy to test objective.  This command will add an
objective that counts the number of times a player has jumped:

```sh
/scoreboard objectives add jumps minecraft.custom:minecraft.jump
```

Once a player has triggered the event they will be added to the scoreboard,
you can then list all players with recorded scores:

```sh
/scoreboard players list
```

View the current scores with a command, substituting your player name:

```sh
/scoreboard players list Etho
```

[rcon]: http://wiki.vg/RCON
[propfile]: https://minecraft.gamepedia.com/Server.properties
[objectives]: https://minecraft.gamepedia.com/Scoreboard#Objectives

## Metrics

- minecraft
  - tags:
    - player
    - port (port of the server)
    - server (hostname:port, deprecated in 1.11; use `source` and `port` tags)
    - source (hostname of the server)
  - fields:
    - `<objective_name>` (integer, count)

## Example Output

```text
minecraft,player=notch,source=127.0.0.1,port=25575 jumps=178i 1498261397000000000
minecraft,player=dinnerbone,source=127.0.0.1,port=25575 deaths=1i,jumps=1999i,cow_kills=1i 1498261397000000000
minecraft,player=jeb,source=127.0.0.1,port=25575 d_pickaxe=1i,damage_dealt=80i,d_sword=2i,hunger=20i,health=20i,kills=1i,level=33i,jumps=264i,armor=15i 1498261397000000000
```
