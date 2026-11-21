import { createElement } from "preact";
import {
  Anchor,
  BookOpen,
  Check,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  ChevronUp,
  Circle,
  CircleAlert,
  CircleQuestionMark,
  ClockArrowDown,
  ClockArrowUp,
  Download,
  Ellipsis,
  ExternalLink,
  Folder,
  FolderMinus,
  FolderPlus,
  Globe,
  Layers,
  LogOut,
  Maximize2,
  Minimize2,
  Plus,
  RotateCw,
  Rss,
  Search,
  SlidersHorizontal,
  SquarePen,
  Star,
  Trash,
  Upload,
  X,
} from "lucide";
import type { IconNode, SVGProps } from "lucide";

const icon = (n: IconNode, a: SVGProps = {}): IconNode => n.map(([t, p]) => [t, { ...p, ...a }]);

const ListAssorted: IconNode = [
  ["path", { d: "M15 5H3", key: "qmgq9p" }],
  ["path", { d: "M3 12h18", key: "1i2n21" }],
  ["path", { d: "M9 19H3", key: "s61nz1" }],
];

const icons = {
  anchor: Anchor,
  "alert-circle": CircleAlert,
  assorted: ListAssorted,
  "book-open": BookOpen,
  check: Check,
  "chevron-down": ChevronDown,
  "chevron-left": ChevronLeft,
  "chevron-right": ChevronRight,
  "chevron-up": ChevronUp,
  circle: Circle,
  "circle-full": icon(Circle, { fill: "currentColor" }),
  download: Download,
  edit: SquarePen,
  "external-link": ExternalLink,
  folder: Folder,
  "folder-minus": FolderMinus,
  "folder-plus": FolderPlus,
  fullscreen: Maximize2,
  "fullscreen-exit": Minimize2,
  globe: Globe,
  "help-circle": CircleQuestionMark,
  layers: Layers,
  "log-out": LogOut,
  "more-horizontal": Ellipsis,
  plus: Plus,
  "rotate-cw": RotateCw,
  rss: Rss,
  search: Search,
  sliders: SlidersHorizontal,
  "sort-new-first": ClockArrowUp,
  "sort-old-first": ClockArrowDown,
  star: Star,
  "star-full": icon(Star, { fill: "currentColor" }),
  trash: Trash,
  upload: Upload,
  x: X,
} satisfies Record<string, IconNode>;

export type IconName = keyof typeof icons;

interface IconProps {
  name: IconName;
  small?: boolean;
  className?: string;
}

export default function Icon({ name, small, className = "" }: IconProps) {
  const cls = `c-icon ${small ? "is-small" : ""} ${className}`.trim();
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      className={cls}>
      {icons[name].map(([tag, attrs]) => createElement(tag, attrs))}
    </svg>
  );
}
