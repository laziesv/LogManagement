# React frontend

React + Vite + TypeScript + Tailwind CSS, with Lucide icons. All application data comes from the Go API; there is no mock fallback.

## Structure

```text
src/
  App.tsx                Session gate: loading, login, or workspace
  main.tsx               React entry point
  layouts/               Workspace shell and page selection
  pages/                 Login, Overview, Logs, Alerts, Sources
  components/
    common/              Brand and empty state
    layout/              Sidebar, topbar, header, notices, footer
    dashboard/           Timeline and rankings
    logs/                Filters, table, event detail dialog
  hooks/                 Session and workspace state/actions
  contexts/              Typed workspace context shared by pages/components
  services/              HTTP client and auth/log/alert API functions
  types/                 Shared domain and navigation types
  constants/             Navigation metadata and source names
  utils/                 Formatting, errors and synthetic sample generation
  styles/                Shared stylesheet
tests/                   Page-render regression checks
```

`useAuth` owns the session check. `Workspace` calls `useWorkspace` once and provides its state through `WorkspaceProvider`. Switching pages preserves the selected filters and loaded data. Pages render state and call actions; the services layer owns endpoints and request payloads. `services/client.ts` handles cookies, JSON and API errors.

API requests run from hooks/services, not from layout components. File upload and sample ingestion are workspace actions. `utils/samples.ts` is only used by the explicit demo buttons. The API remains responsible for authorization; hiding Viewer buttons is not a security boundary.

## Commands

```sh
npm ci
npm run dev
npm test
npm run build
```

Vite binds 127.0.0.1:5173 by default and proxies /api and /ingest to 127.0.0.1:3000. Set backend APP_ORIGIN=http://127.0.0.1:5173 during local development. Production Nginx serves built assets and proxies the same routes.

Tests use Node's built-in runner, TypeScript transpilation and React server rendering; no extra testing dependencies. They verify page rendering and Admin/Viewer controls. They do not replace browser interaction or live API integration tests.

This refactor keeps the existing local tab navigation, styling, endpoints and JSON formats. No new URL router or state-management library was introduced.
