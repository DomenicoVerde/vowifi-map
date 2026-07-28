import type { EpdgMarker } from "../types";
import { isoToFlagEmoji } from "../flagEmoji";

export function EpdgPopup({ marker }: { marker: EpdgMarker }) {
  return (
    <div className="epdg-popup">
      <div className="epdg-popup__title">
        <span className="epdg-popup__flag">{isoToFlagEmoji(marker.iso)}</span>
        <span>{marker.network}</span>
      </div>
      <div className="epdg-popup__subtitle">
        {marker.country} (+{marker.code})
      </div>

      <dl className="epdg-popup__grid">
        <dt>ePDG IP Address:</dt>
        <dd>
          <code>{marker.ip}</code>
        </dd>

        <dt>MCC:</dt>
        <dd>{marker.mcc}</dd>

        <dt>MNC:</dt>
        <dd>{marker.mnc}</dd>

        <dt>Coordinates:</dt>
        <dd>
          {marker.latitude.toFixed(3)}, {marker.longitude.toFixed(3)}
        </dd>

        <dt>IKEv2 Reply:</dt>
        <dd className={`epdg-popup__response--${marker.response}`}>
          {marker.response === "yes" ? "Yes" : "No"}
        </dd>
      </dl>
    </div>
  );
}
