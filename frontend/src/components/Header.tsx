import { useState } from "react";
import type { EpdgMarker } from "../types";

export function Header({ markers }: { markers: EpdgMarker[] }) {
  const [isOpen, setIsOpen] = useState(true);
  const operators = new Set(markers.map((m) => `${m.mcc}-${m.mnc}`)).size;
  const countries = new Set(markers.map((m) => m.iso)).size;
  const responding = markers.filter((m) => m.response === "yes").length;

  if (!isOpen) {
    return (
      <button className="header-toggle" onClick={() => setIsOpen(true)} aria-label="Show info panel">
        <span className="header-toggle__dot" />
        VoWiFi Map
      </button>
    );
  }

  return (
    <div className="header">
      <button className="header__close" onClick={() => setIsOpen(false)} aria-label="Hide info panel">
        ×
      </button>
      <h1>VoWiFi Availability Map</h1>
      <p className="header__subtitle">
        This map displays all mobile network operators (MNOs) worldwide that have deployed VoWiFi
        (also known as Wi-Fi Calling). The location of each ePDG is estimated based on the IP
        address associated with the DNS hostname:{" "}
        <code>epdg.epc.mnc&lt;MNC&gt;.mcc&lt;MCC&gt;.pub.3gppnetwork.org</code>.
      </p>
      <p className="header__legend-title">Marker color legend:</p>
      <ul className="header__legend">
        <li>
          <span className="epdg-popup__response--yes">Green</span>: The ePDG responded to an
          IKEv2 IKE_SA_INIT request.
        </li>
        <li>
          <span className="epdg-popup__response--no">Amber</span>: No response was received.
        </li>
      </ul>
      <div className="header__stats">
        <div className="stat">
          <span className="stat__value">{markers.length}</span>
          <span className="stat__label">EPDGs</span>
        </div>
        <div className="stat">
          <span className="stat__value">{responding}</span>
          <span className="stat__label">Responding</span>
        </div>
        <div className="stat">
          <span className="stat__value">{operators}</span>
          <span className="stat__label">Operators</span>
        </div>
        <div className="stat">
          <span className="stat__value">{countries}</span>
          <span className="stat__label">Countries</span>
        </div>
      </div>
    </div>
  );
}
