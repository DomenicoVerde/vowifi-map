export interface EpdgMarker {
  mcc: string;
  mnc: string;
  iso: string;
  country: string;
  code: string;
  network: string;
  ip: string;
  latitude: number;
  longitude: number;
  response: "yes" | "no";
}
