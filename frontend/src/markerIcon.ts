import L from "leaflet";
import type { EpdgMarker } from "./types";

// A pulsing green dot for EPDGs that answered the IKE_SA_INIT probe — a
// "live" indicator. EPDGs that merely resolved but didn't respond get a
// static amber dot instead, so the two are visually distinct at a glance.
const respondingIcon = L.divIcon({
  className: "epdg-marker",
  html: '<span class="epdg-marker__pulse"></span><span class="epdg-marker__dot"></span>',
  iconSize: [14, 14],
  iconAnchor: [7, 7],
  popupAnchor: [0, -6],
});

const silentIcon = L.divIcon({
  className: "epdg-marker epdg-marker--no",
  html: '<span class="epdg-marker__dot"></span>',
  iconSize: [14, 14],
  iconAnchor: [7, 7],
  popupAnchor: [0, -6],
});

export function epdgIcon(response: EpdgMarker["response"]): L.DivIcon {
  return response === "yes" ? respondingIcon : silentIcon;
}
