# vowifi-scanner (backend)

A Go scanner that, for every Mobile Network Operator (MNO) worldwide, identified by the
MCC/MNC pair, resolves the standard DNS name of its Evolved Packet Data Gateway (**ePDG**)
`epdg.epc.mnc<MNC>.mcc<MCC>.pub.3gppnetwork.org`, geolocates the IPs it finds, and
attempts an IKE_SA_INIT request against every resulting IP to check whether the VoWiFi (or WiFi-Calling)
service is actually active. It produces a JSON report
(`results.json`) which is consumed and displayed by the frontend.

## Requirements

- Go 1.26;
- Two external datasets (see below), **not included in this repository**.

The scanner needs the following files, both shall be downloaded separately and placed under
the `backend/data/` folder. Create it yourself first with:

```bash
mkdir -p backend/data
```

| File | Content                                                    | Source | Expected Path                     |
|---|------------------------------------------------------------|---|-----------------------------------|
| `mcc-mnc.csv` | Worldwide MNO list (MCC, MNC, ISO, country)                | [mcc-mnc.net](https://mcc-mnc.net) | `backend/data/mcc-mnc.csv`        |
| `GeoLite2-City.mmdb` | Database used to geolocate ePDG IPs  (latitude, longitude) | [P3TERX/GeoLite.mmdb releases](https://github.com/P3TERX/GeoLite.mmdb/releases/) (City edition) | `backend/data/GeoLite2-City.mmdb` |

Without these two files the scanner fails immediately on startup
(`Failed to load operators` / `Failed to open GeoIP database`).

## Usage

Run with `backend/src` as the working directory (default paths are relative
to it):

```bash
cd backend/src
go run ./cmd
```
The report is written to `backend/data/results.json` (default folder — this is also
where the frontend's `sync-data.mjs` expects to find it). A full scan may take a few minutes depending on
the timeout values and the number of concurrent workers.

### Flags

| Flag | Default Value | Description |
|---|---------------|---|
| `-csv` | `../data/mcc-mnc.csv` | Path to the operators CSV |
| `-mmdb` | `../data/GeoLite2-City.mmdb` | Path to the GeoLite2-City DB |
| `-output` | `../data/results.json` | Path to the output JSON report |
| `-workers` | `50`          | Number of operators scanned concurrently |
| `-dns-timeout` | `3s`          | Timeout for the ePDG DNS resolution |
| `-ike-timeout` | `3s`          | Timeout waiting for an IKE_SA_INIT response |

## Output
The `results.json` contains a JSON array of operators; each operator has a list of resolved ePDGs.
The list is empty if the DNS name does not resolve at all.

```jsonc
{
  "mcc": "222", "mnc": "001", "iso": "it", "country": "Italy",
  "code": "39", "network": "TIM",
  "epdgs": [
    { "ip": "217.200.184.104", "latitude": 37.5697, "longitude": 14.91, "response": "yes" }
  ]
}
```
`latitude` and `longitude` are `null` if the GeoIP DB has no match for that
IP.

`response` could be `"yes"` or `"no"` depending on whether the ePDG answered or not to the
IKE_SA_INIT request, or `null` if it was never probed (e.g. due to a loopback/sinkhole
IP address in DNS resolution). 
