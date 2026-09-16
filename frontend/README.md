# Frontend

Frontend ของ Logdesk ใช้ React + TypeScript + Vite + Tailwind CSS และ Lucide icons สำหรับ UI

ข้อมูลทั้งหมดมาจาก Go API ไม่มี mock fallback ใน production flow

## คำสั่งที่ใช้บ่อย

```sh
npm ci
npm run dev
npm test
npm run build
```

ระหว่าง local development Vite bind ที่ `127.0.0.1:5173` และ proxy `/api` กับ `/ingest` ไป backend ที่ `127.0.0.1:3000`

ตั้งค่า backend สำหรับ dev

```text
APP_ORIGIN=http://127.0.0.1:5173
```

Production ใช้ Nginx serve built assets และ proxy request ไป backend

## โครงสร้างสำคัญ

```text
src/
  App.tsx                 เลือกสถานะ loading, login หรือ workspace
  main.tsx                entry point ของ React
  layouts/                layout หลักของ workspace
  pages/                  Login, Overview, Logs, Alerts, Sources
  components/common/      brand, empty state และ component ใช้ซ้ำ
  components/layout/      sidebar, topbar, header, notices, footer
  components/dashboard/   timeline และ ranking charts
  components/logs/        filters, table, event detail dialog
  hooks/                  session และ workspace state/actions
  contexts/               workspace context
  services/               HTTP client และ API functions
  types/                  shared domain types
  constants/              navigation metadata และ source names
  utils/                  formatting, error handling และ demo samples
  styles/                 stylesheet กลาง
tests/                    regression checks สำหรับ rendering
```

## การทำงานของ state

- `useAuth` ตรวจ session และสถานะ login
- `useWorkspace` โหลดข้อมูล dashboard/logs/alerts และเก็บ filter state
- `WorkspaceProvider` ส่ง state/actions ให้ pages
- การเปลี่ยนหน้าไม่ล้าง filter และข้อมูลที่โหลดไว้
- API request อยู่ใน hooks/services ไม่อยู่ใน layout components

## สิทธิ์ใน UI

Frontend ซ่อนปุ่มบางอย่างสำหรับ Viewer เพื่อ UX เท่านั้น ส่วน security จริง enforce ที่ backend เสมอ

- Admin: ingest/import, แก้ alert rule, acknowledge alert
- Viewer: อ่าน dashboard/logs/alerts เฉพาะ tenant ของตัวเอง

Tenant filter ถูกถอดออกจาก UI เพราะ tenant ถูก enforce จาก session/API key ที่ backend

## Auto refresh

Overview รีเฟรชข้อมูลอัตโนมัติ และยังมีปุ่ม Refresh สำหรับกดเอง ส่วน Log Explorer ใช้การค้นหา/filter และปุ่ม Refresh ตามการใช้งานปัจจุบัน

## Tests

Frontend tests ใช้ Node built-in runner, TypeScript transpilation และ React server rendering เพื่อเช็ก page rendering และสิทธิ์ปุ่ม Admin/Viewer

Tests เหล่านี้ไม่แทน browser interaction หรือ live API integration tests
