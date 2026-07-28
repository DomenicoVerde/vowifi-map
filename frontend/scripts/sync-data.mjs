// Reads the backend scan output and distills it into the flat, minimal
// marker list the map actually needs: one entry per EPDG that has
// coordinates (i.e. resolved to a real, geolocatable, non-loopback IP),
// whether or not it actually answered the IKE_SA_INIT probe — the map
// tells the two cases apart visually instead of filtering one out.
// Run automatically before dev/build (see package.json) so the site always
// ships whatever backend/data/results.json currently contains.
import { readFileSync, writeFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import path from "node:path";

const here = path.dirname(fileURLToPath(import.meta.url));
const sourcePath = path.resolve(here, "../../backend/data/results.json");
const outPath = path.resolve(here, "../public/data/markers.json");

const operators = JSON.parse(readFileSync(sourcePath, "utf-8"));

const markers = [];
for (const op of operators) {
  for (const epdg of op.epdgs) {
    if (epdg.latitude == null || epdg.longitude == null) continue;
    markers.push({
      mcc: op.mcc,
      mnc: op.mnc,
      iso: op.iso,
      country: op.country,
      code: op.code,
      network: op.network,
      ip: epdg.ip,
      latitude: epdg.latitude,
      longitude: epdg.longitude,
      response: epdg.response,
    });
  }
}

writeFileSync(outPath, JSON.stringify(markers));
console.log(`Wrote ${markers.length} EPDG markers to ${path.relative(process.cwd(), outPath)}`);
