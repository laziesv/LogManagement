import type { Count } from "../../types/domain";
import { number } from "../../utils/format";

export function Ranking({ title, data }: { title: string; data: Count[] }) {
  const max = Math.max(...data.map((d) => d.count), 1);
  return (
    <section className="card ranking">
      <h2>{title}</h2>
      {!data.length && <p className="muted small">No events in this range.</p>}
      {data.map((d, i) => (
        <div className="rank-row" key={d.name}>
          <span className="rank-index">{i + 1}</span>
          <div>
            <div>
              <span title={d.name}>{d.name}</span>
              <strong>{number(d.count)}</strong>
            </div>
            <div className="rank-track">
              <span style={{ width: `${(d.count / max) * 100}%` }} />
            </div>
          </div>
        </div>
      ))}
    </section>
  );
}
