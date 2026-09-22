# vowifi-map
A comprehensive map of all worldwide Mobile Network Operators (MNOs) that have deployed the **VoWiFi**
service (also known as **Wi-Fi Calling**), which allows users to make voice calls over Wi-Fi. The map
is built by probing each operator's Evolved Packet Data Gateway (**ePDG**) publicly exposed
over the internet.

**The map is available online at**: https://domenicoverde.github.io/vowifi-map/.

## Overview
Every MNO worldwide, identified by its MCC/MNC pair, could deploy an ePDG as a gateway to offer
its own VoWiFi service. Each ePDG exposes an IKEv2-based interface over the internet
and has a standard DNS name of the form:

`epdg.epc.mnc<MNC>.mcc<MCC>.pub.3gppnetwork.org`

This project aims to build an interactive VoWiFi map through two main components: 
- **[`backend/`](backend)** resolves the standard ePDG DNS hostname, geolocates the
resulting IPs, and probes an IKEv2 IKE_SA_INIT request to check whether the VoWiFi service is 
actually active;

- **[`frontend/`](frontend)** turns the report produced by the backend into an interactive map.
The map is based on Leaflet and OpenStreetMap.

See each folder's own README for usage or further details.

## License
Distributed under the MIT License — see [LICENSE](LICENSE).

## Contributing
Bug reports and pull requests are welcome.