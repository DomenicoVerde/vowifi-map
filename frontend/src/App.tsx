import { useEffect, useState } from "react";
import "./App.css";
import { Header } from "./components/Header";
import { MapView } from "./components/MapView";
import { ThemeToggle } from "./components/ThemeToggle";
import { useTheme } from "./useTheme";
import type { EpdgMarker } from "./types";

function App() {
  const [markers, setMarkers] = useState<EpdgMarker[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [theme, toggleTheme] = useTheme();

  useEffect(() => {
    fetch(`${import.meta.env.BASE_URL}data/markers.json`)
      .then((res) => {
        if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
        return res.json();
      })
      .then(setMarkers)
      .catch((err) => setError(String(err)));
  }, []);

  return (
    <div className="app">
      <Header markers={markers ?? []} />
      <ThemeToggle theme={theme} onToggle={toggleTheme} />
      {error && <div className="banner banner--error">Failed to load scan data: {error}</div>}
      {markers && <MapView markers={markers} theme={theme} />}
    </div>
  );
}

export default App;
