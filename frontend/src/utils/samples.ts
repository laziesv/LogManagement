export function sample(source: string) {
  return {
    "@timestamp": new Date().toISOString(),
    source,
    event_type:
      source === "ad"
        ? "LogonFailed"
        : source === "crowdstrike"
          ? "malware_detected"
          : source === "aws"
            ? "CreateUser"
            : "connection_allowed",
    severity: source === "crowdstrike" ? 8 : 3,
    src_ip: "203.0.113.7",
    user: "alice",
    host: "demo-host",
    action: source === "crowdstrike" ? "quarantine" : "allow",
  };
}
