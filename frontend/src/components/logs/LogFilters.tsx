import { Search } from "lucide-react";
import { sources } from "../../constants/sources";
import { useWorkspaceContext } from "../../contexts/WorkspaceContext";

export function LogFilters() {
  const {
    source,
    setSource,
    hours,
    setHours,
    query,
    setQuery,
    setSearch,
    setOffset,
    customFrom,
    setCustomFrom,
    customTo,
    setCustomTo,
  } = useWorkspaceContext();
  return (
    <div className="filters">
      <form
        onSubmit={(e) => {
          e.preventDefault();
          setOffset(0);
          setSearch(query);
        }}
      >
        <Search size={17} />
        <input
          aria-label="Search logs"
          placeholder="Search events, IPs, users…"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button type="submit">Search</button>
      </form>
      <label>
        <span className="sr-only">Source</span>
        <select
          value={source}
          onChange={(e) => {
            setSource(e.target.value);
            setOffset(0);
          }}
        >
          <option value="">All sources</option>
          {sources.map((s) => (
            <option key={s}>{s}</option>
          ))}
        </select>
      </label>
      <label>
        <span className="sr-only">Time range</span>
        <select
          value={hours}
          onChange={(e) => {
            setHours(e.target.value);
            setOffset(0);
          }}
        >
          <option value="1">Last hour</option>
          <option value="24">Last 24 hours</option>
          <option value="168">Last 7 days</option>
          <option value="720">Last 30 days</option>
          <option value="custom">Custom range</option>
        </select>
      </label>
      {hours === "custom" && (
        <>
          <label>
            From
            <input
              type="datetime-local"
              value={customFrom}
              onChange={(e) => {
                setCustomFrom(e.target.value);
                setOffset(0);
              }}
            />
          </label>
          <label>
            To
            <input
              type="datetime-local"
              value={customTo}
              onChange={(e) => {
                setCustomTo(e.target.value);
                setOffset(0);
              }}
            />
          </label>
        </>
      )}
    </div>
  );
}
