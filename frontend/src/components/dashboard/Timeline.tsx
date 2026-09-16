import { Activity } from "lucide-react";
import { useState } from "react";
import type { Count } from "../../types/domain";
import { displayTime, number } from "../../utils/format";
import { Empty } from "../common/Empty";

export function Timeline({
  points,
  loading,
}: {
  points: Count[];
  loading: boolean;
}) {
  const [hovered, setHovered] = useState<number | null>(null);
  const [focused, setFocused] = useState<number | null>(null);
  if (!points.length)
    return (
      <Empty
        icon={Activity}
        title={
          loading
            ? "Loading activity…"
            : "Your timeline starts with the first event"
        }
        text="Import a file or send a sample from Data sources to get started."
      />
    );
  const start = new Date(points[0].name).getTime(),
    end = new Date(points.at(-1)!.name).getTime(),
    byHour = new Map(points.map((p) => [p.name, p.count]));
  const count = Math.min(
      744,
      Math.max(1, Math.round((end - start) / 3600000) + 1),
    ),
    buckets: Count[] = Array.from({ length: count }, (_, i) => {
      const name = new Date(start + i * 3600000)
        .toISOString()
        .replace(".000Z", "Z");
      return { name, count: byHour.get(name) || 0 };
    });
  const max = Math.max(...buckets.map((p) => p.count), 1);
  const cursor = hovered ?? focused;
  const active = cursor === null ? null : Math.min(cursor, buckets.length - 1);
  const selected = active === null ? null : buckets[active];
  return (
    <div className="timeline">
      <div className="chart-y">
        <span>{number(max)}</span>
        <span>{number(Math.round(max / 2))}</span>
        <span>0</span>
      </div>
      <div className="plot">
        <div className="chart-grid" />
        <div
          className="bars"
          role="img"
          tabIndex={0}
          onFocus={() => setFocused(0)}
          onBlur={() => setFocused(null)}
          onPointerLeave={() => setHovered(null)}
          onKeyDown={(event) => {
            if (!["ArrowLeft", "ArrowRight", "Home", "End", "Escape"].includes(event.key)) return;
            event.preventDefault();
            setHovered(null);
            if (event.key === "Escape") {
              setFocused(null);
              return;
            }
            setFocused((current) => {
              if (event.key === "Home") return 0;
              if (event.key === "End") return buckets.length - 1;
              const step = event.key === "ArrowRight" ? 1 : -1;
              return Math.max(0, Math.min(buckets.length - 1, (current ?? 0) + step));
            });
          }}
          aria-label={`Use arrow keys to explore hourly event volume. ${points.map((p) => `${p.name}: ${p.count}`).join("; ")}`}
        >
          {buckets.map((p, index) => (
            <div
              key={p.name}
              className={active === index ? "bar-active" : ""}
              onPointerEnter={() => setHovered(index)}
              style={{
                height: `${(p.count / max) * 100}%`,
                minHeight: p.count ? 3 : 0,
              }}
            />
          ))}
        </div>
        <span className="sr-only" role="status">
          {focused !== null && selected
            ? `${displayTime(selected.name)} UTC: ${selected.count} events`
            : ""}
        </span>
        {selected && (
          <div className="chart-tooltip" aria-hidden="true">
            <span>{displayTime(selected.name)} UTC</span>
            <strong>{number(selected.count)} events</strong>
          </div>
        )}
        <div className="chart-x">
          <span>
            {new Date(start).toLocaleString(undefined, {
              timeZone: "UTC",
              month: "short",
              day: "numeric",
              hour: "2-digit",
            })}
          </span>
          <span>
            {new Date(end).toLocaleString(undefined, {
              timeZone: "UTC",
              month: "short",
              day: "numeric",
              hour: "2-digit",
            })}
          </span>
        </div>
      </div>
    </div>
  );
}
