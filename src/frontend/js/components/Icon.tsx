import icons from "../icons";

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
  return <span className={cls} title={title} dangerouslySetInnerHTML={{ __html: svgHtml }} />;
}
