# React frontend

Frontend ใช้ React + Vite + TypeScript + Tailwind CSS และ Lucide icons ข้อมูลทั้งหมดมาจาก Go API ไม่มี mock fallback ใน production flow

## โครงสร้าง

```text
src/
  App.tsx                ตัวกลางเลือกสถานะ loading, login หรือ workspace
  main.tsx               จุดเริ่มต้น React
  layouts/               โครงหน้า workspace และการเลือกหน้า
  pages/                 หน้า Login, Overview, Logs, Alerts, Sources
  components/
    common/              brand และ empty state
    layout/              sidebar, topbar, header, notices, footer
    dashboard/           timeline และ ranking charts
    logs/                filters, table, event detail dialog
  hooks/                 session และ workspace state/actions
  contexts/              workspace context ที่ใช้ร่วมกันระหว่าง pages/components
  services/              HTTP client และฟังก์ชันเรียก auth/log/alert API
  types/                 shared domain และ navigation types
  constants/             navigation metadata และ source names
  utils/                 formatting, errors และ synthetic sample generation
  styles/                stylesheet กลาง
 tests/                  regression checks สำหรับ page rendering
```

`useAuth` รับผิดชอบการตรวจ session ส่วน `Workspace` เรียก `useWorkspace` ครั้งเดียวแล้วส่ง state ผ่าน `WorkspaceProvider` การสลับหน้าไม่ล้าง filter และข้อมูลที่โหลดไว้ Pages มีหน้าที่ render state และเรียก actions ส่วน services layer ดูแล endpoints และ payloads

`services/client.ts` จัดการ cookies, JSON และ API errors

API requests อยู่ใน hooks/services ไม่อยู่ใน layout components การ upload file และ sample ingestion เป็น workspace actions ส่วน `utils/samples.ts` ใช้เฉพาะปุ่ม demo ที่ user กดเองเท่านั้น

Backend ยังเป็นจุด enforce authorization จริง การซ่อนปุ่มของ Viewer ใน UI เป็นแค่ UX ไม่ใช่ security boundary

## คำสั่งที่ใช้บ่อย

```sh
npm ci
npm run dev
npm test
npm run build
```

Vite bind ที่ `127.0.0.1:5173` เป็นค่าเริ่มต้น และ proxy `/api` กับ `/ingest` ไป `127.0.0.1:3000` ระหว่าง local development ให้ตั้ง backend:

```text
APP_ORIGIN=http://127.0.0.1:5173
```

Production ใช้ Nginx serve built assets และ proxy routes เดิม

Tests ใช้ Node built-in runner, TypeScript transpilation และ React server rendering โดยไม่เพิ่ม testing dependency ใหม่ Tests ตรวจ page rendering และปุ่ม/สิทธิ์ Admin/Viewer แต่ไม่แทน browser interaction หรือ live API integration tests

โครงนี้คง tab navigation, styling, endpoints และ JSON formats เดิมไว้ ไม่เพิ่ม URL router หรือ state-management library ใหม่
