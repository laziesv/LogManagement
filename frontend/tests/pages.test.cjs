const assert = require("node:assert/strict");
const fs = require("node:fs");
const test = require("node:test");
const ts = require("typescript");
const React = require("react");
const { renderToStaticMarkup } = require("react-dom/server");

// Compile source modules only in this test process, without generated source files.
for (const extension of [".ts", ".tsx"]) {
  require.extensions[extension] = (module, filename) => {
    const output = ts.transpileModule(fs.readFileSync(filename, "utf8"), {
      compilerOptions: {
        module: ts.ModuleKind.CommonJS,
        target: ts.ScriptTarget.ES2022,
        jsx: ts.JsxEmit.ReactJSX,
        esModuleInterop: true,
      },
      fileName: filename,
    });
    module._compile(output.outputText, filename);
  };
}
const { WorkspaceProvider } = require("../src/contexts/WorkspaceContext.tsx");
const { SourcesPage } = require("../src/pages/SourcesPage.tsx");
const { AlertsPage } = require("../src/pages/AlertsPage.tsx");
const { OverviewPage } = require("../src/pages/OverviewPage.tsx");
const { LogsPage } = require("../src/pages/LogsPage.tsx");
const { LoginPage } = require("../src/pages/LoginPage.tsx");
const noop = () => {};
function state(admin = true) {
  return {
    user: {
      id: "u",
      email: "user@demo.local",
      role: admin ? "admin" : "viewer",
      tenant: "demo-a",
    },
    isAdmin: admin,
    busy: false,
    actionBusy: false,
    tab: "overview",
    total: 1,
    offset: 0,
    logs: [
      {
        id: 1,
        "@timestamp": "2026-09-15T08:00:00Z",
        source: "api",
        event_type: "login_failed",
        src_ip: "203.0.113.7",
        severity: 5,
        user: "alice",
      },
    ],
    stats: {
      total: 1,
      critical: 0,
      sources: 1,
      alerts: 1,
      timeline: [{ name: "2026-09-15T08:00:00Z", count: 1 }],
      top_ips: [],
      top_users: [],
      top_events: [],
    },
    alerts: [
      {
        id: 1,
        src_ip: "203.0.113.7",
        count: 5,
        created_at: "2026-09-15T08:00:00Z",
        status: "open",
      },
    ],
    rule: { enabled: true, threshold: 5, window_minutes: 5 },
    fileInput: { current: null },
    setRule: noop,
    saveRule: noop,
    acknowledge: noop,
    demo: noop,
    setSelected: noop,
    setTab: noop,
    setOffset: noop,
  };
}
function render(Page, value) {
  return renderToStaticMarkup(
    React.createElement(
      WorkspaceProvider,
      { value },
      React.createElement(Page),
    ),
  );
}
test("Admin can see ingest controls; Viewer cannot", () => {
  const admin = render(SourcesPage, state(true));
  const viewer = render(SourcesPage, state(false));
  assert.match(admin, /Send demo batch/);
  assert.match(admin, /Choose JSON file/);
  assert.doesNotMatch(viewer, /Send demo batch|Choose JSON file|Send sample/);
  assert.match(viewer, /AWS CloudTrail/);
});
test("Viewer can read alerts but cannot save a rule or acknowledge", () => {
  const admin = render(AlertsPage, state(true));
  const viewer = render(AlertsPage, state(false));
  assert.match(admin, /Save rule/);
  assert.match(admin, /Acknowledge/);
  assert.match(viewer, /203\.0\.113\.7/);
  assert.doesNotMatch(viewer, /Save rule|Acknowledge/);
  assert.match(viewer, /disabled=""/);
});
test("Overview renders loaded events and empty states without changing page contracts", () => {
  const loaded = render(OverviewPage, state());
  assert.match(loaded, /Event activity/);
  assert.match(loaded, /login_failed/);
  const empty = render(OverviewPage, {
    ...state(),
    logs: [],
    stats: null,
    total: 0,
  });
  assert.match(empty, /No events in this view/);
  assert.match(empty, /Your timeline starts with the first event/);
});
test("Log page renders event details entry and pagination", () => {
  const html = render(LogsPage, { ...state(), tab: "logs" });
  assert.match(html, /Inspect log 1/);
  assert.match(html, /Previous/);
  assert.match(html, /Next/);
});
test("Login keeps a visible connection error and labeled credentials", () => {
  const html = renderToStaticMarkup(
    React.createElement(LoginPage, {
      onLogin: noop,
      retry: noop,
      connectionError: "Offline",
    }),
  );
  assert.match(html, /API unavailable: Offline/);
  assert.match(html, /Email address/);
  assert.match(html, /type="password"/);
  assert.match(html, /Retry/);
});
