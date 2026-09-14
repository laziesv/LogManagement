import { Activity } from "lucide-react";

export function Empty({
  icon: Icon,
  title,
  text,
}: {
  icon: typeof Activity;
  title: string;
  text: string;
}) {
  return (
    <div className="empty">
      <Icon size={27} />
      <strong>{title}</strong>
      <p>{text}</p>
    </div>
  );
}
