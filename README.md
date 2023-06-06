# NTU: No Type Username

Laziers of National Taiwan University need not to type login credentials again.

## Prerequisites

### Target Website

The website that you want to log-in is called "target website" (e.g. cool.ntu.edu.tw). All traffics destined to the target website 
is proxied by `NTU: No Type Username`.

### VPN

Since no authentication is supported now, you have better run this service in a VPN.

You also need to setup routing table in VPN to route traffics to the target website through the host running `NTU: No Type Username` so that
it can proxy the requests.

### NAT Table Setup

Run `scripts/add-nat-rules` to setup NAT tables that redirects all traffics to the target website to `NTU: No Type Username`. This is done by 
DNAT (destination NAT) entries in `iptables` that rewrite the destination to the address of `NTU: No Type Username`.
