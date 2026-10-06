import { useState, useRef, useEffect, useCallback } from "preact/hooks";
import type { ComponentChildren, Ref } from "preact";
import Tooltip from "./Tooltip";

export interface DropdownHandle {
  hide: () => void;
  show: () => void;
  toggle: () => void;
}

interface DropdownProps {
  button: ComponentChildren;
  toggleClass: string;
  drop: "right" | "center";
  title: string;
  disabled?: boolean;
  children: ComponentChildren;
  ref?: Ref<DropdownHandle>;
}

export default function Dropdown({
  button,
  toggleClass,
  drop,
  title,
  disabled,
  children,
  ref,
}: DropdownProps) {
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  const hide = useCallback(() => setOpen(false), []);
  const show = useCallback(() => setOpen(true), []);
  const toggle = useCallback(() => setOpen(prev => !prev), []);

  if (ref) {
    if (typeof ref === "function") {
      ref({ hide, show, toggle });
    } else {
      (ref as { current: DropdownHandle | null }).current = { hide, show, toggle };
    }
  }

  useEffect(() => {
    if (!open) return;

    const clickHandler = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      const menu = menuRef.current;
      if (!menu || !menu.contains(target)) {
        hide();
      } else if (target.closest(".c-dropdown-item") !== null) {
        hide();
      }
    };

    document.addEventListener("click", clickHandler);
    return () => document.removeEventListener("click", clickHandler);
  }, [open, hide]);

  const onButtonClick = (e: MouseEvent) => {
    e.stopPropagation();
    if (!disabled) toggle();
  };

  const menuCls = `c-dropdown position-absolute top-100 z-1 ${
    drop === "center" ? "start-50 translate-middle-x" : "end-0"
  }`;

  return (
    <div className="position-relative d-inline-flex flex-column" ref={containerRef}>
      <Tooltip label={title}>
        <button
          type="button"
          className={`c-button-pill ${toggleClass}`}
          onClick={onButtonClick}
          aria-label={title}
          disabled={disabled}>
          {button}
        </button>
      </Tooltip>
      {open && (
        <div ref={menuRef} className={menuCls}>
          {children}
        </div>
      )}
    </div>
  );
}
