import { Activity } from "lucide-react";
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
          aria-label={`Hourly event volume. ${points.map((p) => `${p.name}: ${p.count}`).join("; ")}`}
        >
          {buckets.map((p) => (
            <div
              key={p.name}
              title={`${displayTime(p.name)} UTC: ${p.count} events`}
              style={{
                height: `${(p.count / max) * 100}%`,
                minHeight: p.count ? 3 : 0,
              }}
            />
          ))}
        </div>
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
