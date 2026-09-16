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

export function abnormalSamples(count = 5, ip = "192.0.2.150") {
  return Array.from({ length: count }, (_, index) => ({
    source: "api",
    event_type: "login_failed",
    severity: 8,
    src_ip: ip,
    user: "demo-admin",
    host: "demo-auth-server",
    action: "login",
    reason: "Repeated incorrect passwords (synthetic sample)",
    attempt: index + 1,
    _tags: ["sample", "suspicious-login"],
  }));
}
