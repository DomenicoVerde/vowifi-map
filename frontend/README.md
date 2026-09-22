# vowifi-map (frontend)

A React + Leaflet single-page app that turns the ePDG scan results produced by the
[backend](../backend) into an interactive world map: one marker per ePDG, colored by whether it
answered the IKE_SA_INIT probe (green or amber), clustered at low zoom levels; the UI supports
both light and dark themes.

No backend/API calls at runtime: the app fetches a single static JSON file
(`public/data/markers.json`), generated at build time from the backend's `results.json`.

## Technology Stack
- **React 19**, **TypeScript** — UI components and app state;
- **Vite** — dev server (HMR) and production bundler;
- **Leaflet** / **react-leaflet** / **react-leaflet-cluster** — the map, markers, and
clustering utilities (basemap tiles are served by [CARTO](https://carto.com/attributions)
on top of OpenStreetMap);
- **oxlint** — linting.

## Requirements
- Node.js ≥ 20, npm.
- A `results.json` report produced by the [backend](../backend) scanner, present at
  `../backend/data/results.json` (relative to this folder).
- A CARTO basemap API key: sign up for free at [carto.com](https://carto.com), create an API key
  under **Your account → API keys**, and restrict it to the domains that will serve the map (e.g.,
  `domenicoverde.github.io`). Copy `.env.example` to `.env` and
  set `VITE_CARTO_API_KEY` to the key — Vite picks it up at dev/build time.

Note: `.env` is ignored by git, so the key will never be committed.

## Usage

```bash
npm install
npm run dev       # start the dev server with hot reload
npm run build     # type-check + production build into dist/
npm run preview   # serve the dist/ build locally, to sanity-check before deploying
```

## Deployment

To deploy the map manually via [`gh-pages`](https://www.npmjs.com/package/gh-pages):

```bash
npm run build
npx gh-pages -d dist
```

GitHub Pages is configured under Settings → Pages → **Deploy from a branch** → `gh-pages` →
`/ (root)`.
