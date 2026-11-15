import icons from "../icons";
import Tooltip from "./Tooltip";

export type IconName = keyof typeof icons;

interface IconProps {
  name: IconName;
  small?: boolean;
  className?: string;
  title?: string;
}

export default function Icon({ name, small, className = "", title }: IconProps) {
  const svgHtml = (icons as Record<string, string>)[name] || "";
  const cls = `c-icon ${small ? "is-small" : ""} ${className}`.trim();
  const span = <span className={cls} dangerouslySetInnerHTML={{ __html: svgHtml }} />;
  if (!title) return span;
  return <Tooltip label={title}>{span}</Tooltip>;
}
