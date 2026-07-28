import "leaflet/dist/leaflet.css";
import "leaflet.markercluster/dist/MarkerCluster.css";
import "leaflet.markercluster/dist/MarkerCluster.Default.css";
import L from "leaflet";
import { MapContainer, TileLayer, Marker, Popup, ZoomControl } from "react-leaflet";
import MarkerClusterGroup from "react-leaflet-cluster";
import type { EpdgMarker } from "../types";
import type { Theme } from "../useTheme";
import { epdgIcon } from "../markerIcon";
import { EpdgPopup } from "./EpdgPopup";

const TILE_URL: Record<Theme, string> = {
  dark: "https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png",
  light: "https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png",
};
const TILE_ATTRIBUTION =
  '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>';

// Colors a cluster bubble by the mix of responding ("yes") vs silent ("no")
// EPDGs inside it, instead of a flat color — a cluster that's all green,
// all amber, or split down the middle should look different at a glance.
// leaflet.markercluster ships no type declarations, so the cluster param is
// typed structurally against just the one method this function calls.
function createClusterIcon(cluster: { getAllChildMarkers(): L.Marker[] }) {
  const childMarkers = cluster.getAllChildMarkers();
  const total = childMarkers.length;
  const respondingCount = childMarkers.filter((marker: L.Marker) => {
    const className = (marker.options.icon as L.DivIcon | undefined)?.options.className ?? "";
    return !String(className).includes("epdg-marker--no");
  }).length;

  const respondingShare = total > 0 ? Math.round((respondingCount / total) * 100) : 100;

  const html =
    `<div class="epdg-cluster" style="background: conic-gradient(var(--accent) 0% ${respondingShare}%, var(--warn) ${respondingShare}% 100%)">` +
    `<span>${total}</span></div>`;

  return L.divIcon({ html, className: "epdg-cluster-wrapper", iconSize: L.point(38, 38, true) });
}

export function MapView({ markers, theme }: { markers: EpdgMarker[]; theme: Theme }) {
  return (
    <MapContainer
      center={[20, 10]}
      zoom={3}
      minZoom={2}
      worldCopyJump
      zoomControl={false}
      className="map"
    >
      <TileLayer key={theme} url={TILE_URL[theme]} attribution={TILE_ATTRIBUTION} />
      <ZoomControl position="bottomright" />
      <MarkerClusterGroup
        chunkedLoading
        maxClusterRadius={40}
        spiderfyOnMaxZoom
        iconCreateFunction={createClusterIcon}
      >
        {markers.map((marker) => (
          <Marker
            key={`${marker.mcc}-${marker.mnc}-${marker.ip}`}
            position={[marker.latitude, marker.longitude]}
            icon={epdgIcon(marker.response)}
          >
            <Popup>
              <EpdgPopup marker={marker} />
            </Popup>
          </Marker>
        ))}
      </MarkerClusterGroup>
    </MapContainer>
  );
}
